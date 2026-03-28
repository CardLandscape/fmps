<template>
  <el-card shadow="never">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <span style="font-weight: 600">案件管理</span>
        <el-button type="primary" @click="openDialog()">
          <el-icon><Plus /></el-icon>
          新建案件
        </el-button>
      </div>
    </template>

    <el-table :data="cases" v-loading="loading" stripe style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="title" label="标题" />
      <el-table-column label="成员" width="120">
        <template #default="{ row }">{{ row.member?.name ?? '-' }}</template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="160">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" plain @click="goDetail(row)">详情</el-button>
          <el-button size="small" type="warning" plain @click="openDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" plain @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <!-- 创建/编辑对话框 -->
  <el-dialog
    v-model="dialogVisible"
    :title="editingId ? '编辑案件' : '新建案件'"
    width="600px"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
      <el-form-item label="标题" prop="title">
        <el-input v-model="form.title" placeholder="案件标题" />
      </el-form-item>
      <el-form-item label="成员" prop="member_id">
        <el-select v-model="form.member_id" placeholder="请选择成员" style="width:100%">
          <el-option
            v-for="m in members"
            :key="m.id"
            :label="m.name"
            :value="m.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="描述">
        <el-input v-model="form.description" type="textarea" :rows="2" placeholder="案件描述（可选）" />
      </el-form-item>
      <el-form-item label="惩罚过程">
        <el-input
          v-model="form.punishment_process"
          type="textarea"
          :rows="8"
          placeholder="每行一条步骤，格式：开始时间|持续分钟|惩罚内容|要求1~5（可空）|扣分规则|扣分值"
        />
        <div style="font-size:12px;color:#909399;margin-top:4px">
          格式：开始时间|持续分钟|惩罚内容|要求1~5（可空）|扣分规则|扣分值
        </div>
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
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCases, createCase, updateCase, deleteCase, getMembers } from '@/utils/api'

const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const cases = ref([])
const members = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const defaultForm = () => ({
  title: '',
  member_id: null,
  description: '',
  punishment_process: ''
})
const form = reactive(defaultForm())

const formRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  member_id: [{ required: true, message: '请选择成员', trigger: 'change' }]
}

function statusLabel(s) {
  return { pending: '待执行', active: '执行中', completed: '已完成' }[s] ?? s
}
function statusTagType(s) {
  return { pending: 'info', active: 'success', completed: '' }[s] ?? 'info'
}
function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

async function load() {
  loading.value = true
  try {
    const [cRes, mRes] = await Promise.all([getCases(), getMembers()])
    cases.value = cRes.data ?? []
    members.value = mRes.data ?? []
  } catch {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

function openDialog(row = null) {
  Object.assign(form, defaultForm())
  if (row) {
    editingId.value = row.id
    form.title = row.title ?? ''
    form.member_id = row.member_id ?? null
    form.description = row.description ?? ''
    form.punishment_process = row.punishment_process ?? ''
  } else {
    editingId.value = null
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
      await updateCase(editingId.value, { ...form })
      ElMessage.success('更新成功')
    } else {
      await createCase({ ...form })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    await load()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '操作失败')
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm(`确定要删除案件"${row.title}"吗？`, '确认删除', {
    type: 'warning', confirmButtonText: '确定', cancelButtonText: '取消'
  })
  try {
    await deleteCase(row.id)
    ElMessage.success('删除成功')
    await load()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '删除失败')
  }
}

function goDetail(row) {
  router.push(`/cases/${row.id}`)
}

onMounted(load)
</script>
