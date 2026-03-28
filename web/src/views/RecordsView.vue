<template>
  <el-card shadow="never">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px">
        <span style="font-weight: 600">惩戒记录</span>
        <div style="display: flex; gap: 12px; align-items: center">
          <el-select
            v-model="filterMemberId"
            placeholder="筛选成员"
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
            添加记录
          </el-button>
        </div>
      </div>
    </template>

    <el-table :data="records" v-loading="loading" stripe style="width: 100%">
      <el-table-column prop="member_name" label="成员" width="120" />
      <el-table-column prop="rule_name" label="违规项目" min-width="140" />
      <el-table-column prop="points" label="分值" width="80">
        <template #default="{ row }">
          <el-tag type="danger">{{ row.points }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="note" label="备注" min-width="150" show-overflow-tooltip />
      <el-table-column prop="occurred_at" label="发生时间" width="160">
        <template #default="{ row }">{{ formatDate(row.occurred_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="danger" plain @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog
    v-model="dialogVisible"
    title="添加惩戒记录"
    width="500px"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="formRules" label-width="90px">
      <el-form-item label="成员" prop="member_id">
        <el-select v-model="form.member_id" placeholder="请选择成员" style="width: 100%">
          <el-option
            v-for="m in members"
            :key="m.id"
            :label="m.name"
            :value="m.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="惩戒规则" prop="rule_id">
        <el-select
          v-model="form.rule_id"
          placeholder="请选择规则"
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
      <el-form-item label="分值" prop="points">
        <el-input-number v-model="form.points" :min="1" :max="1000" style="width: 100%" />
      </el-form-item>
      <el-form-item label="备注" prop="note">
        <el-input v-model="form.note" placeholder="可选备注" />
      </el-form-item>
      <el-form-item label="发生时间" prop="occurred_at">
        <el-date-picker
          v-model="form.occurred_at"
          type="datetime"
          placeholder="选择时间"
          style="width: 100%"
          format="YYYY-MM-DD HH:mm"
          value-format="YYYY-MM-DDTHH:mm:ssZ"
        />
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
import { getRecords, createRecord, deleteRecord, getMembers, getRules } from '@/utils/api'

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

const formRules = {
  member_id: [{ required: true, message: '请选择成员', trigger: 'change' }],
  rule_id: [{ required: true, message: '请选择规则', trigger: 'change' }],
  points: [{ required: true, message: '请输入分值', trigger: 'blur' }]
}

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
    ElMessage.error('加载记录失败')
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
    ElMessage.error('加载数据失败')
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
    ElMessage.success('记录添加成功')
    dialogVisible.value = false
    await loadRecords()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '添加失败')
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm('确定要删除此条记录吗？', '确认删除', {
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  })
  try {
    await deleteRecord(row.id)
    ElMessage.success('删除成功')
    await loadRecords()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '删除失败')
  }
}

onMounted(async () => {
  await loadMembersAndRules()
  await loadRecords()
})
</script>
