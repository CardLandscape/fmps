<template>
  <div>
    <div class="page-header">
      <h2>📝 新建案件</h2>
    </div>
    <el-card>
      <StepBar :active="step" :steps="['选择成员', '选择条款', '选择惩罚', '确认提交']" />

      <!-- Step 0: Select member -->
      <div v-if="step === 0">
        <h3 style="margin-bottom: 16px">选择违规成员</h3>
        <el-radio-group v-model="form.member_id" style="display: flex; flex-direction: column; gap: 12px">
          <el-radio v-for="m in members" :key="m.id" :label="m.id" border :disabled="m.is_protected">
            <span>{{ m.name }} ({{ m.role }})</span>
            <ProtectionBadge :is-protected="m.is_protected" style="margin-left: 8px" />
          </el-radio>
        </el-radio-group>
      </div>

      <!-- Step 1: Select clause -->
      <div v-if="step === 1">
        <h3 style="margin-bottom: 16px">选择违反条款</h3>
        <el-radio-group v-model="form.clause_id" style="display: flex; flex-direction: column; gap: 12px">
          <el-radio v-for="c in clauses" :key="c.id" :label="c.id" border>
            <span>{{ c.code }} - {{ c.title }}</span>
            <el-tag size="small" style="margin-left: 8px">{{ c.category }}</el-tag>
            <el-rate :model-value="c.severity" :max="5" disabled style="display: inline-flex; margin-left: 8px" />
          </el-radio>
        </el-radio-group>
      </div>

      <!-- Step 2: Select template/custom -->
      <div v-if="step === 2">
        <h3 style="margin-bottom: 16px">选择惩罚方式</h3>
        <el-tabs v-model="punishTab">
          <el-tab-pane label="套用模板" name="template">
            <el-radio-group v-model="form.template_id" style="display: flex; flex-direction: column; gap: 12px; margin-top: 12px">
              <el-radio v-for="t in templates" :key="t.id" :label="t.id" border>
                <span>{{ t.name }}</span>
                <el-tag size="small" style="margin-left: 8px">{{ t.punishment_type }}</el-tag>
                <span v-if="t.duration_minutes" style="margin-left: 8px; color: #909399; font-size: 12px">{{ t.duration_minutes }}分钟</span>
              </el-radio>
            </el-radio-group>
          </el-tab-pane>
          <el-tab-pane label="自定义惩罚" name="custom">
            <el-form style="margin-top: 12px">
              <el-form-item label="惩罚内容">
                <el-input v-model="form.punishment_detail" type="textarea" :rows="4" placeholder="详细描述惩罚内容..." />
              </el-form-item>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- Step 3: Confirm -->
      <div v-if="step === 3">
        <h3 style="margin-bottom: 16px">填写事件描述并确认</h3>
        <el-form label-width="100px">
          <el-form-item label="事件描述">
            <el-input v-model="form.incident_description" type="textarea" :rows="4" placeholder="详细描述违规事件..." />
          </el-form-item>
          <el-form-item label="事件时间">
            <el-date-picker v-model="form.incident_time" type="datetime" placeholder="选择时间" />
          </el-form-item>
          <el-form-item label="创建人">
            <el-select v-model="form.created_by" placeholder="选择创建人">
              <el-option v-for="m in members" :key="m.id" :label="`${m.name} (${m.role})`" :value="m.id" />
            </el-select>
          </el-form-item>
        </el-form>

        <el-descriptions title="案件确认信息" :column="2" border style="margin-top: 16px">
          <el-descriptions-item label="违规成员">{{ selectedMember?.name }}</el-descriptions-item>
          <el-descriptions-item label="违规条款">{{ selectedClause?.code }} - {{ selectedClause?.title }}</el-descriptions-item>
          <el-descriptions-item label="惩罚模板">{{ selectedTemplate?.name || '自定义' }}</el-descriptions-item>
          <el-descriptions-item label="惩罚时长">{{ selectedTemplate?.duration_minutes ? selectedTemplate.duration_minutes + '分钟' : '不限' }}</el-descriptions-item>
        </el-descriptions>
      </div>

      <div style="margin-top: 24px; display: flex; justify-content: space-between">
        <el-button v-if="step > 0" @click="step--">上一步</el-button>
        <span v-else />
        <el-button v-if="step < 3" type="primary" @click="nextStep">下一步</el-button>
        <el-button v-else type="primary" @click="submit" :loading="submitting">提交案件</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../utils/api'
import StepBar from '../components/StepBar.vue'
import ProtectionBadge from '../components/ProtectionBadge.vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const step = ref(0)
const punishTab = ref('template')
const submitting = ref(false)

const members = ref([])
const clauses = ref([])
const templates = ref([])

const form = ref({
  member_id: null,
  clause_id: null,
  template_id: null,
  punishment_detail: '',
  incident_description: '',
  incident_time: null,
  created_by: null,
})

const selectedMember = computed(() => members.value.find(m => m.id === form.value.member_id))
const selectedClause = computed(() => clauses.value.find(c => c.id === form.value.clause_id))
const selectedTemplate = computed(() => templates.value.find(t => t.id === form.value.template_id))

const nextStep = () => {
  if (step.value === 0 && !form.value.member_id) { ElMessage.warning('请选择成员'); return }
  if (step.value === 1 && !form.value.clause_id) { ElMessage.warning('请选择条款'); return }
  step.value++
}

const submit = async () => {
  submitting.value = true
  try {
    const payload = {
      member_id: form.value.member_id,
      clause_id: form.value.clause_id,
      template_id: punishTab.value === 'template' ? form.value.template_id : null,
      punishment_detail: form.value.punishment_detail,
      incident_description: form.value.incident_description,
      incident_time: form.value.incident_time ? new Date(form.value.incident_time).toISOString().replace('T', ' ').substring(0, 19) : null,
      created_by: form.value.created_by,
    }
    const created = await api.createCase(payload)
    ElMessage.success('案件创建成功')
    router.push(`/cases/${created.id}`)
  } catch (e) {
    ElMessage.error('创建失败: ' + e)
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  try {
    const [m, c, t] = await Promise.all([api.getAllMembers(), api.getAllClauses(), api.getAllTemplates()])
    members.value = m
    clauses.value = c.filter(cl => cl.is_active)
    templates.value = t.filter(tp => tp.is_active)
  } catch (e) { ElMessage.error('加载数据失败: ' + e) }
})
</script>
