<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import EventTimeline from '../components/EventTimeline.vue'
import InspectorPanel from '../components/InspectorPanel.vue'
import RunControls from '../components/RunControls.vue'
import ScenarioSidebar from '../components/ScenarioSidebar.vue'
import ServiceFlow from '../components/ServiceFlow.vue'
import { scenarios, getScenario } from '../data/scenarios'
import { useLabEngine } from '../composables/useLabEngine'
import { createRun, getCapabilities, sendRunAction, subscribeRun } from '../services/runner'
import type { RunMode, RunnerEvent, RunnerSnapshot } from '../types/scenario'

const route = useRoute()
const scenario = computed(() => getScenario(String(route.params.scenarioId ?? 'normal-withdrawal')))
const engine = useLabEngine(scenario)
const connectedAvailable = ref(false)
const runnerReachable = ref(false)
const runnerRunId = ref<string | null>(null)
const runnerMessage = ref('')
let closeStream: (() => void) | undefined
let provisionVersion = 0

function applyRunnerSnapshot(run: RunnerSnapshot) {
  engine.applySnapshot({
    scenarioId: run.scenarioId,
    currentIndex: run.currentIndex,
    playing: run.playing,
    speed: run.speed,
    mode: run.mode,
  })
}

function applyRunnerEvent(event: RunnerEvent) {
  const index = scenario.value.steps.findIndex(step => step.id === event.step.id)
  if (index < 0) return
  engine.applySnapshot({
    ...engine.snapshot.value,
    scenarioId: event.scenarioId,
    currentIndex: index,
    playing: index < scenario.value.steps.length - 1 ? engine.snapshot.value.playing : false,
    mode: event.mode,
  })
}

async function provisionRun(mode: RunMode = engine.snapshot.value.mode) {
  const version = ++provisionVersion
  closeStream?.()
  closeStream = undefined
  runnerRunId.value = null
  if (!runnerReachable.value) return
  try {
    const run = await createRun(scenario.value.id, mode)
    if (version !== provisionVersion) return
    runnerRunId.value = run.id
    applyRunnerSnapshot(run)
    closeStream = subscribeRun(run.id, applyRunnerEvent, () => {
      runnerMessage.value = 'Runner event stream disconnected'
    })
    runnerMessage.value = 'Runner event stream ready'
  } catch {
    runnerReachable.value = false
    runnerMessage.value = 'Runner unavailable · browser fallback active'
  }
}

onMounted(async () => {
  try {
    const capabilities = await getCapabilities()
    runnerReachable.value = true
    connectedAvailable.value = capabilities.connected
    await provisionRun('simulation')
  } catch {
    runnerReachable.value = false
    runnerMessage.value = 'Browser fallback active'
  }
})

watch(() => scenario.value.id, () => {
  void provisionRun(engine.snapshot.value.mode)
})

onBeforeUnmount(() => {
  provisionVersion++
  closeStream?.()
})

async function selectMode(mode: RunMode) {
  if (mode === 'connected' && !connectedAvailable.value) return
  engine.setMode(mode)
  await provisionRun(mode)
}

async function action(actionName: 'start' | 'pause' | 'step' | 'reset', fallback: () => void) {
  if (!runnerRunId.value) {
    fallback()
    return
  }
  try {
    const run = await sendRunAction(runnerRunId.value, actionName)
    applyRunnerSnapshot(run)
    if (actionName === 'reset') {
      closeStream?.()
      closeStream = subscribeRun(run.id, applyRunnerEvent, () => {
        runnerMessage.value = 'Runner event stream disconnected'
      })
    }
  } catch {
    runnerMessage.value = 'Runner action failed · reset to use browser fallback'
    engine.pause()
  }
}

async function setSpeed(speed: number) {
  if (!runnerRunId.value) {
    engine.setSpeed(speed)
    return
  }
  try {
    applyRunnerSnapshot(await sendRunAction(runnerRunId.value, 'speed', speed))
  } catch {
    runnerMessage.value = 'Runner speed update failed'
  }
}
</script>

<template>
  <div class="lab-page shell-wide" data-screen-label="Lab Cockpit">
    <header class="lab-intro">
      <div><p class="eyebrow">{{ scenario.eyebrow }}</p><h1>{{ scenario.title }}</h1><p>{{ scenario.summary }}</p></div>
      <div class="mode-control" aria-label="运行模式">
        <button type="button" :class="{ active: engine.snapshot.value.mode === 'simulation' }" @click="selectMode('simulation')"><i></i>模拟</button>
        <button type="button" :disabled="!connectedAvailable" :class="{ active: engine.snapshot.value.mode === 'connected' }" :title="connectedAvailable ? '连接本地沙箱' : '设置 LAB_ENABLE_CONNECTED=true 并启动本地 Runner 后开放'" @click="selectMode('connected')"><i></i>本地连接</button>
        <small :title="runnerMessage">{{ runnerReachable ? 'Runner online' : 'Browser engine' }}</small>
      </div>
    </header>

    <div class="lab-grid">
      <ScenarioSidebar :scenarios="scenarios" :active-id="scenario.id" />

      <section class="cockpit-panel">
        <div class="cockpit-heading">
          <div><span class="panel-kicker">Scenario {{ scenario.index }}</span><h2>{{ scenario.shortTitle }}</h2></div>
          <div><small>EXPECTED OUTCOME</small><p>{{ scenario.outcome }}</p></div>
        </div>
        <RunControls
          :playing="engine.snapshot.value.playing"
          :complete="engine.isComplete.value"
          :speed="engine.snapshot.value.speed"
          :progress="engine.progress.value"
          @start="action('start', engine.start)"
          @pause="action('pause', engine.pause)"
          @step="action('step', engine.step)"
          @reset="action('reset', engine.reset)"
          @speed="setSpeed"
        />
        <ServiceFlow
          :visible-steps="engine.visibleSteps.value"
          :current-step="engine.currentStep.value"
          :selected-service="engine.selectedService.value"
          @select="engine.selectedService.value = $event"
        />
        <div class="state-ribbon" :class="engine.currentStep.value ? `tone-${engine.currentStep.value.tone}` : 'tone-neutral'">
          <span>STATE</span><strong>{{ engine.currentStep.value?.state || 'idle' }}</strong><p>{{ engine.currentStep.value?.summary || '运行实验后，这里会显示当前状态及其工程含义。' }}</p><time>{{ engine.currentStep.value?.virtualTime || 'virtual clock ready' }}</time>
        </div>
      </section>

      <InspectorPanel :scenario="scenario" :step="engine.currentStep.value" :density="engine.density.value" @density="engine.density.value = $event" />

      <EventTimeline
        :scenario="scenario"
        :current-index="engine.snapshot.value.currentIndex"
        :selected-index="engine.selectedIndex.value"
        :selected-service="engine.selectedService.value"
        @select="engine.selectEvent"
      />
    </div>
  </div>
</template>
