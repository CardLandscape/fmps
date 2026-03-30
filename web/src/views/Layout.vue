<template>
  <el-container style="height: 100vh">
    <el-aside width="220px" class="sidebar">
      <div class="logo">
        <el-icon :size="28" color="#fff"><House /></el-icon>
        <span>FMPS</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><Odometer /></el-icon>
          <span>{{ t('nav.dashboard') }}</span>
        </el-menu-item>
        <el-menu-item index="/members">
          <el-icon><User /></el-icon>
          <span>{{ t('nav.members') }}</span>
        </el-menu-item>
        <el-menu-item index="/rules">
          <el-icon><List /></el-icon>
          <span>{{ t('nav.rules') }}</span>
        </el-menu-item>
        <el-menu-item index="/records">
          <el-icon><Document /></el-icon>
          <span>{{ t('nav.records') }}</span>
        </el-menu-item>
        <el-menu-item index="/cases">
          <el-icon><Folder /></el-icon>
          <span>{{ t('nav.cases') }}</span>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <span>{{ t('nav.settings') }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <span class="title">{{ t('app.title') }}</span>
        <div style="display: flex; align-items: center; gap: 8px">
          <el-button size="small" plain @click="toggleLang">{{ t('lang.toggle') }}</el-button>
          <el-button type="danger" plain size="small" @click="handleLogout">
            <el-icon><SwitchButton /></el-icon>
            {{ t('nav.logout') }}
          </el-button>
        </div>
      </el-header>
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()

const activeMenu = computed(() => route.path)

function toggleLang() {
  const next = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
  locale.value = next
  localStorage.setItem('fmps_lang', next)
}

function handleLogout() {
  localStorage.removeItem('fmps_token')
  router.push('/login')
}
</script>

<style scoped>
.sidebar {
  background-color: #304156;
  overflow: hidden;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 20px 24px;
  color: #fff;
  font-size: 20px;
  font-weight: bold;
  border-bottom: 1px solid #3f5679;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.main-content {
  background: #f0f2f5;
  padding: 24px;
}

:deep(.el-menu) {
  border-right: none;
}
</style>
