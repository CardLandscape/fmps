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
            {{ row.role === 'child' ? '小孩' : '家长' }}
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
          <el-button size="small" type="primary" plain @click="openEditWithAuth(row)">编辑</el-button>
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
          <el-input v-model="form.name_en" placeholder="请输入英文姓名（必填，例如：ZHANG XIAOMING）" @input="nameEnManuallyEdited = true" />
        </el-form-item>

        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" :disabled="!!editingId" placeholder="请选择角色" style="width: 100%">
            <el-option label="小孩" value="child" />
            <el-option label="家长" value="parent" />
          </el-select>
          <div v-if="editingId" style="font-size:12px;color:#909399;margin-top:2px">成员类型一经确定不可更改</div>
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
            placeholder="请选择国籍（必填）"
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
            placeholder="选择出生日期（必填）"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <!-- 主证件 -->
        <el-divider content-position="left">主证件信息</el-divider>

        <el-form-item label="主证件类型" prop="id_doc_type">
          <el-select
            v-model="form.id_doc_type"
            placeholder="请选择证件类型（必填）"
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
          <el-input v-model="form.id_doc_number" placeholder="请输入证件号码（必填）" />
        </el-form-item>

        <el-form-item label="签发日期" prop="id_issue_date">
          <el-date-picker
            v-model="form.id_issue_date"
            type="date"
            placeholder="选择签发日期（必填）"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="有效期" prop="id_expiry_date">
          <el-date-picker
            v-model="form.id_expiry_date"
            type="date"
            placeholder="选择有效期（必填）"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="签发机关" prop="id_issue_authority">
          <el-input v-model="form.id_issue_authority" placeholder="请输入签发机关（必填）" />
        </el-form-item>

        <!-- 辅助证件（根据主证件类型动态显示） -->
        <template v-if="showAuxDocs">
          <el-divider content-position="left">辅助证件信息</el-divider>

          <!-- 辅助证件1 -->
          <el-form-item :label="aux1Label" prop="aux1_doc_type">
            <el-select
              v-model="form.aux1_doc_type"
              :placeholder="aux1TypePlaceholder"
              :disabled="aux1TypeFixed"
              clearable
              style="width: 100%"
              @change="onAux1DocTypeChange"
            >
              <el-option
                v-for="t in aux1TypeOptions"
                :key="t.code"
                :label="t.name"
                :value="t.code"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="辅助证件1号码" prop="aux1_doc_number" :error="aux1DocNumberError">
            <el-input v-model="form.aux1_doc_number" :placeholder="aux1NumPlaceholder" />
          </el-form-item>

          <!-- 辅助证件2 (type 11 or 21) -->
          <template v-if="showAux2">
            <el-form-item label="辅助证件2类型" prop="aux2_doc_type">
              <el-select
                v-model="form.aux2_doc_type"
                placeholder="请选择辅助证件2类型（必填）"
                style="width: 100%"
                @change="onAux2DocTypeChange"
              >
                <el-option
                  v-for="t in aux2TypeOptions"
                  :key="t.code"
                  :label="t.name"
                  :value="t.code"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="辅助证件2号码" prop="aux2_doc_number" :error="aux2DocNumberError">
              <el-input v-model="form.aux2_doc_number" placeholder="请输入辅助证件2号码（必填）" />
            </el-form-item>
          </template>

          <!-- 主证件04：证明文件补充字段 -->
          <template v-if="form.id_doc_type === '04'">
            <el-form-item label="证明文件类型" prop="proof_doc_type">
              <el-select
                v-model="form.proof_doc_type"
                placeholder="请选择证明文件类型（必填）"
                style="width: 100%"
                @change="onProofDocTypeChange"
              >
                <el-option
                  v-for="t in PROOF_DOC_TYPES"
                  :key="t.code"
                  :label="t.name"
                  :value="t.code"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="签发国家" prop="proof_issue_country">
              <el-select
                v-model="form.proof_issue_country"
                placeholder="请选择签发国家（必填）"
                filterable
                style="width: 100%"
                @change="onProofIssueCountryChange"
              >
                <el-option
                  v-for="c in proofIssueCountryOptions"
                  :key="c.code"
                  :label="`${c.code} - ${c.name}`"
                  :value="c.code"
                />
              </el-select>
            </el-form-item>
          </template>
        </template>

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

  <!-- 授权密码对话框（编辑/删除前验证） -->
  <el-dialog v-model="authDialogVisible" title="需要授权密码" width="360px" :close-on-click-modal="false">
    <el-form @submit.prevent>
      <el-form-item label="授权密码" label-width="80px">
        <el-input
          v-model="authPassword"
          type="password"
          show-password
          placeholder="请输入授权密码"
          @keyup.enter="confirmAuthAction"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="authDialogVisible = false">取 消</el-button>
      <el-button type="primary" :loading="authLoading" @click="confirmAuthAction">确 定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getMembers, createMember, updateMember, deleteMemberWithAuth } from '@/utils/api'
