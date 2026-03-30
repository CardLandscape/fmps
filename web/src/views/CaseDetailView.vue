<template>
  <div v-loading="loading">
    <!-- 案件基本信息 -->
    <el-card shadow="never" style="margin-bottom:16px">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:600">案件详情</span>
          <div>
            <el-button @click="$router.back()" size="small">返回</el-button>
            <!-- Start execution: only when status=pending and all prep items checked -->
            <el-button
              v-if="caseData.status === 'pending'"
              type="success"
              size="small"
              :loading="actionLoading"
              :disabled="!allPrepChecked"
              @click="handleStart"
            >{{ allPrepChecked ? '开始执行' : '请先确认准备物品' }}</el-button>
            <!-- Force complete (legacy / skip) -->
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
        <el-descriptions-item label="惩罚级别">
          <el-tag v-if="caseData.punishment_level" type="danger">{{ caseData.punishment_level }}级</el-tag>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="开始时间">{{ caseData.start_time ? formatDate(caseData.start_time) : '未开始' }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ caseData.description || '-' }}</el-descriptions-item>
        <el-descriptions-item v-if="caseData.status === 'completed'" label="最终成绩" :span="2">
          <el-tag :type="gradeTagType(caseData.final_grade)" size="large" style="font-size:16px">{{ caseData.final_grade || '-' }}</el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- Phase 1: Preparation (status=pending) -->
    <el-card v-if="caseData.status === 'pending' && prepItems.length > 0" shadow="never" style="margin-bottom:16px">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:600">准备阶段</span>
          <span style="font-size:13px;color:#606266">已确认 {{ checkedCount }} / {{ prepItems.length }} 项</span>
        </div>
      </template>
      <div style="padding:4px 0">
        <div
          v-for="(item, idx) in prepItems"
          :key="idx"
          style="display:flex;align-items:center;gap:10px;padding:8px 0;border-bottom:1px solid #f0f0f0"
        >
          <el-checkbox v-model="prepChecked[idx]" />
          <span :style="{ textDecoration: prepChecked[idx] ? 'line-through' : 'none', color: prepChecked[idx] ? '#909399' : '#303133' }">
            {{ item }}
          </span>
        </div>
      </div>
      <div style="margin-top:12px;color:#909399;font-size:13px">
        请逐项确认所有准备物品已就位，全部确认后方可开始执行。
      </div>
    </el-card>

    <!-- Phase 2: Execution (status=active, has parsed steps) -->
    <el-card v-if="caseData.status === 'active' && parsedSteps.length > 0" shadow="never" style="margin-bottom:16px">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:600">执行阶段</span>
          <span style="color:#f56c6c;font-weight:600;font-size:16px">
            已扣：{{ totalDeducted }} 分
            <el-tag :type="gradeTagType(currentGrade)" size="small" style="margin-left:8px">当前成绩：{{ currentGrade }}</el-tag>
          </span>
        </div>
      </template>

      <!-- Progress bar -->
      <div style="margin-bottom:16px">
        <el-progress
          :percentage="stepProgress"
          :format="() => `${Math.min((caseData.current_step_index ?? 0) + 1, parsedSteps.length)} / ${parsedSteps.length}`"
          status="striped"
          striped-flow
        />
      </div>

      <!-- Current step (only one shown at a time) -->
      <el-card shadow="hover" style="border:2px solid #409eff;background:#ecf5ff;margin-bottom:12px">
        <div style="display:flex;align-items:flex-start;gap:12px">
          <el-tag type="primary" size="large" style="flex-shrink:0;font-size:16px;padding:0 12px">
            步骤 {{ (caseData.current_step_index ?? 0) + 1 }}
          </el-tag>
          <div style="flex:1">
            <div style="font-size:15px;font-weight:500;margin-bottom:12px;line-height:1.6">
              {{ parsedSteps[caseData.current_step_index ?? 0] }}
            </div>
            <div style="display:flex;gap:8px;flex-wrap:wrap">
              <!-- Deduct button: quick amounts -->
              <el-button type="danger" size="small" @click="openDeductDialog()">
                <el-icon><Minus /></el-icon>
                扣分
              </el-button>
              <el-button
                type="success"
                size="small"
                :loading="actionLoading"
                @click="handleCompleteStep"
              >
                <el-icon><Check /></el-icon>
                {{ isLastStep ? '完成全部步骤' : '完成此步骤' }}
              </el-button>
            </div>
          </div>
        </div>
      </el-card>

      <!-- Completed steps summary -->
      <div v-if="caseData.current_step_index > 0" style="margin-top:8px">
        <div style="font-size:12px;color:#909399;margin-bottom:6px">已完成步骤：</div>
        <div v-for="i in caseData.current_step_index" :key="i" style="font-size:12px;color:#909399;padding:2px 0">
          <el-icon style="color:#67c23a"><CircleCheck /></el-icon>
          {{ parsedSteps[i - 1] }}
        </div>
      </div>
    </el-card>

    <!-- Phase 2 (legacy): Execution with old step format -->
    <el-card v-if="caseData.status === 'active' && parsedSteps.length === 0 && steps.length > 0" shadow="never" style="margin-bottom:16px">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:600">惩罚过程</span>
          <span style="color:#f56c6c;font-weight:600">
            已扣：{{ totalDeducted }} 分
            <el-tag :type="gradeTagType(currentGrade)" size="small" style="margin-left:8px">{{ currentGrade }}</el-tag>
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
                @click="handleDeductLegacy(step)"
              >扣分</el-button>
            </div>
          </el-card>
        </el-timeline-item>
      </el-timeline>

      <el-empty v-if="steps.length === 0" description="暂无惩罚步骤" />
    </el-card>

    <!-- Phase 3: Completed -->
    <el-card v-if="caseData.status === 'completed'" shadow="never" style="margin-bottom:16px">
      <template #header>
        <span style="font-weight:600">惩罚归档</span>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="惩罚级别">{{ caseData.punishment_level ? caseData.punishment_level + '级' : '-' }}</el-descriptions-item>
        <el-descriptions-item label="总扣分">{{ totalDeducted }} 分</el-descriptions-item>
        <el-descriptions-item label="最终成绩" :span="2">
          <el-tag :type="gradeTagType(caseData.final_grade)" size="large" style="font-size:16px;font-weight:bold">
            {{ caseData.final_grade || '-' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="成绩说明" :span="2">{{ gradeDescription(caseData.final_grade) }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- Deduction records -->
    <el-card v-if="caseData.status === 'active' || caseData.status === 'completed'" shadow="never">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:600">扣分记录</span>
          <span style="color:#f56c6c;font-weight:600">合计扣分：{{ totalDeducted }} 分</span>
        </div>
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

    <!-- Deduct dialog -->
    <el-dialog v-model="deductVisible" title="扣分" width="440px">
      <el-form :model="deductForm" label-width="90px">
        <el-form-item label="扣分原因">
          <el-input v-model="deductForm.ruleText" placeholder="请输入扣分原因" />
        </el-form-item>
        <el-form-item label="扣分分值">
          <div style="display:flex;gap:8px;flex-wrap:wrap;margin-bottom:8px">
            <el-button
              v-for="v in quickDeductValues"
              :key="v"
              size="small"
              type="danger"
              plain
              @click="deductForm.points = v"
            >-{{ v }}</el-button>
          </div>
          <el-input-number
            v-model="deductForm.points"
            :min="1"
            :max="999"
            controls-position="right"
            placeholder="或输入分值"
          />
          <span style="margin-left:8px;color:#909399">分</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="deductVisible = false">取 消</el-button>
        <el-button type="danger" :loading="deductLoading" @click="handleDeduct">确认扣分</el-button>
      </template>
    </el-dialog>

    <!-- Revoke dialog -->
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
import { Minus, Check, CircleCheck } from '@element-plus/icons-vue'
import { getCase, startPunishment, completePunishment, completeStep, addPenalty, revokePenalty } from '@/utils/api'

const route = useRoute()
const caseId = route.params.id

const loading = ref(false)
const actionLoading = ref(false)
const caseData = ref({ status: '', member: null, current_step_index: 0 })
const steps = ref([])        // legacy pipe-format steps
const parsedSteps = ref([])  // new structured steps (array of strings)
const prepItems = ref([])    // preparation items
const prepChecked = ref([])  // checkbox state per prep item
const penalties = ref([])
const totalDeducted = ref(0)
const currentGrade = ref('满分')

const revokeVisible = ref(false)
const revokeLoading = ref(false)
const revokeForm = reactive({ password: '', reason: '', penaltyId: null })

const deductVisible = ref(false)
const deductLoading = ref(false)
const deductForm = reactive({ ruleText: '', points: 1 })

// Quick deduct values
const quickDeductValues = [1, 2, 3, 5, 8, 10, 15, 20]

let timer = null
const now = ref(new Date())

// All prep items checked
const checkedCount = computed(() => prepChecked.value.filter(Boolean).length)
const allPrepChecked = computed(() => prepItems.value.length === 0 || (prepItems.value.length > 0 && checkedCount.value === prepItems.value.length))

// Step progress percentage
const stepProgress = computed(() => {
  const total = parsedSteps.value.length
  if (total === 0) return 0
  const done = (caseData.value.current_step_index ?? 0)
  return Math.round((done / total) * 100)
})

const isLastStep = computed(() => {
  const idx = caseData.value.current_step_index ?? 0
  return idx >= parsedSteps.value.length - 1
})

// ---- Legacy step logic ----
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
    if (!s.duration && currentHM >= startMin) {
      active = i
    }
  }
  return active
})

