<template>
  <el-card shadow="never" style="max-width: 560px">
    <template #header>
      <span style="font-weight: 600">{{ i18n.t('settingsPageTitle') }}</span>
    </template>

    <el-form
      ref="formRef"
      :model="form"
      :rules="formRules"
      label-width="140px"
      v-loading="loading"
    >
      <el-form-item :label="i18n.t('labelAdminUsername')" prop="admin_username">
        <el-input v-model="form.admin_username" :placeholder="i18n.t('placeholderAdminUsername')" />
      </el-form-item>
      <el-form-item :label="i18n.t('labelNewPassword')" prop="new_password">
        <el-input
          v-model="form.new_password"
          type="password"
          :placeholder="i18n.t('placeholderNewPassword')"
          show-password
        />
      </el-form-item>
      <el-form-item :label="i18n.t('labelConfirmPassword')" prop="confirm_password">
        <el-input
          v-model="form.confirm_password"
          type="password"
          :placeholder="i18n.t('placeholderConfirmPassword')"
          show-password
        />
      </el-form-item>

      <el-alert
        :title="i18n.t('hintNoChangePassword')"
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      />

      <el-divider content-position="left">{{ i18n.t('sectionAuthPassword') }}</el-divider>
      <el-form-item :label="i18n.t('labelAuthPasswordSetting')">
        <el-input
          v-model="form.authorization_password"
          type="password"
          show-password
          :placeholder="i18n.t('placeholderAuthPasswordSetting')"
        />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :loading="saving" @click="handleSave">{{ i18n.t('btnSaveSettings') }}</el-button>
        <el-button @click="loadSettings">{{ i18n.t('btnResetSettings') }}</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getSettings, updateSettings } from '@/utils/api'
import { useI18n } from '@/utils/i18n'

const i18n = useI18n()

const loading = ref(false)
const saving = ref(false)
const formRef = ref(null)

const form = reactive({
  admin_username: '',
  new_password: '',
  confirm_password: '',
  authorization_password: ''
})

const formRules = computed(() => {
  const validateConfirmPassword = (_rule, value, callback) => {
    if (form.new_password && value !== form.new_password) {
      callback(new Error(i18n.t('validPasswordsMismatch')))
    } else {
      callback()
    }
  }
  return {
    admin_username: [{ required: true, message: i18n.t('validAdminUsernameRequired'), trigger: 'blur' }],
    confirm_password: [{ validator: validateConfirmPassword, trigger: 'blur' }]
  }
})

async function loadSettings() {
  loading.value = true
  try {
    const res = await getSettings()
    form.admin_username = res.data?.admin_username ?? ''
    form.new_password = ''
    form.confirm_password = ''
    form.authorization_password = ''
  } catch {
    ElMessage.error(i18n.t('loadSettingsFailed'))
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
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
    ElMessage.success(i18n.t('settingsSaved'))
    form.new_password = ''
    form.confirm_password = ''
    form.authorization_password = ''
  } catch (e) {
    ElMessage.error(e.response?.data?.error || i18n.t('saveSettingsFailed'))
  } finally {
    saving.value = false
  }
}

onMounted(loadSettings)
</script>
