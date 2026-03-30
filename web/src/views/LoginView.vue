<template>
  <div class="login-container">
    <el-card class="login-card">
      <!-- Language toggle -->
      <div class="lang-toggle">
        <el-button size="small" plain @click="toggleLang">{{ t('lang.toggle') }}</el-button>
      </div>

      <div class="login-header">
        <el-icon :size="48" color="#409EFF"><House /></el-icon>
        <h1>{{ t('login.title') }}</h1>
        <p class="subtitle">{{ t('login.subtitle') }}</p>
      </div>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        @submit.prevent="handleLogin"
      >
        <el-form-item :label="t('login.usernameLabel')" prop="username">
          <el-input
            v-model="form.username"
            :placeholder="t('login.usernamePlaceholder')"
            prefix-icon="User"
            size="large"
          />
        </el-form-item>
        <el-form-item :label="t('login.passwordLabel')" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            :placeholder="t('login.passwordPlaceholder')"
            prefix-icon="Lock"
            size="large"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-alert
          v-if="errorMsg"
          :title="errorMsg"
          type="error"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            style="width: 100%"
            :loading="loading"
            @click="handleLogin"
          >
            {{ t('login.loginButton') }}
          </el-button>
        </el-form-item>
      </el-form>

      <el-divider />
      <p class="hint">{{ t('login.defaultHint') }}</p>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { login } from '@/utils/api'

const router = useRouter()
const { t, locale } = useI18n()
const formRef = ref(null)
const loading = ref(false)
const errorMsg = ref('')

const form = reactive({
  username: '',
  password: ''
})

const rules = computed(() => ({
  username: [{ required: true, message: t('login.usernameRequired'), trigger: 'blur' }],
  password: [{ required: true, message: t('login.passwordRequired'), trigger: 'blur' }]
}))

function toggleLang() {
  const next = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
  locale.value = next
  localStorage.setItem('fmps_lang', next)
}

async function handleLogin() {
  if (!formRef.value) return
  await formRef.value.validate(async valid => {
    if (!valid) return
    loading.value = true
    errorMsg.value = ''
    try {
      const res = await login(form.username, form.password)
      const token = res.data.token
      localStorage.setItem('fmps_token', token)
      router.push('/dashboard')
    } catch (err) {
      errorMsg.value =
        err.response?.data?.message || err.response?.data?.error || t('login.loginFailed')
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 420px;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  position: relative;
}

.lang-toggle {
  position: absolute;
  top: 16px;
  right: 16px;
  z-index: 1;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-header h1 {
  margin: 12px 0 4px;
  font-size: 22px;
  color: #303133;
}

.subtitle {
  color: #909399;
  font-size: 13px;
  margin: 0;
}

.hint {
  text-align: center;
  color: #909399;
  font-size: 13px;
  margin: 0;
}
</style>
