<script setup lang="ts">
import type { ScenarioDefinition } from '../types/scenario'

defineProps<{ scenarios: ScenarioDefinition[]; activeId: string }>()
</script>

<template>
  <aside class="scenario-sidebar">
    <div class="panel-kicker">Scenario library</div>
    <h2>提现可靠性实验</h2>
    <p>相同主链路，不同故障语义。每个结果都能回到状态和证据。</p>
    <nav class="scenario-list" aria-label="实验场景">
      <router-link
        v-for="scenario in scenarios"
        :key="scenario.id"
        :to="`/lab/${scenario.id}`"
        :class="['scenario-item', `tone-${scenario.accent}`, { active: scenario.id === activeId }]"
      >
        <span>{{ scenario.index }}</span>
        <div><strong>{{ scenario.shortTitle }}</strong><small>{{ scenario.eyebrow }}</small></div>
        <i aria-hidden="true">→</i>
      </router-link>
    </nav>
    <div class="sidebar-note">
      <span class="safe-dot"></span>
      <div><strong>公开安全模式</strong><small>固定数据 · 无真实资金 · 无密钥</small></div>
    </div>
  </aside>
</template>

