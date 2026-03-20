<template>
  <div>
    <div class="page-header">
      <h2>📦 惩罚模板</h2>
      <el-button type="primary" @click="showForm()">
        <el-icon><Plus /></el-icon> 新增模板
      </el-button>
    </div>
    <el-card>
      <el-table :data="templates" stripe>
        <el-table-column prop="name" label="模板名称" />
        <el-table-column prop="punishment_type" label="惩罚类型" width="120" />
        <el-table-column prop="duration_minutes" label="时长(分钟)" width="100" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button link type="primary" @click="showForm(row)">编辑</el-button>
            <el-button link type="danger" @click="deleteTemplate(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑模板' : '新增模板'" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="模板名称" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="惩罚类型" required>
          <el-select v-model="form.punishment_type">
            <el-option v-for="t in types" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <el-form-item label="时长(分钟)">
          <el-input-number v-model="form.duration_minutes" :min="0" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../utils/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const templates = ref([])
const dialogVisible = ref(false)
const types = ['口头警告', '书面检讨', '减少零用钱', '增加家务', '限制娱乐', '其它']
const emptyForm = () => ({ name: '', punishment_type: '口头警告', duration_minutes: 0, description: '' })
const form = ref(emptyForm())

const load = async () => {
  try { templates.value = await api.getAllTemplates() }
  catch (e) { ElMessage.error('加载失败: ' + e) }
}

const showForm = (row = null) => {
  form.value = row ? { ...row } : emptyForm()
  dialogVisible.value = true
}

const save = async () => {
  if (!form.value.name) { ElMessage.warning('请填写模板名称'); return }
  try {
    if (form.value.id) {
      await api.updateTemplate(form.value)
    } else {
      await api.createTemplate(form.value)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    load()
  } catch (e) { ElMessage.error('保存失败: ' + e) }
}

const deleteTemplate = async (row) => {
  try {
    await ElMessageBox.confirm(`确认删除模板「${row.name}」？`, '提示', { type: 'warning' })
    await api.deleteTemplate(row.id)
    ElMessage.success('删除成功')
    load()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败: ' + e)
  }
}

onMounted(load)
</script>
