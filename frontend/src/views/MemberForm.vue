<template>
  <div>
    <div class="page-header">
      <h2>{{ isEdit ? '编辑成员' : '添加成员' }}</h2>
    </div>
    <el-card style="max-width: 500px">
      <el-form :model="form" label-width="80px" @submit.prevent="save">
        <el-form-item label="姓名" required>
          <el-input v-model="form.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="角色" required>
          <el-select v-model="form.role" placeholder="请选择角色">
            <el-option label="父" value="父" />
            <el-option label="母" value="母" />
            <el-option label="子" value="子" />
            <el-option label="女" value="女" />
          </el-select>
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
const form = ref({ name: '', role: '子', avatar: '' })

onMounted(async () => {
  if (isEdit.value) {
    try {
      const m = await api.getMember(Number(route.params.id))
      form.value = { name: m.name, role: m.role, avatar: m.avatar || '', id: m.id }
    } catch (e) { ElMessage.error('加载失败: ' + e) }
  }
})

const save = async () => {
  if (!form.value.name) { ElMessage.warning('请输入姓名'); return }
  try {
    if (isEdit.value) {
      await api.updateMember({ ...form.value, id: Number(route.params.id) })
      ElMessage.success('更新成功')
    } else {
      await api.createMember(form.value)
      ElMessage.success('创建成功')
    }
    router.push('/members')
  } catch (e) { ElMessage.error('保存失败: ' + e) }
}
</script>
