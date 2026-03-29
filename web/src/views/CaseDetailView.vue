<template>
  <div v-loading="loading">
    <!-- 案件基本信息 -->
    <el-card shadow="never" style="margin-bottom:16px">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:600">案件详情</span>
          <div>
            <el-button @click="$router.back()" size="small">返回</el-button>
            <el-button
              v-if="caseData.status === 'pending'"
              type="success"
              size="small"
              :loading="actionLoading"
              @click="handleStart"
            >开始惩罚</el-button>
            <el-button
              v-if="caseData.status === 'active'"
              type="warning"
              size="small"
              :loading="actionLoading"
              @click="handleComplete"
            >结束惩罚</el-button>
          </div>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="标题">{{ caseData.title }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusTagType(caseData.status)">{{ statusLabel(caseData.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="家长">{{ getMemberName(caseData.parent_member) }}</el-descriptions-item>
        <el-descriptions-item label="小孩">{{ getMemberName(caseData.child_member) }}</el-descriptions-item>
        <el-descriptions-item label="开始时间">{{ caseData.start_time ? formatDate(caseData.start_time) : '未开始' }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ caseData.description || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 惩罚过程面板（仅 active/completed 时显示） -->
    <el-card v-if="caseData.status === 'active' || caseData.status === 'completed'" shadow="never" style="margin-bottom:16px">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:600">惩罚过程</span>
          <span v-if="currentStepIndex >= 0 && caseData.status === 'active'" style="color:#67c23a;font-size:14px">
            当前步骤：{{ steps[currentStepIndex]?.punishment_details }}
            <template v-if="countdown > 0">（剩余 {{ countdownDisplay }}）</template>
          </span>
        </div>
      </template>

      <el-timeline>
        <el-timeline-item
          v-for="(step, idx) in steps"
          :key="idx"
          :type="stepTimelineType(idx)"
          :timestamp="step.start_time + (step.duration ? '  共' + step.duration + '分钟' : '')"
          placement="top"
        >
          <el-card
            :class="['step-card', idx === currentStepIndex && caseData.status === 'active' ? 'step-active' : '']"
            shadow="hover"
          >
            <div style="font-weight:600;margin-bottom:8px">{{ step.punishment_details }}</div>
            <div v-if="step.requirements && step.requirements.length" style="margin-bottom:8px">
              <el-tag
                v-for="(req, ri) in step.requirements"
                :key="ri"
                size="small"
                style="margin-right:6px;margin-bottom:4px"
              >{{ req }}</el-tag>
            </div>
            <div v-if="step.deduct_score_rule" style="display:flex;align-items:center;gap:12px;margin-top:8px">
              <span style="color:#e6a23c;font-size:13px">{{ step.deduct_score_rule }}（每次 -{{ step.deduct_score }} 分）</span>
              <el-button
                v-if="caseData.status === 'active'"
                type="danger"
                size="small"
                @click="handleDeduct(step)"
              >扣分</el-button>
            </div>
          </el-card>
        </el-timeline-item>
      </el-timeline>

      <el-empty v-if="steps.length === 0" description="暂无惩罚步骤" />
    </el-card>

    <!-- 扣分记录 -->
    <el-card v-if="caseData.status === 'active' || caseData.status === 'completed'" shadow="never">
      <template #header>
        <span style="font-weight:600">扣分记录（合计：{{ totalScore }} 分）</span>
      </template>
      <el-table :data="penalties" stripe style="width:100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="rule_text" label="扣分规则" />
        <el-table-column prop="score_delta" label="分值" width="80">
          <template #default="{ row }">
            <span :style="{ color: row.revoked ? '#909399' : '#f56c6c' }">
              {{ row.score_delta }}
              <span v-if="row.revoked" style="font-size:12px">（已撤回）</span>
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="备注" />
        <el-table-column prop="created_at" label="时间" width="160">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="!row.revoked && caseData.status === 'active'"
              size="small"
              type="warning"
              plain
              @click="openRevoke(row)"
            >撤回</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="penalties.length === 0" description="暂无扣分记录" />
    </el-card>

    <!-- 撤回对话框 -->
    <el-dialog v-model="revokeVisible" title="撤回扣分" width="400px">
      <el-form :model="revokeForm" label-width="90px">
        <el-form-item label="授权密码">
          <el-input v-model="revokeForm.password" type="password" show-password placeholder="请输入授权密码" />
        </el-form-item>
        <el-form-item label="撤回原因">
          <el-input v-model="revokeForm.reason" placeholder="撤回原因（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="revokeVisible = false">取 消</el-button>
        <el-button type="primary" :loading="revokeLoading" @click="handleRevoke">确 定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCase, startPunishment, completePunishment, addPenalty, revokePenalty } from '@/utils/api'

const route = useRoute()
const caseId = route.params.id

const loading = ref(false)
const actionLoading = ref(false)
const caseData = ref({ status: '', member: null })
const steps = ref([])
const penalties = ref([])
const revokeVisible = ref(false)
const revokeLoading = ref(false)
const revokeForm = reactive({ password: '', reason: '', penaltyId: null })

let timer = null
const now = ref(new Date())

// 当前步骤索引（基于当前时间与步骤的 start_time）
const currentStepIndex = computed(() => {
  if (!caseData.value.start_time || caseData.value.status !== 'active') return -1
  const currentHM = now.value.getHours() * 60 + now.value.getMinutes()
  let active = -1
  for (let i = 0; i < steps.value.length; i++) {
    const s = steps.value[i]
    const [h, m] = (s.start_time || '00:00').split(':').map(Number)
    const startMin = h * 60 + m
    const endMin = startMin + (s.duration || 0)
    if (currentHM >= startMin && currentHM < endMin) {
      active = i
      break
    }
    // 若无 duration，仅判断是否已过开始时间
    if (!s.duration && currentHM >= startMin) {
      active = i
    }
  }
  return active
})

// 当前步骤剩余秒数
const countdown = computed(() => {
  const idx = currentStepIndex.value
  if (idx < 0) return 0
  const s = steps.value[idx]
  if (!s || !s.duration) return 0
  const [h, m] = (s.start_time || '00:00').split(':').map(Number)
  const startSec = (h * 60 + m) * 60
  const endSec = startSec + s.duration * 60
  const currentSec = now.value.getHours() * 3600 + now.value.getMinutes() * 60 + now.value.getSeconds()
  return Math.max(0, endSec - currentSec)
})

const countdownDisplay = computed(() => {
  const s = countdown.value
  const m = Math.floor(s / 60)
  const sec = s % 60
  return `${m}分${sec.toString().padStart(2, '0')}秒`
})

const totalScore = computed(() =>
  penalties.value.filter(p => !p.revoked).reduce((sum, p) => sum + p.score_delta, 0)
)

function statusLabel(s) {
  return { pending: '待执行', active: '执行中', completed: '已完成' }[s] ?? s
}

function getMemberName(m) {
  if (!m) return '-'
  return m.name_cn || m.name_en || m.name || '-'
}
function statusTagType(s) {
  return { pending: 'info', active: 'success', completed: '' }[s] ?? 'info'
}
function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function stepTimelineType(idx) {
  if (caseData.value.status !== 'active') return 'info'
  if (idx === currentStepIndex.value) return 'success'
  const s = steps.value[idx]
  const [h, m] = (s.start_time || '00:00').split(':').map(Number)
  const startMin = h * 60 + m
  const nowMin = now.value.getHours() * 60 + now.value.getMinutes()
  return nowMin > startMin + (s.duration || 0) ? 'info' : 'primary'
}

async function loadData() {
  loading.value = true
  try {
    const res = await getCase(caseId)
    caseData.value = res.data.case
    steps.value = res.data.punishment_steps ?? []
    penalties.value = res.data.penalty_points ?? []
  } catch {
    ElMessage.error('加载案件失败')
  } finally {
    loading.value = false
  }
}

async function handleStart() {
  actionLoading.value = true
  try {
    await startPunishment(caseId)
    ElMessage.success('惩罚已开始')
    await loadData()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '操作失败')
  } finally {
    actionLoading.value = false
  }
}

