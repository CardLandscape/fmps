<template>
  <div v-if="caseData">
    <div class="page-header">
      <h2>📋 案件详情</h2>
      <el-button @click="$router.back()">返回</el-button>
    </div>

    <el-row :gutter="20">
      <el-col :span="16">
        <el-card style="margin-bottom: 16px">
          <el-descriptions :title="caseData.case_no" :column="2" border>
            <el-descriptions-item label="违规成员">{{ caseData.member_name }}</el-descriptions-item>
            <el-descriptions-item label="违规条款">{{ caseData.clause_code }} - {{ caseData.clause_title }}</el-descriptions-item>
            <el-descriptions-item label="案件状态">
              <el-tag :type="statusType(caseData.status)">{{ statusLabels[caseData.status] }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ caseData.created_at }}</el-descriptions-item>
            <el-descriptions-item label="事件描述" :span="2">{{ caseData.incident_description }}</el-descriptions-item>
            <el-descriptions-item label="惩罚内容" :span="2">{{ caseData.punishment_detail }}</el-descriptions-item>
          </el-descriptions>

          <!-- Action buttons -->
          <div style="margin-top: 16px; display: flex; gap: 12px">
            <el-button v-if="caseData.status === 'pending'" type="primary" @click="startPunishment">
              ▶ 开始惩罚
            </el-button>
            <el-button v-if="caseData.status === 'in_progress'" type="success" @click="completePunishment">
              ✅ 完成惩罚
            </el-button>
            <el-button v-if="['pending','in_progress'].includes(caseData.status)" type="danger" @click="cancelCase">
              ✖ 取消案件
            </el-button>
          </div>
        </el-card>

        <!-- Countdown -->
        <el-card v-if="caseData.status === 'in_progress' && template?.duration_minutes" style="margin-bottom: 16px">
          <template #header>⏱ 惩罚倒计时</template>
          <CountdownTimer :started-at="caseData.started_at" :duration-minutes="template.duration_minutes" />
        </el-card>

        <!-- Comments -->
        <el-card style="margin-bottom: 16px">
          <template #header>💬 评论记录</template>
          <div v-for="c in comments" :key="c.id" style="margin-bottom: 12px; padding: 10px; background: #f9f9f9; border-radius: 6px">
            <div style="font-weight: bold; margin-bottom: 4px">{{ c.member_name }} <span style="font-size: 12px; color: #909399; font-weight: normal">{{ c.created_at }}</span></div>
            <div>{{ c.content }}</div>
          </div>
          <el-empty v-if="!comments.length" description="暂无评论" :image-size="60" />
          <div style="margin-top: 16px; display: flex; gap: 12px">
            <el-select v-model="newComment.member_id" placeholder="评论人" style="width: 140px">
              <el-option v-for="m in members" :key="m.id" :label="m.name" :value="m.id" />
            </el-select>
            <el-input v-model="newComment.content" placeholder="添加评论..." style="flex: 1" />
            <el-button type="primary" @click="addComment">发送</el-button>
          </div>
        </el-card>
      </el-col>

      <el-col :span="8">
        <!-- Appeals -->
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>⚖️ 申诉记录</span>
              <el-button v-if="['pending','in_progress'].includes(caseData.status)" size="small" @click="showAppealForm = true">
                发起申诉
              </el-button>
            </div>
          </template>
          <div v-for="a in appeals" :key="a.id" style="margin-bottom: 12px; border-bottom: 1px solid #f0f0f0; padding-bottom: 12px">
            <div style="display: flex; justify-content: space-between">
              <span style="font-weight: bold">{{ a.appellant_name }}</span>
              <el-tag :type="appealStatusType(a.status)" size="small">{{ appealStatusLabels[a.status] }}</el-tag>
            </div>
            <div style="margin: 6px 0; color: #606266">{{ a.reason }}</div>
            <div v-if="a.status === 'pending'" style="display: flex; gap: 8px">
              <el-button size="small" type="success" @click="reviewAppeal(a.id, true)">批准</el-button>
              <el-button size="small" type="danger" @click="reviewAppeal(a.id, false)">驳回</el-button>
            </div>
          </div>
          <el-empty v-if="!appeals.length" description="暂无申诉" :image-size="60" />
        </el-card>
      </el-col>
    </el-row>

    <!-- Appeal dialog -->
    <el-dialog v-model="showAppealForm" title="发起申诉" width="400px">
      <el-form label-width="80px">
        <el-form-item label="申诉人">
          <el-select v-model="appealForm.appellant_id" placeholder="选择申诉人">
            <el-option v-for="m in members" :key="m.id" :label="m.name" :value="m.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="申诉理由">
          <el-input v-model="appealForm.reason" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAppealForm = false">取消</el-button>
        <el-button type="primary" @click="submitAppeal">提交申诉</el-button>
      </template>
    </el-dialog>
  </div>
  <el-empty v-else description="加载中..." />
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '../utils/api'
import CountdownTimer from '../components/CountdownTimer.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const caseData = ref(null)
const template = ref(null)
const comments = ref([])
const appeals = ref([])
const members = ref([])
const showAppealForm = ref(false)
const newComment = ref({ member_id: null, content: '' })
const appealForm = ref({ appellant_id: null, reason: '' })

