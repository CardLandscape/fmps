<template>
  <el-card shadow="never">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center">
        <span style="font-weight: 600">{{ t('rules.title') }}</span>
        <el-button type="primary" @click="openDialog()">
          <el-icon><Plus /></el-icon>
          {{ t('rules.addRule') }}
        </el-button>
      </div>
    </template>

    <el-table :data="rules" v-loading="loading" stripe style="width: 100%">
      <el-table-column prop="name" :label="t('rules.colName')" min-width="120" />
      <el-table-column prop="category" :label="t('rules.colCategory')" width="120" />
      <el-table-column prop="description" :label="t('rules.colDescription')" min-width="200" show-overflow-tooltip />
      <el-table-column prop="points" :label="t('rules.colPoints')" width="80">
        <template #default="{ row }">
          <el-tag type="danger">{{ row.points }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.operation')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" plain @click="openDialog(row)">{{ t('common.edit') }}</el-button>
          <el-button size="small" type="danger" plain @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog
    v-model="dialogVisible"
    :title="editingId ? t('rules.editRule') : t('rules.addRule')"
    width="480px"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="formRules" label-width="80px">
      <el-form-item :label="t('rules.name')" prop="name">
        <el-input v-model="form.name" :placeholder="t('rules.namePlaceholder')" />
      </el-form-item>
      <el-form-item :label="t('rules.category')" prop="category">
        <el-input v-model="form.category" :placeholder="t('rules.categoryPlaceholder')" />
      </el-form-item>
      <el-form-item :label="t('rules.description')" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          :placeholder="t('rules.descriptionPlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('rules.points')" prop="points">
        <el-input-number v-model="form.points" :min="1" :max="100" style="width: 100%" />
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
import { getRules, createRule, updateRule, deleteRule } from '@/utils/api'

const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)
const rules = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const form = reactive({ name: '', category: '', description: '', points: 1 })

const formRules = computed(() => ({
  name: [{ required: true, message: t('rules.nameRequired'), trigger: 'blur' }],
  points: [{ required: true, message: t('rules.pointsRequired'), trigger: 'blur' }]
}))

async function loadRules() {
  loading.value = true
  try {
    const res = await getRules()
    rules.value = res.data ?? []
  } catch {
    ElMessage.error(t('rules.loadFailed'))
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
      ElMessage.success(t('rules.updateSuccess'))
    } else {
      await createRule({ ...form })
      ElMessage.success(t('rules.addSuccess'))
    }
    dialogVisible.value = false
    await loadRules()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || t('rules.saveFailed'))
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm(
    t('rules.deleteConfirm', { name: row.name }),
    t('rules.deleteTitle'),
    { type: 'warning', confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel') }
  )
  try {
    await deleteRule(row.id)
    ElMessage.success(t('rules.deleteSuccess'))
    await loadRules()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || t('rules.deleteFailed'))
  }
}

onMounted(loadRules)
</script>
