<template>
  <div>
    <el-row :gutter="20" class="stat-row">
      <el-col :span="6" v-for="card in statCards" :key="card.key">
        <el-card shadow="hover" class="stat-card">
          <el-statistic :title="t(card.labelKey)" :value="card.value">
            <template #prefix>
              <el-icon :color="card.color" :size="20">
                <component :is="card.icon" />
              </el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never" style="margin-top: 20px">
      <template #header>
        <span style="font-weight: 600">{{ t('dashboard.recentRecords') }}</span>
      </template>
      <el-table
        :data="recentRecords"
        v-loading="loading"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="member_name" :label="t('dashboard.colMemberName')" />
        <el-table-column prop="rule_name" :label="t('dashboard.colViolation')" />
        <el-table-column prop="points" :label="t('dashboard.colPoints')" width="80">
          <template #default="{ row }">
            <el-tag type="danger">{{ row.points }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="note" :label="t('dashboard.colNote')" show-overflow-tooltip />
        <el-table-column prop="occurred_at" :label="t('dashboard.colOccurredAt')" width="160">
          <template #default="{ row }">
            {{ formatDate(row.occurred_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getStats, getRecords } from '@/utils/api'

const { t } = useI18n()
const loading = ref(false)
const recentRecords = ref([])

const memberCount = ref(0)
const ruleCount = ref(0)
const recordCount = ref(0)
const totalPoints = ref(0)

const statCards = computed(() => [
  { key: 'members', labelKey: 'dashboard.memberCount', value: memberCount.value, icon: 'User', color: '#409EFF' },
  { key: 'rules', labelKey: 'dashboard.ruleCount', value: ruleCount.value, icon: 'List', color: '#67C23A' },
  { key: 'records', labelKey: 'dashboard.recordCount', value: recordCount.value, icon: 'Document', color: '#E6A23C' },
  { key: 'points', labelKey: 'dashboard.totalPoints', value: totalPoints.value, icon: 'TrendCharts', color: '#F56C6C' },
])

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

async function loadData() {
  loading.value = true
  try {
    const [statsRes, recordsRes] = await Promise.all([
      getStats(),
      getRecords({ limit: 10 })
    ])
    const s = statsRes.data
    memberCount.value = s.member_count ?? 0
    ruleCount.value = s.rule_count ?? 0
    recordCount.value = s.record_count ?? 0
    totalPoints.value = s.total_points ?? 0
    recentRecords.value = recordsRes.data?.records ?? recordsRes.data ?? []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.stat-row {
  margin-bottom: 4px;
}

.stat-card {
  border-radius: 8px;
}
</style>
