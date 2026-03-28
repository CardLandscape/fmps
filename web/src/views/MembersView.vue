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
      <el-table-column label="姓名">
        <template #default="{ row }">{{ getMemberDisplayName(row) }}</template>
      </el-table-column>
      <el-table-column prop="role" label="角色" width="80">
        <template #default="{ row }">
          <el-tag :type="row.role === 'child' ? 'warning' : 'success'">
            {{ row.role === 'child' ? '孩子' : '成人' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="gender" label="性别" width="70" />
      <el-table-column prop="nationality" label="国籍" width="80" />
      <el-table-column prop="school_name" label="学校" />
      <el-table-column prop="outing_permission" label="外出权限" width="90">
        <template #default="{ row }">
          <el-tag v-if="row.outing_permission" :type="permissionTagType(row.outing_permission)" size="small">
            {{ row.outing_permission }}
          </el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="160">
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
    width="760px"
    @closed="resetForm"
  >
    <div style="max-height: 70vh; overflow-y: auto; padding-right: 8px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="120px">

        <!-- 基本信息 -->
        <el-divider content-position="left">基本信息</el-divider>

        <el-form-item label="中文姓名" prop="name_cn">
          <el-input v-model="form.name_cn" placeholder="请输入中文姓名" />
        </el-form-item>

        <el-form-item label="英文姓名" prop="name_en">
          <el-input v-model="form.name_en" placeholder="请输入英文姓名（可选）" />
        </el-form-item>

        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" placeholder="请选择角色" style="width: 100%">
            <el-option label="孩子" value="child" />
            <el-option label="成人" value="adult" />
          </el-select>
        </el-form-item>

        <el-form-item label="性别" prop="gender">
          <el-radio-group v-model="form.gender">
            <el-radio value="男">男</el-radio>
            <el-radio value="女">女</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="国籍" prop="nationality">
          <el-select
            v-model="form.nationality"
            placeholder="请选择国籍"
            filterable
            style="width: 100%"
            @change="onNationalityChange"
          >
            <el-option
              v-for="c in nationalityOptions"
              :key="c.code"
              :label="`${c.code} - ${c.name}`"
              :value="c.code"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="出生日期" prop="birth_date">
          <el-date-picker
            v-model="form.birth_date"
            type="date"
            placeholder="选择出生日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <!-- 主证件 -->
        <el-divider content-position="left">主证件信息</el-divider>

        <el-form-item label="主证件类型" prop="id_doc_type">
          <el-select
            v-model="form.id_doc_type"
            placeholder="请选择证件类型"
            style="width: 100%"
            @change="onIdDocTypeChange"
          >
            <el-option
              v-for="t in ID_DOC_TYPES"
              :key="t.code"
              :label="t.name"
              :value="t.code"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="主证件号码" prop="id_doc_number" :error="idDocNumberError">
          <el-input v-model="form.id_doc_number" placeholder="请输入证件号码" />
        </el-form-item>

        <el-form-item label="签发日期" prop="id_issue_date">
          <el-date-picker
            v-model="form.id_issue_date"
            type="date"
            placeholder="选择签发日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="有效期" prop="id_expiry_date">
          <el-date-picker
            v-model="form.id_expiry_date"
            type="date"
            placeholder="选择有效期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="签发机关" prop="id_issue_authority">
          <el-input v-model="form.id_issue_authority" placeholder="请输入签发机关" />
        </el-form-item>

        <!-- 辅助证件 -->
        <el-divider content-position="left">辅助证件信息</el-divider>

        <el-form-item label="辅助证件类型" prop="aux_doc_type">
          <el-select
            v-model="form.aux_doc_type"
            placeholder="请选择辅助证件类型（可选）"
            clearable
            style="width: 100%"
            @change="onAuxDocTypeChange"
          >
            <el-option
              v-for="t in ID_DOC_TYPES"
              :key="t.code"
              :label="t.name"
              :value="t.code"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="辅助证件号码" prop="aux_doc_number" :error="auxDocNumberError">
          <el-input v-model="form.aux_doc_number" placeholder="请输入辅助证件号码（可选）" />
        </el-form-item>

        <!-- 学校信息 -->
        <el-divider content-position="left">学校信息</el-divider>

        <el-form-item label="就读学校" prop="school_name">
          <el-input v-model="form.school_name" placeholder="请输入就读学校" />
        </el-form-item>

        <el-form-item label="年级" prop="grade">
          <el-select v-model="form.grade" placeholder="请选择年级" clearable style="width: 100%">
            <el-option v-for="g in GRADES" :key="g" :label="g" :value="g" />
          </el-select>
        </el-form-item>

        <el-form-item label="班级" prop="class_name">
          <el-select v-model="form.class_name" placeholder="请选择班级" clearable style="width: 100%">
            <el-option v-for="c in CLASSES" :key="c" :label="`${c}班`" :value="c" />
          </el-select>
        </el-form-item>

        <el-form-item label="班主任姓名" prop="class_teacher_name">
          <el-input v-model="form.class_teacher_name" placeholder="请输入班主任姓名" />
        </el-form-item>

        <el-form-item label="班主任电话" prop="class_teacher_phone">
          <el-input v-model="form.class_teacher_phone" placeholder="请输入班主任电话" />
        </el-form-item>

        <!-- 外出权限 -->
        <el-divider content-position="left">外出权限</el-divider>

        <el-form-item label="外出权限" prop="outing_permission">
          <el-radio-group v-model="form.outing_permission">
            <el-radio value="许可">许可</el-radio>
            <el-radio value="不许可">不许可</el-radio>
            <el-radio value="受限">受限</el-radio>
          </el-radio-group>
        </el-form-item>

        <template v-if="form.outing_permission === '受限'">
          <el-form-item label="允许外出日期" prop="outing_dates">
            <el-date-picker
              v-model="outingDatesArray"
              type="dates"
              placeholder="选择允许外出的日期"
              value-format="YYYY-MM-DD"
              style="width: 100%"
              @change="onOutingDatesChange"
            />
          </el-form-item>

          <el-form-item label="允许外出时段" prop="outing_time_ranges">
            <div style="width: 100%">
              <div
                v-for="(range, idx) in outingTimeRanges"
                :key="idx"
                style="display: flex; align-items: center; gap: 8px; margin-bottom: 8px"
              >
                <el-time-picker
                  v-model="range.start"
                  placeholder="开始时间"
                  value-format="HH:mm"
                  format="HH:mm"
                  style="width: 140px"
                  @change="syncTimeRangesToForm"
                />
                <span>至</span>
                <el-time-picker
                  v-model="range.end"
                  placeholder="结束时间"
                  value-format="HH:mm"
                  format="HH:mm"
                  style="width: 140px"
                  @change="syncTimeRangesToForm"
                />
                <el-button type="danger" size="small" plain @click="removeTimeRange(idx)">删除</el-button>
              </div>
              <el-button size="small" @click="addTimeRange">+ 添加时间段</el-button>
            </div>
          </el-form-item>
        </template>

      </el-form>
    </div>

    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getMembers, createMember, updateMember, deleteMember } from '@/utils/api'
import {
  COUNTRIES,
  ID_DOC_TYPES,
  GRADES,
  CLASSES,
  getAvailableCountries,
  validateIDNumber,
  validateNationalityDocType
} from '@/utils/constants'

// ─── State ───────────────────────────────────────────────────────────────────

const members = ref([])
const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const idDocNumberError = ref('')
const auxDocNumberError = ref('')

const outingDatesArray = ref([])
const outingTimeRanges = ref([])

const form = reactive({
  name_cn: '',
  name_en: '',
  role: '',
  gender: '',
  nationality: '',
  birth_date: '',
  id_doc_type: '',
  id_doc_number: '',
  id_issue_date: '',
  id_expiry_date: '',
  id_issue_authority: '',
  aux_doc_type: '',
  aux_doc_number: '',
  school_name: '',
  grade: '',
  class_name: '',
  class_teacher_name: '',
  class_teacher_phone: '',
  outing_permission: '',
  outing_dates: '',
  outing_time_ranges: ''
})

// ─── Computed ─────────────────────────────────────────────────────────────────

const nationalityOptions = computed(() => {
  const fromMain = form.id_doc_type ? getAvailableCountries(form.id_doc_type) : null
  const fromAux = form.aux_doc_type ? getAvailableCountries(form.aux_doc_type) : null

  if (!fromMain && !fromAux) return COUNTRIES
  if (fromMain && !fromAux) return fromMain
  if (!fromMain && fromAux) return fromAux

  // Both set: intersection (codes present in both lists)
  const auxCodes = new Set(fromAux.map(c => c.code))
  return fromMain.filter(c => auxCodes.has(c.code))
})

const formRules = computed(() => ({
  name_cn: form.nationality === 'CHN' ? [{ required: true, message: '中文姓名为必填项', trigger: 'blur' }] : [],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }]
}))

