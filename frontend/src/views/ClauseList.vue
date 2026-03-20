<template>
  <div>
    <div class="page-header">
      <h2>📜 惩罚条款</h2>
      <el-button type="primary" @click="$router.push('/clauses/create')">
        <el-icon><Plus /></el-icon> 新增条款
      </el-button>
    </div>
    <el-card>
      <el-table :data="clauses" stripe>
        <el-table-column prop="code" label="编号" width="80" />
        <el-table-column prop="title" label="标题" />
        <el-table-column prop="category" label="分类" width="80">
          <template #default="{ row }">
            <el-tag size="small" effect="plain">{{ row.category }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="severity" label="严重等级" width="100">
          <template #default="{ row }">
            <el-rate :model-value="row.severity" :max="5" disabled />
          </template>
        </el-table-column>
        <el-table-column prop="is_active" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'info'" size="small">{{ row.is_active ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button link type="primary" @click="$router.push(`/clauses/${row.id}/edit`)">编辑</el-button>
            <el-button link type="danger" @click="deleteClause(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../utils/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const clauses = ref([])
const load = async () => {
  try { clauses.value = await api.getAllClauses() }
  catch (e) { ElMessage.error('加载失败: ' + e) }
}
const deleteClause = async (row) => {
  try {
    await ElMessageBox.confirm(`确认删除条款「${row.title}」？`, '提示', { type: 'warning' })
    await api.deleteClause(row.id)
    ElMessage.success('删除成功')
    load()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败: ' + e)
  }
}
onMounted(load)
</script>
