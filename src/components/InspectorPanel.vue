<script setup lang="ts">
import { computed } from 'vue'
import { getEvidence } from '../data/scenarios'
import type { LabDensity, ScenarioDefinition, ScenarioStep } from '../types/scenario'

const props = defineProps<{
  scenario: ScenarioDefinition
  step: ScenarioStep | null
  density: LabDensity
}>()

const emit = defineEmits<{ density: [value: LabDensity] }>()
const evidence = computed(() => getEvidence(props.step?.evidenceIds ?? []))
</script>

<template>
  <aside class="inspector-panel">
    <div class="density-switch" aria-label="信息密度">
      <button type="button" :class="{ active: density === 'brief' }" @click="emit('density', 'brief')">速览</button>
      <button type="button" :class="{ active: density === 'developer' }" @click="emit('density', 'developer')">开发者</button>
    </div>
    <div v-if="step" class="inspector-content">
      <div class="inspector-meta"><span>{{ step.phase }}</span><time>{{ step.virtualTime }}</time></div>
      <h2>{{ step.title }}</h2>
      <p class="inspector-summary">{{ step.summary }}</p>
      <div class="judgement-block">
        <span>Engineering judgement</span>
        <p>{{ step.reasoning }}</p>
      </div>
      <dl class="detail-list">
        <div v-for="detail in step.details" :key="detail.label"><dt>{{ detail.label }}</dt><dd>{{ detail.value }}</dd></div>
      </dl>
      <template v-if="density === 'developer'">
        <div class="log-block"><span>Sanitized log</span><code>{{ step.log }}</code></div>
        <div class="evidence-mini-list">
          <span>Evidence</span>
          <article v-for="item in evidence" :key="item.id">
            <div><strong>{{ item.testName }}</strong><small>{{ item.sourceRole }}</small></div>
            <div class="evidence-actions"><em :class="item.status">{{ item.status === 'verified' ? '已验证' : '目标设计' }}</em><a :href="item.sourceUrl" target="_blank" rel="noreferrer" :aria-label="`${item.sourceLabel}：${item.testName}`">源码 ↗</a></div>
          </article>
        </div>
      </template>
    </div>
    <div v-else class="inspector-empty">
      <span>00</span>
      <h2>准备运行实验</h2>
      <p>{{ scenario.summary }}</p>
      <small>点击“开始运行”自动播放，或使用“单步”检查每次状态变化。</small>
    </div>
  </aside>
</template>
