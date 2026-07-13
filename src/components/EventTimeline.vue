<script setup lang="ts">
import type { ScenarioDefinition, ServiceId } from '../types/scenario'

defineProps<{
  scenario: ScenarioDefinition
  currentIndex: number
  selectedIndex: number
  selectedService: ServiceId | null
}>()

const emit = defineEmits<{ select: [index: number] }>()
</script>

<template>
  <section class="timeline-panel">
    <div class="timeline-heading">
      <div><span class="panel-kicker">Deterministic event stream</span><h2>事件时间线</h2></div>
      <p>虚拟时间固定，重复运行会得到完全相同的状态与证据。</p>
    </div>
    <div class="timeline-scroll">
      <button
        v-for="(step, index) in scenario.steps"
        :key="step.id"
        type="button"
        :disabled="index > currentIndex"
        :class="['timeline-event', `tone-${step.tone}`, { current: index === selectedIndex, filtered: selectedService && step.service !== selectedService }]"
        @click="emit('select', index)"
      >
        <span class="event-sequence">{{ String(index + 1).padStart(2, '0') }}</span>
        <span class="event-marker"><i></i></span>
        <span class="event-copy"><small>{{ step.virtualTime }} · {{ step.service }}</small><strong>{{ step.title }}</strong><em>{{ step.state }}</em></span>
      </button>
    </div>
  </section>
</template>

