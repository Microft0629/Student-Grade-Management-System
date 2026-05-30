<!-- 学生管理页面 -->
<script setup>
import { ref, onMounted, computed } from 'vue'
import {
  CreateStudent, GetAllStudents, DeleteStudent,
  SearchStudents, GetStudentsByPage,
} from '../../wailsjs/go/api/StudentAPI'

const students = ref([])
const keyword = ref('')
const total = ref(0)
const currentPage = ref(1)
const pageSize = 10
const showForm = ref(false)

const form = ref({ StudentID: '', Name: '', Gender: '', ClassName: '', Major: '' })

const totalPages = computed(() => Math.ceil(total.value / pageSize))

async function loadStudents() {
  if (keyword.value.trim() !== '') {
    students.value = await SearchStudents(keyword.value)
    total.value = students.value.length
    return
  }
  const result = await GetStudentsByPage(currentPage.value, pageSize)
  students.value = result.List || []
  total.value = result.Total
}

async function handleCreate() {
  if (!form.value.StudentID || !form.value.Name) {
    alert('学号和姓名不能为空')
    return
  }
  try {
    await CreateStudent(form.value)
    form.value = { StudentID: '', Name: '', Gender: '', ClassName: '', Major: '' }
    showForm.value = false
    await loadStudents()
    alert('新增成功')
  } catch (error) {
    alert(error)
  }
}

async function handleDelete(id) {
  if (!confirm('确认删除该学生吗？关联成绩也会一并删除。')) return
  await DeleteStudent(id)
  await loadStudents()
}

async function handleSearch() {
  currentPage.value = 1
  await loadStudents()
}

async function handleReset() {
  keyword.value = ''
  currentPage.value = 1
  await loadStudents()
}

async function prevPage() {
  if (currentPage.value > 1) { currentPage.value--; await loadStudents() }
}
async function nextPage() {
  if (currentPage.value < totalPages.value) { currentPage.value++; await loadStudents() }
}

onMounted(() => { loadStudents() })
</script>

<template>
  <div>
    <div class="card">
      <div class="card-title">学生管理</div>

      <div class="form-row">
        <input v-model="keyword" placeholder="搜索学生姓名" @keyup.enter="handleSearch" />
        <button class="btn-primary" @click="handleSearch">搜索</button>
        <button class="btn-default" @click="handleReset">重置</button>
        <button class="btn-success" @click="showForm = !showForm">
          {{ showForm ? '收起' : '+ 新增学生' }}
        </button>
      </div>

      <div v-if="showForm" class="form-row" style="background:#fafafa;padding:16px;border-radius:8px;">
        <input v-model="form.StudentID" placeholder="学号 *" />
        <input v-model="form.Name" placeholder="姓名 *" />
        <select v-model="form.Gender">
          <option value="">性别</option>
          <option value="男">男</option>
          <option value="女">女</option>
        </select>
        <input v-model="form.ClassName" placeholder="班级" />
        <input v-model="form.Major" placeholder="专业" />
        <button class="btn-primary" @click="handleCreate">确认新增</button>
      </div>

      <table class="data-table">
        <thead>
          <tr>
            <th>ID</th><th>学号</th><th>姓名</th><th>性别</th><th>班级</th><th>专业</th><th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in students" :key="s.ID">
            <td>{{ s.ID }}</td>
            <td>{{ s.StudentID }}</td>
            <td>{{ s.Name }}</td>
            <td>{{ s.Gender }}</td>
            <td>{{ s.ClassName }}</td>
            <td>{{ s.Major }}</td>
            <td>
              <button class="btn-danger btn-sm" @click="handleDelete(s.ID)">删除</button>
            </td>
          </tr>
          <tr v-if="students.length === 0">
            <td colspan="7" style="text-align:center;color:#999;padding:24px;">暂无数据</td>
          </tr>
        </tbody>
      </table>

      <div style="display:flex;justify-content:center;align-items:center;gap:16px;margin-top:16px;">
        <button class="btn-default btn-sm" @click="prevPage">上一页</button>
        <span style="font-size:13px;color:#666;">第 {{ currentPage }} / {{ totalPages || 1 }} 页（共 {{ total }} 条）</span>
        <button class="btn-default btn-sm" @click="nextPage">下一页</button>
      </div>
    </div>
  </div>
</template>
