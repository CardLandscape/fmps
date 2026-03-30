<template>
  <el-card shadow="never">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <span style="font-weight: 600">{{ t('cases.title') }}</span>
        <el-button type="primary" @click="openDialog()">
          <el-icon><Plus /></el-icon>
          {{ t('cases.newCase') }}
        </el-button>
      </div>
    </template>

    <el-table :data="cases" v-loading="loading" stripe style="width: 100%">
      <el-table-column prop="id" :label="t('cases.colId')" width="60" />
      <el-table-column prop="title" :label="t('cases.colTitle')" />
      <el-table-column :label="t('cases.colParent')" width="120">
        <template #default="{ row }">
          {{ getMemberName(row.parent_member) }}
        </template>
      </el-table-column>
      <el-table-column :label="t('cases.colChild')" width="120">
        <template #default="{ row }">
          {{ getMemberName(row.child_member) }}
        </template>
      </el-table-column>
      <el-table-column :label="t('cases.colLevel')" width="70">
        <template #default="{ row }">
          <el-tag v-if="row.punishment_level" size="small" type="danger">{{ row.punishment_level }}{{ t('cases.levelSuffix') }}</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" :label="t('cases.colStatus')" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('cases.colGrade')" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.final_grade" :type="gradeTagType(row.final_grade)" size="small">{{ row.final_grade }}</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" :label="t('cases.colCreatedAt')" width="160">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.operation')" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" plain @click="goDetail(row)">{{ t('cases.detail') }}</el-button>
          <el-button size="small" type="warning" plain @click="openDialog(row)">{{ t('common.edit') }}</el-button>
          <el-button size="small" type="danger" plain @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <!-- 创建/编辑对话框 -->
  <el-dialog
    v-model="dialogVisible"
    :title="editingId ? t('cases.editCase') : t('cases.newCase')"
    width="640px"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
      <el-form-item :label="t('cases.title2')" prop="title">
        <el-input v-model="form.title" :placeholder="t('cases.titlePlaceholder')" />
      </el-form-item>
      <el-form-item :label="t('cases.parent')" prop="parent_member_id">
        <el-select v-model="form.parent_member_id" :placeholder="t('cases.parentPlaceholder')" style="width:100%">
          <el-option
            v-for="m in parentMembers"
            :key="m.id"
            :label="getMemberName(m)"
            :value="m.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('cases.child')" prop="child_member_id">
        <el-select v-model="form.child_member_id" :placeholder="t('cases.childPlaceholder')" style="width:100%">
          <el-option
            v-for="m in childMembers"
            :key="m.id"
            :label="getMemberName(m)"
            :value="m.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('cases.description')">
        <el-input v-model="form.description" type="textarea" :rows="2" :placeholder="t('cases.descriptionPlaceholder')" />
      </el-form-item>

      <!-- Punishment workflow section -->
      <el-divider content-position="left">{{ t('cases.punishmentSection') }}</el-divider>

      <el-form-item :label="t('cases.punishmentLevel')" prop="punishment_level">
        <el-radio-group v-model="form.punishment_level">
          <el-radio-button value="A">A</el-radio-button>
          <el-radio-button value="B">B</el-radio-button>
          <el-radio-button value="C">C</el-radio-button>
          <el-radio-button value="D">D</el-radio-button>
        </el-radio-group>
      </el-form-item>

      <el-form-item :label="t('cases.importTxt')">
        <div style="width:100%">
          <div style="display:flex;gap:8px;align-items:center;margin-bottom:6px">
            <el-button size="small" @click="triggerTxtImport" :disabled="!form.punishment_level">
              <el-icon><Upload /></el-icon>
              {{ t('cases.selectTxt') }}
            </el-button>
            <span v-if="form.txt_filename" style="font-size:12px;color:#606266">{{ form.txt_filename }}</span>
            <span v-if="!form.punishment_level" style="font-size:12px;color:#f56c6c">{{ t('cases.selectLevelFirst') }}</span>
            <input
              ref="txtFileInput"
              type="file"
              accept=".txt"
              style="display:none"
              @change="handleTxtImport"
            />
          </div>

          <!-- Parsed preview -->
          <div v-if="parsedPreview.prepItems.length || parsedPreview.steps.length">
            <div v-if="parsedPreview.prepItems.length" style="margin-bottom:8px">
              <div style="font-weight:600;font-size:13px;margin-bottom:4px">{{ t('cases.parsedPrepPrefix', { n: parsedPreview.prepItems.length }) }}</div>
              <el-tag
                v-for="(item, i) in parsedPreview.prepItems"
                :key="i"
                size="small"
                style="margin-right:6px;margin-bottom:4px"
              >{{ item }}</el-tag>
            </div>
            <div v-if="parsedPreview.steps.length">
              <div style="font-weight:600;font-size:13px;margin-bottom:4px">{{ t('cases.parsedStepsPrefix', { n: parsedPreview.steps.length }) }}</div>
              <div v-for="(step, i) in parsedPreview.steps" :key="i" style="font-size:12px;color:#606266;margin-bottom:2px">
                {{ i + 1 }}. {{ step }}
              </div>
            </div>
          </div>
          <div style="font-size:12px;color:#909399;margin-top:4px">{{ t('cases.txtHint') }}</div>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" :loading="saving" @click="handleSave">{{ t('common.confirm') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Upload } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { getCases, createCase, updateCase, deleteCase, getMembers, parseCaseTxt } from '@/utils/api'

