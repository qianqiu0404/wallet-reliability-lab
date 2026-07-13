import { computed, onBeforeUnmount, ref, watch, type Ref } from 'vue'
import type { LabDensity, RunMode, RunSnapshot, ScenarioDefinition, ServiceId } from '../types/scenario'

export function advanceSnapshot(snapshot: RunSnapshot, stepCount: number): RunSnapshot {
  if (stepCount <= 0) return { ...snapshot, currentIndex: -1, playing: false }
  if (snapshot.currentIndex >= stepCount - 1) return { ...snapshot, playing: false }
  const nextIndex = snapshot.currentIndex + 1
  return { ...snapshot, currentIndex: nextIndex, playing: nextIndex < stepCount - 1 ? snapshot.playing : false }
}

export function createInitialSnapshot(scenarioId: string): RunSnapshot {
  return { scenarioId, currentIndex: -1, playing: false, speed: 1, mode: 'simulation' }
}

export function useLabEngine(scenario: Ref<ScenarioDefinition>) {
  const snapshot = ref<RunSnapshot>(createInitialSnapshot(scenario.value.id))
  const density = ref<LabDensity>('brief')
  const selectedService = ref<ServiceId | null>(null)
  const selectedIndex = ref(-1)
  let timer: ReturnType<typeof window.setTimeout> | undefined

  const currentStep = computed(() => {
    const index = selectedIndex.value >= 0 ? selectedIndex.value : snapshot.value.currentIndex
    return scenario.value.steps[index] ?? null
  })
  const visibleSteps = computed(() => scenario.value.steps.slice(0, snapshot.value.currentIndex + 1))
  const progress = computed(() => {
    if (snapshot.value.currentIndex < 0) return 0
    return ((snapshot.value.currentIndex + 1) / scenario.value.steps.length) * 100
  })
  const isComplete = computed(() => snapshot.value.currentIndex === scenario.value.steps.length - 1)

  function clearTimer() {
    if (timer !== undefined) window.clearTimeout(timer)
    timer = undefined
  }

  function scheduleNext() {
    clearTimer()
    if (!snapshot.value.playing) return
    timer = window.setTimeout(() => {
      step()
      scheduleNext()
    }, 1350 / snapshot.value.speed)
  }

  function start() {
    if (isComplete.value) reset()
    snapshot.value = { ...snapshot.value, playing: true }
    if (snapshot.value.currentIndex < 0) step()
    scheduleNext()
  }

  function pause() {
    snapshot.value = { ...snapshot.value, playing: false }
    clearTimer()
  }

  function step() {
    snapshot.value = advanceSnapshot(snapshot.value, scenario.value.steps.length)
    selectedIndex.value = snapshot.value.currentIndex
    if (isComplete.value) clearTimer()
  }

  function reset() {
    clearTimer()
    snapshot.value = { ...createInitialSnapshot(scenario.value.id), speed: snapshot.value.speed, mode: snapshot.value.mode }
    selectedIndex.value = -1
    selectedService.value = null
  }

  function setSpeed(speed: number) {
    snapshot.value = { ...snapshot.value, speed }
    if (snapshot.value.playing) scheduleNext()
  }

  function setMode(mode: RunMode) {
    snapshot.value = { ...snapshot.value, mode }
  }

  function selectEvent(index: number) {
    if (index <= snapshot.value.currentIndex) selectedIndex.value = index
  }

  function applySnapshot(next: Pick<RunSnapshot, 'scenarioId' | 'currentIndex' | 'playing' | 'speed' | 'mode'>) {
    clearTimer()
    snapshot.value = { ...next }
    selectedIndex.value = next.currentIndex
  }

  watch(() => scenario.value.id, id => {
    clearTimer()
    snapshot.value = createInitialSnapshot(id)
    selectedIndex.value = -1
    selectedService.value = null
  })

  onBeforeUnmount(clearTimer)

  return {
    snapshot,
    density,
    selectedService,
    selectedIndex,
    currentStep,
    visibleSteps,
    progress,
    isComplete,
    start,
    pause,
    step,
    reset,
    setSpeed,
    setMode,
    selectEvent,
    applySnapshot,
  }
}