async function handleComplete() {
  actionLoading.value = true
  try {
    await completePunishment(caseId)
    ElMessage.success('惩罚已结束')
    await loadData()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '操作失败')
  } finally {
    actionLoading.value = false
  }
}

async function handleDeduct(step) {
  try {
    await addPenalty(caseId, {
      rule_text: step.deduct_score_rule,
      score_delta: -(step.deduct_score || 1),
      reason: step.punishment_details
    })
    ElMessage.success(`已扣除 ${step.deduct_score} 分`)
    await loadData()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '扣分失败')
  }
}

function openRevoke(row) {
  revokeForm.password = ''
  revokeForm.reason = ''
  revokeForm.penaltyId = row.id
  revokeVisible.value = true
}

async function handleRevoke() {
  if (!revokeForm.password) {
    ElMessage.warning('请输入授权密码')
    return
  }
  revokeLoading.value = true
  try {
    await revokePenalty(revokeForm.penaltyId, {
      password: revokeForm.password,
      reason: revokeForm.reason
    })
    ElMessage.success('撤回成功')
    revokeVisible.value = false
    await loadData()
  } catch (e) {
    const msg = e.response?.data?.message || '撤回失败'
    ElMessage.error(msg)
    if (e.response?.status === 403) {
      // 密码错误，已自动扣分，刷新记录并关闭
      await loadData()
      revokeVisible.value = false
    }
    // 其他错误保持对话框打开，让用户重试
  } finally {
    revokeLoading.value = false
  }
}

onMounted(() => {
  loadData()
  timer = setInterval(() => { now.value = new Date() }, 1000)
})
onUnmounted(() => { if (timer) clearInterval(timer) })
</script>

<style scoped>
.step-card {
  transition: border 0.2s;
}
.step-active {
  border: 2px solid #67c23a;
  background: #f0f9eb;
}
</style>
