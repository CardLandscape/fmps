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
      <el-table-column prop="name" label="姓名" />
      <el-table-column prop="role" label="角色" width="80">
        <template #default="{ row }">
          <el-tag :type="row.role === 'child' ? 'warning' : 'success'">
            {{ row.role === 'child' ? '孩子' : '成人' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="gender" label="性别" width="70" />
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
    width="680px"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="formRules" label-width="110px">
      <el-divider content-position="left">基本信息</el-divider>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="姓名" prop="name">
            <el-input v-model="form.name" placeholder="请输入姓名" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="角色" prop="role">
            <el-select v-model="form.role" style="width: 100%">
              <el-option label="孩子" value="child" />
              <el-option label="成人" value="adult" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="性别">
            <el-radio-group v-model="form.gender">
              <el-radio value="男">男</el-radio>
              <el-radio value="女">女</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="国籍">
            <el-input v-model="form.nationality" placeholder="如：中国" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="出生日期">
            <el-date-picker v-model="form.birth_date" type="date" value-format="YYYY-MM-DD" style="width:100%" placeholder="选择日期" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-divider content-position="left">主要证件</el-divider>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="证件类型">
            <el-input v-model="form.id_doc_type" placeholder="如：身份证" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="证件号码">
            <el-input v-model="form.id_doc_number" placeholder="证件号码" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="签发日期">
            <el-date-picker v-model="form.id_issue_date" type="date" value-format="YYYY-MM-DD" style="width:100%" placeholder="选择日期" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="有效期">
            <el-date-picker v-model="form.id_expiry_date" type="date" value-format="YYYY-MM-DD" style="width:100%" placeholder="选择日期" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item label="签发机关">
        <el-input v-model="form.id_issue_authority" placeholder="证件签发机关" />
      </el-form-item>

      <el-divider content-position="left">辅助证件</el-divider>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="辅助证件类型">
            <el-input v-model="form.aux_doc_type" placeholder="如：护照" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="辅助证件号码">
            <el-input v-model="form.aux_doc_number" placeholder="辅助证件号码" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-divider content-position="left">学籍信息</el-divider>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="就读学校">
            <el-input v-model="form.school_name" placeholder="学校名称" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="年级">
            <el-input v-model="form.grade" placeholder="如：三年级" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="班级">
            <el-input v-model="form.class_name" placeholder="如：二班" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="班主任姓名">
            <el-input v-model="form.class_teacher_name" placeholder="班主任姓名" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item label="班主任电话">
        <el-input v-model="form.class_teacher_phone" placeholder="班主任联系电话" />
      </el-form-item>

      <el-divider content-position="left">外出权限</el-divider>
      <el-form-item label="外出权限">
        <el-radio-group v-model="form.outing_permission">
          <el-radio value="许可">许可</el-radio>
          <el-radio value="不许可">不许可</el-radio>
          <el-radio value="受限">受限</el-radio>
        </el-radio-group>
      </el-form-item>

      <template v-if="form.outing_permission === '受限'">
        <el-form-item label="允许外出日期">
          <el-input
            v-model="form.outing_dates"
            type="textarea"
            :rows="3"
            placeholder='JSON数组，如：["2024-01-01","2024-01-02"]'
          />
        </el-form-item>
        <el-form-item label="允许外出时段">
          <el-input
            v-model="form.outing_time_ranges"
            type="textarea"
            :rows="3"
            placeholder='JSON数组，如：[{"start":"09:00","end":"12:00"}]'
          />
        </el-form-item>
      </template>
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
import { getMembers, createMember, updateMember, deleteMember } from '@/utils/api'

const loading = ref(false)
const saving = ref(false)
const members = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const defaultForm = () => ({
  name: '', role: 'child', gender: '', nationality: '', birth_date: '',
  id_doc_type: '', id_doc_number: '', id_issue_date: '', id_expiry_date: '',
  id_issue_authority: '', aux_doc_type: '', aux_doc_number: '',
  school_name: '', grade: '', class_name: '',
  class_teacher_name: '', class_teacher_phone: '',
  outing_permission: '', outing_dates: '', outing_time_ranges: ''
})

const form = reactive(defaultForm())

const formRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }]
}

function permissionTagType(p) {
  if (p === '许可') return 'success'
  if (p === '不许可') return 'danger'
  return 'warning'
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

async function loadMembers() {
  loading.value = true
  try {
    const res = await getMembers()
    members.value = res.data ?? []
  } catch {
    ElMessage.error('加载成员失败')
  } finally {
    loading.value = false
  }
}

function openDialog(row = null) {
  Object.assign(form, defaultForm())
  if (row) {
    editingId.value = row.id
    Object.keys(defaultForm()).forEach(k => {
      if (row[k] != null) form[k] = row[k]
    })
  } else {
    editingId.value = null
  }
  dialogVisible.value = true
}

function resetForm() {
  formRef.value?.resetFields()
  editingId.value = null
}

async function handleSave() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (editingId.value) {
      await updateMember(editingId.value, { ...form })
      ElMessage.success('更新成功')
    } else {
      await createMember({ ...form })
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    await loadMembers()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '操作失败')
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm(`确定要删除成员"${row.name}"吗？`, '确认删除', {
    type: 'warning', confirmButtonText: '确定', cancelButtonText: '取消'
  })
  try {
    await deleteMember(row.id)
    ElMessage.success('删除成功')
    await loadMembers()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '删除失败')
  }
}

onMounted(loadMembers)
</script>
