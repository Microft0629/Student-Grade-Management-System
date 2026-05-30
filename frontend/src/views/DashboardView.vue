<!-- 首页概览 -->
<script setup>
import { ref, onMounted } from 'vue'
import { GetAllStudents } from '../../wailsjs/go/api/StudentAPI'
import { GetAllCourses } from '../../wailsjs/go/api/CourseAPI'
import { GetAllGrades } from '../../wailsjs/go/api/GradeAPI'
import { GetStudentStatistics } from '../../wailsjs/go/api/StatisticsAPI'

const studentCount = ref(0)
const courseCount = ref(0)
const gradeCount = ref(0)
const avgGpa = ref(0)

onMounted(async () => {
  try {
    const students = await GetAllStudents()
    studentCount.value = students.length
  } catch (_) {}
  try {
    const courses = await GetAllCourses()
    courseCount.value = courses.length
  } catch (_) {}
  try {
    const grades = await GetAllGrades()
    gradeCount.value = grades.length
  } catch (_) {}
  try {
    const stats = await GetStudentStatistics()
    if (stats.length > 0) {
      let total = 0
      stats.forEach(s => { total += s.GPA })
      avgGpa.value = (total / stats.length).toFixed(2)
    }
  } catch (_) {}
})
</script>

<template>
  <div>
    <div class="stat-cards">
      <div class="stat-card">
        <div class="stat-value">{{ studentCount }}</div>
        <div class="stat-label">学生总数</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ courseCount }}</div>
        <div class="stat-label">课程总数</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ gradeCount }}</div>
        <div class="stat-label">成绩记录</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ avgGpa }}</div>
        <div class="stat-label">整体平均绩点</div>
      </div>
    </div>

    <div class="card">
      <div class="card-title">欢迎使用成绩管理系统</div>
      <p style="color:#666;line-height:1.8;">
        本系统支持学生管理、课程管理、成绩录入与查询、绩点自动计算、
        统计分析、批量操作、数据备份恢复及操作日志追溯等功能。
        请通过左侧菜单导航到各功能模块。
      </p>
    </div>
  </div>
</template>
