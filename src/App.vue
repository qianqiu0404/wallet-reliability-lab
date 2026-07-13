<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const menuOpen = ref(false)
const xiuqiuSiteUrl = import.meta.env.VITE_XIUQIU_SITE_URL || 'https://xiuqiu-site.vercel.app'
</script>

<template>
  <header class="site-header">
    <nav class="site-nav shell" aria-label="Main navigation">
      <router-link class="brand-lockup" to="/" @click="menuOpen = false">
        <span class="brand-mark" aria-hidden="true"><i></i><i></i><i></i></span>
        <span><strong>Wallet Reliability Lab</strong><small>by xiuqiu</small></span>
      </router-link>
      <button class="nav-toggle" type="button" :aria-expanded="menuOpen" aria-label="打开导航" @click="menuOpen = !menuOpen">
        <span></span><span></span>
      </button>
      <div class="nav-links" :class="{ open: menuOpen }">
        <router-link to="/lab/normal-withdrawal" @click="menuOpen = false">实验室</router-link>
        <router-link to="/architecture" @click="menuOpen = false">系统边界</router-link>
        <router-link to="/evidence" @click="menuOpen = false">证据索引</router-link>
        <a :href="xiuqiuSiteUrl" target="_blank" rel="noreferrer">xiuqiu-site ↗</a>
      </div>
      <span v-if="route.name === 'lab'" class="header-live"><i></i> Simulation safe</span>
    </nav>
  </header>

  <main>
    <router-view />
  </main>

  <footer class="site-footer">
    <div class="shell footer-inner">
      <span>Wallet Reliability Lab · deterministic simulation</span>
      <a :href="xiuqiuSiteUrl" target="_blank" rel="noreferrer">by xiuqiu</a>
    </div>
  </footer>
</template>