function stepTimelineType(idx) {
  if (caseData.value.status !== 'active') return 'info'
  if (idx === currentStepIndex.value) return 'success'
  const s = steps.value[idx]
  const [h, m] = (s.start_time || '00:00').split(':').map(Number)
  const startMin = h * 60 + m
  const nowMin = now.value.getHours() * 60 + now.value.getMinutes()
  return nowMin > startMin + (s.duration || 0) ? 'info' : 'primary'
}
// ---- End legacy ----

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
function gradeTagType(g) {
  const map = { '满分': 'success', '优': 'success', '良': '', '达标': 'warning', '不达标': 'danger', '态度不端正': 'danger' }
  return map[g] ?? 'info'
}
function gradeDescription(g) {
  const map = {
    '满分': '本次惩罚未发生任何扣分，表现完美。',
    '优': '本次惩罚扣分1-5分，成绩为优。需在第二日重复执行下一级别惩罚。',
    '良': '本次惩罚扣分6-15分，成绩为良。需在第二日重复执行同等级惩罚。',
    '达标': '本次惩罚扣分16-19分，成绩为达标。需再执行三次该惩罚。',
    '不达标': '本次惩罚扣分20分以上，成绩为不达标。需继续执行一周。',
    '态度不端正': '本次惩罚扣分40分以上，成绩为态度不端正。需执行至态度端正为止。'
  }
  return map[g] ?? '-'
}
function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