import {
  COUNTRIES,
  ID_DOC_TYPES,
  AUX_DOC_TYPES,
  PROOF_DOC_TYPES,
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

// Auth password dialog state
const authDialogVisible = ref(false)
const authPassword = ref('')
const authLoading = ref(false)
const authAction = ref(null)   // { type: 'edit'|'delete', row: ... }

const idDocNumberError = ref('')
const aux1DocNumberError = ref('')
const aux2DocNumberError = ref('')

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
  aux1_doc_type: '',
  aux1_doc_number: '',
  aux2_doc_type: '',
  aux2_doc_number: '',
  proof_doc_type: '',
  proof_issue_country: '',
  school_name: '',
  grade: '',
  class_name: '',
  class_teacher_name: '',
  class_teacher_phone: '',
  outing_permission: '',
  outing_dates: '',
  outing_time_ranges: ''
})

// ─── Aux doc logic ────────────────────────────────────────────────────────────

// Whether to show aux doc section at all
const showAuxDocs = computed(() => {
  const t = form.id_doc_type
  return t !== '' && t !== '01' && t !== '91' && t !== '05'
})

// Whether to show aux doc 2 (only for type 11 and 21)
const showAux2 = computed(() => ['11', '21'].includes(form.id_doc_type))

// Aux1 type options based on primary doc type
const aux1TypeOptions = computed(() => {
  const nat = form.nationality
  switch (form.id_doc_type) {
    case '11': return AUX_DOC_TYPES.filter(t => t.code === '02')
    case '21': return AUX_DOC_TYPES.filter(t => t.code === '03')
    case '04': return AUX_DOC_TYPES.filter(t => t.code === '94')
    case '52': return AUX_DOC_TYPES.filter(t => ['95', '98'].includes(t.code))
    case '31': return AUX_DOC_TYPES.filter(t => t.code === '05')
    case '02':
      if (nat === 'HKG') return AUX_DOC_TYPES.filter(t => ['90', '92'].includes(t.code))
      if (nat === 'MAC') return AUX_DOC_TYPES.filter(t => ['96', '97'].includes(t.code))
      return AUX_DOC_TYPES.filter(t => ['90', '92', '96', '97'].includes(t.code))
    case '03': return AUX_DOC_TYPES.filter(t => t.code === '93')
    default: return AUX_DOC_TYPES
  }
})

// Aux1 type is fixed (only one option)
const aux1TypeFixed = computed(() => ['11', '21', '04'].includes(form.id_doc_type))

const aux1Label = computed(() => {
  if (form.id_doc_type === '04') return '辅助证件类型'
  return '辅助证件1类型'
})

const aux1NumPlaceholder = computed(() => {
  const req = ['11', '21', '04'].includes(form.id_doc_type)
  return req ? '请输入辅助证件号码（必填）' : '请输入辅助证件号码（可选）'
})

