<template>
  <div>
    <div class="page-header"><h2>⚙️ 系统设置</h2></div>
    <el-card style="max-width: 600px">
      <el-form :model="settings" label-width="160px">
        <el-form-item label="系统名称">
          <el-input v-model="settings.system_name" placeholder="FMPS" />
        </el-form-item>
        <el-form-item label="案件完成后保护天数">
          <el-input-number v-model="settings.protection_days" :min="1" :max="365" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">保存设置</el-button>
        </el-form-item>
      </el-form>
      <el-divider />
      <div style="color: #909399; font-size: 13px">
        <p>版本: FMPS v1.0.0</p>
        <p>技术栈: Go + Wails v2 + Vue 3 + SQLite</p>
        <p>数据库: 与程序同目录下 data.db</p>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../utils/api'
import { ElMessage } from 'element-plus'

const settings = ref({ system_name: 'FMPS', protection_days: 30 })

onMounted(async () => {
  try {
    const s = await api.getAllSettings()
    if (s.system_name) settings.value.system_name = s.system_name
    if (s.protection_days) settings.value.protection_days = Number(s.protection_days)
  } catch (e) {}
})

const save = async () => {
  try {
    await Promise.all([
      api.setSetting('system_name', settings.value.system_name),
      api.setSetting('protection_days', String(settings.value.protection_days)),
    ])
    ElMessage.success('设置已保存')
  } catch (e) { ElMessage.error('保存失败: ' + e) }
}
</script>
