package lab

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrRunNotFound = errors.New("run not found")

type Store struct {
	mu          sync.RWMutex
	runs        map[string]*RunSnapshot
	subscribers map[string]map[chan struct{}]struct{}
	nextID      int64
}

func NewStore() *Store {
	return &Store{runs: map[string]*RunSnapshot{}, subscribers: map[string]map[chan struct{}]struct{}{}}
}

func (s *Store) Create(scenarioID string, mode RunMode) (RunSnapshot, error) {
	if _, ok := FindScenario(scenarioID); !ok {
		return RunSnapshot{}, fmt.Errorf("unknown scenario %q", scenarioID)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	id := fmt.Sprintf("run_%06d", s.nextID)
	run := &RunSnapshot{ID: id, ScenarioID: scenarioID, Mode: mode, CurrentIndex: -1, Speed: 1, Events: []RunEvent{}}
	s.runs[id] = run
	return cloneRun(run), nil
}

func (s *Store) Get(id string) (RunSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	run, ok := s.runs[id]
	if !ok {
		return RunSnapshot{}, ErrRunNotFound
	}
	return cloneRun(run), nil
}

func (s *Store) Step(id string, connectorOK bool) (RunSnapshot, error) {
	s.mu.Lock()
	result, advanced, err := s.stepLocked(id, connectorOK)
	if err != nil {
		s.mu.Unlock()
		return RunSnapshot{}, err
	}
	subscribers := subscriberList(s.subscribers[id])
	s.mu.Unlock()
	if advanced {
		notifySubscribers(subscribers)
	}
	return result, nil
}

func (s *Store) StepIfPlaying(id string, generation uint64, connectorOK bool) (RunSnapshot, bool, error) {
	s.mu.Lock()
	run, ok := s.runs[id]
	if !ok {
		s.mu.Unlock()
		return RunSnapshot{}, false, ErrRunNotFound
	}
	if !run.Playing || run.generation != generation {
		result := cloneRun(run)
		s.mu.Unlock()
		return result, false, nil
	}
	result, advanced, err := s.stepLocked(id, connectorOK)
	subscribers := subscriberList(s.subscribers[id])
	s.mu.Unlock()
	if advanced {
		notifySubscribers(subscribers)
	}
	return result, advanced, err
}

func (s *Store) stepLocked(id string, connectorOK bool) (RunSnapshot, bool, error) {
	run, ok := s.runs[id]
	if !ok {
		return RunSnapshot{}, false, ErrRunNotFound
	}
	scenario, _ := FindScenario(run.ScenarioID)
	if run.CurrentIndex >= len(scenario.Steps)-1 {
		run.Playing = false
		run.Complete = true
		return cloneRun(run), false, nil
	}
	run.CurrentIndex++
	step := scenario.Steps[run.CurrentIndex]
	run.nextSequence++
	event := RunEvent{
		RunID: run.ID, Sequence: run.nextSequence, ScenarioID: run.ScenarioID, Mode: run.Mode,
		Step: step, EmittedAt: time.Now().UTC(), RequestID: "wd_20260713_0018", Sanitized: true, ConnectorOK: connectorOK,
	}
	if step.State == "signed" || step.State == "broadcasted" || step.State == "wallet_done" || step.State == "notified" {
		event.TxHash = "0x8fe4…c91a"
	}
	run.Events = append(run.Events, event)
	if run.CurrentIndex == len(scenario.Steps)-1 {
		run.Complete = true
		run.Playing = false
	}
	return cloneRun(run), true, nil
}

func (s *Store) Start(id string) (RunSnapshot, uint64, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[id]
	if !ok {
		return RunSnapshot{}, 0, false, ErrRunNotFound
	}
	if run.Playing || run.Complete {
		return cloneRun(run), run.generation, false, nil
	}
	run.generation++
	run.Playing = true
	return cloneRun(run), run.generation, true, nil
}

func (s *Store) Pause(id string) (RunSnapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[id]
	if !ok {
		return RunSnapshot{}, ErrRunNotFound
	}
	run.generation++
	run.Playing = false
	return cloneRun(run), nil
}

func (s *Store) SetSpeed(id string, speed int) (RunSnapshot, error) {
	if speed != 1 && speed != 2 && speed != 4 {
		return RunSnapshot{}, fmt.Errorf("speed must be 1, 2 or 4")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[id]
	if !ok {
		return RunSnapshot{}, ErrRunNotFound
	}
	run.Speed = speed
	return cloneRun(run), nil
}

func (s *Store) Reset(id string) (RunSnapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[id]
	if !ok {
		return RunSnapshot{}, ErrRunNotFound
	}
	run.CurrentIndex, run.Playing, run.Complete = -1, false, false
	run.generation++
	run.Events = []RunEvent{}
	return cloneRun(run), nil
}

func (s *Store) SubscribeFrom(id string, afterSequence int) ([]RunEvent, <-chan struct{}, func(), error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[id]
	if !ok {
		return nil, nil, nil, ErrRunNotFound
	}
	channel := make(chan struct{}, 1)
	if s.subscribers[id] == nil {
		s.subscribers[id] = map[chan struct{}]struct{}{}
	}
	s.subscribers[id][channel] = struct{}{}
	replay := eventsAfter(run.Events, afterSequence)
	cancel := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if _, ok := s.subscribers[id][channel]; ok {
			delete(s.subscribers[id], channel)
		}
	}
	return replay, channel, cancel, nil
}

func (s *Store) EventsAfter(id string, afterSequence int) ([]RunEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	run, ok := s.runs[id]
	if !ok {
		return nil, ErrRunNotFound
	}
	return eventsAfter(run.Events, afterSequence), nil
}

func cloneRun(run *RunSnapshot) RunSnapshot {
	clone := *run
	clone.Events = append([]RunEvent(nil), run.Events...)
	return clone
}

func eventsAfter(events []RunEvent, afterSequence int) []RunEvent {
	result := make([]RunEvent, 0, len(events))
	for _, event := range events {
		if event.Sequence > afterSequence {
			result = append(result, event)
		}
	}
	return result
}

func subscriberList(items map[chan struct{}]struct{}) []chan struct{} {
	result := make([]chan struct{}, 0, len(items))
	for item := range items {
		result = append(result, item)
	}
	return result
}

func notifySubscribers(subscribers []chan struct{}) {
	for _, subscriber := range subscribers {
		select {
		case subscriber <- struct{}{}:
		default:
		}
	}
}
