<!-- 用户管理页面（仅管理员） -->
<script setup>
import { ref, onMounted, computed } from 'vue'
import { CreateTeacher, DeleteUser, GetAllTeachers } from '../../wailsjs/go/api/UserAPI'
import { useNotify } from '../composables/useNotify'

const notify = useNotify()
const teachers = ref([])
const showForm = ref(false)
const newUsername = ref('')
const newPassword = ref('')
const searchKeyword = ref('')

const filteredTeachers = computed(() => {
  const kw = searchKeyword.value.trim().toLowerCase()
  if (!kw) return teachers.value
  return teachers.value.filter(t => t.Username.toLowerCase().includes(kw))
})

async function loadTeachers() {
  try { teachers.value = await GetAllTeachers() || [] } catch (_) { teachers.value = [] }
}

async function handleCreate() {
  if (!newUsername.value) { await notify.info('请输入用户名（7位数字工号）'); return }
  if (!/^\d{7}$/.test(newUsername.value)) { await notify.info('用户名必须为7位数字'); return }
  if (!newPassword.value) { await notify.info('请输入密码'); return }
  if (newPassword.value.length < 8 || newPassword.value.length > 12) { await notify.info('密码长度须为8-12位'); return }
  try {
    await CreateTeacher(newUsername.value, newPassword.value)
    newUsername.value = ''
    newPassword.value = ''
    showForm.value = false
    await loadTeachers()
    await notify.success('老师账号创建成功')
  } catch (error) { await notify.error(String(error)) }
}

async function handleDelete(username) {
  if (!await notify.confirm(`确认删除老师账号"${username}"吗？`)) return
  try {
    await DeleteUser(username)
    await loadTeachers()
    await notify.success('账号已删除')
  } catch (error) { await notify.error(String(error)) }
}

onMounted(() => { loadTeachers() })
</script>

<template>
  <div>
    <div class="card">
      <div class="card-title">用户管理</div>

      <div class="form-row">
        <input v-model="searchKeyword" placeholder="搜索用户名..." @input="searchKeyword" />
        <button class="btn-success" @click="showForm = !showForm">
          {{ showForm ? '收起' : '+ 新增老师账号' }}
        </button>
      </div>

      <div v-if="showForm" style="background:#fafafa;padding:16px;border-radius:8px;margin-bottom:16px;">
        <div class="form-row">
          <input v-model="newUsername" placeholder="用户名（7位数字工号）" />
          <input v-model="newPassword" type="password" placeholder="密码（8-12位）" />
          <button class="btn-primary" @click="handleCreate">确认新增</button>
        </div>
      </div>

      <table class="data-table">
        <thead>
          <tr><th>用户名</th><th>角色</th><th>操作</th></tr>
        </thead>
        <tbody>
          <tr v-for="t in filteredTeachers" :key="t.ID">
            <td>{{ t.Username }}</td>
            <td><span class="tag tag-green">老师</span></td>
            <td><button class="btn-danger btn-sm" @click="handleDelete(t.Username)">删除</button></td>
          </tr>
          <tr v-if="filteredTeachers.length === 0">
            <td colspan="3" style="text-align:center;color:#999;padding:24px;">暂无老师账号</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