const { t } = useI18n()
const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const cases = ref([])
const members = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref(null)
const txtFileInput = ref(null)

const parsedPreview = reactive({ prepItems: [], steps: [] })

const defaultForm = () => ({
  title: '',
  parent_member_id: null,
  child_member_id: null,
  description: '',
  punishment_level: '',
  prep_items: '',
  parsed_steps: '',
  txt_filename: '',
  punishment_process: ''
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
  title: [{ required: true, message: t('cases.titleRequired'), trigger: 'blur' }],
  parent_member_id: [{ required: true, message: t('cases.parentRequired'), trigger: 'change' }],
  child_member_id: [
    { required: true, message: t('cases.childRequired'), trigger: 'change' },
    {
      validator: (rule, value, callback) => {
        if (value && form.parent_member_id) {
          if (value === form.parent_member_id) {
            callback(new Error(t('cases.sameParentChildError')))
            return
          }
        }
        callback()
      },
      trigger: 'change'
    }
  ],
  punishment_level: [{ required: true, message: t('cases.levelRequired'), trigger: 'change' }]
}))

function getMemberName(m) {
  if (!m) return '-'
  return m.name_cn || m.name_en || m.name || '-'
}

function statusLabel(s) {
  return {
    pending: t('cases.statusPending'),
    active: t('cases.statusActive'),
    completed: t('cases.statusCompleted')
  }[s] ?? s
}
function statusTagType(s) {
  return { pending: 'info', active: 'success', completed: '' }[s] ?? 'info'
}
function gradeTagType(g) {
  const map = { '满分': 'success', '优': 'success', '良': '', '达标': 'warning', '不达标': 'danger', '态度不端正': 'danger' }
  return map[g] ?? 'info'
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
    ElMessage.error(t('cases.loadFailed'))
  } finally {
    loading.value = false
  }
}

function openDialog(row = null) {
  Object.assign(form, defaultForm())
  parsedPreview.prepItems = []
  parsedPreview.steps = []
  if (row) {
    editingId.value = row.id
    form.title = row.title ?? ''
    form.parent_member_id = row.parent_member_id || null
    form.child_member_id = row.child_member_id || null
    form.description = row.description ?? ''
    form.punishment_level = row.punishment_level ?? ''
    form.prep_items = row.prep_items ?? ''
    form.parsed_steps = row.parsed_steps ?? ''
    form.txt_filename = row.txt_filename ?? ''
    form.punishment_process = row.punishment_process ?? ''
    if (form.prep_items) {
      try { parsedPreview.prepItems = JSON.parse(form.prep_items) } catch { /* ignore */ }
    }
    if (form.parsed_steps) {
      try { parsedPreview.steps = JSON.parse(form.parsed_steps) } catch { /* ignore */ }
    }
  } else {
    editingId.value = null
  }
  dialogVisible.value = true
}

function resetForm() {
  formRef.value?.resetFields()
  editingId.value = null
  parsedPreview.prepItems = []
  parsedPreview.steps = []
}

function triggerTxtImport() {
  txtFileInput.value?.click()
}

function handleTxtImport(event) {
  const file = event.target.files?.[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = async (e) => {
    const content = e.target?.result ?? ''
    form.txt_filename = file.name
    form.punishment_process = content
    try {
      const res = await parseCaseTxt({ content, level: form.punishment_level, txt_filename: file.name })
      const data = res.data
      const prepItems = data.prep_items ?? []
      const steps = data.steps ?? []
      parsedPreview.prepItems = prepItems
      parsedPreview.steps = steps
      form.prep_items = JSON.stringify(prepItems)
      form.parsed_steps = JSON.stringify(steps)
      ElMessage.success(t('cases.parseSuccess', { prep: prepItems.length, steps: steps.length }))
    } catch (err) {
      ElMessage.error(err.response?.data?.message || t('cases.parseFailed'))
    }
  }
  reader.onerror = () => {
    ElMessage.error(t('cases.parseFailed'))
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
    if (editingId.value) {
      await updateCase(editingId.value, { ...form })
      ElMessage.success(t('cases.updateSuccess'))
    } else {
      await createCase({ ...form })
      ElMessage.success(t('cases.createSuccess'))
    }
    dialogVisible.value = false
    await load()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || t('cases.saveFailed'))
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm(
    t('cases.deleteConfirm', { title: row.title }),
    t('cases.deleteTitle'),
    { type: 'warning', confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel') }
  )
  try {
    await deleteCase(row.id)
    ElMessage.success(t('cases.deleteSuccess'))
    await load()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || t('cases.deleteFailed'))
  }
}

function goDetail(row) {
  router.push(`/cases/${row.id}`)
}

onMounted(load)
</script>
