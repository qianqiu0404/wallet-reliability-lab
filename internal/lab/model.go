package lab

import "time"

type RunMode string

const (
	ModeSimulation RunMode = "simulation"
	ModeConnected  RunMode = "connected"
)

type Step struct {
	ID          string `json:"id"`
	Service     string `json:"service"`
	Phase       string `json:"phase"`
	State       string `json:"state"`
	Tone        string `json:"tone"`
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	VirtualTime string `json:"virtualTime"`
}

type Scenario struct {
	ID         string `json:"id"`
	Index      string `json:"index"`
	Title      string `json:"title"`
	ShortTitle string `json:"shortTitle"`
	Summary    string `json:"summary"`
	Steps      []Step `json:"steps"`
}

type RunEvent struct {
	RunID       string    `json:"runId"`
	Sequence    int       `json:"sequence"`
	ScenarioID  string    `json:"scenarioId"`
	Mode        RunMode   `json:"mode"`
	Step        Step      `json:"step"`
	EmittedAt   time.Time `json:"emittedAt"`
	RequestID   string    `json:"requestId"`
	TxHash      string    `json:"txHash,omitempty"`
	Sanitized   bool      `json:"sanitized"`
	ConnectorOK bool      `json:"connectorOk"`
}

type RunSnapshot struct {
	ID           string     `json:"id"`
	ScenarioID   string     `json:"scenarioId"`
	Mode         RunMode    `json:"mode"`
	CurrentIndex int        `json:"currentIndex"`
	Playing      bool       `json:"playing"`
	Complete     bool       `json:"complete"`
	Speed        int        `json:"speed"`
	Events       []RunEvent `json:"events"`
	generation   uint64
	nextSequence int
}

type Capabilities struct {
	Simulation      bool                     `json:"simulation"`
	Connected       bool                     `json:"connected"`
	ConnectedReason string                   `json:"connectedReason"`
	Services        map[string]ServiceStatus `json:"services,omitempty"`
}

type ServiceStatus struct {
	Target    string `json:"-"`
	Reachable bool   `json:"reachable"`
	Message   string `json:"message"`
}
