<template>
  <div>
    <el-row :gutter="20" class="stat-row">
      <el-col :span="6" v-for="card in statCards" :key="card.labelKey">
        <el-card shadow="hover" class="stat-card">
          <el-statistic :title="i18n.t(card.labelKey)" :value="card.value">
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
        <span style="font-weight: 600">{{ i18n.t('recentRecords') }}</span>
      </template>
      <el-table
        :data="recentRecords"
        v-loading="loading"
        stripe
        style="width: 100%"
      >
        <el-table-column :label="i18n.t('colMemberName')" prop="member_name" />
        <el-table-column :label="i18n.t('colViolation')" prop="rule_name" />
        <el-table-column :label="i18n.t('colPoints')" prop="points" width="80">
          <template #default="{ row }">
            <el-tag type="danger">{{ row.points }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="i18n.t('colNote')" prop="note" show-overflow-tooltip />
        <el-table-column :label="i18n.t('colOccurredAt')" prop="occurred_at" width="160">
          <template #default="{ row }">
            {{ formatDate(row.occurred_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getStats, getRecords } from '@/utils/api'
import { useI18n } from '@/utils/i18n'

const i18n = useI18n()

const loading = ref(false)
const stats = ref({ member_count: 0, rule_count: 0, record_count: 0, total_points: 0 })
const recentRecords = ref([])

const statCards = ref([
  { labelKey: 'statMembers', value: 0, icon: 'User', color: '#409EFF' },
  { labelKey: 'statRules', value: 0, icon: 'List', color: '#67C23A' },
  { labelKey: 'statRecords', value: 0, icon: 'Document', color: '#E6A23C' },
  { labelKey: 'statPoints', value: 0, icon: 'TrendCharts', color: '#F56C6C' }
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
    statCards.value[0].value = s.member_count ?? 0
    statCards.value[1].value = s.rule_count ?? 0
    statCards.value[2].value = s.record_count ?? 0
    statCards.value[3].value = s.total_points ?? 0
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
