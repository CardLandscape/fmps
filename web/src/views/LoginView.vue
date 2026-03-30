<template>
  <div class="login-container">
    <el-card class="login-card">
      <div class="lang-toggle">
        <el-button size="small" text @click="toggleLang">{{ i18n.t('switchLang') }}</el-button>
      </div>

      <div class="login-header">
        <el-icon :size="48" color="#409EFF"><House /></el-icon>
        <h1>{{ i18n.t('appTitle') }}</h1>
        <p class="subtitle">{{ i18n.t('appSubtitle') }}</p>
      </div>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        @submit.prevent="handleLogin"
      >
        <el-form-item :label="i18n.t('loginUsername')" prop="username">
          <el-input
            v-model="form.username"
            :placeholder="i18n.t('loginUsernamePlaceholder')"
            prefix-icon="User"
            size="large"
          />
        </el-form-item>
        <el-form-item :label="i18n.t('loginPassword')" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            :placeholder="i18n.t('loginPasswordPlaceholder')"
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
            {{ i18n.t('loginButton') }}
          </el-button>
        </el-form-item>
      </el-form>

      <el-divider />
      <p class="hint">{{ i18n.t('loginDefaultHint') }}</p>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '@/utils/api'
import { useI18n, setLang, LANGUAGES } from '@/utils/i18n'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)
const errorMsg = ref('')

const i18n = useI18n()

function toggleLang() {
  const other = LANGUAGES.find(l => l.code !== i18n.lang.value)
  if (other) setLang(other.code)
}

const form = reactive({
  username: '',
  password: ''
})

const rules = computed(() => ({
  username: [{ required: true, message: i18n.t('loginUsernameRequired'), trigger: 'blur' }],
  password: [{ required: true, message: i18n.t('loginPasswordRequired'), trigger: 'blur' }]
}))

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
        err.response?.data?.message || err.response?.data?.error || i18n.t('loginFailed')
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
}

.lang-toggle {
  display: flex;
  justify-content: flex-end;
  margin-bottom: -8px;
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
