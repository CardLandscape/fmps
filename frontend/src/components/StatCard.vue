<template>
  <el-card class="stat-card" shadow="hover">
    <div class="stat-content">
      <div class="stat-info">
        <div class="stat-label">{{ label }}</div>
        <div class="stat-value">{{ value }}</div>
        <div v-if="subtitle" class="stat-subtitle">{{ subtitle }}</div>
      </div>
      <div class="stat-icon" :style="{ background: iconBg }">
        <el-icon :size="28" :color="iconColor">
          <component :is="icon" />
        </el-icon>
      </div>
    </div>
  </el-card>
</template>

<script setup>
defineProps({
  label: String,
  value: [String, Number],
  subtitle: String,
  icon: Object,
  iconBg: { type: String, default: '#fff3eb' },
  iconColor: { type: String, default: '#e6722e' },
})
</script>

<style scoped>
.stat-card { cursor: default; }
.stat-content { display: flex; justify-content: space-between; align-items: center; }
.stat-label { font-size: 13px; color: #909399; margin-bottom: 8px; }
.stat-value { font-size: 32px; font-weight: bold; color: #303133; }
.stat-subtitle { font-size: 12px; color: #c0c4cc; margin-top: 4px; }
.stat-icon { width: 60px; height: 60px; border-radius: 12px; display: flex; align-items: center; justify-content: center; }
</style>