// ─── Watchers ─────────────────────────────────────────────────────────────────

watch(
  () => [form.id_doc_number, form.id_doc_type, form.nationality],
  ([number, docType, nationality]) => {
    idDocNumberError.value = validateIDNumber(docType, number, nationality) || ''
  }
)

watch(
  () => [form.aux_doc_number, form.aux_doc_type, form.nationality],
  ([number, docType, nationality]) => {
    auxDocNumberError.value = validateIDNumber(docType, number, nationality) || ''
  }
)

// ─── Helpers ──────────────────────────────────────────────────────────────────

function getMemberDisplayName(member) {
  return member.name_cn || member.name || '-'
}

function formatDate(val) {
  if (!val) return '-'
  return new Date(val).toLocaleString('zh-CN', { hour12: false })
}

function permissionTagType(p) {
  if (p === '许可') return 'success'
  if (p === '不许可') return 'danger'
  if (p === '受限') return 'warning'
  return 'info'
}

// ─── Nationality / doc type linkage ──────────────────────────────────────────

function onNationalityChange(val) {
  // Validate compatibility with main doc type
  if (form.id_doc_type) {
    const err = validateNationalityDocType(form.id_doc_type, val)
    if (err) ElMessage.warning(err)
  }
  // Validate compatibility with aux doc type
  if (form.aux_doc_type) {
    const err = validateNationalityDocType(form.aux_doc_type, val)
    if (err) ElMessage.warning(err)
  }
}

