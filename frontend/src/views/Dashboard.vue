<template>
  <div>
    <div class="page-header">
      <h2>🏠 仪表盘</h2>
      <el-button type="primary" @click="$router.push('/cases/create')">
        <el-icon><Plus /></el-icon> 新建案件
      </el-button>
    </div>

    <el-row :gutter="20" style="margin-bottom: 24px">
      <el-col :span="6">
        <StatCard label="家庭成员" :value="stats.total_members" icon="User" icon-bg="#fff3eb" icon-color="#e6722e" />
      </el-col>
      <el-col :span="6">
        <StatCard label="活跃条款" :value="stats.active_clauses" icon="Document" icon-bg="#ecf5ff" icon-color="#409eff" />
      </el-col>
      <el-col :span="6">
        <StatCard label="本月案件" :value="stats.month_cases" icon="Folder" icon-bg="#fdf6ec" icon-color="#e6a23c" />
      </el-col>
      <el-col :span="6">
        <StatCard label="待处理案件" :value="stats.pending_cases" icon="Warning" icon-bg="#fef0f0" icon-color="#f56c6c" />
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>📋 最近案件</span>
          </template>
          <el-table :data="recentCases" stripe style="width: 100%">
            <el-table-column prop="case_no" label="案件编号" width="200" />
            <el-table-column prop="member_name" label="成员" width="100" />
            <el-table-column prop="clause_title" label="违规条款" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="statusType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="80">
              <template #default="{ row }">
                <el-button link type="primary" @click="$router.push(`/cases/${row.id}`)">详情</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header><span>📊 案件状态分布</span></template>
          <div v-for="(item, i) in statusChart" :key="i" style="margin-bottom: 12px">
            <div style="display:flex;justify-content:space-between;margin-bottom:4px">
              <span style="font-size:13px">{{ item.label }}</span>
              <span style="font-size:13px;font-weight:bold">{{ item.count }}</span>
            </div>
            <el-progress :percentage="item.pct" :color="item.color" :show-text="false" />
          </div>
          <el-empty v-if="!statusChart.length" description="暂无案件" :image-size="60" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { api } from '../utils/api'
import StatCard from '../components/StatCard.vue'
import { ElMessage } from 'element-plus'

const stats = ref({ total_members: 0, active_clauses: 0, month_cases: 0, pending_cases: 0, status_counts: {} })
const recentCases = ref([])

const statusLabels = { pending: '待处理', in_progress: '进行中', completed: '已完成', cancelled: '已取消', appealed: '申诉中' }
const statusTypes = { pending: 'info', in_progress: 'warning', completed: 'success', cancelled: 'danger', appealed: 'warning' }
const statusColors = { pending: '#909399', in_progress: '#e6a23c', completed: '#67c23a', cancelled: '#f56c6c', appealed: '#e6722e' }

const statusLabel = (s) => statusLabels[s] || s
const statusType = (s) => statusTypes[s] || 'info'

const statusChart = computed(() => {
  const counts = stats.value.status_counts || {}
  const total = Object.values(counts).reduce((a, b) => a + b, 0) || 1
  return Object.entries(counts).map(([s, count]) => ({
    label: statusLabels[s] || s,
    count,
    pct: Math.round((count / total) * 100),
    color: statusColors[s] || '#909399',
  }))
})

onMounted(async () => {
  try {
    const [s, r] = await Promise.all([api.getCaseStats(), api.getRecentCases(5)])
    stats.value = s
    recentCases.value = r
  } catch (e) {
    ElMessage.error('加载数据失败: ' + e)
  }
})
</script>
