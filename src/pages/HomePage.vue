<script setup lang="ts">
import { scenarios } from '../data/scenarios'
</script>

<template>
  <div data-screen-label="Home">
    <section class="home-hero shell">
      <div class="hero-copy">
        <p class="eyebrow">Wallet state machines · failure recovery · evidence</p>
        <h1>把钱包故障<br />变成可以运行的证据</h1>
        <p class="hero-lede">Wallet Reliability Lab 用确定性实验拆开提现状态、广播失败与补偿恢复。先看工程判断，再深入请求、日志和测试。</p>
        <div class="hero-actions">
          <router-link class="button-primary" to="/lab/normal-withdrawal">运行正常提现基线 <span>→</span></router-link>
          <router-link class="button-secondary" to="/architecture">先看系统边界</router-link>
        </div>
        <div class="hero-assurances">
          <span><i></i> 无真实资金</span><span><i></i> 固定脱敏数据</span><span><i></i> 证据可回溯</span>
        </div>
      </div>
      <div class="hero-console" aria-label="实验工作台预览">
        <div class="console-top"><span>LIVE PREVIEW / 01</span><small>simulation</small></div>
        <div class="console-state">
          <span>Current state</span>
          <strong>broadcasted</strong>
          <p>节点接受交易，等待链上确认。</p>
        </div>
        <div class="mini-services">
          <article class="done"><span>01</span><strong>wallet-service</strong><small>signed</small></article>
          <i></i>
          <article class="done"><span>02</span><strong>risk-service</strong><small>approved</small></article>
          <i></i>
          <article class="active"><span>03</span><strong>wallet-api</strong><small>broadcasted</small></article>
          <i></i>
          <article><span>04</span><strong>wallet-sign</strong><small>signed</small></article>
        </div>
        <div class="console-log"><span>09:42:16.090</span><code>broadcast accepted · 1 of 1 items classified as success</code></div>
        <div class="console-progress"><span></span></div>
      </div>
    </section>

    <section class="home-manifesto">
      <div class="shell manifesto-grid">
        <p class="section-index">01 / POSITION</p>
        <h2>不是演示“成功路径”，<br />而是解释结果不确定时为什么这样判断。</h2>
        <p>每个场景同时展示状态、边界、失败语义和脱敏测试证据。面试官能快速理解，钱包开发者可以继续向下钻取。</p>
      </div>
    </section>

    <section class="scenario-section shell">
      <div class="section-heading-row">
        <div><p class="section-index">02 / FIRST RELEASE</p><h2>三条实验，一套资金语义</h2></div>
        <p>首版只覆盖提现，把一个闭环做深。</p>
      </div>
      <div class="home-scenario-grid">
        <router-link v-for="scenario in scenarios" :key="scenario.id" :to="`/lab/${scenario.id}`" :class="`tone-${scenario.accent}`">
          <span>{{ scenario.index }}</span>
          <small>{{ scenario.eyebrow }}</small>
          <h3>{{ scenario.title }}</h3>
          <p>{{ scenario.summary }}</p>
          <div><em>{{ scenario.steps.length }} events</em><strong>运行实验 →</strong></div>
        </router-link>
      </div>
    </section>

    <section class="home-boundary shell">
      <div class="boundary-copy"><p class="section-index">03 / SYSTEM BOUNDARY</p><h2>四个服务，各自只守住一类责任</h2><p>资金编排、风险审批、链交互和密钥签名不混在同一个进程里。实验会高亮当前责任方，也会明确最终事实来自哪里。</p><router-link to="/architecture">查看完整系统边界 →</router-link></div>
      <div class="boundary-list">
        <div><span>01</span><strong>wallet-service</strong><small>业务状态、幂等、worker、补偿</small></div>
        <div><span>02</span><strong>risk-service</strong><small>交易校验与审批凭证</small></div>
        <div><span>03</span><strong>wallet-api</strong><small>多链资源、构建、查询与广播</small></div>
        <div><span>04</span><strong>wallet-sign</strong><small>策略复核、密钥隔离与签名</small></div>
      </div>
    </section>
  </div>
</template>

