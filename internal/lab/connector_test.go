package lab

import "testing"

func TestIsLoopbackTarget(t *testing.T) {
	tests := []struct {
		target string
		want   bool
	}{
		{"127.0.0.1:50051", true},
		{"localhost:50051", true},
		{"[::1]:50051", true},
		{"10.0.0.5:50051", false},
		{"wallet.internal:50051", false},
		{"missing-port", false},
	}
	for _, item := range tests {
		if got := IsLoopbackTarget(item.target); got != item.want {
			t.Fatalf("IsLoopbackTarget(%q) = %v, want %v", item.target, got, item.want)
		}
	}
}

func TestParseTargetsRejectsInvalidShape(t *testing.T) {
	if _, err := ParseTargets("wallet-service"); err == nil {
		t.Fatal("expected invalid target error")
	}
	targets, err := ParseTargets("wallet-service=127.0.0.1:1,wallet-api=localhost:2")
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 2 {
		t.Fatalf("targets = %d, want 2", len(targets))
	}
}
