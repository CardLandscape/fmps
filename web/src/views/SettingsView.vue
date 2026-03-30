<template>
  <el-card shadow="never" style="max-width: 560px">
    <template #header>
      <span style="font-weight: 600">{{ t('settings.title') }}</span>
    </template>

    <el-form
      ref="formRef"
      :model="form"
      :rules="formRules"
      label-width="120px"
      v-loading="loading"
    >
      <el-form-item :label="t('settings.adminUsername')" prop="admin_username">
        <el-input v-model="form.admin_username" :placeholder="t('settings.adminUsernamePlaceholder')" />
      </el-form-item>
      <el-form-item :label="t('settings.newPassword')" prop="new_password">
        <el-input
          v-model="form.new_password"
          type="password"
          :placeholder="t('settings.newPasswordPlaceholder')"
          show-password
        />
      </el-form-item>
      <el-form-item :label="t('settings.confirmPassword')" prop="confirm_password">
        <el-input
          v-model="form.confirm_password"
          type="password"
          :placeholder="t('settings.confirmPasswordPlaceholder')"
          show-password
        />
      </el-form-item>

      <el-alert
        :title="t('settings.passwordHint')"
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      />

      <el-divider content-position="left">{{ t('settings.authPasswordSection') }}</el-divider>
      <el-form-item :label="t('settings.authPassword')">
        <el-input
          v-model="form.authorization_password"
          type="password"
          show-password
          :placeholder="t('settings.authPasswordPlaceholder')"
        />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :loading="saving" @click="handleSave">{{ t('settings.saveBtn') }}</el-button>
        <el-button @click="loadSettings">{{ t('settings.resetBtn') }}</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { getSettings, updateSettings } from '@/utils/api'

const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)
const formRef = ref(null)

const form = reactive({
  admin_username: '',
  new_password: '',
  confirm_password: '',
  authorization_password: ''
})

const validateConfirmPassword = (_rule, value, callback) => {
  if (form.new_password && value !== form.new_password) {
    callback(new Error(t('settings.passwordMismatch')))
  } else {
    callback()
  }
}

const formRules = computed(() => ({
  admin_username: [{ required: true, message: t('settings.adminUsernameRequired'), trigger: 'blur' }],
  confirm_password: [{ validator: validateConfirmPassword, trigger: 'blur' }]
}))

async function loadSettings() {
  loading.value = true
  try {
    const res = await getSettings()
    form.admin_username = res.data?.admin_username ?? ''
    form.new_password = ''
    form.confirm_password = ''
    form.authorization_password = ''
  } catch {
    ElMessage.error(t('settings.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  if (form.new_password && form.new_password !== form.confirm_password) {
    ElMessage.error(t('settings.passwordMismatch'))
    return
  }

  saving.value = true
  try {
    const payload = { admin_username: form.admin_username }
    if (form.new_password) {
      payload.admin_password = form.new_password
    }
    if (form.authorization_password) {
      payload.authorization_password = form.authorization_password
    }
    await updateSettings(payload)
    ElMessage.success(t('settings.saveSuccess'))
    form.new_password = ''
    form.confirm_password = ''
  } catch (e) {
    ElMessage.error(e.response?.data?.error || t('settings.saveFailed'))
  } finally {
    saving.value = false
  }
}

onMounted(loadSettings)
</script>
