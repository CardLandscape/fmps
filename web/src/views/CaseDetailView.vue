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
            >强制结束</el-button>
          </div>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="标题">{{ caseData.title }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusTagType(caseData.status)">{{ statusLabel(caseData.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="惩罚级别">
          <el-tag v-if="caseData.punishment_level" type="warning">{{ caseData.punishment_level }}级</el-tag>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="最终成绩">
          <el-tag v-if="caseData.final_grade" :type="gradeTagType(caseData.final_grade)" size="default">
            {{ caseData.final_grade }}
          </el-tag>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="家长">{{ getMemberName(caseData.parent_member) }}</el-descriptions-item>
        <el-descriptions-item label="小孩">{{ getMemberName(caseData.child_member) }}</el-descriptions-item>
        <el-descriptions-item label="开始时间">{{ caseData.start_time ? formatDate(caseData.start_time) : '未开始' }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ caseData.description || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- ===== 新惩罚流程：准备阶段（active 且 current_step = -1）===== -->
    <el-card
      v-if="caseData.status === 'active' && caseData.current_step === -1 && prepItems.length > 0"
      shadow="never"
      style="margin-bottom:16px"
    >
      <template #header>
        <span style="font-weight:600">准备物品确认</span>
      </template>
      <div style="margin-bottom:12px;color:#606266;font-size:14px">
        请逐一确认以下物品已准备完毕，全部勾选后方可开始执行。
      </div>
      <div v-for="(item, idx) in prepItems" :key="idx" style="margin-bottom:8px">
        <el-checkbox v-model="prepChecked[idx]" :label="item" />
      </div>
      <div style="margin-top:16px">
        <el-button
          type="primary"
          :disabled="!allPrepChecked"
          :loading="actionLoading"
          @click="handleStartExecution"
        >
          {{ allPrepChecked ? '所有物品已准备完毕，开始执行' : `请先勾选所有准备物品（${checkedCount}/${prepItems.length}）` }}
        </el-button>
      </div>
    </el-card>

    <!-- ===== 新惩罚流程：执行阶段（active 且 current_step >= 0，有 exec_steps）===== -->
    <el-card
      v-if="caseData.status === 'active' && caseData.current_step >= 0 && execSteps.length > 0"
      shadow="never"
      style="margin-bottom:16px"
    >
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:600">
            执行步骤（{{ caseData.current_step + 1 }} / {{ execSteps.length }}）
          </span>
          <span style="color:#f56c6c;font-size:15px;font-weight:600">
            已扣分：{{ totalDeducted }} 分
          </span>
        </div>
      </template>

      <!-- Current step card -->
      <el-card shadow="hover" style="border:2px solid #409eff;background:#f0f8ff;margin-bottom:16px">
        <div style="font-size:16px;font-weight:600;margin-bottom:10px;color:#303133">
          步骤 {{ caseData.current_step + 1 }}：{{ execSteps[caseData.current_step] }}
        </div>

        <!-- Quick deduct buttons -->
        <div style="display:flex;flex-wrap:wrap;gap:8px;margin-bottom:12px">
          <el-button
            v-for="pts in [1,2,3,5,8,10,15,20]"
            :key="pts"
            size="small"
            type="danger"
            plain
            @click="quickDeduct(pts)"
          >-{{ pts }} 分</el-button>
          <el-button size="small" type="danger" @click="openCustomDeduct">自定义扣分</el-button>
        </div>

        <div style="display:flex;justify-content:flex-end">
          <el-button
            type="success"
            :loading="actionLoading"
            @click="handleAdvanceStep"
          >
            {{ caseData.current_step + 1 < execSteps.length ? '完成此步骤，进入下一步 →' : '完成最后步骤，归档结果' }}
          </el-button>
        </div>
      </el-card>

      <!-- Upcoming steps (greyed out) -->
      <div v-if="caseData.current_step + 1 < execSteps.length">
        <div style="font-size:13px;color:#909399;margin-bottom:8px">后续步骤（完成当前步骤后逐一解锁）：</div>
        <ol :start="caseData.current_step + 2" style="margin:0;padding-left:20px;color:#c0c4cc;font-size:13px">
          <li
            v-for="(step, i) in execSteps.slice(caseData.current_step + 1)"
            :key="i"
            style="margin-bottom:4px"
          >{{ step }}</li>
        </ol>
      </div>
    </el-card>

    <!-- ===== 旧格式兼容：timeline 展示（active/completed，无 exec_steps 但有 punishment_steps）===== -->
    <el-card
      v-if="(caseData.status === 'active' || caseData.status === 'completed') && execSteps.length === 0 && legacySteps.length > 0"
      shadow="never"
      style="margin-bottom:16px"
    >
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:600">惩罚过程</span>
          <span v-if="currentStepIndex >= 0 && caseData.status === 'active'" style="color:#67c23a;font-size:14px">
            当前步骤：{{ legacySteps[currentStepIndex]?.punishment_details }}
            <template v-if="countdown > 0">（剩余 {{ countdownDisplay }}）</template>
          </span>
        </div>
      </template>

      <el-timeline>
        <el-timeline-item
          v-for="(step, idx) in legacySteps"
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
                @click="handleLegacyDeduct(step)"
              >扣分</el-button>
            </div>
          </el-card>
        </el-timeline-item>
      </el-timeline>

      <el-empty v-if="legacySteps.length === 0" description="暂无惩罚步骤" />
    </el-card>

    <!-- ===== 完成归档展示 ===== -->
    <el-card
      v-if="caseData.status === 'completed'"
      shadow="never"
      style="margin-bottom:16px;border:2px solid #67c23a"
    >
      <template #header>
        <span style="font-weight:600;color:#67c23a">🎉 惩罚已完成 · 归档结果</span>
      </template>
      <el-descriptions :column="3" border>
        <el-descriptions-item label="惩罚级别">
          <el-tag v-if="caseData.punishment_level" type="warning">{{ caseData.punishment_level }}级</el-tag>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="总扣分">
          <span style="color:#f56c6c;font-size:18px;font-weight:700">{{ totalDeducted }} 分</span>
        </el-descriptions-item>
        <el-descriptions-item label="最终成绩">
          <el-tag :type="gradeTagType(caseData.final_grade)" size="large" style="font-size:16px;padding:8px 16px">
            {{ caseData.final_grade || '-' }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>
      <div v-if="caseData.final_grade" style="margin-top:12px;color:#606266;font-size:13px">
        <span v-if="caseData.final_grade === '满分'">🏆 零扣分，完美表现！</span>
        <span v-else-if="caseData.final_grade === '优'">扣分 1–5 分，成绩优秀。</span>
        <span v-else-if="caseData.final_grade === '良'">扣分 6–15 分，成绩良好。</span>
        <span v-else-if="caseData.final_grade === '达标'">扣分 16–19 分，勉强达标。</span>
        <span v-else-if="caseData.final_grade === '不达标'">扣分 20+ 分，成绩不达标。</span>
        <span v-else-if="caseData.final_grade === '态度不端正'">扣分 40+ 分，态度不端正。</span>
      </div>
    </el-card>

    <!-- 扣分记录 -->
    <el-card v-if="caseData.status === 'active' || caseData.status === 'completed'" shadow="never">
      <template #header>
        <span style="font-weight:600">扣分记录（已扣：{{ totalDeducted }} 分）</span>
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

    <!-- 自定义扣分对话框 -->
    <el-dialog v-model="deductVisible" title="自定义扣分" width="380px">
      <el-form :model="deductForm" label-width="80px">
        <el-form-item label="扣分值">
          <el-input-number v-model="deductForm.points" :min="1" :max="9999" placeholder="扣分分值" style="width:100%" />
        </el-form-item>
        <el-form-item label="扣分规则">
          <el-input v-model="deductForm.ruleText" placeholder="扣分规则说明（可选）" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="deductForm.reason" placeholder="备注（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="deductVisible = false">取 消</el-button>
        <el-button type="danger" :loading="deductLoading" @click="handleCustomDeduct">确认扣分</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCase, startPunishment, completePunishment, advanceStep, addPenalty, revokePenalty } from '@/utils/api'

const route = useRoute()
const caseId = route.params.id

const loading = ref(false)
const actionLoading = ref(false)
const caseData = ref({ status: '', current_step: -1, member: null })
const legacySteps = ref([])
const penalties = ref([])
const revokeVisible = ref(false)
const revokeLoading = ref(false)
const revokeForm = reactive({ password: '', reason: '', penaltyId: null })
const deductVisible = ref(false)
const deductLoading = ref(false)
const deductForm = reactive({ points: 1, ruleText: '', reason: '' })

// Preparation items
const prepItems = computed(() => {
  try {
    return caseData.value.prep_items ? JSON.parse(caseData.value.prep_items) : []
  } catch { return [] }
})
const execSteps = computed(() => {
  try {
    return caseData.value.exec_steps ? JSON.parse(caseData.value.exec_steps) : []
  } catch { return [] }
})

// Preparation checkboxes
const prepChecked = ref([])
const checkedCount = computed(() => prepChecked.value.filter(Boolean).length)
const allPrepChecked = computed(() => prepItems.value.length > 0 && checkedCount.value === prepItems.value.length)

// Total deducted (positive number = points lost)
const totalDeducted = computed(() =>
  penalties.value.filter(p => !p.revoked && p.score_delta < 0)
    .reduce((sum, p) => sum + (-p.score_delta), 0)
)

// Legacy countdown support
let timer = null
const now = ref(new Date())

const currentStepIndex = computed(() => {
  if (!caseData.value.start_time || caseData.value.status !== 'active') return -1
  const currentHM = now.value.getHours() * 60 + now.value.getMinutes()
  let active = -1
  for (let i = 0; i < legacySteps.value.length; i++) {
    const s = legacySteps.value[i]
    const [h, m] = (s.start_time || '00:00').split(':').map(Number)
    const startMin = h * 60 + m
    const endMin = startMin + (s.duration || 0)
    if (currentHM >= startMin && currentHM < endMin) { active = i; break }
    if (!s.duration && currentHM >= startMin) active = i
  }
  return active
})

const countdown = computed(() => {
  const idx = currentStepIndex.value
  if (idx < 0) return 0
  const s = legacySteps.value[idx]
  if (!s || !s.duration) return 0
  const [h, m] = (s.start_time || '00:00').split(':').map(Number)
  const startSec = (h * 60 + m) * 60
  const endSec = startSec + s.duration * 60
  const currentSec = now.value.getHours() * 3600 + now.value.getMinutes() * 60 + now.value.getSeconds()
  return Math.max(0, endSec - currentSec)
})
const countdownDisplay = computed(() => {
  const s = countdown.value
  return `${Math.floor(s / 60)}分${(s % 60).toString().padStart(2, '0')}秒`
})

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
  if (g === '满分' || g === '优') return 'success'
  if (g === '良' || g === '达标') return 'warning'
  return 'danger'
}
function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' })
}
function stepTimelineType(idx) {
  if (caseData.value.status !== 'active') return 'info'
  if (idx === currentStepIndex.value) return 'success'
  const s = legacySteps.value[idx]
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
    legacySteps.value = res.data.punishment_steps ?? []
    penalties.value = res.data.penalty_points ?? []
    // Initialize prep checkboxes
    const items = res.data.case.prep_items
    const count = items ? JSON.parse(items).length : 0
    prepChecked.value = Array(count).fill(false)
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

async function handleStartExecution() {
  actionLoading.value = true
  try {
    await advanceStep(caseId)
    ElMessage.success('准备完成，开始执行第一步')
    await loadData()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '操作失败')
  } finally {
    actionLoading.value = false
  }
}

