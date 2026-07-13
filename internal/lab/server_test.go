package lab

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type connectorStub struct{ statuses map[string]ServiceStatus }

func (c connectorStub) Probe(_ context.Context) map[string]ServiceStatus { return c.statuses }

func TestCapabilitiesKeepConnectedModeBehindReachability(t *testing.T) {
	server := NewServer(NewStore(), connectorStub{statuses: map[string]ServiceStatus{
		"wallet-service": {Target: "127.0.0.1:1", Reachable: true},
		"wallet-api":     {Target: "127.0.0.1:2", Reachable: false},
	}}, true)
	request := httptest.NewRequest(http.MethodGet, "/api/capabilities", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	var result Capabilities
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Connected {
		t.Fatal("connected mode should be unavailable while a service is unreachable")
	}
	if strings.Contains(response.Body.String(), "127.0.0.1") {
		t.Fatalf("capabilities leaked local target: %s", response.Body.String())
	}
}

func TestRepeatedStartDoesNotAdvanceTwice(t *testing.T) {
	server := NewServer(NewStore(), nil, false)
	handler := server.Handler()
	createResponse := httptest.NewRecorder()
	handler.ServeHTTP(createResponse, httptest.NewRequest(http.MethodPost, "/api/runs", bytes.NewBufferString(`{"scenarioId":"normal-withdrawal","mode":"simulation"}`)))
	var run RunSnapshot
	if err := json.NewDecoder(createResponse.Body).Decode(&run); err != nil {
		t.Fatal(err)
	}
	for index := 0; index < 2; index++ {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/runs/"+run.ID+"/actions", bytes.NewBufferString(`{"action":"start"}`)))
		if response.Code != http.StatusOK {
			t.Fatalf("start %d status = %d", index+1, response.Code)
		}
	}
	current, err := server.store.Get(run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(current.Events) != 1 || current.CurrentIndex != 0 {
		t.Fatalf("duplicate start advanced more than once: %+v", current)
	}
	_, _ = server.store.Pause(run.ID)
}

func TestWriteSSEIncludesReconnectableEventID(t *testing.T) {
	event := RunEvent{Sequence: 7, RunID: "run_1"}
	response := httptest.NewRecorder()
	if err := writeSSE(response, event); err != nil {
		t.Fatal(err)
	}
	body := response.Body.String()
	if !strings.Contains(body, "id: 7\n") || !strings.Contains(body, "event: run-event\n") {
		t.Fatalf("SSE frame is missing id/event: %q", body)
	}
}

func TestCreateAndStepSimulationRun(t *testing.T) {
	server := NewServer(NewStore(), nil, false)
	handler := server.Handler()
	createBody := bytes.NewBufferString(`{"scenarioId":"normal-withdrawal","mode":"simulation"}`)
	createRequest := httptest.NewRequest(http.MethodPost, "/api/runs", createBody)
	createResponse := httptest.NewRecorder()
	handler.ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body=%s", createResponse.Code, createResponse.Body.String())
	}
	var run RunSnapshot
	if err := json.NewDecoder(createResponse.Body).Decode(&run); err != nil {
		t.Fatal(err)
	}
	stepRequest := httptest.NewRequest(http.MethodPost, "/api/runs/"+run.ID+"/actions", bytes.NewBufferString(`{"action":"step"}`))
	stepResponse := httptest.NewRecorder()
	handler.ServeHTTP(stepResponse, stepRequest)
	if stepResponse.Code != http.StatusOK {
		t.Fatalf("step status = %d, body=%s", stepResponse.Code, stepResponse.Body.String())
	}
	if err := json.NewDecoder(stepResponse.Body).Decode(&run); err != nil {
		t.Fatal(err)
	}
	if run.CurrentIndex != 0 || len(run.Events) != 1 || !run.Events[0].Sanitized {
		t.Fatalf("unexpected run: %+v", run)
	}
}

func TestConnectedRunIsForbiddenWhenDisabled(t *testing.T) {
	server := NewServer(NewStore(), nil, false)
	request := httptest.NewRequest(http.MethodPost, "/api/runs", bytes.NewBufferString(`{"scenarioId":"normal-withdrawal","mode":"connected"}`))
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusForbidden)
	}
}
