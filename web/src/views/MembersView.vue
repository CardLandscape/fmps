<template>
  <el-card shadow="never">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <span style="font-weight: 600">{{ i18n.t('membersPageTitle') }}</span>
        <el-button type="primary" @click="openDialog()">
          <el-icon><Plus /></el-icon>
          {{ i18n.t('btnAddMember') }}
        </el-button>
      </div>
    </template>

    <el-table :data="members" v-loading="loading" stripe style="width: 100%">
      <el-table-column :label="i18n.t('colName')">
        <template #default="{ row }">{{ getMemberDisplayName(row) }}</template>
      </el-table-column>
      <el-table-column prop="role" :label="i18n.t('colRole')" width="80">
        <template #default="{ row }">
          <el-tag :type="row.role === 'child' ? 'warning' : 'success'">
            {{ row.role === 'child' ? i18n.t('roleChild') : i18n.t('roleParent') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="gender" :label="i18n.t('colGender')" width="70" />
      <el-table-column prop="nationality" :label="i18n.t('colNationality')" width="80" />
      <el-table-column prop="school_name" :label="i18n.t('colSchool')" />
      <el-table-column prop="outing_permission" :label="i18n.t('colOutingPerm')" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.outing_permission" :type="permissionTagType(row.outing_permission)" size="small">
            {{ outingPermLabel(row.outing_permission) }}
          </el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" :label="i18n.t('colCreatedAt')" width="160">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column :label="i18n.t('colActions')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" plain @click="openEditWithAuth(row)">{{ i18n.t('btnEdit') }}</el-button>
          <el-button size="small" type="danger" plain @click="handleDelete(row)">{{ i18n.t('btnDelete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog
    v-model="dialogVisible"
    :title="editingId ? i18n.t('dialogEditMember') : i18n.t('dialogAddMember')"
    width="760px"
    @closed="resetForm"
  >
    <div style="max-height: 70vh; overflow-y: auto; padding-right: 8px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="140px">

        <!-- Basic info -->
        <el-divider content-position="left">{{ i18n.t('sectionBasicInfo') }}</el-divider>

        <el-form-item :label="i18n.t('labelNameCn')" prop="name_cn">
          <el-input v-model="form.name_cn" :placeholder="i18n.t('placeholderNameCn')" />
        </el-form-item>

        <el-form-item :label="i18n.t('labelNameEn')" prop="name_en">
          <el-input v-model="form.name_en" :placeholder="i18n.t('placeholderNameEn')" @input="nameEnManuallyEdited = true" />
        </el-form-item>

        <el-form-item :label="i18n.t('labelRole')" prop="role">
          <el-select v-model="form.role" :disabled="!!editingId" :placeholder="i18n.t('placeholderRole')" style="width: 100%">
            <el-option :label="i18n.t('roleChild')" value="child" />
            <el-option :label="i18n.t('roleParent')" value="parent" />
          </el-select>
          <div v-if="editingId" style="font-size:12px;color:#909399;margin-top:2px">{{ i18n.t('roleFixedHint') }}</div>
        </el-form-item>

        <el-form-item :label="i18n.t('labelGender')" prop="gender" :error="genderIdError">
          <el-radio-group v-model="form.gender">
            <el-radio value="男">男</el-radio>
            <el-radio value="女">女</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item :label="i18n.t('labelNationality')" prop="nationality">
          <el-select
            v-model="form.nationality"
            :placeholder="i18n.t('placeholderNationality')"
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

        <el-form-item :label="i18n.t('labelBirthDate')" prop="birth_date" :error="birthDateError">
          <el-date-picker
            v-model="form.birth_date"
            type="date"
            :placeholder="i18n.t('placeholderBirthDate')"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <!-- Primary document -->
        <el-divider content-position="left">{{ i18n.t('sectionPrimaryDoc') }}</el-divider>

        <el-form-item :label="i18n.t('labelDocType')" prop="id_doc_type">
          <el-select
            v-model="form.id_doc_type"
            :placeholder="i18n.t('placeholderDocType')"
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

        <el-form-item :label="i18n.t('labelDocNumber')" prop="id_doc_number" :error="idDocNumberError">
          <el-input v-model="form.id_doc_number" :placeholder="i18n.t('placeholderDocNumber')" />
        </el-form-item>

        <el-form-item :label="i18n.t('labelIssueDate')" prop="id_issue_date" :error="issueDateError">
          <el-date-picker
            v-model="form.id_issue_date"
            type="date"
            :placeholder="i18n.t('placeholderIssueDate')"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item :label="i18n.t('labelExpiryDate')" prop="id_expiry_date" :error="expiryDateError">
          <el-date-picker
            v-model="form.id_expiry_date"
            type="date"
            :placeholder="i18n.t('placeholderExpiryDate')"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item :label="i18n.t('labelIssueAuth')" prop="id_issue_authority">
          <el-input v-model="form.id_issue_authority" :placeholder="i18n.t('placeholderIssueAuth')" />
        </el-form-item>

        <!-- Auxiliary documents (shown based on primary doc type) -->
        <template v-if="showAuxDocs">
          <el-divider content-position="left">{{ i18n.t('sectionAuxDoc') }}</el-divider>

          <!-- Aux doc 1 -->
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

          <el-form-item :label="i18n.t('labelAux1Number')" prop="aux1_doc_number" :error="aux1DocNumberError">
            <el-input v-model="form.aux1_doc_number" :placeholder="aux1NumPlaceholder" />
          </el-form-item>

          <!-- Aux doc 2 (type 11 or 21 only) -->
          <template v-if="showAux2">
            <el-form-item :label="i18n.t('labelAux2Type')" prop="aux2_doc_type">
              <el-select
                v-model="form.aux2_doc_type"
                :placeholder="i18n.t('validAux2TypeRequired')"
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

            <el-form-item :label="i18n.t('labelAux2Number')" prop="aux2_doc_number" :error="aux2DocNumberError">
              <el-input v-model="form.aux2_doc_number" :placeholder="i18n.t('validAux2NumberRequired')" />
            </el-form-item>
          </template>

          <!-- Primary doc type 04: proof document fields -->
          <template v-if="form.id_doc_type === '04'">
            <el-form-item :label="i18n.t('labelProofDocType')" prop="proof_doc_type">
              <el-select
                v-model="form.proof_doc_type"
                :placeholder="i18n.t('validProofDocTypeRequired')"
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

            <el-form-item :label="i18n.t('labelProofIssueCountry')" prop="proof_issue_country">
              <el-select
                v-model="form.proof_issue_country"
                :placeholder="i18n.t('validProofIssueCountryRequired')"
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

        <!-- School info -->
        <el-divider content-position="left">{{ i18n.t('sectionSchool') }}</el-divider>

        <el-form-item :label="i18n.t('labelSchoolName')" prop="school_name">
          <el-input v-model="form.school_name" :placeholder="i18n.t('placeholderSchoolName')" />
        </el-form-item>

        <el-form-item :label="i18n.t('labelGrade')" prop="grade">
          <el-select v-model="form.grade" placeholder="-" clearable style="width: 100%">
            <el-option v-for="g in GRADES" :key="g" :label="g" :value="g" />
          </el-select>
        </el-form-item>

        <el-form-item :label="i18n.t('labelClass')" prop="class_name">
          <el-select v-model="form.class_name" placeholder="-" clearable style="width: 100%">
            <el-option v-for="c in CLASSES" :key="c" :label="`${c}班`" :value="c" />
          </el-select>
        </el-form-item>

        <el-form-item :label="i18n.t('labelClassTeacher')" prop="class_teacher_name">
          <el-input v-model="form.class_teacher_name" :placeholder="i18n.t('placeholderClassTeacher')" />
        </el-form-item>

        <el-form-item :label="i18n.t('labelClassTeacherPhone')" prop="class_teacher_phone">
          <el-input v-model="form.class_teacher_phone" :placeholder="i18n.t('placeholderClassTeacherPhone')" />
        </el-form-item>

        <!-- Outing permission -->
        <el-divider content-position="left">{{ i18n.t('sectionOuting') }}</el-divider>

        <el-form-item :label="i18n.t('labelOutingPerm')" prop="outing_permission">
          <el-radio-group v-model="form.outing_permission">
            <el-radio value="许可">{{ i18n.t('outingAllowed') }}</el-radio>
            <el-radio value="不许可">{{ i18n.t('outingDenied') }}</el-radio>
            <el-radio value="受限">{{ i18n.t('outingRestricted') }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <template v-if="form.outing_permission === '受限'">
          <el-form-item :label="i18n.t('labelOutingDates')" prop="outing_dates">
            <el-date-picker
              v-model="outingDatesArray"
              type="dates"
              placeholder="-"
              value-format="YYYY-MM-DD"
              style="width: 100%"
              @change="onOutingDatesChange"
            />
          </el-form-item>

          <el-form-item :label="i18n.t('labelOutingTimes')" prop="outing_time_ranges">
            <div style="width: 100%">
              <div
                v-for="(range, idx) in outingTimeRanges"
                :key="idx"
                style="display: flex; align-items: center; gap: 8px; margin-bottom: 8px"
              >
                <el-time-picker
                  v-model="range.start"
                  placeholder="--:--"
                  value-format="HH:mm"
                  format="HH:mm"
                  style="width: 140px"
                  @change="syncTimeRangesToForm"
                />
                <span>{{ i18n.t('timeTo') }}</span>
                <el-time-picker
                  v-model="range.end"
                  placeholder="--:--"
                  value-format="HH:mm"
                  format="HH:mm"
                  style="width: 140px"
                  @change="syncTimeRangesToForm"
                />
                <el-button type="danger" size="small" plain @click="removeTimeRange(idx)">{{ i18n.t('btnRemoveTimeRange') }}</el-button>
              </div>
              <el-button size="small" @click="addTimeRange">{{ i18n.t('btnAddTimeRange') }}</el-button>
            </div>
          </el-form-item>
        </template>

      </el-form>
    </div>

    <template #footer>
      <el-button @click="dialogVisible = false">{{ i18n.t('btnCancel') }}</el-button>
      <el-button type="primary" :loading="saving" @click="handleSave">{{ i18n.t('btnSave') }}</el-button>
    </template>
  </el-dialog>

  <!-- Authorization password dialog (before edit/delete) -->
  <el-dialog v-model="authDialogVisible" :title="i18n.t('authDialogTitle')" width="360px" :close-on-click-modal="false">
    <el-form @submit.prevent>
      <el-form-item :label="i18n.t('labelAuthPassword')" label-width="120px">
        <el-input
          v-model="authPassword"
          type="password"
          show-password
          :placeholder="i18n.t('placeholderAuthPassword')"
          @keyup.enter="confirmAuthAction"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="authDialogVisible = false">{{ i18n.t('btnCancel') }}</el-button>
      <el-button type="primary" :loading="authLoading" @click="confirmAuthAction">{{ i18n.t('btnConfirm') }}</el-button>
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
  validateNationalityDocType,
  validateBirthDate,
  validateIssueDate,
  validateExpiryDate,
  validateIDConsistency
} from '@/utils/constants'
import { useI18n } from '@/utils/i18n'

// ─── State ───────────────────────────────────────────────────────────────────

const i18n = useI18n()

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

// Date range validation error refs (translated strings)
const birthDateRangeError = ref('')
const issueDateError = ref('')
const expiryDateError = ref('')

// ID consistency error refs (translated strings)
// These show on the birth_date / gender fields when the entered value
// contradicts what is encoded in the ID number.
const birthDateIdError = ref('')
const genderIdError = ref('')

// Combined birth date error: range first, then ID-consistency
const birthDateError = computed(() => birthDateRangeError.value || birthDateIdError.value)

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

const aux1Label = computed(() => i18n.t('labelAux1Type'))

const aux1NumPlaceholder = computed(() => {
  const req = ['11', '21', '04'].includes(form.id_doc_type)
  return req ? i18n.t('validAux1NumberRequired') : i18n.t('labelAux1Number')
})

const aux1TypePlaceholder = computed(() => {
  const req = ['11', '21', '04'].includes(form.id_doc_type)
  return req ? i18n.t('validAux1TypeRequired') : i18n.t('labelAux1Type')
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
    name_cn: requireChinese ? [{ required: true, message: i18n.t('validNameCnRequired'), trigger: 'blur' }] : [],
    name_en: [{ required: true, message: i18n.t('validNameEnRequired'), trigger: 'blur' }],
    role: [{ required: true, message: i18n.t('validRoleRequired'), trigger: 'change' }],
    gender: [{ required: true, message: i18n.t('validGenderRequired'), trigger: 'change' }],
    nationality: [{ required: true, message: i18n.t('validNationalityRequired'), trigger: 'change' }],
    birth_date: [{ required: true, message: i18n.t('validBirthDateRequired'), trigger: 'change' }],
    id_doc_type: [{ required: true, message: i18n.t('validDocTypeRequired'), trigger: 'change' }],
    id_doc_number: [{ required: true, message: i18n.t('validDocNumberRequired'), trigger: 'blur' }],
    id_issue_date: [{ required: true, message: i18n.t('validIssueDateRequired'), trigger: 'change' }],
    id_expiry_date: [{ required: true, message: i18n.t('validExpiryDateRequired'), trigger: 'change' }],
    id_issue_authority: [{ required: true, message: i18n.t('validIssueAuthRequired'), trigger: 'blur' }],
    aux1_doc_type: needAux1 ? [{ required: true, message: i18n.t('validAux1TypeRequired'), trigger: 'change' }] : [],
    aux1_doc_number: needAux1 ? [{ required: true, message: i18n.t('validAux1NumberRequired'), trigger: 'blur' }] : [],
    aux2_doc_type: needAux2 ? [{ required: true, message: i18n.t('validAux2TypeRequired'), trigger: 'change' }] : [],
    aux2_doc_number: needAux2 ? [{ required: true, message: i18n.t('validAux2NumberRequired'), trigger: 'blur' }] : [],
    proof_doc_type: mainType === '04' ? [{ required: true, message: i18n.t('validProofDocTypeRequired'), trigger: 'change' }] : [],
    proof_issue_country: mainType === '04' ? [{ required: true, message: i18n.t('validProofIssueCountryRequired'), trigger: 'change' }] : [],
    school_name: [{ validator: validateSchoolNameRule, trigger: 'blur' }]
  }
})

// ─── Watchers ─────────────────────────────────────────────────────────────────

// ID number format-only validation (does NOT check birthdate/gender consistency).
// Consistency is validated separately below so mismatches are reported on the
// correct field rather than being shown as a generic "document number" error.
watch(
  () => [form.id_doc_number, form.id_doc_type, form.nationality],
  ([number, docType, nationality]) => {
    idDocNumberError.value = validateIDNumber(docType, number, nationality) || ''
  }
)

// ID-number ↔ birthdate / gender consistency watcher.
// When the value encoded in the ID differs from the separately entered field,
// show the error on THAT field rather than on the document number field.
watch(
  () => [form.id_doc_number, form.id_doc_type, form.birth_date, form.gender],
  ([number, docType, birthDate, gender]) => {
    const result = validateIDConsistency(docType, number, gender, birthDate)
    if (result?.field === 'birth_date') {
      // Only apply if the date itself has passed basic range validation
      if (!validateBirthDate(birthDate)) {
        birthDateIdError.value = i18n.t(result.key)
      } else {
        birthDateIdError.value = ''
      }
      genderIdError.value = ''
    } else if (result?.field === 'gender') {
      genderIdError.value = i18n.t(result.key)
      birthDateIdError.value = ''
    } else {
      birthDateIdError.value = ''
      genderIdError.value = ''
    }
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

// Date validation watchers
watch(
  () => form.birth_date,
  (val) => {
    const key = validateBirthDate(val)
    birthDateRangeError.value = key ? i18n.t(key) : ''
    // Clear consistency error when birth date changes (consistency watcher re-evaluates separately)
  }
)

watch(
  () => form.id_issue_date,
  (val) => {
    const key = validateIssueDate(val)
    issueDateError.value = key ? i18n.t(key) : ''
    // Re-validate expiry whenever issue date changes
    _revalidateExpiry()
  }
)

watch(
  () => [form.id_expiry_date, form.id_doc_type, form.birth_date, form.id_issue_date, form.role],
  () => { _revalidateExpiry() }
)

function _revalidateExpiry() {
  const key = validateExpiryDate(
    form.id_expiry_date,
    form.id_doc_type,
    form.birth_date,
    form.id_issue_date,
    form.role
  )
  expiryDateError.value = key ? i18n.t(key) : ''
}

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
    callback(new Error(i18n.t('schoolKeywordError')))
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

function outingPermLabel(p) {
  if (p === '许可') return i18n.t('outingAllowed')
  if (p === '不许可') return i18n.t('outingDenied')
  if (p === '受限') return i18n.t('outingRestricted')
  return p
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
      ElMessage.warning(i18n.t('nationalityIncompatible'))
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
        i18n.t('highRiskCountryMsg'),
        i18n.t('highRiskCountryTitle'),
        {
          confirmButtonText: i18n.t('highRiskConfirm'),
          cancelButtonText: i18n.t('btnCancel'),
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
    ElMessage.error(i18n.t('fetchMembersFailed'))
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
    ElMessage.warning(i18n.t('authPasswordRequired'))
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
      ElMessage.success(i18n.t('memberDeleted'))
      authDialogVisible.value = false
      await fetchMembers()
    }
  } catch (e) {
    const msg = e.response?.data?.message || i18n.t('saveFailed')
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
  birthDateRangeError.value = ''
  birthDateIdError.value = ''
  genderIdError.value = ''
  issueDateError.value = ''
  expiryDateError.value = ''
  nameEnManuallyEdited.value = false
  pendingAuthPassword.value = ''
  formRef.value?.resetFields()
}

async function handleSave() {
  try {
    await formRef.value.validate()
  } catch {
    ElMessage.error(i18n.t('formError'))
    return
  }

  // Block save if there are field-level validation errors.
  // Check each field individually so the toast identifies the exact issue.
  if (birthDateError.value) {
    ElMessage.error(birthDateError.value)
    return
  }
  if (genderIdError.value) {
    ElMessage.error(genderIdError.value)
    return
  }
  if (idDocNumberError.value) {
    ElMessage.error(i18n.t('idNumberError'))
    return
  }
  if (aux1DocNumberError.value) {
    ElMessage.error(i18n.t('aux1NumberError'))
    return
  }
  if (aux2DocNumberError.value) {
    ElMessage.error(i18n.t('aux2NumberError'))
    return
  }
  if (issueDateError.value) {
    ElMessage.error(issueDateError.value)
    return
  }
  if (expiryDateError.value) {
    ElMessage.error(expiryDateError.value)
    return
  }

  saving.value = true
  try {
    const payload = { ...form, name: form.name_cn || form.name_en }
    if (editingId.value) {
      // Include auth password for server-side validation
      payload.auth_password = pendingAuthPassword.value
      await updateMember(editingId.value, payload)
      ElMessage.success(i18n.t('memberUpdated'))
    } else {
      // For new parent records, include auth_password if provided (same-person child exists)
      payload.auth_password = pendingAuthPassword.value
      await createMember(payload)
      ElMessage.success(i18n.t('memberAdded'))
    }
    dialogVisible.value = false
    pendingAuthPassword.value = ''
    await fetchMembers()
  } catch (e) {
    const msg = e.response?.data?.message || i18n.t('saveFailed')
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
      i18n.t('confirmDeleteMember').replace('{name}', getMemberDisplayName(row)),
      i18n.t('confirmDeleteTitle'),
      { type: 'warning', confirmButtonText: i18n.t('confirmDeleteBtn'), cancelButtonText: i18n.t('btnCancel') }
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