const aux1TypePlaceholder = computed(() => {
  const req = ['11', '21', '04'].includes(form.id_doc_type)
  return req ? '请选择辅助证件类型（必填）' : '请选择辅助证件类型（可选）'
})

// Aux2 type options based on primary doc type
const aux2TypeOptions = computed(() => {
  switch (form.id_doc_type) {
    case '11': return AUX_DOC_TYPES.filter(t => ['90', '92', '96', '97'].includes(t.code))
    case '21': return AUX_DOC_TYPES.filter(t => t.code === '93')
    default: return AUX_DOC_TYPES
  }
})

// Proof issue country options: exclude CHN/HKG/MAC/TWN unless proofDocType=94NP
const RESTRICTED = new Set(['CHN', 'HKG', 'MAC', 'TWN'])
const proofIssueCountryOptions = computed(() => {
  if (form.proof_doc_type === '94NP') {
    return COUNTRIES.filter(c => c.code === 'CHN')
  }
  return COUNTRIES.filter(c => !RESTRICTED.has(c.code))
})

// ─── Computed ─────────────────────────────────────────────────────────────────

const nationalityOptions = computed(() => {
  const fromMain = form.id_doc_type ? getAvailableCountries(form.id_doc_type) : null
  if (!fromMain) return COUNTRIES
  return fromMain
})

const formRules = computed(() => {
  const requireChinese = ['CHN', 'HKG', 'MAC', 'TWN'].includes(form.nationality)
  const mainType = form.id_doc_type
  const needAux1 = ['11', '21', '04'].includes(mainType)
  const needAux2 = ['11', '21'].includes(mainType)

  return {
    name_cn: requireChinese ? [{ required: true, message: '中文姓名为必填项（CHN/HKG/MAC/TWN国籍）', trigger: 'blur' }] : [],
    name_en: [{ required: true, message: '英文姓名为必填项', trigger: 'blur' }],
    role: [{ required: true, message: '请选择角色', trigger: 'change' }],
    gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
    nationality: [{ required: true, message: '国籍为必填项', trigger: 'change' }],
    birth_date: [{ required: true, message: '出生日期为必填项', trigger: 'change' }],
    id_doc_type: [{ required: true, message: '主证件类型为必填项', trigger: 'change' }],
    id_doc_number: [{ required: true, message: '主证件号码为必填项', trigger: 'blur' }],
    id_issue_date: [{ required: true, message: '签发日期为必填项', trigger: 'change' }],
    id_expiry_date: [{ required: true, message: '有效期为必填项', trigger: 'change' }],
    id_issue_authority: [{ required: true, message: '签发机关为必填项', trigger: 'blur' }],
    aux1_doc_type: needAux1 ? [{ required: true, message: '辅助证件1类型为必填项', trigger: 'change' }] : [],
    aux1_doc_number: needAux1 ? [{ required: true, message: '辅助证件1号码为必填项', trigger: 'blur' }] : [],
    aux2_doc_type: needAux2 ? [{ required: true, message: '辅助证件2类型为必填项', trigger: 'change' }] : [],
    aux2_doc_number: needAux2 ? [{ required: true, message: '辅助证件2号码为必填项', trigger: 'blur' }] : [],
    proof_doc_type: mainType === '04' ? [{ required: true, message: '证明文件类型为必填项', trigger: 'change' }] : [],
    proof_issue_country: mainType === '04' ? [{ required: true, message: '签发国家为必填项', trigger: 'change' }] : [],
    school_name: [{ validator: validateSchoolNameRule, trigger: 'blur' }]
  }
})

// ─── Watchers ─────────────────────────────────────────────────────────────────

watch(
  () => [form.id_doc_number, form.id_doc_type, form.nationality, form.gender, form.birth_date],
  ([number, docType, nationality, gender, birthDate]) => {
    idDocNumberError.value = validateIDNumber(docType, number, nationality, { gender, birthDate }) || ''
  }
)

