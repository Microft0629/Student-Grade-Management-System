<!-- 课程管理页面 -->
<script setup>
import { ref, onMounted } from 'vue'
import { CreateCourse, GetAllCourses, DeleteCourse, SearchCourses } from '../../wailsjs/go/api/CourseAPI'
import { useNotify } from '../composables/useNotify'
import { useAuthStore } from '../store/auth'

const notify = useNotify()
const authStore = useAuthStore()

function canDelete(course) {
  if (authStore.isAdmin()) return true
  return course.CreatorName === authStore.user?.Username
}

const courses = ref([])
const keyword = ref('')
const showForm = ref(false)
const form = ref({ CourseCode: '', CourseName: '', Term: '', Credit: 0, Teacher: '' })
const terms = ['2024-2025-1', '2024-2025-2', '2025-2026-1', '2025-2026-2']

async function loadCourses() {
  courses.value = await GetAllCourses()
}

async function handleCreate() {
  if (!form.value.CourseCode) { await notify.info('请输入课程代码'); return }
  if (!form.value.CourseName) { await notify.info('请输入课程名称'); return }
  if (!form.value.Term) { await notify.info('请选择学期'); return }
  if (!form.value.Credit || form.value.Credit <= 0) { await notify.info('请输入有效的学分'); return }
  if (!form.value.Teacher) { await notify.info('请输入任课教师'); return }
  try {
    await CreateCourse(form.value)
    form.value = { CourseCode: '', CourseName: '', Term: '', Credit: 0, Teacher: '' }
    showForm.value = false
    await loadCourses()
    await notify.success('新增成功')
  } catch (error) { await notify.error(String(error)) }
}

async function handleDelete(id) {
  if (!await notify.confirm('确认删除该课程吗？关联成绩也会一并删除。')) return
  try {
    await DeleteCourse(id)
    await loadCourses()
    await notify.success('删除成功')
  } catch (error) { await notify.error(String(error)) }
}

async function handleSearch() {
  if (!keyword.value.trim()) { await loadCourses(); return }
  courses.value = await SearchCourses(keyword.value)
}

onMounted(() => { loadCourses() })
</script>

<template>
  <div>
    <div class="card">
      <div class="card-title">课程管理</div>

      <div class="form-row">
        <input v-model="keyword" placeholder="搜索课程名称" @keyup.enter="handleSearch" />
        <button class="btn-primary" @click="handleSearch">搜索</button>
        <button class="btn-default" @click="keyword='';loadCourses()">重置</button>
        <button class="btn-success" @click="showForm = !showForm">
          {{ showForm ? '收起' : '+ 新增课程' }}
        </button>
      </div>

      <div v-if="showForm" class="form-row" style="background:#fafafa;padding:16px;border-radius:8px;">
        <input v-model="form.CourseCode" placeholder="课程编号 *" />
        <input v-model="form.CourseName" placeholder="课程名称 *" />
        <select v-model="form.Term">
          <option value="">选择学期 *</option>
          <option v-for="t in terms" :key="t" :value="t">{{ t }}</option>
        </select>
        <input v-model.number="form.Credit" type="number" placeholder="学分 *" min="0" step="0.5" style="width:100px;" />
        <input v-model="form.Teacher" placeholder="任课教师" />
        <button class="btn-primary" @click="handleCreate">确认新增</button>
      </div>

      <table class="data-table">
        <thead>
          <tr>
            <th>课程代码</th><th>课程名称</th><th>学期</th><th class="col-center">学分</th><th>教师</th><th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in courses" :key="c.ID">
            <td>{{ c.CourseCode }}</td>
            <td>{{ c.CourseName }}</td>
            <td>{{ c.Term }}</td>
            <td class="col-center">{{ c.Credit }}</td>
            <td>{{ c.Teacher }}</td>
            <td>
              <button v-if="canDelete(c)" class="btn-danger btn-sm" @click="handleDelete(c.ID)">删除</button>
            </td>
          </tr>
          <tr v-if="courses.length === 0">
            <td colspan="6" style="text-align:center;color:#999;padding:24px;">暂无数据</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
