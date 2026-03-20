<template>
  <div>
    <div class="page-header">
      <h2>{{ isEdit ? '编辑条款' : '新增条款' }}</h2>
    </div>
    <el-card style="max-width: 600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="条款编号" required>
          <el-input v-model="form.code" placeholder="如 C011" />
        </el-form-item>
        <el-form-item label="标题" required>
          <el-input v-model="form.title" placeholder="条款标题" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="form.category">
            <el-option v-for="c in categories" :key="c" :label="c" :value="c" />
          </el-select>
        </el-form-item>
        <el-form-item label="严重等级">
          <el-rate v-model="form.severity" :max="5" />
        </el-form-item>
        <el-form-item v-if="isEdit" label="状态">
          <el-switch v-model="form.is_active" active-text="启用" inactive-text="禁用" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">保存</el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../utils/api'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const categories = ['行为', '态度', '学习', '家务', '其它']
const form = ref({ code: '', title: '', description: '', severity: 1, category: '行为', is_active: true })

onMounted(async () => {
  if (isEdit.value) {
    try {
      const c = await api.getClause(Number(route.params.id))
      form.value = { ...c }
    } catch (e) { ElMessage.error('加载失败: ' + e) }
  }
})

const save = async () => {
  if (!form.value.code || !form.value.title) { ElMessage.warning('请填写必填项'); return }
  try {
    if (isEdit.value) {
      await api.updateClause({ ...form.value, id: Number(route.params.id) })
    } else {
      await api.createClause(form.value)
    }
    ElMessage.success('保存成功')
    router.push('/clauses')
  } catch (e) { ElMessage.error('保存失败: ' + e) }
}
</script>