watch(
  () => [form.aux1_doc_number, form.aux1_doc_type, form.nationality, form.birth_date, form.proof_doc_type, form.gender],
  ([number, docType, nationality, birthDate, proofDocType, gender]) => {
    aux1DocNumberError.value = validateIDNumber(docType, number, nationality, { birthDate, proofDocType, gender }) || ''
  }
)

watch(
  () => [form.aux2_doc_number, form.aux2_doc_type, form.nationality, form.birth_date, form.gender],
  ([number, docType, nationality, birthDate, gender]) => {
    aux2DocNumberError.value = validateIDNumber(docType, number, nationality, { birthDate, gender }) || ''
  }
)

// Track whether name_en was manually edited to avoid overwriting user input
const nameEnManuallyEdited = ref(false)

// Auto-fill name_en from name_cn (uppercase ASCII) when not manually edited.
// For purely Chinese names this produces an empty string; the user should manually
// enter the romanized form (e.g. ZHANG XIAOMING) following the displayed placeholder.
watch(() => form.name_cn, (newCn) => {
  if (!nameEnManuallyEdited.value) {
    const suggested = (newCn || '').replace(/[^\x00-\x7F]/g, '').trim().toUpperCase()
    form.name_en = suggested
  }
})
watch(() => form.id_doc_type, (newType) => {
  if (newType === '11') {
    form.aux1_doc_type = '02'
  } else if (newType === '21') {
    form.aux1_doc_type = '03'
  } else if (newType === '04') {
    form.aux1_doc_type = '94'
  } else {
    form.aux1_doc_type = ''
  }
  form.aux2_doc_type = ''
  form.aux1_doc_number = ''
  form.aux2_doc_number = ''
  form.proof_doc_type = ''
  form.proof_issue_country = ''
})

// ─── Helpers ──────────────────────────────────────────────────────────────────

function validateSchoolNameRule(rule, value, callback) {
  if (!value || !value.trim()) {
    callback()
    return
  }
  const keywords = ['小学', '中学', '大学', '学院']
  if (keywords.some(kw => value.includes(kw))) {
    callback()
  } else {
    callback(new Error('就读学校名称须包含"小学"、"中学"、"大学"或"学院"之一'))
  }
}

