<template>
  <el-steps :active="active" finish-status="success" align-center style="margin-bottom: 30px">
    <el-step v-for="step in steps" :key="step" :title="step" />
  </el-steps>
</template>

<script setup>
defineProps({
  active: { type: Number, default: 0 },
  steps: { type: Array, default: () => [] },
})
</script>
