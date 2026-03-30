<template>
  <el-card shadow="never">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <span style="font-weight: 600">{{ i18n.t('rulesPageTitle') }}</span>
        <el-button type="primary" @click="openDialog()">
          <el-icon><Plus /></el-icon>
          {{ i18n.t('btnAddRule') }}
        </el-button>
      </div>
    </template>

    <el-table :data="rules" v-loading="loading" stripe style="width: 100%">
      <el-table-column prop="name" :label="i18n.t('colRuleName')" min-width="120" />
      <el-table-column prop="category" :label="i18n.t('colCategory')" width="120" />
      <el-table-column prop="description" :label="i18n.t('colDescription')" min-width="200" show-overflow-tooltip />
      <el-table-column prop="points" :label="i18n.t('colPoints')" width="80">
        <template #default="{ row }">
          <el-tag type="danger">{{ row.points }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="i18n.t('colActions')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" plain @click="openDialog(row)">{{ i18n.t('btnEdit') }}</el-button>
          <el-button size="small" type="danger" plain @click="handleDelete(row)">{{ i18n.t('btnDelete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog
    v-model="dialogVisible"
    :title="editingId ? i18n.t('dialogEditRule') : i18n.t('dialogAddRule')"
    width="480px"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="formRules" label-width="90px">
      <el-form-item :label="i18n.t('labelRuleName')" prop="name">
        <el-input v-model="form.name" :placeholder="i18n.t('placeholderRuleName')" />
      </el-form-item>
      <el-form-item :label="i18n.t('labelRuleCategory')" prop="category">
        <el-input v-model="form.category" :placeholder="i18n.t('placeholderRuleCategory')" />
      </el-form-item>
      <el-form-item :label="i18n.t('labelRuleDescription')" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          :placeholder="i18n.t('placeholderRuleDescription')"
        />
      </el-form-item>
      <el-form-item :label="i18n.t('labelRulePoints')" prop="points">
        <el-input-number v-model="form.points" :min="1" :max="100" style="width: 100%" />
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getRules, createRule, updateRule, deleteRule } from '@/utils/api'
import { useI18n } from '@/utils/i18n'

const i18n = useI18n()

const loading = ref(false)
const saving = ref(false)
const rules = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const form = reactive({ name: '', category: '', description: '', points: 1 })

const formRules = computed(() => ({
  name: [{ required: true, message: i18n.t('validRuleNameRequired'), trigger: 'blur' }],
  points: [{ required: true, message: i18n.t('validRulePointsRequired'), trigger: 'blur' }]
}))

async function loadRules() {
  loading.value = true
  try {
    const res = await getRules()
    rules.value = res.data ?? []
  } catch {
    ElMessage.error(i18n.t('loadRulesFailed'))
  } finally {
    loading.value = false
  }
}

function openDialog(row = null) {
  if (row) {
    editingId.value = row.id
    form.name = row.name
    form.category = row.category || ''
    form.description = row.description || ''
    form.points = row.points
  } else {
    editingId.value = null
    form.name = ''
    form.category = ''
    form.description = ''
    form.points = 1
  }
  dialogVisible.value = true
}

function resetForm() {
  formRef.value?.resetFields()
  editingId.value = null
}

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (editingId.value) {
      await updateRule(editingId.value, { ...form })
      ElMessage.success(i18n.t('updateSuccess'))
    } else {
      await createRule({ ...form })
      ElMessage.success(i18n.t('addSuccess'))
    }
    dialogVisible.value = false
    await loadRules()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || i18n.t('saveFailed'))
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm(
    i18n.t('confirmDeleteRule').replace('{name}', row.name),
    i18n.t('confirmDeleteTitle'),
    {
      type: 'warning',
      confirmButtonText: i18n.t('confirmDeleteBtn'),
      cancelButtonText: i18n.t('btnCancel')
    }
  )
  try {
    await deleteRule(row.id)
    ElMessage.success(i18n.t('deleteSuccess'))
    await loadRules()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || i18n.t('deleteFailed'))
  }
}

onMounted(loadRules)
</script>
