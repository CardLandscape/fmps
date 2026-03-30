<template>
  <el-card shadow="never">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <span style="font-weight: 600">{{ t('members.title') }}</span>
        <el-button type="primary" @click="openDialog()">
          <el-icon><Plus /></el-icon>
          {{ t('members.addMember') }}
        </el-button>
      </div>
    </template>

    <el-table :data="members" v-loading="loading" stripe style="width: 100%">
      <el-table-column :label="t('members.colName')">
        <template #default="{ row }">{{ getMemberDisplayName(row) }}</template>
      </el-table-column>
      <el-table-column prop="role" :label="t('members.colRole')" width="80">
        <template #default="{ row }">
          <el-tag :type="row.role === 'child' ? 'warning' : 'success'">
            {{ row.role === 'child' ? t('members.roleChild') : t('members.roleParent') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="gender" :label="t('members.colGender')" width="70" />
      <el-table-column prop="nationality" :label="t('members.colNationality')" width="80" />
      <el-table-column prop="school_name" :label="t('members.colSchool')" />
      <el-table-column prop="outing_permission" :label="t('members.colOutingPerm')" width="90">
        <template #default="{ row }">
          <el-tag v-if="row.outing_permission" :type="permissionTagType(row.outing_permission)" size="small">
            {{ row.outing_permission }}
          </el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" :label="t('members.colCreatedAt')" width="160">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column :label="t('members.colOperation')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" plain @click="openEditWithAuth(row)">{{ t('common.edit') }}</el-button>
          <el-button size="small" type="danger" plain @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog
    v-model="dialogVisible"
    :title="editingId ? t('members.editMember') : t('members.addMember')"
    width="760px"
    @closed="resetForm"
  >
    <div style="max-height: 70vh; overflow-y: auto; padding-right: 8px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="120px">

        <!-- 基本信息 -->
        <el-divider content-position="left">{{ t('members.sectionBasic') }}</el-divider>

        <el-form-item :label="t('members.nameCn')" prop="name_cn">
          <el-input v-model="form.name_cn" placeholder="请输入中文姓名" />
        </el-form-item>

        <el-form-item :label="t('members.nameEn')" prop="name_en">
          <el-input v-model="form.name_en" :placeholder="t('members.nameEnPlaceholder')" @input="nameEnManuallyEdited = true" />
        </el-form-item>

        <el-form-item :label="t('members.role')" prop="role">
          <el-select v-model="form.role" :disabled="!!editingId" :placeholder="t('members.rolePlaceholder')" style="width: 100%">
            <el-option :label="t('members.roleChild')" value="child" />
            <el-option :label="t('members.roleParent')" value="parent" />
          </el-select>
          <div v-if="editingId" style="font-size:12px;color:#909399;margin-top:2px">{{ t('members.roleImmutable') }}</div>
        </el-form-item>

        <el-form-item :label="t('members.gender')" prop="gender">
          <el-radio-group v-model="form.gender">
            <el-radio value="男">{{ t('members.genderMale') }}</el-radio>
            <el-radio value="女">{{ t('members.genderFemale') }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item :label="t('members.nationality')" prop="nationality">
          <el-select
            v-model="form.nationality"
            :placeholder="t('members.nationalityPlaceholder')"
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

        <el-form-item :label="t('members.birthDate')" prop="birth_date">
          <el-date-picker
            v-model="form.birth_date"
            type="date"
            :placeholder="t('members.birthDatePlaceholder')"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <!-- 主证件 -->
        <el-divider content-position="left">{{ t('members.sectionId') }}</el-divider>

        <el-form-item :label="t('members.idDocType')" prop="id_doc_type">
          <el-select
            v-model="form.id_doc_type"
            :placeholder="t('members.idDocTypePlaceholder')"
            style="width: 100%"
            @change="onIdDocTypeChange"
          >
            <el-option
              v-for="dt in ID_DOC_TYPES"
              :key="dt.code"
              :label="dt.name"
              :value="dt.code"
            />
          </el-select>
        </el-form-item>

        <el-form-item :label="t('members.idDocNumber')" prop="id_doc_number" :error="idDocNumberError">
          <el-input v-model="form.id_doc_number" :placeholder="t('members.idDocNumberPlaceholder')" />
        </el-form-item>

        <el-form-item :label="t('members.issueDate')" prop="id_issue_date">
          <el-date-picker
            v-model="form.id_issue_date"
            type="date"
            :placeholder="t('members.issueDatePlaceholder')"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item :label="t('members.expiryDate')" prop="id_expiry_date">
          <template v-if="isLongTermDoc">
            <div style="display: flex; align-items: center; gap: 12px; width: 100%">
              <el-checkbox v-model="longTermEnabled" @change="onLongTermChange">{{ t('members.longTermCheckbox') }}</el-checkbox>
              <span v-if="longTermEnabled" style="color: #E6A23C; font-weight: 600">{{ t('members.longTerm') }}（2099-12-31）</span>
              <el-date-picker
                v-if="!longTermEnabled"
                v-model="form.id_expiry_date"
                type="date"
                :placeholder="t('members.expiryDatePlaceholder')"
                value-format="YYYY-MM-DD"
                style="flex: 1"
              />
            </div>
          </template>
          <template v-else>
            <div style="width: 100%">
              <el-date-picker
                v-model="form.id_expiry_date"
                type="date"
                :placeholder="t('members.expiryDatePlaceholder')"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
              <div v-if="calcExpiryDate" style="font-size: 12px; color: #909399; margin-top: 4px">
                {{ t('members.expectedExpiry') }}: <span style="color: #409EFF; font-weight: 600">{{ calcExpiryDate }}</span>
              </div>
            </div>
          </template>
        </el-form-item>

        <el-form-item :label="t('members.issueAuthority')" prop="id_issue_authority">
          <el-input v-model="form.id_issue_authority" :placeholder="t('members.issueAuthorityPlaceholder')" />
        </el-form-item>

        <!-- 辅助证件（根据主证件类型动态显示） -->
        <template v-if="showAuxDocs">
          <el-divider content-position="left">{{ t('members.sectionAux') }}</el-divider>

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
                v-for="dt in aux1TypeOptions"
                :key="dt.code"
                :label="dt.name"
                :value="dt.code"
              />
            </el-select>
          </el-form-item>

          <el-form-item :label="t('members.aux1Number')" prop="aux1_doc_number" :error="aux1DocNumberError">
            <el-input v-model="form.aux1_doc_number" :placeholder="aux1NumPlaceholder" />
          </el-form-item>

          <!-- 辅助证件2 (type 11 or 21) -->
          <template v-if="showAux2">
            <el-form-item :label="t('members.aux2Type')" prop="aux2_doc_type">
              <el-select
                v-model="form.aux2_doc_type"
                :placeholder="t('members.aux2TypePlaceholder')"
                style="width: 100%"
                @change="onAux2DocTypeChange"
              >
                <el-option
                  v-for="dt in aux2TypeOptions"
                  :key="dt.code"
                  :label="dt.name"
                  :value="dt.code"
                />
              </el-select>
            </el-form-item>

            <el-form-item :label="t('members.aux2Number')" prop="aux2_doc_number" :error="aux2DocNumberError">
              <el-input v-model="form.aux2_doc_number" :placeholder="t('members.aux2NumberPlaceholder')" />
            </el-form-item>
          </template>

          <!-- 主证件04：证明文件补充字段 -->
          <template v-if="form.id_doc_type === '04'">
            <el-form-item :label="t('members.proofDocType')" prop="proof_doc_type">
              <el-select
                v-model="form.proof_doc_type"
                :placeholder="t('members.proofDocTypePlaceholder')"
                style="width: 100%"
                @change="onProofDocTypeChange"
              >
                <el-option
                  v-for="pt in PROOF_DOC_TYPES"
                  :key="pt.code"
                  :label="pt.name"
                  :value="pt.code"
                />
              </el-select>
            </el-form-item>

            <el-form-item :label="t('members.proofIssueCountry')" prop="proof_issue_country">
              <el-select
                v-model="form.proof_issue_country"
                :placeholder="t('members.proofIssueCountryPlaceholder')"
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
        <el-divider content-position="left">{{ t('members.sectionSchool') }}</el-divider>

        <el-form-item :label="t('members.school')" prop="school_name">
          <el-input v-model="form.school_name" :placeholder="t('members.schoolPlaceholder')" />
        </el-form-item>

        <el-form-item :label="t('members.grade')" prop="grade">
          <el-select v-model="form.grade" :placeholder="t('members.gradePlaceholder')" clearable style="width: 100%">
            <el-option v-for="g in GRADES" :key="g" :label="g" :value="g" />
          </el-select>
        </el-form-item>

        <el-form-item :label="t('members.className')" prop="class_name">
          <el-select v-model="form.class_name" :placeholder="t('members.classPlaceholder')" clearable style="width: 100%">
            <el-option v-for="c in CLASSES" :key="c" :label="`${c}${t('members.classSuffix')}`" :value="c" />
          </el-select>
        </el-form-item>

        <el-form-item :label="t('members.teacherName')" prop="class_teacher_name">
          <el-input v-model="form.class_teacher_name" :placeholder="t('members.teacherNamePlaceholder')" />
        </el-form-item>

        <el-form-item :label="t('members.teacherPhone')" prop="class_teacher_phone">
          <el-input v-model="form.class_teacher_phone" :placeholder="t('members.teacherPhonePlaceholder')" />
        </el-form-item>

        <!-- 外出权限 -->
        <el-divider content-position="left">{{ t('members.sectionOuting') }}</el-divider>

        <el-form-item :label="t('members.outingPermission')" prop="outing_permission">
          <el-radio-group v-model="form.outing_permission">
            <el-radio value="许可">{{ t('members.outingPermAllowed') }}</el-radio>
            <el-radio value="不许可">{{ t('members.outingPermDenied') }}</el-radio>
            <el-radio value="受限">{{ t('members.outingPermRestricted') }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <template v-if="form.outing_permission === '受限'">
          <el-form-item :label="t('members.outingDates')" prop="outing_dates">
            <el-date-picker
              v-model="outingDatesArray"
              type="dates"
              placeholder="选择允许外出的日期"
              value-format="YYYY-MM-DD"
              style="width: 100%"
              @change="onOutingDatesChange"
            />
          </el-form-item>

          <el-form-item :label="t('members.outingTimeRanges')" prop="outing_time_ranges">
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
                <span>{{ t('members.outingTimeTo') }}</span>
                <el-time-picker
                  v-model="range.end"
                  placeholder="结束时间"
                  value-format="HH:mm"
                  format="HH:mm"
                  style="width: 140px"
                  @change="syncTimeRangesToForm"
                />
                <el-button type="danger" size="small" plain @click="removeTimeRange(idx)">{{ t('members.deleteTimeRange') }}</el-button>
              </div>
              <el-button size="small" @click="addTimeRange">{{ t('members.addTimeRange') }}</el-button>
            </div>
          </el-form-item>
        </template>

      </el-form>
    </div>

    <template #footer>
      <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" :loading="saving" @click="handleSave">{{ t('common.save') }}</el-button>
    </template>
  </el-dialog>

  <!-- 授权密码对话框（编辑/删除前验证） -->
  <el-dialog v-model="authDialogVisible" :title="t('members.authDialogTitle')" width="360px" :close-on-click-modal="false">
    <el-form @submit.prevent>
      <el-form-item :label="t('members.authPasswordLabel')" label-width="80px">
        <el-input
          v-model="authPassword"
          type="password"
          show-password
          :placeholder="t('members.authPasswordPlaceholder')"
          @keyup.enter="confirmAuthAction"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="authDialogVisible = false">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" :loading="authLoading" @click="confirmAuthAction">{{ t('common.confirm') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
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

const { t } = useI18n()

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
  if (form.id_doc_type === '04') return t('members.aux1TypeFixed')
  return t('members.aux1Type')
})

const aux1NumPlaceholder = computed(() => {
  const req = ['11', '21', '04'].includes(form.id_doc_type)
  return req ? t('members.aux1NumberRequired') : t('members.aux1NumberOptional')
})

const aux1TypePlaceholder = computed(() => {
  const req = ['11', '21', '04'].includes(form.id_doc_type)
  return req ? t('members.aux1TypeRequired') : t('members.aux1TypeOptional')
})

// Aux2 type options based on primary doc type
const aux2TypeOptions = computed(() => {
  switch (form.id_doc_type) {
    case '11': return AUX_DOC_TYPES.filter(dt => ['90', '92', '96', '97'].includes(dt.code))
    case '21': return AUX_DOC_TYPES.filter(dt => dt.code === '93')
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
    name_cn: requireChinese ? [{ required: true, message: t('members.nameCnRequired'), trigger: 'blur' }] : [],
    name_en: [{ required: true, message: t('members.nameEnRequired'), trigger: 'blur' }],
    role: [{ required: true, message: t('members.roleRequired'), trigger: 'change' }],
    gender: [{ required: true, message: t('members.genderRequired'), trigger: 'change' }],
    nationality: [{ required: true, message: t('members.nationalityRequired'), trigger: 'change' }],
    birth_date: [
      { required: true, message: t('members.birthDateRequired'), trigger: 'change' },
      { validator: validateBirthDateRule, trigger: 'change' }
    ],
    id_doc_type: [{ required: true, message: t('members.idDocTypeRequired'), trigger: 'change' }],
    id_doc_number: [{ required: true, message: t('members.idDocNumberRequired'), trigger: 'blur' }],
    id_issue_date: [
      { required: true, message: t('members.issueDateRequired'), trigger: 'change' },
      { validator: validateIssueDateRule, trigger: 'change' }
    ],
    id_expiry_date: [
      { required: true, message: t('members.expiryDateRequired'), trigger: 'change' },
      { validator: validateExpiryDateRule, trigger: 'change' }
    ],
    id_issue_authority: [{ required: true, message: t('members.issueAuthorityRequired'), trigger: 'blur' }],
    aux1_doc_type: needAux1 ? [{ required: true, message: t('members.aux1TypeRequired2'), trigger: 'change' }] : [],
    aux1_doc_number: needAux1 ? [{ required: true, message: t('members.aux1NumberRequired2'), trigger: 'blur' }] : [],
    aux2_doc_type: needAux2 ? [{ required: true, message: t('members.aux2TypeRequired'), trigger: 'change' }] : [],
    aux2_doc_number: needAux2 ? [{ required: true, message: t('members.aux2NumberRequired'), trigger: 'blur' }] : [],
    proof_doc_type: mainType === '04' ? [{ required: true, message: t('members.proofDocTypeRequired'), trigger: 'change' }] : [],
    proof_issue_country: mainType === '04' ? [{ required: true, message: t('members.proofIssueCountryRequired'), trigger: 'change' }] : [],
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
  longTermEnabled.value = false
})

// ─── Document date helpers ────────────────────────────────────────────────────

// Whether this doc type allows long-term expiry (type 01, age ≥ 46)
const isLongTermDoc = computed(() => form.id_doc_type === '01')

// Long-term checkbox state (type 01, age ≥ 46)
const longTermEnabled = ref(false)

function onLongTermChange(val) {
  if (val) {
    form.id_expiry_date = '2099-12-31'
  } else {
    // Clear and let calcExpiryDate repopulate if possible
    form.id_expiry_date = ''
    const exp = calcExpiryDate.value
    if (exp && exp !== '2099-12-31') {
      form.id_expiry_date = exp
    }
  }
}

// Helper: parse a YYYY-MM-DD string into a Date (UTC midnight)
function parseDate(s) {
  if (!s) return null
  const d = new Date(s + 'T00:00:00Z')
  return isNaN(d.getTime()) ? null : d
}

// Helper: isLeapYear
function isLeapYearJS(y) {
  return (y % 4 === 0 && y % 100 !== 0) || y % 400 === 0
}

// Helper: addSameDayYears (mirrors Go logic)
function addSameDayYearsJS(d, years) {
  const targetYear = d.getUTCFullYear() + years
  const month = d.getUTCMonth() // 0-based
  const day = d.getUTCDate()
  // Feb 29 → if target year not leap → Mar 1
  if (month === 1 && day === 29 && !isLeapYearJS(targetYear)) {
    return new Date(Date.UTC(targetYear, 2, 1)) // March 1
  }
  return new Date(Date.UTC(targetYear, month, day))
}

// Helper: addYearsMinusOneDay (mirrors Go logic)
function addYearsMinusOneDayJS(d, years) {
  const result = new Date(d.getTime())
  result.setUTCFullYear(result.getUTCFullYear() + years)
  result.setUTCDate(result.getUTCDate() - 1)
  return result
}

// Helper: ageAtDate (mirrors Go logic)
function ageAtDateJS(birthDate, targetDate) {
  let years = targetDate.getUTCFullYear() - birthDate.getUTCFullYear()
  let effectiveBirthday
  const bMonth = birthDate.getUTCMonth()
  const bDay = birthDate.getUTCDate()
  const tYear = targetDate.getUTCFullYear()
  if (bMonth === 1 && bDay === 29) {
    // Feb 29 birthday → always treated as Mar 1
    effectiveBirthday = new Date(Date.UTC(tYear, 2, 1))
  } else {
    effectiveBirthday = new Date(Date.UTC(tYear, bMonth, bDay))
  }
  if (targetDate.getTime() < effectiveBirthday.getTime()) {
    years--
  }
  return years
}

// Format date as YYYY-MM-DD
function fmtDate(d) {
  if (!d) return ''
  return d.toISOString().slice(0, 10)
}

// Reactive computed: expected expiry date based on doc type + birth + issue + role
const calcExpiryDate = computed(() => {
  const docType = form.id_doc_type
  if (!docType) return ''

  const birthDate = parseDate(form.birth_date)
  const issueDate = parseDate(form.id_issue_date)
  if (!issueDate) return ''

  switch (docType) {
    case '01': {
      if (!birthDate) return ''
      const ageAtIssue = ageAtDateJS(birthDate, issueDate)
      if (ageAtIssue >= 46) return '2099-12-31'
      if (ageAtIssue >= 26) return fmtDate(addSameDayYearsJS(issueDate, 20))
      if (ageAtIssue >= 16) return fmtDate(addSameDayYearsJS(issueDate, 10))
      return fmtDate(addSameDayYearsJS(issueDate, 5))
    }
    case '91': {
      if (!birthDate) return ''
      // Expiry = 16th birthday
      const bYear = birthDate.getUTCFullYear()
      const bMonth = birthDate.getUTCMonth()
      const bDay = birthDate.getUTCDate()
      if (bMonth === 1 && bDay === 29 && !isLeapYearJS(bYear + 16)) {
        return fmtDate(new Date(Date.UTC(bYear + 16, 2, 1)))
      }
      return fmtDate(new Date(Date.UTC(bYear + 16, bMonth, bDay)))
    }
    case '11':
    case '21':
      return fmtDate(addSameDayYearsJS(issueDate, 5))
    case '31': {
      if (!birthDate) return ''
      const ageAtIssue = ageAtDateJS(birthDate, issueDate)
      return fmtDate(addYearsMinusOneDayJS(issueDate, ageAtIssue < 18 ? 5 : 10))
    }
    case '02': {
      if (!birthDate) return ''
      const ageAtIssue = ageAtDateJS(birthDate, issueDate)
      return fmtDate(addYearsMinusOneDayJS(issueDate, ageAtIssue < 18 ? 5 : 10))
    }
    case '03':
    case '52':
      return fmtDate(addYearsMinusOneDayJS(issueDate, 5))
    case '04': {
      if (!birthDate) return ''
      const ageAtIssue = ageAtDateJS(birthDate, issueDate)
      return fmtDate(addYearsMinusOneDayJS(issueDate, ageAtIssue < 16 ? 5 : 10))
    }
    case '05':
      // No fixed expiry, just max 10yr – return max
      return fmtDate(addSameDayYearsJS(issueDate, 10))
    default:
      return ''
  }
})

// Auto-fill expiry date when calcExpiryDate changes (for most types)
watch(calcExpiryDate, (newExpiry) => {
  if (!newExpiry) return
  const docType = form.id_doc_type
  // For type 05 don't auto-fill (user chooses any date ≤ max)
  if (docType === '05') return
  // For type 01 with long-term, handle separately
  if (docType === '01' && newExpiry === '2099-12-31') {
    longTermEnabled.value = true
    form.id_expiry_date = '2099-12-31'
    return
  }
  if (docType === '01' && longTermEnabled.value) return
  form.id_expiry_date = newExpiry
})

// ─── Date validation rules ────────────────────────────────────────────────────

function validateBirthDateRule(_rule, value, callback) {
  if (!value) { callback(); return }
  const d = parseDate(value)
  if (!d) { callback(); return }
  const today = new Date()
  today.setUTCHours(0, 0, 0, 0)
  if (d >= today) {
    callback(new Error(t('members.birthBeforeToday')))
    return
  }
  const hundredYearsAgo = new Date(today)
  hundredYearsAgo.setUTCFullYear(hundredYearsAgo.getUTCFullYear() - 100)
  if (d < hundredYearsAgo) {
    callback(new Error(t('members.birthWithin100Years')))
    return
  }
  callback()
}

function validateIssueDateRule(_rule, value, callback) {
  if (!value) { callback(); return }
  const d = parseDate(value)
  if (!d) { callback(); return }
  const today = new Date()
  today.setUTCHours(0, 0, 0, 0)
  if (d > today) {
    callback(new Error(t('members.issueDateNotFuture')))
    return
  }
  const twentyYearsAgo = new Date(today)
  twentyYearsAgo.setUTCFullYear(twentyYearsAgo.getUTCFullYear() - 20)
  if (d < twentyYearsAgo) {
    callback(new Error(t('members.issueDateWithin20Years')))
    return
  }
  callback()
}

function validateExpiryDateRule(_rule, value, callback) {
  if (!value) { callback(); return }
  // Long-term is always valid
  if (value === '2099-12-31') { callback(); return }
  const d = parseDate(value)
  if (!d) { callback(); return }
  const today = new Date()
  today.setUTCHours(0, 0, 0, 0)
  if (d <= today) {
    callback(new Error(t('members.expiryAfterToday')))
    return
  }
  callback()
}

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
    callback(new Error(t('members.schoolKeywordError')))
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
      ElMessage.warning(t('members.natMismatchWarning'))
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
        t('members.highRiskCountryConfirm'),
        t('members.highRiskCountryTitle'),
        {
          confirmButtonText: t('members.highRiskConfirmBtn'),
          cancelButtonText: t('common.cancel'),
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
    ElMessage.error(t('members.loadFailed'))
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
    ElMessage.warning(t('members.authPasswordRequired'))
    return
  }
  authLoading.value = true
  try {
    const { type, row } = authAction.value
    if (type === 'edit') {
      pendingAuthPassword.value = authPassword.value
      authDialogVisible.value = false
      if (row) {
        openDialog(row)
      }
    } else if (type === 'delete') {
      await deleteMemberWithAuth(row.id, authPassword.value)
      ElMessage.success(t('members.deletedSuccess'))
      authDialogVisible.value = false
      await fetchMembers()
    }
  } catch (e) {
    const msg = e.response?.data?.message || t('common.failed')
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
    // Set long-term checkbox based on expiry date
    longTermEnabled.value = form.id_expiry_date === '2099-12-31'
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
  longTermEnabled.value = false
  nameEnManuallyEdited.value = false
  pendingAuthPassword.value = ''
  formRef.value?.resetFields()
}

async function handleSave() {
  try {
    await formRef.value.validate()
  } catch {
    ElMessage.error(t('members.formError'))
    return
  }

  // Block save if there are ID number validation errors
  if (idDocNumberError.value) {
    ElMessage.error(t('members.mainDocNumberError'))
    return
  }
  if (aux1DocNumberError.value) {
    ElMessage.error(t('members.aux1NumberError'))
    return
  }
  if (aux2DocNumberError.value) {
    ElMessage.error(t('members.aux2NumberError'))
    return
  }

  saving.value = true
  try {
    const payload = { ...form, name: form.name_cn || form.name_en }
    if (editingId.value) {
      payload.auth_password = pendingAuthPassword.value
      await updateMember(editingId.value, payload)
      ElMessage.success(t('members.saveSuccess'))
    } else {
      payload.auth_password = pendingAuthPassword.value
      await createMember(payload)
      ElMessage.success(t('members.createSuccess'))
    }
    dialogVisible.value = false
    pendingAuthPassword.value = ''
    await fetchMembers()
  } catch (e) {
    const msg = e.response?.data?.message || t('members.saveFailed')
    ElMessage.error(msg)
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
      t('members.deleteConfirmMsg', { name: getMemberDisplayName(row) }),
      t('members.deleteConfirmTitle'),
      { type: 'warning', confirmButtonText: t('members.deleteBtn'), cancelButtonText: t('common.cancel') }
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
