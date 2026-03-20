<template>
  <el-container class="layout-container">
    <el-aside width="220px" class="sidebar">
      <div class="logo">
        <span class="logo-icon">⚖️</span>
        <span class="logo-text">FMPS</span>
      </div>
      <el-menu
        :default-active="currentPath"
        router
        background-color="#2c3e50"
        text-color="#bdc3c7"
        active-text-color="#e6722e"
      >
        <el-menu-item index="/">
          <el-icon><Odometer /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/members">
          <el-icon><User /></el-icon>
          <span>家庭成员</span>
        </el-menu-item>
        <el-menu-item index="/clauses">
          <el-icon><Document /></el-icon>
          <span>惩罚条款</span>
        </el-menu-item>
        <el-menu-item index="/templates">
          <el-icon><Collection /></el-icon>
          <span>惩罚模板</span>
        </el-menu-item>
        <el-menu-item index="/cases">
          <el-icon><Folder /></el-icon>
          <span>案件管理</span>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <span>系统设置</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-main class="main-content">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const currentPath = computed(() => route.path)
</script>

<style scoped>
.layout-container { height: 100vh; }
.sidebar {
  background: #2c3e50;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border-bottom: 1px solid #3d5166;
}
.logo-icon { font-size: 24px; }
.logo-text {
  font-size: 20px;
  font-weight: bold;
  color: #e6722e;
  letter-spacing: 2px;
}
.el-menu { border-right: none; flex: 1; }
.main-content { padding: 24px; overflow-y: auto; background: #f5f6fa; }
</style>
