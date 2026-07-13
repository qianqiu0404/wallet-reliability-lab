<script setup lang="ts">
import { computed, ref } from 'vue'
import { evidenceRecords, scenarios } from '../data/scenarios'

const filter = ref<'all' | 'verified' | 'target'>('all')
const filtered = computed(() => filter.value === 'all' ? evidenceRecords : evidenceRecords.filter(item => item.status === filter.value))

function scenarioName(id: string) {
  return scenarios.find(item => item.id === id)?.shortTitle ?? id
}
</script>

<template>
  <div class="content-page shell" data-screen-label="Evidence">
    <header class="content-hero evidence-hero"><div><p class="eyebrow">Evidence / sanitized by default</p><h1>判断必须回到证据</h1><p>只公开测试名称、代码角色、验证结论和适用边界；私有路径、密钥、RPC 地址与原始日志不会进入浏览器。</p></div><div class="evidence-count"><strong>{{ evidenceRecords.filter(item => item.status === 'verified').length }}</strong><span>verified records</span><small>+ {{ evidenceRecords.filter(item => item.status === 'target').length }} target design</small></div></header>
    <div class="evidence-toolbar"><div><button type="button" :class="{ active: filter === 'all' }" @click="filter = 'all'">全部</button><button type="button" :class="{ active: filter === 'verified' }" @click="filter = 'verified'">已验证</button><button type="button" :class="{ active: filter === 'target' }" @click="filter = 'target'">目标设计</button></div><span>事实与目标态分开显示</span></div>
    <section class="evidence-grid">
      <article v-for="(item, index) in filtered" :key="item.id">
        <header><span>{{ String(index + 1).padStart(2, '0') }}</span><em :class="item.status">{{ item.status === 'verified' ? '已验证' : '目标设计' }}</em></header>
        <p class="evidence-label">{{ item.label }}</p><h2>{{ item.testName }}</h2><small>{{ item.sourceRole }}</small><p>{{ item.summary }}</p>
        <footer><router-link v-for="scenarioId in item.relatedScenarioIds" :key="scenarioId" :to="`/lab/${scenarioId}`">{{ scenarioName(scenarioId) }}</router-link></footer>
      </article>
    </section>
  </div>
</template>

