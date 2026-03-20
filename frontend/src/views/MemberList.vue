<template>
  <div>
    <div class="page-header">
      <h2>👥 家庭成员</h2>
      <el-button type="primary" @click="$router.push('/members/create')">
        <el-icon><Plus /></el-icon> 添加成员
      </el-button>
    </div>

    <el-card>
      <el-table :data="members" stripe>
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="row.role === '父' || row.role === '母' ? 'primary' : 'success'">{{ row.role }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <ProtectionBadge :is-protected="row.is_protected" />
          </template>
        </el-table-column>
        <el-table-column prop="protected_until" label="保护期至" />
        <el-table-column prop="created_at" label="创建时间" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button link type="primary" @click="$router.push(`/members/${row.id}/edit`)">编辑</el-button>
            <el-button link type="warning" @click="showProtection(row)">设置保护期</el-button>
            <el-button link type="danger" @click="deleteMember(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Protection Dialog -->
    <el-dialog v-model="protectionVisible" title="设置保护期" width="400px">
      <el-form>
        <el-form-item label="保护时长（小时）">
          <el-input-number v-model="protectionHours" :min="1" :max="8760" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="protectionVisible = false">取消</el-button>
        <el-button type="primary" @click="setProtection">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../utils/api'
import ProtectionBadge from '../components/ProtectionBadge.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const members = ref([])
const protectionVisible = ref(false)
const protectionHours = ref(720)
const selectedMember = ref(null)

const loadMembers = async () => {
  try { members.value = await api.getAllMembers() }
  catch (e) { ElMessage.error('加载失败: ' + e) }
}

const deleteMember = async (row) => {
  try {
    await ElMessageBox.confirm(`确认删除成员「${row.name}」？`, '提示', { type: 'warning' })
    await api.deleteMember(row.id)
    ElMessage.success('删除成功')
    loadMembers()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败: ' + e)
  }
}

const showProtection = (row) => {
  selectedMember.value = row
  protectionVisible.value = true
}

const setProtection = async () => {
  try {
    await api.setProtection(selectedMember.value.id, protectionHours.value)
    ElMessage.success('保护期设置成功')
    protectionVisible.value = false
    loadMembers()
  } catch (e) {
    ElMessage.error('设置失败: ' + e)
  }
}

onMounted(loadMembers)
</script>
