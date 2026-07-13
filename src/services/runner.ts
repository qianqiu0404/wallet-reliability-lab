import type { RunMode, RunnerCapabilities, RunnerEvent, RunnerSnapshot } from '../types/scenario'

const REQUEST_TIMEOUT_MS = 1600

async function requestJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...init?.headers },
    signal: AbortSignal.timeout(REQUEST_TIMEOUT_MS),
  })
  const contentType = response.headers.get('content-type') ?? ''
  if (!response.ok || !contentType.includes('application/json')) {
    throw new Error(`Runner request failed: ${response.status}`)
  }
  return response.json() as Promise<T>
}

export function getCapabilities(): Promise<RunnerCapabilities> {
  return requestJSON('/api/capabilities')
}

export function createRun(scenarioId: string, mode: RunMode): Promise<RunnerSnapshot> {
  return requestJSON('/api/runs', { method: 'POST', body: JSON.stringify({ scenarioId, mode }) })
}

export function sendRunAction(runId: string, action: string, speed?: number): Promise<RunnerSnapshot> {
  return requestJSON(`/api/runs/${encodeURIComponent(runId)}/actions`, {
    method: 'POST',
    body: JSON.stringify({ action, ...(speed ? { speed } : {}) }),
  })
}

export function subscribeRun(
  runId: string,
  onEvent: (event: RunnerEvent) => void,
  onUnavailable: () => void,
): () => void {
  const source = new EventSource(`/api/runs/${encodeURIComponent(runId)}/events`)
  source.addEventListener('run-event', raw => {
    onEvent(JSON.parse((raw as MessageEvent<string>).data) as RunnerEvent)
  })
  source.onerror = () => {
    if (source.readyState === EventSource.CLOSED) onUnavailable()
  }
  return () => source.close()
}
