<template>
  <el-card shadow="never">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px">
        <span style="font-weight: 600">{{ t('records.title') }}</span>
        <div style="display: flex; gap: 12px; align-items: center">
          <el-select
            v-model="filterMemberId"
            :placeholder="t('records.filterMember')"
            clearable
            style="width: 160px"
            @change="loadRecords"
          >
            <el-option
              v-for="m in members"
              :key="m.id"
              :label="m.name"
              :value="m.id"
            />
          </el-select>
          <el-button type="primary" @click="openDialog">
            <el-icon><Plus /></el-icon>
            {{ t('records.addRecord') }}
          </el-button>
        </div>
      </div>
    </template>

    <el-table :data="records" v-loading="loading" stripe style="width: 100%">
      <el-table-column prop="member_name" :label="t('records.colMember')" width="120" />
      <el-table-column prop="rule_name" :label="t('records.colViolation')" min-width="140" />
      <el-table-column prop="points" :label="t('records.colPoints')" width="80">
        <template #default="{ row }">
          <el-tag type="danger">{{ row.points }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="note" :label="t('records.colNote')" min-width="150" show-overflow-tooltip />
      <el-table-column prop="occurred_at" :label="t('records.colOccurredAt')" width="160">
        <template #default="{ row }">{{ formatDate(row.occurred_at) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.operation')" width="100" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="danger" plain @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog
    v-model="dialogVisible"
    :title="t('records.dialogTitle')"
    width="500px"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="formRules" label-width="90px">
      <el-form-item :label="t('records.member')" prop="member_id">
        <el-select v-model="form.member_id" :placeholder="t('records.memberPlaceholder')" style="width: 100%">
          <el-option
            v-for="m in members"
            :key="m.id"
            :label="m.name"
            :value="m.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('records.rule')" prop="rule_id">
        <el-select
          v-model="form.rule_id"
          :placeholder="t('records.rulePlaceholder')"
          style="width: 100%"
          @change="onRuleChange"
        >
          <el-option
            v-for="r in rulesList"
            :key="r.id"
            :label="r.name"
            :value="r.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('records.points')" prop="points">
        <el-input-number v-model="form.points" :min="1" :max="1000" style="width: 100%" />
      </el-form-item>
      <el-form-item :label="t('records.note')" prop="note">
        <el-input v-model="form.note" :placeholder="t('records.notePlaceholder')" />
      </el-form-item>
      <el-form-item :label="t('records.occurredAt')" prop="occurred_at">
        <el-date-picker
          v-model="form.occurred_at"
          type="datetime"
          :placeholder="t('records.occurredAtPlaceholder')"
          style="width: 100%"
          format="YYYY-MM-DD HH:mm"
          value-format="YYYY-MM-DDTHH:mm:ssZ"
        />
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { getRecords, createRecord, deleteRecord, getMembers, getRules } from '@/utils/api'

const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)
const records = ref([])
const members = ref([])
const rulesList = ref([])
const dialogVisible = ref(false)
const filterMemberId = ref(null)
const formRef = ref(null)

const form = reactive({
  member_id: null,
  rule_id: null,
  points: 1,
  note: '',
  occurred_at: new Date().toISOString()
})

const formRules = computed(() => ({
  member_id: [{ required: true, message: t('records.memberRequired'), trigger: 'change' }],
  rule_id: [{ required: true, message: t('records.ruleRequired'), trigger: 'change' }],
  points: [{ required: true, message: t('records.pointsRequired'), trigger: 'blur' }]
}))

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

async function loadRecords() {
  loading.value = true
  try {
    const params = {}
    if (filterMemberId.value) params.member_id = filterMemberId.value
    const res = await getRecords(params)
    records.value = res.data?.records ?? res.data ?? []
  } catch {
    ElMessage.error(t('records.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function loadMembersAndRules() {
  try {
    const [mRes, rRes] = await Promise.all([getMembers(), getRules()])
    members.value = mRes.data ?? []
    rulesList.value = rRes.data ?? []
  } catch {
    ElMessage.error(t('records.loadDataFailed'))
  }
}

function onRuleChange(ruleId) {
  const rule = rulesList.value.find(r => r.id === ruleId)
  if (rule) form.points = rule.points
}

function openDialog() {
  form.member_id = null
  form.rule_id = null
  form.points = 1
  form.note = ''
  form.occurred_at = new Date().toISOString()
  dialogVisible.value = true
}

function resetForm() {
  formRef.value?.resetFields()
}

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    await createRecord({ ...form })
    ElMessage.success(t('records.addSuccess'))
    dialogVisible.value = false
    await loadRecords()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || t('records.addFailed'))
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm(
    t('records.deleteConfirm'),
    t('records.deleteTitle'),
    { type: 'warning', confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel') }
  )
  try {
    await deleteRecord(row.id)
    ElMessage.success(t('records.deleteSuccess'))
    await loadRecords()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || t('records.deleteFailed'))
  }
}

onMounted(async () => {
  await loadMembersAndRules()
  await loadRecords()
})
</script>
