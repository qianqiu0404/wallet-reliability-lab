<script setup lang="ts">
import { computed } from 'vue'
import type { ScenarioStep, ServiceId } from '../types/scenario'

const props = defineProps<{
  visibleSteps: ScenarioStep[]
  currentStep: ScenarioStep | null
  selectedService: ServiceId | null
}>()

const emit = defineEmits<{ select: [service: ServiceId | null] }>()

const services: Array<{ id: ServiceId; index: string; name: string; role: string }> = [
  { id: 'wallet-service', index: '01', name: 'wallet-service', role: '资金编排与状态' },
  { id: 'risk-service', index: '02', name: 'risk-service', role: '校验与审批凭证' },
  { id: 'wallet-api', index: '03', name: 'wallet-api', role: '链资源与广播' },
  { id: 'wallet-sign', index: '04', name: 'wallet-sign', role: '策略校验与签名' },
]

function latestFor(service: ServiceId) {
  return [...props.visibleSteps].reverse().find(step => step.service === service)
}

function nodeClass(service: ServiceId) {
  const latest = latestFor(service)
  return [
    latest ? `tone-${latest.tone}` : 'tone-neutral',
    { active: props.currentStep?.service === service, dimmed: props.selectedService && props.selectedService !== service },
  ]
}

const edgeProgress = computed(() => {
  const index = services.findIndex(service => service.id === props.currentStep?.service)
  return index < 0 ? 0 : ((index + 0.5) / services.length) * 100
})
</script>

<template>
  <section class="service-flow" aria-label="钱包服务调用链">
    <div class="flow-context"><span>BUSINESS REQUEST</span><i></i><span>CHAIN FACT</span></div>
    <div class="flow-track"><span :style="{ width: `${edgeProgress}%` }"></span></div>
    <div class="service-nodes">
      <button
        v-for="service in services"
        :key="service.id"
        type="button"
        :class="['service-node', nodeClass(service.id)]"
        @click="emit('select', selectedService === service.id ? null : service.id)"
      >
        <span class="service-index">{{ service.index }}</span>
        <span class="service-signal"><i></i></span>
        <strong>{{ service.name }}</strong>
        <small>{{ service.role }}</small>
        <em>{{ latestFor(service.id)?.state || 'waiting' }}</em>
      </button>
    </div>
    <div class="secondary-boundaries">
      <span :class="{ active: currentStep?.service === 'chain' }"><i></i> chain / confirmations</span>
      <span :class="{ active: currentStep?.service === 'notifier' }"><i></i> notifier / callback</span>
    </div>
  </section>
</template>

