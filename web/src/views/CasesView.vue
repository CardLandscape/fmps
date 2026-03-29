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
      <el-table-column label="级别" width="80">
        <template #default="{ row }">
          <el-tag v-if="row.punishment_level" size="small" type="warning">{{ row.punishment_level }}级</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="家长" width="120">
        <template #default="{ row }">
          {{ getMemberName(row.parent_member) }}
        </template>
      </el-table-column>
      <el-table-column label="小孩" width="120">
        <template #default="{ row }">
          {{ getMemberName(row.child_member) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="最终成绩" width="110">
        <template #default="{ row }">
          <el-tag v-if="row.final_grade" :type="gradeTagType(row.final_grade)" size="small">{{ row.final_grade }}</el-tag>
          <span v-else>-</span>
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
    width="620px"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
      <el-form-item label="标题" prop="title">
        <el-input v-model="form.title" placeholder="案件标题" />
      </el-form-item>
      <el-form-item label="家长" prop="parent_member_id">
        <el-select v-model="form.parent_member_id" placeholder="请选择家长成员" style="width:100%">
          <el-option
            v-for="m in parentMembers"
            :key="m.id"
            :label="getMemberName(m)"
            :value="m.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="小孩" prop="child_member_id">
        <el-select v-model="form.child_member_id" placeholder="请选择小孩成员" style="width:100%">
          <el-option
            v-for="m in childMembers"
            :key="m.id"
            :label="getMemberName(m)"
            :value="m.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="描述">
        <el-input v-model="form.description" type="textarea" :rows="2" placeholder="案件描述（可选）" />
      </el-form-item>
      <el-form-item label="惩罚级别" prop="punishment_level">
        <el-select v-model="form.punishment_level" placeholder="请选择惩罚级别" style="width:100%">
          <el-option label="A 级（最轻）" value="A" />
          <el-option label="B 级" value="B" />
          <el-option label="C 级" value="C" />
          <el-option label="D 级（最重）" value="D" />
        </el-select>
      </el-form-item>
      <el-form-item label="惩罚流程">
        <div style="width:100%">
          <div style="display:flex;gap:8px;margin-bottom:6px;align-items:center">
            <el-button size="small" @click="triggerTxtImport">
              <el-icon><Upload /></el-icon>
              从TXT导入
            </el-button>
            <input
              ref="txtFileInput"
              type="file"
              accept=".txt"
              style="display:none"
              @change="handleTxtImport"
            />
            <span v-if="form.prep_items.length || form.exec_steps.length" style="color:#67c23a;font-size:13px">
              已导入：{{ form.prep_items.length }} 项准备物品，{{ form.exec_steps.length }} 个执行步骤
            </span>
          </div>
          <div style="font-size:12px;color:#909399;margin-bottom:8px">
            TXT 格式：用 <code>[准备物品]</code> 和 <code>[执行步骤]</code> 分隔区块，每行一条（空行忽略）
          </div>
          <!-- Preview prep items -->
          <div v-if="form.prep_items.length" style="margin-bottom:8px">
            <div style="font-size:13px;font-weight:600;margin-bottom:4px">准备物品（{{ form.prep_items.length }}）：</div>
            <el-tag
              v-for="(item, i) in form.prep_items"
              :key="i"
              size="small"
              style="margin:2px"
            >{{ item }}</el-tag>
          </div>
          <!-- Preview exec steps -->
          <div v-if="form.exec_steps.length">
            <div style="font-size:13px;font-weight:600;margin-bottom:4px">执行步骤（{{ form.exec_steps.length }}）：</div>
            <ol style="margin:0;padding-left:20px;font-size:13px;color:#606266">
              <li v-for="(step, i) in form.exec_steps" :key="i" style="margin-bottom:2px">{{ step }}</li>
            </ol>
          </div>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Upload } from '@element-plus/icons-vue'
import { getCases, createCase, updateCase, deleteCase, getMembers } from '@/utils/api'

const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const cases = ref([])
const members = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref(null)
const txtFileInput = ref(null)

const defaultForm = () => ({
  title: '',
  parent_member_id: null,
  child_member_id: null,
  description: '',
  punishment_level: '',
  prep_items: [],
  exec_steps: []
})
const form = reactive(defaultForm())

// Filtered member lists for dropdowns
const parentMembers = computed(() =>
  members.value.filter(m => m.role === 'parent' || m.role === 'adult')
)
const childMembers = computed(() =>
  members.value.filter(m => m.role === 'child')
)

const formRules = computed(() => ({
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  parent_member_id: [{ required: true, message: '请选择家长成员', trigger: 'change' }],
  child_member_id: [
    { required: true, message: '请选择小孩成员', trigger: 'change' },
    {
      validator: (rule, value, callback) => {
        if (value && form.parent_member_id) {
          if (value === form.parent_member_id) {
            callback(new Error('家长和小孩不能是同一条记录'))
            return
          }
        }
        callback()
      },
      trigger: 'change'
    }
  ]
}))

function getMemberName(m) {
  if (!m) return '-'
  return m.name_cn || m.name_en || m.name || '-'
}

function statusLabel(s) {
  return { pending: '待执行', active: '执行中', completed: '已完成' }[s] ?? s
}
function statusTagType(s) {
  return { pending: 'info', active: 'success', completed: '' }[s] ?? 'info'
}
function gradeTagType(g) {
  if (g === '满分' || g === '优') return 'success'
  if (g === '良' || g === '达标') return 'warning'
  return 'danger'
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

function parseTxtSections(text) {
  const lines = text.replace(/\r\n/g, '\n').split('\n')
  const prepItems = []
  const execSteps = []
  let section = null
  for (const rawLine of lines) {
    const line = rawLine.trim()
    if (line === '[准备物品]') { section = 'prep'; continue }
    if (line === '[执行步骤]') { section = 'exec'; continue }
    if (!line) continue
    if (section === 'prep') prepItems.push(line)
    else if (section === 'exec') execSteps.push(line)
  }
  return { prepItems, execSteps }
}

function openDialog(row = null) {
  Object.assign(form, defaultForm())
  if (row) {
    editingId.value = row.id
    form.title = row.title ?? ''
    form.parent_member_id = row.parent_member_id || null
    form.child_member_id = row.child_member_id || null
    form.description = row.description ?? ''
    form.punishment_level = row.punishment_level ?? ''
    try { form.prep_items = row.prep_items ? JSON.parse(row.prep_items) : [] } catch { form.prep_items = [] }
    try { form.exec_steps = row.exec_steps ? JSON.parse(row.exec_steps) : [] } catch { form.exec_steps = [] }
  } else {
    editingId.value = null
  }
  dialogVisible.value = true
}

function resetForm() {
  formRef.value?.resetFields()
  editingId.value = null
}

function triggerTxtImport() {
  txtFileInput.value?.click()
}

function handleTxtImport(event) {
  const file = event.target.files?.[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = (e) => {
    const text = e.target?.result ?? ''
    const { prepItems, execSteps } = parseTxtSections(text)
    // Additive import: append to existing items (dedup by value)
    const newPrep = [...new Set([...form.prep_items, ...prepItems])]
    const newExec = [...form.exec_steps, ...execSteps]
    form.prep_items = newPrep
    form.exec_steps = newExec
    ElMessage.success(`已导入 ${prepItems.length} 项准备物品，${execSteps.length} 个执行步骤`)
  }
  reader.onerror = () => {
    ElMessage.error('文件读取失败')
  }
  reader.readAsText(file, 'utf-8')
  event.target.value = ''
}

async function handleSave() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    const payload = {
      title: form.title,
      parent_member_id: form.parent_member_id,
      child_member_id: form.child_member_id,
      description: form.description,
      punishment_level: form.punishment_level,
      prep_items: JSON.stringify(form.prep_items),
      exec_steps: JSON.stringify(form.exec_steps)
    }
    if (editingId.value) {
      await updateCase(editingId.value, payload)
      ElMessage.success('更新成功')
    } else {
      await createCase(payload)
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
