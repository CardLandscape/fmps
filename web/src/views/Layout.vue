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
          <span>{{ i18n.t('navDashboard') }}</span>
        </el-menu-item>
        <el-menu-item index="/members">
          <el-icon><User /></el-icon>
          <span>{{ i18n.t('navMembers') }}</span>
        </el-menu-item>
        <el-menu-item index="/rules">
          <el-icon><List /></el-icon>
          <span>{{ i18n.t('navRules') }}</span>
        </el-menu-item>
        <el-menu-item index="/records">
          <el-icon><Document /></el-icon>
          <span>{{ i18n.t('navRecords') }}</span>
        </el-menu-item>
        <el-menu-item index="/cases">
          <el-icon><Folder /></el-icon>
          <span>{{ i18n.t('navCases') }}</span>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <span>{{ i18n.t('navSettings') }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <span class="title">{{ i18n.t('appTitle') }}</span>
        <div style="display: flex; align-items: center; gap: 12px">
          <el-button size="small" text @click="toggleLang">{{ i18n.t('switchLang') }}</el-button>
          <el-button type="danger" plain size="small" @click="handleLogout">
            <el-icon><SwitchButton /></el-icon>
            {{ i18n.t('logout') }}
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
import { useI18n, setLang, LANGUAGES } from '@/utils/i18n'

const route = useRoute()
const router = useRouter()

const i18n = useI18n()

const activeMenu = computed(() => route.path)

function toggleLang() {
  const other = LANGUAGES.find(l => l.code !== i18n.lang.value)
  if (other) setLang(other.code)
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
