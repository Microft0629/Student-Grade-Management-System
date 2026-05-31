<!-- 成绩统计页面 -->
<script setup>
import { ref, onMounted, computed } from 'vue'
import {
  GetStudentStatistics, GetStudentRanking,
} from '../../wailsjs/go/api/StatisticsAPI'
import { GetAllCourses } from '../../wailsjs/go/api/CourseAPI'
import { GetAllStudents } from '../../wailsjs/go/api/StudentAPI'
import { ExportCourseStats, ExportStudentStats } from '../../wailsjs/go/api/ExcelAPI'
import { useNotify } from '../composables/useNotify'
import SearchableSelect from '../components/SearchableSelect.vue'

const notify = useNotify()

// 课程选项：支持按课程名、课程代码搜索
const courseOptions = computed(() => courses.value.map(c => ({
  value: c.ID,
  label: c.CourseName + ' · ' + c.Term,
  searchText: c.CourseName + ' ' + c.CourseCode,
  term: c.Term,
})))

// 学生选项：支持按姓名、学号搜索
const studentOptions = computed(() => students.value.map(s => ({
  value: s.ID,
  label: s.Name + ' · ' + s.StudentID,
  searchText: s.Name + ' ' + s.StudentID,
})))

const studentStats = ref([])
const studentRankings = ref([])
const courses = ref([])
const students = ref([])
const activeTab = ref('student')

// 课程统计导出
const exportCourseId = ref(0)
const courseExporting = ref(false)

// 学生统计导出
const exportStudentId = ref(0)
const studentExporting = ref(false)

async function loadData() {
  try { studentStats.value = await GetStudentStatistics() } catch (_) {}
  try { studentRankings.value = await GetStudentRanking() } catch (_) {}
  courses.value = await GetAllCourses()
  students.value = await GetAllStudents()
}

async function handleExportCourse() {
  if (exportCourseId.value === 0) { await notify.info('请选择课程'); return }
  courseExporting.value = true
  try {
    const path = await ExportCourseStats(exportCourseId.value)
    await notify.success('报表已导出：' + path)
  } catch (error) { await notify.error(String(error)) }
  finally { courseExporting.value = false }
}

async function handleExportStudent() {
  if (exportStudentId.value === 0) { await notify.info('请选择学生'); return }
  studentExporting.value = true
  try {
    const path = await ExportStudentStats(exportStudentId.value)
    await notify.success('报表已导出：' + path)
  } catch (error) { await notify.error(String(error)) }
  finally { studentExporting.value = false }
}

onMounted(() => { loadData() })
</script>

<template>
  <div>
    <div style="display:flex;gap:8px;margin-bottom:16px;">
      <button :class="activeTab==='student'?'btn-primary':'btn-default'" @click="activeTab='student'">学生统计</button>
      <button :class="activeTab==='ranking'?'btn-primary':'btn-default'" @click="activeTab='ranking'">综合排名</button>
      <button :class="activeTab==='course'?'btn-primary':'btn-default'" @click="activeTab='course'">课程报表</button>
      <button :class="activeTab==='studentReport'?'btn-primary':'btn-default'" @click="activeTab='studentReport'">学生报表</button>
    </div>

    <!-- 学生统计 -->
    <div v-if="activeTab==='student'" class="card">
      <div class="card-title">学生成绩统计</div>
      <table class="data-table">
        <thead>
          <tr><th>学生</th><th class="col-center">平均分</th><th class="col-center">GPA</th><th class="col-center">总学分</th><th class="col-center">课程数</th></tr>
        </thead>
        <tbody>
          <tr v-for="item in studentStats" :key="item.StudentName">
            <td><strong>{{ item.StudentName }}</strong></td>
            <td class="col-center">{{ item.AverageScore?.toFixed(1) }}</td>
            <td class="col-center">
              <span class="tag" :class="item.GPA >= 3.0 ? 'tag-blue' : item.GPA >= 2.0 ? 'tag-orange' : 'tag-red'">
                {{ item.GPA?.toFixed(2) }}
              </span>
            </td>
            <td class="col-center">{{ item.TotalCredits?.toFixed(1) }}</td>
            <td class="col-center">{{ item.CourseCount }}</td>
          </tr>
          <tr v-if="studentStats.length === 0">
            <td colspan="5" style="text-align:center;color:#999;padding:24px;">暂无数据</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 综合排名 -->
    <div v-if="activeTab==='ranking'" class="card">
      <div class="card-title">学生综合排名（按平均绩点降序）</div>
      <table class="data-table">
        <thead>
          <tr><th class="col-center">排名</th><th>学号</th><th>姓名</th><th class="col-center">总分</th><th class="col-center">平均分</th><th class="col-center">绩点</th><th class="col-center">课程数</th></tr>
        </thead>
        <tbody>
          <tr v-for="r in studentRankings" :key="r.StudentID">
            <td class="col-center">
              <span v-if="r.Rank === 1" style="font-size:18px;">🥇</span>
              <span v-else-if="r.Rank === 2" style="font-size:18px;">🥈</span>
              <span v-else-if="r.Rank === 3" style="font-size:18px;">🥉</span>
              <span v-else>{{ r.Rank }}</span>
            </td>
            <td>{{ r.StudentID }}</td>
            <td><strong>{{ r.StudentName }}</strong></td>
            <td class="col-center">{{ r.TotalScore?.toFixed(1) }}</td>
            <td class="col-center">{{ r.AverageScore?.toFixed(1) }}</td>
            <td class="col-center"><span class="tag tag-blue">{{ r.GPA?.toFixed(2) }}</span></td>
            <td class="col-center">{{ r.CourseCount }}</td>
          </tr>
          <tr v-if="studentRankings.length === 0">
            <td colspan="7" style="text-align:center;color:#999;padding:24px;">暂无数据</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 课程报表（导出Excel） -->
    <div v-if="activeTab==='course'" class="card">
      <div class="card-title">单课程统计报表</div>
      <p style="color:#888;font-size:13px;margin-bottom:20px;">
        选择一门课程，导出包含平均分、及格率、分数段分布的 Excel 统计报表。
      </p>
      <div style="display:flex;gap:12px;align-items:flex-start;">
        <SearchableSelect v-model="exportCourseId" :options="courseOptions" placeholder="请选择要统计的课程（可输入关键词搜索）" />
        <button class="btn-primary" :disabled="courseExporting" @click="handleExportCourse" style="height:38px;min-width:120px;">
          {{ courseExporting ? '导出中...' : '导出 Excel' }}
        </button>
      </div>
    </div>

    <!-- 学生报表（导出Excel） -->
    <div v-if="activeTab==='studentReport'" class="card">
      <div class="card-title">单学生统计报表</div>
      <p style="color:#888;font-size:13px;margin-bottom:20px;">
        选择一名学生，导出包含各科成绩明细、总分、平均绩点及排名的 Excel 统计报表。
      </p>
      <div style="display:flex;gap:12px;align-items:flex-start;">
        <SearchableSelect v-model="exportStudentId" :options="studentOptions" placeholder="请选择要统计的学生（可输入学号或姓名搜索）" />
        <button class="btn-primary" :disabled="studentExporting" @click="handleExportStudent" style="height:38px;min-width:120px;">
          {{ studentExporting ? '导出中...' : '导出 Excel' }}
        </button>
      </div>
    </div>
  </div>
</template>