async function handleAdvanceStep() {
  actionLoading.value = true
  try {
    const res = await advanceStep(caseId)
    const updated = res.data
    if (updated.status === 'completed') {
      ElMessage.success(`所有步骤完成！最终成绩：${updated.final_grade}`)
    } else {
      ElMessage.success(`步骤完成，进入第 ${updated.current_step + 1} 步`)
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
    ElMessage.success(`惩罚已结束，最终成绩：${res.data.final_grade}`)
    await loadData()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '操作失败')
  } finally {
    actionLoading.value = false
  }
}

async function quickDeduct(pts) {
  try {
    const stepDesc = execSteps.value[caseData.value.current_step] || ''
    await addPenalty(caseId, {
      rule_text: `步骤${caseData.value.current_step + 1}扣分`,
      score_delta: -pts,
      reason: stepDesc
    })
    ElMessage.success(`已扣除 ${pts} 分`)
    await loadData()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '扣分失败')
  }
}

function openCustomDeduct() {
  deductForm.points = 1
  deductForm.ruleText = ''
  deductForm.reason = ''
  deductVisible.value = true
}

async function handleCustomDeduct() {
  if (!deductForm.points || deductForm.points <= 0) {
    ElMessage.warning('请输入有效的扣分值')
    return
  }
  deductLoading.value = true
  try {
    const stepDesc = execSteps.value[caseData.value.current_step] || ''
    await addPenalty(caseId, {
      rule_text: deductForm.ruleText || `步骤${caseData.value.current_step + 1}扣分`,
      score_delta: -deductForm.points,
      reason: deductForm.reason || stepDesc
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

// Legacy deduct for old-format steps
async function handleLegacyDeduct(step) {
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
