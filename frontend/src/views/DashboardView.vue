<!-- 首页概览 -->
<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { GetAllStudents } from '../../wailsjs/go/api/StudentAPI'
import { GetAllCourses } from '../../wailsjs/go/api/CourseAPI'
import { GetAllGrades } from '../../wailsjs/go/api/GradeAPI'
import { GetStudentStatistics } from '../../wailsjs/go/api/StatisticsAPI'

const router = useRouter()
const studentCount = ref(0)
const courseCount = ref(0)
const gradeCount = ref(0)
const avgGpa = ref(0)
const today = new Date().toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' })

const statCards = [
  { label: '学生总数', value: studentCount, icon: '👤', color: '#4a90d9', bg: '#ecf3fc' },
  { label: '课程总数', value: courseCount, icon: '📚', color: '#52c41a', bg: '#edf9e8' },
  { label: '成绩记录', value: gradeCount, icon: '📝', color: '#faad14', bg: '#fffbe6' },
  { label: '平均绩点', value: avgGpa, icon: '🎯', color: '#ff6b6b', bg: '#fff1f0' },
]

const shortcuts = [
  { label: '学生管理', desc: '添加、搜索、管理学生信息', icon: '👤', path: '/main/students', color: '#4a90d9' },
  { label: '课程管理', desc: '管理课程、学期及教师信息', icon: '📚', path: '/main/courses', color: '#52c41a' },
  { label: '成绩管理', desc: '录入成绩、批量导入与调整', icon: '📝', path: '/main/grades', color: '#faad14' },
  { label: '统计分析', desc: '课程统计、排名与报表导出', icon: '📈', path: '/main/statistics', color: '#ff6b6b' },
]

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
    <!-- 欢迎横幅 -->
    <div class="welcome-banner">
      <div class="welcome-text">
        <h1>欢迎回来 👋</h1>
        <p>{{ today }}</p>
      </div>
      <div class="welcome-decoration">
        <span>🎓</span>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="dash-stat-row">
      <div v-for="card in statCards" :key="card.label" class="dash-stat-card">
        <div class="dash-stat-icon" :style="{ background: card.bg, color: card.color }">
          {{ card.icon }}
        </div>
        <div class="dash-stat-body">
          <div class="dash-stat-value" :style="{ color: card.color }">{{ card.value }}</div>
          <div class="dash-stat-label">{{ card.label }}</div>
        </div>
      </div>
    </div>

    <!-- 快捷入口 -->
    <div class="dash-shortcut-title">快捷入口</div>
    <div class="dash-shortcut-row">
      <div v-for="s in shortcuts" :key="s.label" class="dash-shortcut-card" @click="router.push(s.path)">
        <div class="dash-shortcut-icon" :style="{ color: s.color }">{{ s.icon }}</div>
        <div class="dash-shortcut-text">
          <div class="dash-shortcut-label">{{ s.label }}</div>
          <div class="dash-shortcut-desc">{{ s.desc }}</div>
        </div>
        <div class="dash-shortcut-arrow">→</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 欢迎横幅 */
.welcome-banner {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  border-radius: 16px;
  padding: 36px 40px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  color: #fff;
}
.welcome-text h1 {
  margin: 0 0 8px;
  font-size: 26px;
  font-weight: 700;
}
.welcome-text p {
  margin: 0;
  font-size: 14px;
  color: rgba(255,255,255,0.6);
}
.welcome-decoration span {
  font-size: 56px;
  opacity: 0.9;
}

/* 统计卡片 */
.dash-stat-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 28px;
}
.dash-stat-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  transition: all 0.25s;
  cursor: default;
}
.dash-stat-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(0,0,0,0.1);
}
.dash-stat-icon {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}
.dash-stat-body {
  flex: 1;
}
.dash-stat-value {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
}
.dash-stat-label {
  font-size: 13px;
  color: #999;
  margin-top: 4px;
}

/* 快捷入口 */
.dash-shortcut-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 14px;
}
.dash-shortcut-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 14px;
}
.dash-shortcut-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.06);
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid transparent;
}
.dash-shortcut-card:hover {
  border-color: #e0e0e0;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
  transform: translateX(4px);
}
.dash-shortcut-icon {
  font-size: 28px;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8f9fa;
  border-radius: 12px;
  flex-shrink: 0;
}
.dash-shortcut-text {
  flex: 1;
}
.dash-shortcut-label {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}
.dash-shortcut-desc {
  font-size: 12px;
  color: #999;
  margin-top: 3px;
}
.dash-shortcut-arrow {
  font-size: 18px;
  color: #ccc;
  transition: all 0.2s;
}
.dash-shortcut-card:hover .dash-shortcut-arrow {
  color: #4a90d9;
  transform: translateX(2px);
}

@media (max-width: 900px) {
  .dash-stat-row { grid-template-columns: repeat(2, 1fr); }
  .dash-shortcut-row { grid-template-columns: 1fr; }
}
</style>
