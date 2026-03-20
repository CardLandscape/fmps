<template>
  <el-tag v-if="isProtected" type="warning" size="small" effect="light">
    🛡️ 保护期
  </el-tag>
  <el-tag v-else type="success" size="small" effect="light">
    ✅ 正常
  </el-tag>
</template>

<script setup>
defineProps({ isProtected: Boolean })
</script>
