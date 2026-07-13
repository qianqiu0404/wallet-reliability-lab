package lab

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	store            *Store
	connector        Connector
	connectedEnabled bool
	connectedReason  string
}

func NewServerFromEnvironment() (*Server, error) {
	targets, err := ParseTargets(os.Getenv("LAB_SERVICE_TARGETS"))
	if err != nil {
		return nil, err
	}
	requested := strings.EqualFold(os.Getenv("LAB_ENABLE_CONNECTED"), "true")
	allLoopback := len(targets) > 0
	for _, target := range targets {
		allLoopback = allLoopback && IsLoopbackTarget(target)
	}
	reason := "connected mode is disabled; set LAB_ENABLE_CONNECTED=true with loopback targets"
	if requested && len(targets) == 0 {
		reason = "connected mode needs LAB_SERVICE_TARGETS"
	}
	if requested && len(targets) > 0 && !allLoopback {
		reason = "one or more service targets are not loopback"
	}
	if requested && allLoopback {
		reason = "loopback targets configured; reachability is checked at request time"
	}
	return &Server{
		store: NewStore(), connector: LoopbackConnector{Targets: targets, Timeout: 300 * time.Millisecond},
		connectedEnabled: requested && allLoopback, connectedReason: reason,
	}, nil
}

func NewServer(store *Store, connector Connector, connectedEnabled bool) *Server {
	return &Server{store: store, connector: connector, connectedEnabled: connectedEnabled, connectedReason: "configured by caller"}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/capabilities", s.handleCapabilities)
	mux.HandleFunc("GET /api/scenarios", s.handleScenarios)
	mux.HandleFunc("POST /api/runs", s.handleCreateRun)
	mux.HandleFunc("GET /api/runs/{id}", s.handleGetRun)
	mux.HandleFunc("GET /api/runs/{id}/events", s.handleEvents)
	mux.HandleFunc("POST /api/runs/{id}/actions", s.handleAction)
	return s.withHeaders(mux)
}

func (s *Server) withHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("X-Content-Type-Options", "nosniff")
		response.Header().Set("Cache-Control", "no-store")
		response.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		if request.Method == http.MethodOptions {
			response.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(response, request)
	})
}

func (s *Server) capabilities(ctx context.Context) Capabilities {
	capability := Capabilities{Simulation: true, ConnectedReason: s.connectedReason}
	if s.connector == nil || !s.connectedEnabled {
		return capability
	}
	statuses := s.connector.Probe(ctx)
	capability.Services = statuses
	capability.Connected = len(statuses) > 0
	for _, status := range statuses {
		capability.Connected = capability.Connected && status.Reachable
	}
	if !capability.Connected {
		capability.ConnectedReason = "connected mode configured, but one or more local services are unreachable"
	}
	return capability
}

func (s *Server) handleCapabilities(response http.ResponseWriter, request *http.Request) {
	writeJSON(response, http.StatusOK, s.capabilities(request.Context()))
}

func (s *Server) handleScenarios(response http.ResponseWriter, _ *http.Request) {
	writeJSON(response, http.StatusOK, BuiltinScenarios())
}