function getMemberDisplayName(member) {
  return member.name_cn || member.name_en || member.name || '-'
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
  if (form.id_doc_type) {
    const err = validateNationalityDocType(form.id_doc_type, val)
    if (err) ElMessage.warning(err)
  }
  // When primary=02, nationality change may narrow aux options
  if (form.id_doc_type === '02' && form.aux1_doc_type) {
    const hkgAux = new Set(['90', '92'])
    const macAux = new Set(['96', '97'])
    if (val === 'HKG' && !hkgAux.has(form.aux1_doc_type)) {
      form.aux1_doc_type = ''
      form.aux1_doc_number = ''
    } else if (val === 'MAC' && !macAux.has(form.aux1_doc_type)) {
      form.aux1_doc_type = ''
      form.aux1_doc_number = ''
    }
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

function onAux1DocTypeChange() {
  form.aux1_doc_number = ''
  aux1DocNumberError.value = ''
}

function onAux2DocTypeChange() {
  form.aux2_doc_number = ''
  aux2DocNumberError.value = ''
}

function onProofDocTypeChange(val) {
  // If 94NP, force proof_issue_country to CHN; otherwise clear if it's CHN
  if (val === '94NP') {
    form.proof_issue_country = 'CHN'
  } else if (form.proof_issue_country === 'CHN') {
    form.proof_issue_country = ''
  }
}

async function onProofIssueCountryChange(val) {
  const warningCountries = new Set(['VUT', 'GIN', 'GNB'])
  if (warningCountries.has(val)) {
    try {
      await ElMessageBox.confirm(
        '此人近2年内是否实际居住在该国家超过18个月？如否，建议加强核实。',
        '⚠️ 高风险国家确认',
        {
          confirmButtonText: '确认继续',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
      // User confirmed, continue
    } catch {
      // User cancelled, clear the selection
      form.proof_issue_country = ''
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

// Open auth dialog before editing
function openEditWithAuth(row) {
  authAction.value = { type: 'edit', row }
  authPassword.value = ''
  authDialogVisible.value = true
}

async function confirmAuthAction() {
  if (!authAction.value) return
  if (!authPassword.value) {
    ElMessage.warning('请输入授权密码')
    return
  }
  authLoading.value = true
  try {
    const { type, row } = authAction.value
    if (type === 'edit') {
      // Store the auth password; it will be included when the user saves
      pendingAuthPassword.value = authPassword.value
      authDialogVisible.value = false
      if (row) {
        openDialog(row)
      }
      // If row is null this was called from the create-parent flow; the save will retry
    } else if (type === 'delete') {
      await deleteMemberWithAuth(row.id, authPassword.value)
      ElMessage.success('已删除')
      authDialogVisible.value = false
      await fetchMembers()
    }
  } catch (e) {
    const msg = e.response?.data?.message || '操作失败'
    ElMessage.error(msg)
  } finally {
    authLoading.value = false
  }
}

// Stores the auth password entered before opening the edit dialog
const pendingAuthPassword = ref('')

function openDialog(row = null) {
  if (row) {
    editingId.value = row.id
    nameEnManuallyEdited.value = true // existing record: don't overwrite name_en
    Object.keys(form).forEach(k => {
      form[k] = row[k] ?? ''
    })
    // Backward compatibility: use row.name if name_cn is absent
    if (!form.name_cn && row.name) {
      form.name_cn = row.name
    }
    // Normalize legacy "adult" role
    if (form.role === 'adult') {
      form.role = 'parent'
    }
    // Backward compatibility: migrate old aux_doc fields to aux1
    if (!form.aux1_doc_type && row.aux_doc_type) {
      form.aux1_doc_type = row.aux_doc_type
    }
    if (!form.aux1_doc_number && row.aux_doc_number) {
      form.aux1_doc_number = row.aux_doc_number
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
    pendingAuthPassword.value = ''
    resetForm()
  }
  dialogVisible.value = true
}

function resetForm() {
  Object.keys(form).forEach(k => { form[k] = '' })
  outingDatesArray.value = []
  outingTimeRanges.value = []
  idDocNumberError.value = ''
  aux1DocNumberError.value = ''
  aux2DocNumberError.value = ''
  nameEnManuallyEdited.value = false
  pendingAuthPassword.value = ''
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
  if (aux1DocNumberError.value) {
    ElMessage.error('辅助证件1号码格式有误，请修正后保存')
    return
  }
  if (aux2DocNumberError.value) {
    ElMessage.error('辅助证件2号码格式有误，请修正后保存')
    return
  }

  saving.value = true
  try {
    const payload = { ...form, name: form.name_cn || form.name_en }
    if (editingId.value) {
      // Include auth password for server-side validation
      payload.auth_password = pendingAuthPassword.value
      await updateMember(editingId.value, payload)
      ElMessage.success('成员已更新')
    } else {
      // For new parent records, include auth_password if provided (same-person child exists)
      payload.auth_password = pendingAuthPassword.value
      await createMember(payload)
      ElMessage.success('成员已添加')
    }
    dialogVisible.value = false
    pendingAuthPassword.value = ''
    await fetchMembers()
  } catch (e) {
    const msg = e.response?.data?.message || '保存失败，请稍后重试'
    ElMessage.error(msg)
    // If server requires auth password for new parent (same-person child exists),
    // show the auth dialog so the user can supply it and retry
    if (e.response?.status === 403 && !editingId.value) {
      authAction.value = { type: 'edit', row: null }
      authPassword.value = ''
      authDialogVisible.value = true
    }
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
    // All deletions require auth password
    authAction.value = { type: 'delete', row }
    authPassword.value = ''
    authDialogVisible.value = true
  } catch {
    // User cancelled the confirmation – do nothing
  }
}

// ─── Init ─────────────────────────────────────────────────────────────────────

onMounted(fetchMembers)
</script>
