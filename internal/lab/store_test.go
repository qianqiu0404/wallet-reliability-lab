package lab

import (
	"sync"
	"testing"
)

func TestStoreRunAdvancesDeterministically(t *testing.T) {
	store := NewStore()
	run, err := store.Create("retryable-broadcast", ModeSimulation)
	if err != nil {
		t.Fatal(err)
	}
	for index := 0; index < 4; index++ {
		run, err = store.Step(run.ID, true)
		if err != nil {
			t.Fatal(err)
		}
	}
	if run.CurrentIndex != 3 {
		t.Fatalf("current index = %d, want 3", run.CurrentIndex)
	}
	if got := run.Events[2].Step.State; got != "signed" {
		t.Fatalf("retryable failure state = %q, want signed", got)
	}
	if got := run.Events[3].Step.State; got != "signed" {
		t.Fatalf("classified failure state = %q, want signed", got)
	}
}

func TestStartIsIdempotentAndUsesSingleGeneration(t *testing.T) {
	store := NewStore()
	run, _ := store.Create("normal-withdrawal", ModeSimulation)
	first, generation, started, err := store.Start(run.ID)
	if err != nil || !started || !first.Playing {
		t.Fatalf("first start = %+v, generation=%d, started=%v, err=%v", first, generation, started, err)
	}
	second, secondGeneration, startedAgain, err := store.Start(run.ID)
	if err != nil || startedAgain || secondGeneration != generation || !second.Playing {
		t.Fatalf("duplicate start = %+v, generation=%d, started=%v, err=%v", second, secondGeneration, startedAgain, err)
	}
}

func TestPauseAndResetInvalidatePendingPlayer(t *testing.T) {
	store := NewStore()
	run, _ := store.Create("normal-withdrawal", ModeSimulation)
	_, generation, _, _ := store.Start(run.ID)
	_, _ = store.Pause(run.ID)
	paused, advanced, err := store.StepIfPlaying(run.ID, generation, true)
	if err != nil || advanced || paused.CurrentIndex != -1 {
		t.Fatalf("paused player advanced: %+v, advanced=%v, err=%v", paused, advanced, err)
	}
	_, generation, _, _ = store.Start(run.ID)
	_, _ = store.Reset(run.ID)
	reset, advanced, err := store.StepIfPlaying(run.ID, generation, true)
	if err != nil || advanced || reset.CurrentIndex != -1 {
		t.Fatalf("reset player advanced: %+v, advanced=%v, err=%v", reset, advanced, err)
	}
}

func TestEventSequenceRemainsMonotonicAcrossReset(t *testing.T) {
	store := NewStore()
	run, _ := store.Create("normal-withdrawal", ModeSimulation)
	run, _ = store.Step(run.ID, true)
	firstSequence := run.Events[0].Sequence
	_, _ = store.Reset(run.ID)
	run, _ = store.Step(run.ID, true)
	if run.Events[0].Sequence <= firstSequence {
		t.Fatalf("sequence after reset = %d, want > %d", run.Events[0].Sequence, firstSequence)
	}
}

func TestSubscriberCancellationIsSafeDuringConcurrentSteps(t *testing.T) {
	store := NewStore()
	run, _ := store.Create("normal-withdrawal", ModeSimulation)
	var wait sync.WaitGroup
	for index := 0; index < 100; index++ {
		_, _, cancel, err := store.SubscribeFrom(run.ID, 0)
		if err != nil {
			t.Fatal(err)
		}
		wait.Add(2)
		go func() { defer wait.Done(); cancel() }()
		go func() { defer wait.Done(); _, _ = store.Step(run.ID, true) }()
	}
	wait.Wait()
}

func TestSubscribeFromReplaysOnlyMissingEvents(t *testing.T) {
	store := NewStore()
	run, _ := store.Create("normal-withdrawal", ModeSimulation)
	for index := 0; index < 3; index++ {
		run, _ = store.Step(run.ID, true)
	}
	replay, _, cancel, err := store.SubscribeFrom(run.ID, run.Events[0].Sequence)
	defer cancel()
	if err != nil {
		t.Fatal(err)
	}
	if len(replay) != 2 || replay[0].Sequence != run.Events[1].Sequence || replay[1].Sequence != run.Events[2].Sequence {
		t.Fatalf("unexpected replay: %+v", replay)
	}
}

func TestStoreResetPreservesModeAndSpeed(t *testing.T) {
	store := NewStore()
	run, _ := store.Create("normal-withdrawal", ModeConnected)
	run, _ = store.SetSpeed(run.ID, 4)
	run, _ = store.Step(run.ID, true)
	run, _ = store.Reset(run.ID)
	if run.Mode != ModeConnected || run.Speed != 4 || run.CurrentIndex != -1 || len(run.Events) != 0 {
		t.Fatalf("unexpected reset snapshot: %+v", run)
	}
}

func TestCompensationScenarioCrossesThresholdBeforeRecovery(t *testing.T) {
	scenario, _ := FindScenario("compensation-recovery")
	if scenario.Steps[1].VirtualTime != "11:29:59.999" || scenario.Steps[2].VirtualTime != "11:30:00.001" {
		t.Fatalf("unexpected compensation boundary: %s -> %s", scenario.Steps[1].VirtualTime, scenario.Steps[2].VirtualTime)
	}
	if scenario.Steps[5].State != "wallet_done" {
		t.Fatalf("recovery state = %s", scenario.Steps[5].State)
	}
}
