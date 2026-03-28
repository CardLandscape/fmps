<template>
  <el-card shadow="never">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <span style="font-weight: 600">惩戒规则管理</span>
        <el-button type="primary" @click="openDialog()">
          <el-icon><Plus /></el-icon>
          添加规则
        </el-button>
      </div>
    </template>

    <el-table :data="rules" v-loading="loading" stripe style="width: 100%">
      <el-table-column prop="name" label="名称" min-width="120" />
      <el-table-column prop="category" label="分类" width="120" />
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column prop="points" label="分值" width="80">
        <template #default="{ row }">
          <el-tag type="danger">{{ row.points }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" plain @click="openDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" plain @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog
    v-model="dialogVisible"
    :title="editingId ? '编辑规则' : '添加规则'"
    width="480px"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="formRules" label-width="80px">
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入规则名称" />
      </el-form-item>
      <el-form-item label="分类" prop="category">
        <el-input v-model="form.category" placeholder="如：学习、行为、礼貌" />
      </el-form-item>
      <el-form-item label="描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="规则详细描述"
        />
      </el-form-item>
      <el-form-item label="分值" prop="points">
        <el-input-number v-model="form.points" :min="1" :max="100" style="width: 100%" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取 消</el-button>
      <el-button type="primary" :loading="saving" @click="handleSave">确 定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getRules, createRule, updateRule, deleteRule } from '@/utils/api'

const loading = ref(false)
const saving = ref(false)
const rules = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const form = reactive({ name: '', category: '', description: '', points: 1 })

const formRules = {
  name: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
  points: [{ required: true, message: '请输入分值', trigger: 'blur' }]
}

async function loadRules() {
  loading.value = true
  try {
    const res = await getRules()
    rules.value = res.data ?? []
  } catch {
    ElMessage.error('加载规则失败')
  } finally {
    loading.value = false
  }
}

function openDialog(row = null) {
  if (row) {
    editingId.value = row.id
    form.name = row.name
    form.category = row.category || ''
    form.description = row.description || ''
    form.points = row.points
  } else {
    editingId.value = null
    form.name = ''
    form.category = ''
    form.description = ''
    form.points = 1
  }
  dialogVisible.value = true
}

function resetForm() {
  formRef.value?.resetFields()
  editingId.value = null
}

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (editingId.value) {
      await updateRule(editingId.value, { ...form })
      ElMessage.success('更新成功')
    } else {
      await createRule({ ...form })
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    await loadRules()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '操作失败')
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm(`确定要删除规则"${row.name}"吗？`, '确认删除', {
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  })
  try {
    await deleteRule(row.id)
    ElMessage.success('删除成功')
    await loadRules()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '删除失败')
  }
}

onMounted(loadRules)
</script>
