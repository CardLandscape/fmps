<template>
  <el-card shadow="never">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <span style="font-weight: 600">{{ i18n.t('casesPageTitle') }}</span>
        <el-button type="primary" @click="openDialog()">
          <el-icon><Plus /></el-icon>
          {{ i18n.t('btnNewCase') }}
        </el-button>
      </div>
    </template>

    <el-table :data="cases" v-loading="loading" stripe style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="title" :label="i18n.t('colTitle')" />
      <el-table-column :label="i18n.t('colParent')" width="120">
        <template #default="{ row }">
          {{ getMemberName(row.parent_member) }}
        </template>
      </el-table-column>
      <el-table-column :label="i18n.t('colChild')" width="120">
        <template #default="{ row }">
          {{ getMemberName(row.child_member) }}
        </template>
      </el-table-column>
      <el-table-column :label="i18n.t('colLevel')" width="70">
        <template #default="{ row }">
          <el-tag v-if="row.punishment_level" size="small" type="danger">{{ row.punishment_level }}</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" :label="i18n.t('colStatus')" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="i18n.t('colGrade')" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.final_grade" :type="gradeTagType(row.final_grade)" size="small">{{ row.final_grade }}</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" :label="i18n.t('colCreatedAt')" width="160">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column :label="i18n.t('colActions')" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" plain @click="goDetail(row)">{{ i18n.t('btnDetail') }}</el-button>
          <el-button size="small" type="warning" plain @click="openDialog(row)">{{ i18n.t('btnEdit') }}</el-button>
          <el-button size="small" type="danger" plain @click="handleDelete(row)">{{ i18n.t('btnDelete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <!-- 创建/编辑对话框 -->
  <el-dialog
    v-model="dialogVisible"
    :title="editingId ? i18n.t('dialogEditCase') : i18n.t('dialogNewCase')"
    width="640px"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
      <el-form-item :label="i18n.t('labelCaseTitle')" prop="title">
        <el-input v-model="form.title" :placeholder="i18n.t('placeholderCaseTitle')" />
      </el-form-item>
      <el-form-item :label="i18n.t('labelCaseParent')" prop="parent_member_id">
        <el-select v-model="form.parent_member_id" :placeholder="i18n.t('placeholderCaseParent')" style="width:100%">
          <el-option
            v-for="m in parentMembers"
            :key="m.id"
            :label="getMemberName(m)"
            :value="m.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="i18n.t('labelCaseChild')" prop="child_member_id">
        <el-select v-model="form.child_member_id" :placeholder="i18n.t('placeholderCaseChild')" style="width:100%">
          <el-option
            v-for="m in childMembers"
            :key="m.id"
            :label="getMemberName(m)"
            :value="m.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="i18n.t('labelCaseDescription')">
        <el-input v-model="form.description" type="textarea" :rows="2" :placeholder="i18n.t('placeholderCaseDescription')" />
      </el-form-item>

      <!-- Punishment workflow section -->
      <el-divider content-position="left">{{ i18n.t('sectionPunishment') }}</el-divider>

      <!-- TXT import first to detect levels -->
      <el-form-item :label="i18n.t('labelImportTxt')">
        <div style="width:100%">
          <div style="display:flex;gap:8px;align-items:center;margin-bottom:6px">
            <el-button size="small" @click="triggerTxtImport">
              <el-icon><Upload /></el-icon>
              {{ i18n.t('btnImportTxt') }}
            </el-button>
            <span v-if="form.txt_filename" style="font-size:12px;color:#606266">{{ form.txt_filename }}</span>
            <input
              ref="txtFileInput"
              type="file"
              accept=".txt"
              style="display:none"
              @change="handleTxtImport"
            />
          </div>
          <div style="font-size:12px;color:#909399">{{ i18n.t('txtImportHint') }}</div>
        </div>
      </el-form-item>

      <!-- Dynamic level selector -->
      <el-form-item :label="i18n.t('labelPunishLevel')" prop="punishment_level">
        <div style="width:100%">
          <!-- Detected levels as radio buttons -->
          <div v-if="detectedLevels.length > 0" style="margin-bottom:8px">
            <div style="font-size:12px;color:#67c23a;margin-bottom:6px">
              {{ i18n.t('detectedLevels').replace('{n}', detectedLevels.length) }}
            </div>
            <div style="display:flex;flex-wrap:wrap;gap:6px;margin-bottom:8px">
              <el-tag
                v-for="(lv, idx) in detectedLevels"
                :key="idx"
                :type="form.punishment_level === lv ? 'primary' : 'info'"
                style="cursor:pointer;user-select:none"
                @click="selectLevel(lv)"
              >
                <span>{{ lv }}级</span>
                <el-icon
                  style="margin-left:4px;vertical-align:middle"
                  @click.stop="removeDetectedLevel(idx)"
                ><Close /></el-icon>
              </el-tag>
              <el-tag
                type="success"
                style="cursor:pointer"
                @click="showAddLevelInput = true"
                v-if="!showAddLevelInput"
              >+ {{ i18n.t('btnAddLevel') }}</el-tag>
            </div>
            <!-- Inline add level -->
            <div v-if="showAddLevelInput" style="display:flex;gap:6px;margin-bottom:6px">
              <el-input
                v-model="newLevelName"
                size="small"
                :placeholder="i18n.t('placeholderCustomLevel')"
                style="width:160px"
                @keyup.enter="addCustomLevel"
              />
              <el-button size="small" type="primary" @click="addCustomLevel">{{ i18n.t('btnConfirm') }}</el-button>
              <el-button size="small" @click="showAddLevelInput = false; newLevelName = ''">{{ i18n.t('btnCancel') }}</el-button>
            </div>
          </div>
          <!-- No levels detected yet: show manual input -->
          <div v-else style="margin-bottom:8px">
            <div style="font-size:12px;color:#909399;margin-bottom:6px">
              {{ form.txt_filename ? i18n.t('noLevelsDetected') : i18n.t('importTxtFirst') }}
            </div>
          </div>
          <!-- Always show manual level input -->
          <div style="display:flex;align-items:center;gap:8px">
            <el-input
              v-model="form.punishment_level"
              size="small"
              :placeholder="i18n.t('levelSelector')"
              style="width:200px"
              @change="onLevelChange"
            />
            <span style="font-size:13px;color:#606266">级</span>
            <span v-if="form.punishment_level" style="font-size:12px;color:#409eff">
              → {{ form.punishment_level }}级
            </span>
          </div>
          <div v-if="!form.punishment_level" style="font-size:12px;color:#f56c6c;margin-top:4px">
            {{ i18n.t('validCaseLevelRequired') }}
          </div>
        </div>
      </el-form-item>

      <!-- Parsed preview with editable items -->
      <el-form-item v-if="parsedPreview.prepItems.length || parsedPreview.steps.length" :label="i18n.t('prepItemsLabel')">
        <div style="width:100%">
          <div v-if="parsedPreview.prepItems.length" style="margin-bottom:12px">
            <div style="font-weight:600;font-size:13px;margin-bottom:6px">{{ i18n.t('prepItemsLabel') }}（{{ parsedPreview.prepItems.length }}）：</div>
            <div style="display:flex;flex-wrap:wrap;gap:6px;margin-bottom:6px">
              <el-tag
                v-for="(item, i) in parsedPreview.prepItems"
                :key="i"
                closable
                size="small"
                @close="removePrepItem(i)"
              >{{ item }}</el-tag>
            </div>
            <div style="display:flex;gap:6px">
              <el-input
                v-model="newPrepItem"
                size="small"
                :placeholder="'+ 添加准备物品'"
                style="width:180px"
                @keyup.enter="addPrepItem"
              />
              <el-button size="small" @click="addPrepItem">{{ i18n.t('btnConfirm') }}</el-button>
            </div>
          </div>
          <div v-if="parsedPreview.steps.length">
            <div style="font-weight:600;font-size:13px;margin-bottom:4px">{{ i18n.t('execStepsLabel') }}（{{ parsedPreview.steps.length }}）：</div>
            <div v-for="(step, i) in parsedPreview.steps" :key="i" style="font-size:12px;color:#606266;margin-bottom:2px;display:flex;align-items:center;gap:6px">
              <span style="color:#409eff;font-weight:600;min-width:20px">{{ i + 1 }}.</span>
              <span>{{ step }}</span>
            </div>
          </div>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">{{ i18n.t('btnCancel') }}</el-button>
      <el-button type="primary" :loading="saving" @click="handleSave">{{ i18n.t('btnConfirm') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from '@/utils/i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Upload, Close } from '@element-plus/icons-vue'
import { getCases, createCase, updateCase, deleteCase, getMembers, parseCaseTxt, parseCaseTxtLevels } from '@/utils/api'

const router = useRouter()
const i18n = useI18n()
const loading = ref(false)
const saving = ref(false)
const cases = ref([])
const members = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref(null)
const txtFileInput = ref(null)

const parsedPreview = reactive({ prepItems: [], steps: [], postures: [] })
const detectedLevels = ref([])
const showAddLevelInput = ref(false)
const newLevelName = ref('')
const newPrepItem = ref('')
let rawTxtContent = ''

const defaultForm = () => ({
  title: '',
  parent_member_id: null,
  child_member_id: null,
  description: '',
  punishment_level: '',
  prep_items: '',
  parsed_steps: '',
  postures: '',
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
  title: [{ required: true, message: i18n.t('validCaseTitleRequired'), trigger: 'blur' }],
  parent_member_id: [{ required: true, message: i18n.t('validCaseParentRequired'), trigger: 'change' }],
  child_member_id: [
    { required: true, message: i18n.t('validCaseChildRequired'), trigger: 'change' },
    {
      validator: (rule, value, callback) => {
        if (value && form.parent_member_id) {
          if (value === form.parent_member_id) {
            callback(new Error(i18n.t('validCaseSamePersonError')))
            return
          }
        }
        callback()
      },
      trigger: 'change'
    }
  ],
  punishment_level: [{ required: true, message: i18n.t('validCaseLevelRequired'), trigger: 'change' }]
}))

function getMemberName(m) {
  if (!m) return '-'
  return m.name_cn || m.name_en || m.name || '-'
}

function statusLabel(s) {
  return { pending: i18n.t('statusPending'), active: i18n.t('statusActive'), completed: i18n.t('statusCompleted') }[s] ?? s
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
    ElMessage.error(i18n.t('loadCasesFailed'))
  } finally {
    loading.value = false
  }
}

function openDialog(row = null) {
  Object.assign(form, defaultForm())
  parsedPreview.prepItems = []
  parsedPreview.steps = []
  parsedPreview.postures = []
  detectedLevels.value = []
  showAddLevelInput.value = false
  newLevelName.value = ''
  newPrepItem.value = ''
  rawTxtContent = ''
  if (row) {
    editingId.value = row.id
    form.title = row.title ?? ''
    form.parent_member_id = row.parent_member_id || null
    form.child_member_id = row.child_member_id || null
    form.description = row.description ?? ''
    form.punishment_level = row.punishment_level ?? ''
    form.prep_items = row.prep_items ?? ''
    form.parsed_steps = row.parsed_steps ?? ''
    form.postures = row.postures ?? ''
    form.txt_filename = row.txt_filename ?? ''
    form.punishment_process = row.punishment_process ?? ''
    // Restore preview from saved data
    if (form.prep_items) {
      try { parsedPreview.prepItems = JSON.parse(form.prep_items) } catch { /* ignore */ }
    }
    if (form.parsed_steps) {
      try { parsedPreview.steps = JSON.parse(form.parsed_steps) } catch { /* ignore */ }
    }
    if (form.postures) {
      try { parsedPreview.postures = JSON.parse(form.postures) } catch { /* ignore */ }
    }
    // If there's a saved level, show it as a detected level option
    if (form.punishment_level) {
      detectedLevels.value = [form.punishment_level]
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
  parsedPreview.postures = []
  detectedLevels.value = []
  showAddLevelInput.value = false
  newLevelName.value = ''
  newPrepItem.value = ''
  rawTxtContent = ''
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
    rawTxtContent = content
    form.txt_filename = file.name
    // Store raw content for legacy compatibility
    form.punishment_process = content
    // First detect levels from TXT
    try {
      const lvlRes = await parseCaseTxtLevels({ content })
      const levels = lvlRes.data?.levels ?? []
      detectedLevels.value = levels
      if (levels.length > 0 && !form.punishment_level) {
        // Auto-select first level
        form.punishment_level = levels[0]
        await parseForLevel(content, levels[0])
      } else if (form.punishment_level) {
        await parseForLevel(content, form.punishment_level)
      } else if (levels.length === 0) {
        ElMessage.warning(i18n.t('noLevelsDetected'))
      }
    } catch (err) {
      ElMessage.error(err.response?.data?.message || i18n.t('loadFailed'))
    }
  }
  reader.onerror = () => {
    ElMessage.error(i18n.t('loadFailed'))
  }
  reader.readAsText(file, 'utf-8')
  event.target.value = ''
}

async function parseForLevel(content, level) {
  try {
    const res = await parseCaseTxt({ content, level, txt_filename: form.txt_filename })
    const data = res.data
    const prepItems = data.prep_items ?? []
    const steps = data.steps ?? []
    const postures = data.postures ?? []
    parsedPreview.prepItems = prepItems
    parsedPreview.steps = steps
    parsedPreview.postures = postures
    form.prep_items = JSON.stringify(prepItems)
    form.parsed_steps = JSON.stringify(steps)
    form.postures = JSON.stringify(postures)
    ElMessage.success(`${i18n.t('prepItemsLabel')}: ${prepItems.length}, ${i18n.t('execStepsLabel')}: ${steps.length}`)
  } catch (err) {
    ElMessage.error(err.response?.data?.message || i18n.t('loadFailed'))
  }
}

async function selectLevel(lv) {
  form.punishment_level = lv
  if (rawTxtContent) {
    await parseForLevel(rawTxtContent, lv)
  }
}

async function onLevelChange(val) {
  if (rawTxtContent && val) {
    await parseForLevel(rawTxtContent, val)
  }
}

function removeDetectedLevel(idx) {
  const removed = detectedLevels.value[idx]
  detectedLevels.value.splice(idx, 1)
  if (form.punishment_level === removed) {
    form.punishment_level = detectedLevels.value[0] ?? ''
    if (form.punishment_level && rawTxtContent) {
      parseForLevel(rawTxtContent, form.punishment_level)
    }
  }
}

function addCustomLevel() {
  const name = newLevelName.value.trim()
  if (!name) return
  if (!detectedLevels.value.includes(name)) {
    detectedLevels.value.push(name)
  }
  form.punishment_level = name
  showAddLevelInput.value = false
  newLevelName.value = ''
  if (rawTxtContent) {
    parseForLevel(rawTxtContent, name)
  }
}

function removePrepItem(idx) {
  parsedPreview.prepItems.splice(idx, 1)
  form.prep_items = JSON.stringify(parsedPreview.prepItems)
}

function addPrepItem() {
  const item = newPrepItem.value.trim()
  if (!item) return
  parsedPreview.prepItems.push(item)
  form.prep_items = JSON.stringify(parsedPreview.prepItems)
  newPrepItem.value = ''
}

async function handleSave() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (editingId.value) {
      await updateCase(editingId.value, { ...form })
      ElMessage.success(i18n.t('updateSuccess'))
    } else {
      await createCase({ ...form })
      ElMessage.success(i18n.t('addSuccess'))
    }
    dialogVisible.value = false
    await load()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || i18n.t('saveFailed'))
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm(
    i18n.t('confirmDeleteCase'),
    i18n.t('confirmDeleteTitle'),
    { type: 'warning', confirmButtonText: i18n.t('confirmDeleteBtn'), cancelButtonText: i18n.t('btnCancel') }
  )
  try {
    await deleteCase(row.id)
    ElMessage.success(i18n.t('deleteSuccess'))
    await load()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || i18n.t('deleteFailed'))
  }
}

function goDetail(row) {
  router.push(`/cases/${row.id}`)
}

onMounted(load)
</script>