async function loadData() {
  loading.value = true
  try {
    const res = await getCase(caseId)
    caseData.value = res.data.case
    steps.value = res.data.punishment_steps ?? []
    parsedSteps.value = res.data.parsed_steps ?? []
    prepItems.value = res.data.prep_items ?? []
    penalties.value = res.data.penalty_points ?? []
    totalDeducted.value = res.data.total_deducted ?? 0
    currentGrade.value = res.data.current_grade ?? '满分'
    // Initialize prep checkboxes (preserve existing state)
    if (prepChecked.value.length !== prepItems.value.length) {
      prepChecked.value = prepItems.value.map(() => false)
    }
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

async function handleCompleteStep() {
  actionLoading.value = true
  try {
    const res = await completeStep(caseId)
    const data = res.data
    if (data.finished) {
      ElMessage.success(`所有步骤已完成！最终成绩：${data.final_grade}`)
    } else {
      ElMessage.success('步骤已完成，进入下一步')
    }
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
    const res = await completePunishment(caseId)
    const grade = res.data.final_grade
    ElMessage.success(`惩罚已结束，最终成绩：${grade}`)
    await loadData()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '操作失败')
  } finally {
    actionLoading.value = false
  }
}

function openDeductDialog() {
  deductForm.ruleText = ''
  deductForm.points = 1
  deductVisible.value = true
}

async function handleDeduct() {
  if (!deductForm.points || deductForm.points < 1) {
    ElMessage.warning('请输入扣分分值')
    return
  }
  deductLoading.value = true
  try {
    await addPenalty(caseId, {
      rule_text: deductForm.ruleText || '违规扣分',
      score_delta: -(deductForm.points),
      reason: ''
    })
    ElMessage.success(`已扣除 ${deductForm.points} 分`)
    deductVisible.value = false
    await loadData()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '扣分失败')
  } finally {
    deductLoading.value = false
  }
}

async function handleDeductLegacy(step) {
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
      await loadData()
      revokeVisible.value = false
    }
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

