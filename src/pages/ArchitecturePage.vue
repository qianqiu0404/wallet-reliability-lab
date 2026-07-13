<script setup lang="ts">
const services = [
  { index: '01', name: 'wallet-service', title: '资金编排层', owns: '业务状态、request_id 幂等、数据库、worker、余额与补偿', avoids: '不直接接触私钥，不把链 SDK 差异带进业务状态机' },
  { index: '02', name: 'risk-service', title: '风险审批层', owns: '交易内容校验、离线交易一致性、risk approval 与重复请求识别', avoids: '不构建交易，不负责签名，不推进资金终态' },
  { index: '03', name: 'wallet-api', title: '多链网关层', owns: '链资源查询、交易构建、signed raw transaction 拼装、广播与链上查询', avoids: '不保存业务余额，不拥有私钥，不决定业务完成' },
  { index: '04', name: 'wallet-sign', title: '签名隔离层', owns: '双重审批校验、密钥存储边界、摘要签名与多链签名适配', avoids: '不接受盲签，不连接业务数据库，不广播交易' },
]
</script>

<template>
  <div class="content-page shell" data-screen-label="Architecture">
    <header class="content-hero"><p class="eyebrow">Architecture / explicit ownership</p><h1>边界清楚，故障才有归属</h1><p>Wallet Reliability Lab 不把四个仓库包装成一个黑盒。每次实验都指出当前责任方、可观察证据和最终事实来源。</p></header>
    <section class="architecture-flow">
      <div class="external-node"><span>BUSINESS</span><strong>提现请求</strong></div>
      <template v-for="service in services" :key="service.name">
        <i class="architecture-arrow">→</i>
        <div class="architecture-node"><span>{{ service.index }}</span><strong>{{ service.name }}</strong><small>{{ service.title }}</small></div>
      </template>
      <i class="architecture-arrow">→</i>
      <div class="external-node"><span>CHAIN</span><strong>最终事实</strong></div>
    </section>
    <section class="boundary-table">
      <article v-for="service in services" :key="service.name">
        <header><span>{{ service.index }}</span><div><h2>{{ service.name }}</h2><p>{{ service.title }}</p></div></header>
        <dl><div><dt>OWNS</dt><dd>{{ service.owns }}</dd></div><div><dt>DOES NOT OWN</dt><dd>{{ service.avoids }}</dd></div></dl>
      </article>
    </section>
    <section class="state-machine-section">
      <div><p class="section-index">STATE MACHINE</p><h2>主路径与恢复路径重新汇合</h2><p>故障处理不会发明第二套终态。重试和补偿最终都回到同一个资金状态机。</p></div>
      <div class="state-machine-line">
        <span>review_pending</span><i>→</i><span>create_unsign</span><i>→</i><span>signed</span><i>→</i><span>broadcasted</span><i>→</i><span>wallet_done</span><i>→</i><span>notified</span>
      </div>
      <div class="recovery-lines"><p><strong>Retryable broadcast</strong><span>signed ↺ signed</span><small>复用同一 raw transaction</small></p><p><strong>Stale broadcast</strong><span>broadcasted ⇢ wallet_done</span><small>以链上查询恢复本地事实</small></p></div>
    </section>
  </div>
</template>

