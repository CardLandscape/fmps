<template>
  <el-card shadow="never">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <span style="font-weight: 600">家庭成员管理</span>
        <el-button type="primary" @click="openDialog()">
          <el-icon><Plus /></el-icon>
          添加成员
        </el-button>
      </div>
    </template>

    <el-table :data="members" v-loading="loading" stripe style="width: 100%">
      <el-table-column prop="name" label="姓名" />
      <el-table-column prop="role" label="角色" width="100">
        <template #default="{ row }">
          <el-tag :type="row.role === 'child' ? 'warning' : 'success'">
            {{ row.role === 'child' ? '孩子' : '成人' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
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
    :title="editingId ? '编辑成员' : '添加成员'"
    width="440px"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
      <el-form-item label="姓名" prop="name">
        <el-input v-model="form.name" placeholder="请输入姓名" />
      </el-form-item>
      <el-form-item label="角色" prop="role">
        <el-select v-model="form.role" placeholder="请选择角色" style="width: 100%">
          <el-option label="孩子" value="child" />
          <el-option label="成人" value="adult" />
        </el-select>
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
import { getMembers, createMember, updateMember, deleteMember } from '@/utils/api'

const loading = ref(false)
const saving = ref(false)
const members = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const form = reactive({ name: '', role: 'child' })

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }]
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

async function loadMembers() {
  loading.value = true
  try {
    const res = await getMembers()
    members.value = res.data ?? []
  } catch {
    ElMessage.error('加载成员失败')
  } finally {
    loading.value = false
  }
}

function openDialog(row = null) {
  if (row) {
    editingId.value = row.id
    form.name = row.name
    form.role = row.role
  } else {
    editingId.value = null
    form.name = ''
    form.role = 'child'
  }
  dialogVisible.value = true
}

function resetForm() {
  formRef.value?.resetFields()
  editingId.value = null
}

async function handleSave() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (editingId.value) {
      await updateMember(editingId.value, { ...form })
      ElMessage.success('更新成功')
    } else {
      await createMember({ ...form })
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    await loadMembers()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '操作失败')
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm(`确定要删除成员"${row.name}"吗？`, '确认删除', {
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  })
  try {
    await deleteMember(row.id)
    ElMessage.success('删除成功')
    await loadMembers()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '删除失败')
  }
}

onMounted(loadMembers)
</script>
