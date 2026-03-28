<template>
  <el-card shadow="never" style="max-width: 560px">
    <template #header>
      <span style="font-weight: 600">系统设置</span>
    </template>

    <el-form
      ref="formRef"
      :model="form"
      :rules="formRules"
      label-width="120px"
      v-loading="loading"
    >
      <el-form-item label="管理员用户名" prop="admin_username">
        <el-input v-model="form.admin_username" placeholder="管理员用户名" />
      </el-form-item>
      <el-form-item label="新密码" prop="new_password">
        <el-input
          v-model="form.new_password"
          type="password"
          placeholder="留空则不修改密码"
          show-password
        />
      </el-form-item>
      <el-form-item label="确认新密码" prop="confirm_password">
        <el-input
          v-model="form.confirm_password"
          type="password"
          placeholder="再次输入新密码"
          show-password
        />
      </el-form-item>

      <el-alert
        title="留空密码字段则不修改密码"
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      />

      <el-form-item>
        <el-button type="primary" :loading="saving" @click="handleSave">保存设置</el-button>
        <el-button @click="loadSettings">重 置</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getSettings, updateSettings } from '@/utils/api'

const loading = ref(false)
const saving = ref(false)
const formRef = ref(null)

const form = reactive({
  admin_username: '',
  new_password: '',
  confirm_password: ''
})

const validateConfirmPassword = (_rule, value, callback) => {
  if (form.new_password && value !== form.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const formRules = {
  admin_username: [{ required: true, message: '请输入管理员用户名', trigger: 'blur' }],
  confirm_password: [{ validator: validateConfirmPassword, trigger: 'blur' }]
}

async function loadSettings() {
  loading.value = true
  try {
    const res = await getSettings()
    form.admin_username = res.data?.admin_username ?? ''
    form.new_password = ''
    form.confirm_password = ''
  } catch {
    ElMessage.error('加载设置失败')
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  if (form.new_password && form.new_password !== form.confirm_password) {
    ElMessage.error('两次输入的密码不一致')
    return
  }

  saving.value = true
  try {
    const payload = { admin_username: form.admin_username }
    if (form.new_password) {
      payload.admin_password = form.new_password
    }
    await updateSettings(payload)
    ElMessage.success('设置保存成功')
    form.new_password = ''
    form.confirm_password = ''
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(loadSettings)
</script>