const statusLabels = { pending: '待处理', in_progress: '进行中', completed: '已完成', cancelled: '已取消', appealed: '申诉中' }
const statusTypes = { pending: 'info', in_progress: 'warning', completed: 'success', cancelled: 'danger', appealed: 'warning' }
const appealStatusLabels = { pending: '待审核', approved: '已批准', rejected: '已驳回' }
const appealStatusType = (s) => ({ pending: 'warning', approved: 'success', rejected: 'danger' }[s] || 'info')
const statusType = (s) => statusTypes[s] || 'info'

const load = async () => {
  const id = Number(route.params.id)
  try {
    const [c, co, a, m] = await Promise.all([
      api.getCase(id),
      api.getCaseComments(id),
      api.getAppealsByCase(id),
      api.getAllMembers(),
    ])
    caseData.value = c
    comments.value = co
    appeals.value = a
    members.value = m
    if (c.template_id) {
      template.value = await api.getTemplate(c.template_id)
    }
  } catch (e) { ElMessage.error('加载失败: ' + e) }
}

const startPunishment = async () => {
  try { await api.startPunishment(caseData.value.id); ElMessage.success('惩罚已开始'); load() }
  catch (e) { ElMessage.error('操作失败: ' + e) }
}

const completePunishment = async () => {
  try { await api.completePunishment(caseData.value.id); ElMessage.success('惩罚已完成'); load() }
  catch (e) { ElMessage.error('操作失败: ' + e) }
}

const cancelCase = async () => {
  try {
    await ElMessageBox.confirm('确认取消此案件？', '提示', { type: 'warning' })
    await api.updateCaseStatus(caseData.value.id, 'cancelled')
    ElMessage.success('案件已取消')
    load()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('操作失败: ' + e)
  }
}

const addComment = async () => {
  if (!newComment.value.member_id || !newComment.value.content) { ElMessage.warning('请选择评论人并填写内容'); return }
  try {
    await api.addCaseComment({ case_id: caseData.value.id, ...newComment.value })
    newComment.value.content = ''
    ElMessage.success('评论已添加')
    const co = await api.getCaseComments(caseData.value.id)
    comments.value = co
  } catch (e) { ElMessage.error('添加失败: ' + e) }
}

const submitAppeal = async () => {
  if (!appealForm.value.appellant_id || !appealForm.value.reason) { ElMessage.warning('请填写完整信息'); return }
  try {
    await api.createAppeal({ case_id: caseData.value.id, ...appealForm.value })
    ElMessage.success('申诉已提交')
    showAppealForm.value = false
    load()
  } catch (e) { ElMessage.error('提交失败: ' + e) }
}

const reviewAppeal = async (appealId, approved) => {
  try {
    await api.reviewAppeal(appealId, 0, approved, '')
    ElMessage.success(approved ? '已批准申诉' : '已驳回申诉')
    load()
  } catch (e) { ElMessage.error('操作失败: ' + e) }
}

onMounted(load)
</script>
