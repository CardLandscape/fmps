<template>
  <div v-loading="loading">
    <!-- 案件基本信息 -->
    <el-card shadow="never" style="margin-bottom:16px">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:600">{{ i18n.t('caseDetailTitle') }}</span>
          <div>
            <el-button @click="$router.back()" size="small">{{ i18n.t('btnBack') }}</el-button>
            <!-- Start execution: only when status=pending and all prep items checked -->
            <el-button
              v-if="caseData.status === 'pending'"
              type="success"
              size="small"
              :loading="actionLoading"
              :disabled="!allPrepChecked"
              @click="handleStart"
            >{{ allPrepChecked ? i18n.t('btnStartExec') : i18n.t('confirmPrepFirst') }}</el-button>
            <!-- Force complete (legacy / skip) -->
            <el-button
              v-if="caseData.status === 'active'"
              type="warning"
              size="small"
              :loading="actionLoading"
              @click="handleComplete"
            >{{ i18n.t('btnEndPunishment') }}</el-button>
          </div>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="i18n.t('labelCaseTitle')">{{ caseData.title }}</el-descriptions-item>
        <el-descriptions-item :label="i18n.t('labelCaseStatusHeader')">
          <el-tag :type="statusTagType(caseData.status)">{{ statusLabel(caseData.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="i18n.t('labelCaseParentHeader')">{{ getMemberName(caseData.parent_member) }}</el-descriptions-item>
        <el-descriptions-item :label="i18n.t('labelCaseChildHeader')">{{ getMemberName(caseData.child_member) }}</el-descriptions-item>
        <el-descriptions-item :label="i18n.t('labelPunishLevelHeader')">
          <el-tag v-if="caseData.punishment_level" type="danger">{{ caseData.punishment_level }}级</el-tag>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item :label="i18n.t('labelStartTime')">{{ caseData.start_time ? formatDate(caseData.start_time) : i18n.t('notStarted') }}</el-descriptions-item>
        <el-descriptions-item :label="i18n.t('labelDescription')" :span="2">{{ caseData.description || '-' }}</el-descriptions-item>
        <el-descriptions-item v-if="caseData.status === 'completed'" :label="i18n.t('labelFinalGrade')" :span="2">
          <el-tag :type="gradeTagType(caseData.final_grade)" size="large" style="font-size:16px">{{ caseData.final_grade || '-' }}</el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- Phase 1: Preparation (status=pending) -->
    <el-card v-if="caseData.status === 'pending' && prepItems.length > 0" shadow="never" style="margin-bottom:16px">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:600">{{ i18n.t('prepPhaseTitle') }}</span>
          <span style="font-size:13px;color:#606266">{{ i18n.t('prepProgress').replace('{done}', checkedCount).replace('{total}', prepItems.length) }}</span>
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
        {{ i18n.t('prepConfirmAllHint') }}
      </div>
    </el-card>

    <!-- Phase 2: Execution (status=active, has parsed steps) -->
    <el-card v-if="caseData.status === 'active' && parsedSteps.length > 0" shadow="never" style="margin-bottom:16px">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center;flex-wrap:wrap;gap:8px">
          <span style="font-weight:600">{{ i18n.t('execPhaseTitle') }}</span>
          <div style="display:flex;align-items:center;gap:16px;flex-wrap:wrap">
            <!-- Real-time clock -->
            <span style="font-size:13px;color:#409eff;font-family:monospace;background:#ecf5ff;padding:2px 8px;border-radius:4px">
              🕐 {{ formatDateTime(now) }}
            </span>
            <span style="color:#f56c6c;font-weight:600;font-size:16px">
              {{ i18n.t('totalDeducted').replace('{n}', totalDeducted) }}
              <el-tag :type="gradeTagType(currentGrade)" size="small" style="margin-left:8px">{{ i18n.t('currentGrade') }}{{ currentGrade }}</el-tag>
            </span>
          </div>
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
            {{ i18n.t('stepLabel').replace('{n}', (caseData.current_step_index ?? 0) + 1) }}
          </el-tag>
          <div style="flex:1">
            <div style="font-size:15px;font-weight:500;margin-bottom:12px;line-height:1.6">
              {{ parsedSteps[caseData.current_step_index ?? 0] }}
            </div>

            <!-- Posture info panel: shown when current step references a known posture -->
            <div
              v-if="currentStepPostures.length > 0"
              style="margin-bottom:12px"
            >
              <div
                v-for="posture in currentStepPostures"
                :key="posture.name"
                style="background:#fffbe6;border:1px solid #ffe58f;border-radius:6px;padding:10px 14px;margin-bottom:8px"
              >
                <div style="font-weight:600;color:#d48806;margin-bottom:6px;font-size:14px">
                  📐 {{ posture.name }}
                </div>
                <div v-if="posture.requirements" style="font-size:13px;color:#595959;margin-bottom:6px;line-height:1.6">
                  <span style="font-weight:600">{{ i18n.t('postureRequirements') }}：</span>{{ posture.requirements }}
                </div>
                <div v-if="posture.deduct_rule" style="display:flex;align-items:center;gap:8px;flex-wrap:wrap">
                  <span style="font-size:13px;color:#cf1322">
                    <span style="font-weight:600">{{ i18n.t('postureDeductRule') }}：</span>{{ posture.deduct_rule }}
                  </span>
                  <el-button
                    v-if="posture.deduct_points > 0"
                    type="danger"
                    size="small"
                    @click="handlePostureDeduct(posture)"
                  >
                    -{{ posture.deduct_points }} {{ i18n.t('postureQuickDeduct') }}
                  </el-button>
                </div>
              </div>
            </div>

            <div style="display:flex;gap:8px;flex-wrap:wrap">
              <!-- Deduct button: quick amounts -->
              <el-button type="danger" size="small" @click="openDeductDialog()">
                <el-icon><Minus /></el-icon>
                {{ i18n.t('btnDeduct') }}
              </el-button>
              <el-button
                type="success"
                size="small"
                :loading="actionLoading"
                @click="handleCompleteStep"
              >
                <el-icon><Check /></el-icon>
                {{ isLastStep ? i18n.t('btnCompleteAll') : i18n.t('btnCompleteStep') }}
              </el-button>
            </div>
          </div>
        </div>
      </el-card>

      <!-- Completed steps summary -->
      <div v-if="caseData.current_step_index > 0" style="margin-top:8px">
        <div style="font-size:12px;color:#909399;margin-bottom:6px">{{ i18n.t('completedSteps') }}</div>
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
          <span style="font-weight:600">{{ i18n.t('legacyPhaseTitle') }}</span>
          <span style="color:#f56c6c;font-weight:600">
            {{ i18n.t('totalDeducted').replace('{n}', totalDeducted) }}
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
              >{{ i18n.t('btnDeduct') }}</el-button>
            </div>
          </el-card>
        </el-timeline-item>
      </el-timeline>

      <el-empty v-if="steps.length === 0" description="暂无惩罚步骤" />
    </el-card>

    <!-- Phase 3: Completed -->
    <el-card v-if="caseData.status === 'completed'" shadow="never" style="margin-bottom:16px">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:600">{{ i18n.t('archivePhaseTitle') }}</span>
          <el-button type="primary" size="small" @click="openReportDialog">
            📄 {{ i18n.t('btnGenerateReport') }}
          </el-button>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="i18n.t('labelPunishLevelArchive')">{{ caseData.punishment_level ? caseData.punishment_level + '级' : '-' }}</el-descriptions-item>
        <el-descriptions-item :label="i18n.t('labelTotalDeducted')">{{ totalDeducted }}</el-descriptions-item>
        <el-descriptions-item :label="i18n.t('labelFinalGradeArchive')" :span="2">
          <el-tag :type="gradeTagType(caseData.final_grade)" size="large" style="font-size:16px;font-weight:bold">
            {{ caseData.final_grade || '-' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="i18n.t('labelGradeDesc')" :span="2">{{ gradeDescription(caseData.final_grade) }}</el-descriptions-item>
        <el-descriptions-item v-if="caseData.start_time" :label="i18n.t('execStartedAt')">{{ formatDate(caseData.start_time) }}</el-descriptions-item>
        <el-descriptions-item v-if="caseData.end_time" :label="i18n.t('execEndedAt')">{{ formatDate(caseData.end_time) }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- Deduction records -->
    <el-card v-if="caseData.status === 'active' || caseData.status === 'completed'" shadow="never">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:600">{{ i18n.t('penaltiesTitle') }}</span>
          <span style="color:#f56c6c;font-weight:600">{{ i18n.t('totalDeductedSummary').replace('{n}', totalDeducted) }}</span>
        </div>
      </template>
      <el-table :data="penalties" stripe style="width:100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="rule_text" :label="i18n.t('colPenaltyRule')" />
        <el-table-column prop="score_delta" :label="i18n.t('colPenaltyScore')" width="80">
          <template #default="{ row }">
            <span :style="{ color: row.revoked ? '#909399' : '#f56c6c' }">
              {{ row.score_delta }}
              <span v-if="row.revoked" style="font-size:12px">{{ i18n.t('penaltyRevoked') }}</span>
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="reason" :label="i18n.t('colPenaltyReason')" />
        <el-table-column prop="created_at" :label="i18n.t('colPenaltyTime')" width="160">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column :label="i18n.t('colActions')" width="100" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="!row.revoked && caseData.status === 'active'"
              size="small"
              type="warning"
              plain
              @click="openRevoke(row)"
            >{{ i18n.t('btnRevoke') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="penalties.length === 0" :description="i18n.t('noPenalties')" />
    </el-card>

    <!-- Deduct dialog -->
    <el-dialog v-model="deductVisible" :title="i18n.t('deductDialogTitle')" width="440px">
      <el-form :model="deductForm" label-width="90px">
        <el-form-item :label="i18n.t('labelDeductReason')">
          <el-input v-model="deductForm.ruleText" :placeholder="i18n.t('placeholderDeductReason')" />
        </el-form-item>
        <el-form-item :label="i18n.t('labelDeductPoints')">
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
        <el-button @click="deductVisible = false">{{ i18n.t('btnCancel') }}</el-button>
        <el-button type="danger" :loading="deductLoading" @click="handleDeduct">{{ i18n.t('btnConfirmDeduct') }}</el-button>
      </template>
    </el-dialog>

    <!-- Revoke dialog -->
    <el-dialog v-model="revokeVisible" :title="i18n.t('revokeDialogTitle')" width="400px">
      <el-form :model="revokeForm" label-width="90px">
        <el-form-item :label="i18n.t('labelRevokePassword')">
          <el-input v-model="revokeForm.password" type="password" show-password :placeholder="i18n.t('placeholderRevokePassword')" />
        </el-form-item>
        <el-form-item :label="i18n.t('labelRevokeReason')">
          <el-input v-model="revokeForm.reason" :placeholder="i18n.t('placeholderRevokeReason')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="revokeVisible = false">{{ i18n.t('btnCancel') }}</el-button>
        <el-button type="primary" :loading="revokeLoading" @click="handleRevoke">{{ i18n.t('btnConfirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- Score Report Dialog -->
    <el-dialog v-model="reportVisible" :title="i18n.t('reportTitle')" width="760px" class="report-dialog">
      <div style="margin-bottom:16px;display:flex;align-items:center;gap:12px;flex-wrap:wrap">
        <span style="font-weight:600">{{ i18n.t('reportLang') }}：</span>
        <el-radio-group v-model="reportLang" size="small">
          <el-radio-button value="zh">{{ i18n.t('reportLangZh') }}</el-radio-button>
          <el-radio-button value="en">{{ i18n.t('reportLangEn') }}</el-radio-button>
          <el-radio-button value="bilingual">{{ i18n.t('reportLangBilingual') }}</el-radio-button>
        </el-radio-group>
      </div>

      <!-- Report content (printable) -->
      <div id="score-report-print" class="score-report">
        <div class="report-header">
          <div class="report-sys-title">{{ rLabel('家庭惩戒管理系统', 'Family Management & Penalty System') }}</div>
          <div class="report-main-title">{{ rLabel('惩罚成绩单', 'Punishment Report Card') }}</div>
        </div>

        <table class="report-table" style="width:100%;border-collapse:collapse;margin-bottom:12px">
          <tr>
            <td class="report-label">{{ rLabel('案件标题', 'Case Title') }}</td>
            <td colspan="3" class="report-value">{{ caseData.title }}</td>
          </tr>
          <tr>
            <td class="report-label">{{ rLabel('家长', 'Parent') }}</td>
            <td class="report-value">{{ getMemberName(caseData.parent_member) }}</td>
            <td class="report-label">{{ rLabel('小孩', 'Child') }}</td>
            <td class="report-value">{{ getMemberName(caseData.child_member) }}</td>
          </tr>
          <tr>
            <td class="report-label">{{ rLabel('惩罚级别', 'Punishment Level') }}</td>
            <td class="report-value">{{ caseData.punishment_level ? caseData.punishment_level + '级' : '-' }}</td>
            <td class="report-label">{{ rLabel('案件编号', 'Case ID') }}</td>
            <td class="report-value">#{{ caseData.id }}</td>
          </tr>
          <tr>
            <td class="report-label">{{ rLabel('开始时间', 'Start Time') }}</td>
            <td class="report-value">{{ caseData.start_time ? formatDate(caseData.start_time) : '-' }}</td>
            <td class="report-label">{{ rLabel('结束时间', 'End Time') }}</td>
            <td class="report-value">{{ caseData.end_time ? formatDate(caseData.end_time) : '-' }}</td>
          </tr>
        </table>

        <!-- Steps executed -->
        <div v-if="parsedSteps.length > 0" style="margin-bottom:12px">
          <div class="report-section-title">{{ rLabel('执行步骤', 'Steps Executed') }}</div>
          <ol style="margin:4px 0;padding-left:20px;font-size:13px">
            <li v-for="(step, idx) in parsedSteps" :key="idx" style="margin-bottom:3px">{{ step }}</li>
          </ol>
        </div>

        <!-- Deduction log -->
        <div style="margin-bottom:12px">
          <div class="report-section-title">{{ rLabel('扣分记录', 'Deduction Log') }}</div>
          <div v-if="activePenalties.length === 0" style="font-size:13px;color:#909399;padding:4px 0">
            {{ rLabel('本次惩罚无扣分', 'No points were deducted.') }}
          </div>
          <table v-else style="width:100%;border-collapse:collapse;font-size:13px">
            <thead>
              <tr style="background:#f5f7fa">
                <th class="report-th">#</th>
                <th class="report-th">{{ rLabel('规则', 'Rule') }}</th>
                <th class="report-th">{{ rLabel('分值', 'Points') }}</th>
                <th class="report-th">{{ rLabel('备注', 'Notes') }}</th>
                <th class="report-th">{{ rLabel('时间', 'Time') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(p, idx) in activePenalties" :key="p.id">
                <td class="report-td">{{ idx + 1 }}</td>
                <td class="report-td">{{ p.rule_text }}</td>
                <td class="report-td" style="color:#f56c6c">{{ p.score_delta }}</td>
                <td class="report-td">{{ p.reason || '-' }}</td>
                <td class="report-td">{{ formatDate(p.created_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Summary -->
        <table class="report-table" style="width:100%;border-collapse:collapse;margin-bottom:16px">
          <tr>
            <td class="report-label" style="width:30%">{{ rLabel('总扣分', 'Total Deducted') }}</td>
            <td class="report-value" style="color:#f56c6c;font-weight:bold;font-size:16px">{{ totalDeducted }} {{ rLabel('分', 'pts') }}</td>
          </tr>
          <tr>
            <td class="report-label">{{ rLabel('最终成绩', 'Final Grade') }}</td>
            <td class="report-value" style="font-weight:bold;font-size:18px">{{ caseData.final_grade || '-' }}</td>
          </tr>
          <tr v-if="caseData.description">
            <td class="report-label">{{ rLabel('备注', 'Notes') }}</td>
            <td class="report-value">{{ caseData.description }}</td>
          </tr>
        </table>

        <!-- Grade description -->
        <div style="margin-bottom:16px;font-size:12px;color:#606266;border:1px solid #eee;padding:8px;border-radius:4px;background:#fafafa">
          {{ gradeDescription(caseData.final_grade) }}
        </div>

        <!-- Signature area -->
        <table style="width:100%;border-collapse:collapse">
          <tr>
            <td style="width:50%;padding:8px 0;font-size:13px">
              {{ rLabel('家长签名', "Parent's Signature") }}：
              <span style="display:inline-block;width:120px;border-bottom:1px solid #333;margin-left:8px">&nbsp;</span>
            </td>
            <td style="width:50%;padding:8px 0;font-size:13px">
              {{ rLabel('日期', 'Date') }}：
              <span style="display:inline-block;width:120px;border-bottom:1px solid #333;margin-left:8px">&nbsp;</span>
            </td>
          </tr>
        </table>

        <div style="text-align:right;font-size:11px;color:#909399;margin-top:12px;border-top:1px solid #eee;padding-top:8px">
          {{ rLabel('成绩单生成时间', 'Report generated at') }}：{{ formatDate(new Date()) }}
        </div>
      </div>

      <template #footer>
        <el-button @click="reportVisible = false">{{ i18n.t('btnCancel') }}</el-button>
        <el-button type="primary" @click="printReport">🖨️ {{ i18n.t('btnPrintReport') }}</el-button>
        <el-button type="success" @click="exportPdf">📄 {{ i18n.t('btnExportPdf') }}</el-button>
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
import { useI18n } from '@/utils/i18n'

const route = useRoute()
const i18n = useI18n()
const caseId = route.params.id

const loading = ref(false)
const actionLoading = ref(false)
const caseData = ref({ status: '', member: null, current_step_index: 0 })
const steps = ref([])        // legacy pipe-format steps
const parsedSteps = ref([])  // new structured steps (array of strings)
const prepItems = ref([])    // preparation items
const prepChecked = ref([])  // checkbox state per prep item
const postures = ref([])     // punishment postures parsed from TXT
const penalties = ref([])
const totalDeducted = ref(0)
const currentGrade = ref('满分')

const revokeVisible = ref(false)
const revokeLoading = ref(false)
const revokeForm = reactive({ password: '', reason: '', penaltyId: null })

const deductVisible = ref(false)
const deductLoading = ref(false)
const deductForm = reactive({ ruleText: '', points: 1 })

// Score report state
const reportVisible = ref(false)
const reportLang = ref('zh')

// Active (non-revoked) penalties for report
const activePenalties = computed(() => penalties.value.filter(p => !p.revoked))

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

// Postures relevant to the current execution step.
// A posture is considered "active" if its name appears as a substring of the current step text.
const currentStepPostures = computed(() => {
  if (!postures.value.length || !parsedSteps.value.length) return []
  const stepIdx = caseData.value.current_step_index ?? 0
  const stepText = parsedSteps.value[stepIdx] ?? ''
  return postures.value.filter(p => p.name && stepText.includes(p.name))
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
  return { pending: i18n.t('statusPending'), active: i18n.t('statusActive'), completed: i18n.t('statusCompleted') }[s] ?? s
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
  const keyMap = {
    '满分': 'gradeDescPerfect',
    '优': 'gradeDescExcellent',
    '良': 'gradeDescGood',
    '达标': 'gradeDescPass',
    '不达标': 'gradeDescFail',
    '态度不端正': 'gradeDescBadAttitude'
  }
  return keyMap[g] ? i18n.t(keyMap[g]) : '-'
}
function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function formatDateTime(d) {
  if (!d) return '-'
  return new Date(d).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', weekday: 'short', hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false })
}

function openReportDialog() {
  reportLang.value = 'zh'
  reportVisible.value = true
}

/**
 * Returns a label string appropriate for the current report language.
 * - 'zh': returns zh only
 * - 'en': returns en only
 * - 'bilingual': returns "zh / en"
 */
function rLabel(zh, en) {
  if (reportLang.value === 'en') return en
  if (reportLang.value === 'bilingual') return `${zh} / ${en}`
  return zh
}

function printReport() {
  const printContent = document.getElementById('score-report-print')
  if (!printContent) return
  const win = window.open('', '_blank', 'width=800,height=700')
  if (!win) return
  win.document.write(`
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="utf-8">
      <title>${i18n.t('reportTitle')}</title>
      <style>
        body { font-family: 'SimSun', 'Arial', sans-serif; padding: 20px; color: #333; font-size: 13px; }
        .report-header { text-align: center; margin-bottom: 20px; border-bottom: 2px solid #333; padding-bottom: 12px; }
        .report-sys-title { font-size: 14px; color: #666; margin-bottom: 6px; }
        .report-main-title { font-size: 22px; font-weight: bold; }
        table { width: 100%; border-collapse: collapse; }
        td, th { border: 1px solid #ccc; padding: 6px 10px; }
        .report-label { background: #f5f7fa; font-weight: 600; width: 25%; }
        .report-value { background: #fff; }
        .report-section-title { font-weight: 600; font-size: 14px; border-left: 3px solid #409eff; padding-left: 8px; margin: 10px 0 6px; }
        .report-th { background: #f5f7fa; font-weight: 600; text-align: center; }
        .report-td { text-align: center; }
        @media print { body { padding: 0; } }
      </style>
    </head>
    <body>
      ${printContent.innerHTML}
    </body>
    </html>
  `)
  win.document.close()
  win.focus()
  // Small delay allows the browser to finish rendering before triggering print
  setTimeout(() => { win.print() }, 300)
}

function exportPdf() {
  // Opens the system print dialog; choose "Save as PDF" as the printer destination
  printReport()
}

async function loadData() {
  loading.value = true
  try {
    const res = await getCase(caseId)
    caseData.value = res.data.case
    steps.value = res.data.punishment_steps ?? []
    parsedSteps.value = res.data.parsed_steps ?? []
    prepItems.value = res.data.prep_items ?? []
    postures.value = res.data.postures ?? []
    penalties.value = res.data.penalty_points ?? []
    totalDeducted.value = res.data.total_deducted ?? 0
    currentGrade.value = res.data.current_grade ?? '满分'
    // Initialize prep checkboxes (preserve existing state)
    if (prepChecked.value.length !== prepItems.value.length) {
      prepChecked.value = prepItems.value.map(() => false)
    }
  } catch {
    ElMessage.error(i18n.t('loadCaseFailed'))
  } finally {
    loading.value = false
  }
}

async function handleStart() {
  actionLoading.value = true
  try {
    await startPunishment(caseId)
    ElMessage.success(i18n.t('punishStarted'))
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
      ElMessage.success(i18n.t('allStepsCompleted').replace('{grade}', data.final_grade))
    } else {
      ElMessage.success(i18n.t('stepCompleted'))
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
    ElMessage.success(i18n.t('punishEnded').replace('{grade}', grade))
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

/**
 * Quick deduction triggered from posture panel — deducts the posture's default points.
 */
async function handlePostureDeduct(posture) {
  const pts = posture.deduct_points || 1
  try {
    await addPenalty(caseId, {
      rule_text: i18n.t('postureViolation').replace('{name}', posture.name),
      score_delta: -pts,
      reason: posture.deduct_rule || ''
    })
    ElMessage.success(i18n.t('deductSuccess').replace('{n}', pts))
    await loadData()
  } catch (e) {
    ElMessage.error(e.response?.data?.message || i18n.t('saveFailed'))
  }
}

async function handleDeduct() {
  if (!deductForm.points || deductForm.points < 1) {
    ElMessage.warning(i18n.t('deductPointsRequired'))
    return
  }
  deductLoading.value = true
  try {
    await addPenalty(caseId, {
      rule_text: deductForm.ruleText || '违规扣分',
      score_delta: -(deductForm.points),
      reason: ''
    })
    ElMessage.success(i18n.t('deductSuccess').replace('{n}', deductForm.points))
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
    ElMessage.warning(i18n.t('revokePasswordRequired'))
    return
  }
  revokeLoading.value = true
  try {
    await revokePenalty(revokeForm.penaltyId, {
      password: revokeForm.password,
      reason: revokeForm.reason
    })
    ElMessage.success(i18n.t('revokeSuccess'))
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

/* Score report styles */
.score-report {
  background: #fff;
  padding: 8px;
  font-family: 'SimSun', 'Arial', sans-serif;
  font-size: 13px;
  color: #333;
}
.score-report .report-header {
  text-align: center;
  margin-bottom: 16px;
  border-bottom: 2px solid #333;
  padding-bottom: 10px;
}
.score-report .report-sys-title {
  font-size: 13px;
  color: #666;
  margin-bottom: 4px;
}
.score-report .report-main-title {
  font-size: 20px;
  font-weight: bold;
}
.score-report .report-section-title {
  font-weight: 600;
  font-size: 14px;
  border-left: 3px solid #409eff;
  padding-left: 8px;
  margin: 10px 0 6px;
}
.score-report .report-table td,
.score-report table td,
.score-report table th {
  border: 1px solid #ddd;
  padding: 6px 10px;
}
.score-report .report-label {
  background: #f5f7fa;
  font-weight: 600;
  width: 25%;
}
.score-report .report-value {
  background: #fff;
}
.score-report .report-th {
  background: #f5f7fa;
  font-weight: 600;
  text-align: center;
}
.score-report .report-td {
  text-align: center;
}
</style>

