<template>
  <div v-if="endTime" class="countdown">
    <el-statistic title="剩余时间" :value="remaining">
      <template #suffix>秒</template>
    </el-statistic>
    <el-progress :percentage="progress" :color="progressColor" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  startedAt: String,
  durationMinutes: { type: Number, default: 0 },
})

const now = ref(Date.now())
let timer

onMounted(() => { timer = setInterval(() => { now.value = Date.now() }, 1000) })
onUnmounted(() => clearInterval(timer))

const endTime = computed(() => {
  if (!props.startedAt || !props.durationMinutes) return null
  const start = new Date(props.startedAt.replace(' ', 'T')).getTime()
  return start + props.durationMinutes * 60 * 1000
})

const remaining = computed(() => {
  if (!endTime.value) return 0
  return Math.max(0, Math.floor((endTime.value - now.value) / 1000))
})

const progress = computed(() => {
  if (!props.durationMinutes) return 0
  const total = props.durationMinutes * 60
  const elapsed = total - remaining.value
  return Math.min(100, Math.floor((elapsed / total) * 100))
})

const progressColor = computed(() => {
  if (remaining.value === 0) return '#67c23a'
  if (progress.value > 70) return '#e6a23c'
  return '#e6722e'
})
</script>

<style scoped>
.countdown { padding: 12px 0; }
</style>
