<template>
  <div>
    <div class="page-header">
      <h2>📁 案件列表</h2>
      <el-button type="primary" @click="$router.push('/cases/create')">
        <el-icon><Plus /></el-icon> 新建案件
      </el-button>
    </div>
    <el-card>
      <div style="margin-bottom: 16px; display: flex; gap: 12px">
        <el-select v-model="filterStatus" placeholder="筛选状态" clearable @change="load" style="width: 140px">
          <el-option v-for="(label, key) in statusLabels" :key="key" :label="label" :value="key" />
        </el-select>
      </div>
      <el-table :data="filteredCases" stripe>
        <el-table-column prop="case_no" label="案件编号" width="200" />
        <el-table-column prop="member_name" label="成员" width="90" />
        <el-table-column prop="clause_code" label="条款" width="70" />
        <el-table-column prop="clause_title" label="违规行为" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" size="small">{{ statusLabels[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button link type="primary" @click="$router.push(`/cases/${row.id}`)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { api } from '../utils/api'
import { ElMessage } from 'element-plus'

const cases = ref([])
const filterStatus = ref('')
const statusLabels = { pending: '待处理', in_progress: '进行中', completed: '已完成', cancelled: '已取消', appealed: '申诉中' }
const statusTypes = { pending: 'info', in_progress: 'warning', completed: 'success', cancelled: 'danger', appealed: 'warning' }
const statusType = (s) => statusTypes[s] || 'info'

const filteredCases = computed(() => {
  if (!filterStatus.value) return cases.value
  return cases.value.filter(c => c.status === filterStatus.value)
})

const load = async () => {
  try { cases.value = await api.getAllCases() }
  catch (e) { ElMessage.error('加载失败: ' + e) }
}
onMounted(load)
</script>