func (s *Server) handleCreateRun(response http.ResponseWriter, request *http.Request) {
	var input struct {
		ScenarioID string  `json:"scenarioId"`
		Mode       RunMode `json:"mode"`
	}
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(response, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if input.Mode == "" {
		input.Mode = ModeSimulation
	}
	if input.Mode != ModeSimulation && input.Mode != ModeConnected {
		writeError(response, http.StatusBadRequest, "mode must be simulation or connected")
		return
	}
	if input.Mode == ModeConnected && !s.capabilities(request.Context()).Connected {
		writeError(response, http.StatusForbidden, "connected mode is not available")
		return
	}
	run, err := s.store.Create(input.ScenarioID, input.Mode)
	if err != nil {
		writeError(response, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(response, http.StatusCreated, run)
}

func (s *Server) handleGetRun(response http.ResponseWriter, request *http.Request) {
	run, err := s.store.Get(request.PathValue("id"))
	if err != nil {
		s.writeStoreError(response, err)
		return
	}
	writeJSON(response, http.StatusOK, run)
}

func (s *Server) handleAction(response http.ResponseWriter, request *http.Request) {
	var input struct {
		Action string `json:"action"`
		Speed  int    `json:"speed"`
	}
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(response, http.StatusBadRequest, "invalid JSON body")
		return
	}
	id := request.PathValue("id")
	var run RunSnapshot
	var err error
	switch input.Action {
	case "step":
		run, err = s.store.Step(id, s.connectorHealthy(request.Context(), id))
	case "start":
		var generation uint64
		var started bool
		run, generation, started, err = s.store.Start(id)
		if err == nil && started {
			run, _, err = s.store.StepIfPlaying(id, generation, s.connectorHealthy(request.Context(), id))
			if err == nil && !run.Complete {
				go s.autoAdvance(id, generation)
			}
		}
	case "pause":
		run, err = s.store.Pause(id)
	case "reset":
		run, err = s.store.Reset(id)
	case "speed":
		run, err = s.store.SetSpeed(id, input.Speed)
	default:
		writeError(response, http.StatusBadRequest, "action must be start, pause, step, reset or speed")
		return
	}
	if err != nil {
		s.writeStoreError(response, err)
		return
	}
	writeJSON(response, http.StatusOK, run)
}

func (s *Server) autoAdvance(id string, generation uint64) {
	for {
		run, err := s.store.Get(id)
		if err != nil || !run.Playing || run.Complete {
			return
		}
		delay := 1200 * time.Millisecond / time.Duration(run.Speed)
		timer := time.NewTimer(delay)
		<-timer.C
		current, advanced, err := s.store.StepIfPlaying(id, generation, s.connectorHealthy(context.Background(), id))
		if err != nil || !advanced || current.Complete {
			return
		}
	}
}

func (s *Server) connectorHealthy(ctx context.Context, runID string) bool {
	run, err := s.store.Get(runID)
	if err != nil || run.Mode == ModeSimulation {
		return true
	}
	return s.capabilities(ctx).Connected
}

func (s *Server) handleEvents(response http.ResponseWriter, request *http.Request) {
	flusher, ok := response.(http.Flusher)
	if !ok {
		writeError(response, http.StatusInternalServerError, "streaming is not supported")
		return
	}
	response.Header().Set("Content-Type", "text/event-stream")
	response.Header().Set("Connection", "keep-alive")
	response.Header().Set("Content-Security-Policy", "default-src 'none'")
	afterSequence := 0
	if value := request.Header.Get("Last-Event-ID"); value != "" {
		parsed, parseErr := strconv.Atoi(value)
		if parseErr != nil || parsed < 0 {
			writeError(response, http.StatusBadRequest, "Last-Event-ID must be a non-negative integer")
			return
		}
		afterSequence = parsed
	}
	replay, channel, cancel, err := s.store.SubscribeFrom(request.PathValue("id"), afterSequence)
	if err != nil {
		s.writeStoreError(response, err)
		return
	}
	defer cancel()
	lastSequence := afterSequence
	for _, event := range replay {
		if err := writeSSE(response, event); err != nil {
			return
		}
		lastSequence = event.Sequence
	}
	flusher.Flush()
	for {
		select {
		case <-request.Context().Done():
			return
		case <-channel:
			events, eventsErr := s.store.EventsAfter(request.PathValue("id"), lastSequence)
			if eventsErr != nil {
				return
			}
			for _, event := range events {
				if err := writeSSE(response, event); err != nil {
					return
				}
				lastSequence = event.Sequence
			}
			flusher.Flush()
		}
	}
}

func (s *Server) writeStoreError(response http.ResponseWriter, err error) {
	if errors.Is(err, ErrRunNotFound) {
		writeError(response, http.StatusNotFound, err.Error())
		return
	}
	writeError(response, http.StatusBadRequest, err.Error())
}

func writeSSE(response http.ResponseWriter, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	sequence := 0
	if event, ok := value.(RunEvent); ok {
		sequence = event.Sequence
	}
	_, err = fmt.Fprintf(response, "id: %d\nevent: run-event\ndata: %s\n\n", sequence, payload)
	return err
}

func writeJSON(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(value)
}

func writeError(response http.ResponseWriter, status int, message string) {
	writeJSON(response, status, map[string]string{"error": message})
}
