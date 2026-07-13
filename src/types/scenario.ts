export type ServiceId =
  | 'business'
  | 'wallet-service'
  | 'risk-service'
  | 'wallet-api'
  | 'wallet-sign'
  | 'chain'
  | 'notifier'

export type EventTone = 'neutral' | 'active' | 'success' | 'warning' | 'danger' | 'recovered'
export type LabDensity = 'brief' | 'developer'
export type RunMode = 'simulation' | 'connected'

export interface EventDetail {
  label: string
  value: string
}

export interface ScenarioStep {
  id: string
  service: ServiceId
  phase: string
  state: string
  tone: EventTone
  title: string
  summary: string
  reasoning: string
  virtualTime: string
  details: EventDetail[]
  log: string
  evidenceIds: string[]
}

export interface ScenarioDefinition {
  id: string
  index: string
  title: string
  shortTitle: string
  eyebrow: string
  summary: string
  outcome: string
  duration: string
  accent: EventTone
  steps: ScenarioStep[]
}

export interface EvidenceRecord {
  id: string
  kind: 'test' | 'design'
  label: string
  testName: string
  sourceRole: string
  summary: string
  status: 'verified' | 'target'
  relatedScenarioIds: string[]
}

export interface RunSnapshot {
  scenarioId: string
  currentIndex: number
  playing: boolean
  speed: number
  mode: RunMode
}

export interface RunnerCapabilities {
  simulation: boolean
  connected: boolean
  connectedReason: string
  services?: Record<string, { reachable: boolean; message: string }>
}

export interface RunnerEvent {
  runId: string
  sequence: number
  scenarioId: string
  mode: RunMode
  step: Pick<ScenarioStep, 'id' | 'service' | 'phase' | 'state' | 'tone' | 'title' | 'summary' | 'virtualTime'>
  emittedAt: string
  requestId: string
  txHash?: string
  sanitized: boolean
  connectorOk: boolean
}

export interface RunnerSnapshot {
  id: string
  scenarioId: string
  mode: RunMode
  currentIndex: number
  playing: boolean
  complete: boolean
  speed: number
  events: RunnerEvent[]
}
