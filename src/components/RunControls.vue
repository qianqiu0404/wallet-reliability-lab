<script setup lang="ts">
defineProps<{
  playing: boolean
  complete: boolean
  speed: number
  progress: number
}>()

const emit = defineEmits<{
  start: []
  pause: []
  step: []
  reset: []
  speed: [value: number]
}>()
</script>

<template>
  <div class="run-controls">
    <div class="run-actions">
      <button v-if="!playing" class="control-button primary" type="button" @click="emit('start')">
        <span aria-hidden="true">▶</span>{{ complete ? '重新运行' : '开始运行' }}
      </button>
      <button v-else class="control-button primary" type="button" @click="emit('pause')"><span aria-hidden="true">Ⅱ</span>暂停</button>
      <button class="control-button" type="button" :disabled="playing || complete" @click="emit('step')">单步 →</button>
      <button class="icon-button" type="button" aria-label="重置实验" title="重置实验" @click="emit('reset')">↺</button>
    </div>
    <div class="speed-control" aria-label="运行速度">
      <span>Speed</span>
      <button v-for="item in [1, 2, 4]" :key="item" type="button" :class="{ active: speed === item }" @click="emit('speed', item)">{{ item }}×</button>
    </div>
    <div class="progress-track" aria-label="实验进度"><span :style="{ width: `${progress}%` }"></span></div>
  </div>
</template>