function onIdDocTypeChange(docType) {
  if (form.nationality) {
    const available = getAvailableCountries(docType)
    const stillValid = available.some(c => c.code === form.nationality)
    if (!stillValid) {
      form.nationality = ''
      ElMessage.warning('当前国籍与所选证件类型不兼容，已清空国籍')
    }
  }
}

function onAuxDocTypeChange(docType) {
  if (!docType) return
  if (form.nationality) {
    const available = getAvailableCountries(docType)
    const stillValid = available.some(c => c.code === form.nationality)
    if (!stillValid) {
      form.nationality = ''
      ElMessage.warning('当前国籍与所选辅助证件类型不兼容，已清空国籍')
    }
  }
}

// ─── Outing dates / time ranges ───────────────────────────────────────────────

function onOutingDatesChange(dates) {
  form.outing_dates = dates && dates.length ? JSON.stringify(dates) : ''
}

function addTimeRange() {
  outingTimeRanges.value.push({ start: '', end: '' })
}

function removeTimeRange(idx) {
  outingTimeRanges.value.splice(idx, 1)
  syncTimeRangesToForm()
}

function syncTimeRangesToForm() {
  form.outing_time_ranges = outingTimeRanges.value.length
    ? JSON.stringify(outingTimeRanges.value)
    : ''
}

// ─── CRUD ─────────────────────────────────────────────────────────────────────

async function fetchMembers() {
  loading.value = true
  try {
    const res = await getMembers()
    members.value = res.data || []
  } catch {
    ElMessage.error('获取成员列表失败')
  } finally {
    loading.value = false
  }
}

function openDialog(row = null) {
  if (row) {
    editingId.value = row.id
    Object.keys(form).forEach(k => {
      form[k] = row[k] ?? ''
    })
    // Backward compatibility: use row.name if name_cn is absent
    if (!form.name_cn && row.name) {
      form.name_cn = row.name
    }
    // Parse outing_dates
    try {
      outingDatesArray.value = form.outing_dates ? JSON.parse(form.outing_dates) : []
    } catch {
      outingDatesArray.value = []
    }
    // Parse outing_time_ranges
    try {
      outingTimeRanges.value = form.outing_time_ranges ? JSON.parse(form.outing_time_ranges) : []
    } catch {
      outingTimeRanges.value = []
    }
  } else {
    editingId.value = null
    resetForm()
  }
  dialogVisible.value = true
}

function resetForm() {
  Object.keys(form).forEach(k => { form[k] = '' })
  outingDatesArray.value = []
  outingTimeRanges.value = []
  idDocNumberError.value = ''
  auxDocNumberError.value = ''
  formRef.value?.resetFields()
}

async function handleSave() {
  try {
    await formRef.value.validate()
  } catch {
    ElMessage.error('请检查表单填写是否正确')
    return
  }

  // Block save if there are ID number validation errors
  if (idDocNumberError.value) {
    ElMessage.error('主证件号码格式有误，请修正后保存')
    return
  }
  if (auxDocNumberError.value) {
    ElMessage.error('辅助证件号码格式有误，请修正后保存')
    return
  }

  saving.value = true
  try {
    const payload = { ...form, name: form.name_cn }
    if (editingId.value) {
      await updateMember(editingId.value, payload)
      ElMessage.success('成员已更新')
    } else {
      await createMember(payload)
      ElMessage.success('成员已添加')
    }
    dialogVisible.value = false
    await fetchMembers()
  } catch {
    ElMessage.error('保存失败，请稍后重试')
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(
      `确定要删除成员「${getMemberDisplayName(row)}」吗？`,
      '删除确认',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
    await deleteMember(row.id)
    ElMessage.success('已删除')
    await fetchMembers()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败')
  }
}

// ─── Init ─────────────────────────────────────────────────────────────────────

onMounted(fetchMembers)
</script>
