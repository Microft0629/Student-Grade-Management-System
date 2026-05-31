<!-- 成绩管理页面 -->
<script setup>
import { ref, onMounted, computed } from 'vue'
import {
  CreateGrade, UpdateGrade, GetAllGrades, DeleteGrade, SearchGrades,
  BatchImportGrades, BatchAdjustScores, AggregateGrades,
} from '../../wailsjs/go/api/GradeAPI'
import { ExportTranscript } from '../../wailsjs/go/api/ExcelAPI'
import { GetAllStudents } from '../../wailsjs/go/api/StudentAPI'
import { GetAllCourses } from '../../wailsjs/go/api/CourseAPI'
import { useNotify } from '../composables/useNotify'
import { useAuthStore } from '../store/auth'

const notify = useNotify()
const authStore = useAuthStore()

// 当前用户是否能修改/删除该成绩
function canModify(grade) {
  if (authStore.isAdmin()) return true
  return grade.CreatorName === authStore.user?.Username
}

const grades = ref([])
const students = ref([])
const courses = ref([])
const activeTab = ref('list')

// 新增/编辑
const form = ref({ StudentID: 0, CourseID: 0, Score: 0 })
const editId = ref(0)
const editScore = ref(0)

// 查询
const searchForm = ref({ StudentKeyword: '', CourseKeyword: '', Term: '' })

// 批量导入
const importText = ref('')
const importResults = ref([])

// 批量调整
const batchForm = ref({ CourseID: 0, MinScore: 0, MaxScore: 100, Delta: 0 })
const batchResults = ref(null)

// 成绩汇总
const aggTerm = ref('')
const aggKeyword = ref('')
const aggData = ref([])

// 成绩单
const transcriptTerm = ref('')

const filteredCourses = computed(() => {
  if (!searchForm.value.Term) return courses.value
  return courses.value.filter(c => c.Term === searchForm.value.Term)
})

async function loadData() {
  grades.value = await GetAllGrades()
  students.value = await GetAllStudents()
  courses.value = await GetAllCourses()
}

// 录入/修改
async function handleCreate() {
  if (form.value.StudentID === 0) { await notify.info('请选择学生'); return }
  if (form.value.CourseID === 0) { await notify.info('请选择课程'); return }
  if (form.value.Score < 0 || form.value.Score > 100) { await notify.info('成绩必须在0-100之间'); return }
  try {
    await CreateGrade(form.value)
    form.value = { StudentID: 0, CourseID: 0, Score: 0 }
    await loadData()
  } catch (error) { await notify.error(String(error)) }
}

function openEdit(grade) {
  editId.value = grade.ID
  editScore.value = grade.Score
}
async function handleUpdate() {
  if (editScore.value < 0 || editScore.value > 100) { await notify.info('成绩必须在0-100之间'); return }
  try {
    await UpdateGrade(editId.value, editScore.value)
    editId.value = 0
    await loadData()
    await notify.success('修改成功')
  } catch (error) { await notify.error(String(error)) }
}
function cancelEdit() { editId.value = 0 }

async function handleDelete(id) {
  if (!await notify.confirm('确认删除该成绩记录吗？')) return
  try {
    await DeleteGrade(id)
    await loadData()
    await notify.success('删除成功')
  } catch (error) { await notify.error(String(error)) }
}

// 查询
async function handleSearch() {
  grades.value = await SearchGrades(
    searchForm.value.StudentKeyword,
    searchForm.value.CourseKeyword,
    searchForm.value.Term,
  )
}
async function handleReset() {
  searchForm.value = { StudentKeyword: '', CourseKeyword: '', Term: '' }
  await loadData()
}

// 批量导入
async function handleBatchImport() {
  if (!importText.value.trim()) { await notify.info('请粘贴成绩数据'); return }
  // 重新加载最新的学生和课程数据，避免因在其他页面新增后数据过期
  await loadData()
  const lines = importText.value.trim().split('\n')
  const gradeList = []
  const parseErrors = []

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i].trim()
    if (!line) continue
    const lineNum = i + 1

    const parts = line.split(/[,\t ]+/)
    if (parts.length < 3) {
      parseErrors.push(`第${lineNum}行：格式错误，需要"学号,课程代码,分数"（用逗号/Tab/空格分隔）`)
      continue
    }

    const studentID = parts[0].trim()
    const courseCode = parts[1].trim()
    const scoreStr = parts[2].trim()

    const sid = students.value.find(s => s.StudentID === studentID)
    if (!sid) {
      parseErrors.push(`第${lineNum}行：学号"${studentID}"不存在`)
      continue
    }

    const cid = courses.value.find(c => c.CourseCode === courseCode)
    if (!cid) {
      parseErrors.push(`第${lineNum}行：课程代码"${courseCode}"不存在`)
      continue
    }

    const score = parseFloat(scoreStr)
    if (isNaN(score)) {
      parseErrors.push(`第${lineNum}行：分数"${scoreStr}"不是有效数字`)
      continue
    }
    if (score < 0 || score > 100) {
      parseErrors.push(`第${lineNum}行：分数${score}不在0-100范围内`)
      continue
    }

    gradeList.push({
      StudentID: sid.ID,
      CourseID: cid.ID,
      Score: score,
      GradePoint: 0,
      Student: sid,
      Course: cid,
    })
  }

  const allResults = [...parseErrors]

  if (gradeList.length === 0) {
    importResults.value = allResults
    await notify.error('没有解析到有效数据，详见下方错误列表')
    return
  }

  try {
    const [successCount, errors] = await BatchImportGrades(gradeList)
    allResults.push(`✓ 成功导入 ${successCount} 条`)
    for (const e of errors || []) {
      allResults.push('✗ ' + e)
    }
    importResults.value = allResults
    await loadData()
  } catch (error) {
    await notify.error('导入失败: ' + String(error))
    importResults.value = allResults
  }
}

// 批量调整
async function handleBatchAdjust() {
  if (batchForm.value.CourseID === 0) { await notify.info('请选择课程'); return }
  const d = batchForm.value.Delta
  if (d === 0) { await notify.info('调整分值不能为0'); return }
  const course = courses.value.find(c => c.ID === batchForm.value.CourseID)
  if (!await notify.confirm(`确认将课程"${course.CourseName}"中分数 [${batchForm.value.MinScore}, ${batchForm.value.MaxScore}] 范围的所有成绩 ${d > 0 ? '+' + d : d} 分？`)) return
  try {
    batchResults.value = await BatchAdjustScores(batchForm.value.CourseID, batchForm.value.MinScore, batchForm.value.MaxScore, d)
    await loadData()
  } catch (error) { await notify.error(String(error)) }
}

// 汇总
async function handleAggregate() {
  aggData.value = await AggregateGrades(aggTerm.value, aggKeyword.value)
}

// 导出成绩单
async function handleExport() {
  try {
    const path = await ExportTranscript(transcriptTerm.value)
    await notify.success('成绩单已导出：' + path)
  } catch (error) { await notify.error(String(error)) }
}

onMounted(() => { loadData() })
</script>

<template>
  <div>
    <!-- Tab 切换 -->
    <div style="display:flex;gap:8px;margin-bottom:16px;">
      <button :class="activeTab==='list'?'btn-primary':'btn-default'" @click="activeTab='list'">成绩列表</button>
      <button :class="activeTab==='import'?'btn-primary':'btn-default'" @click="activeTab='import'">批量导入</button>
      <button :class="activeTab==='adjust'?'btn-primary':'btn-default'" @click="activeTab='adjust'">批量调整</button>
      <button :class="activeTab==='aggregate'?'btn-primary':'btn-default'" @click="activeTab='aggregate'">数据汇总</button>
      <button :class="activeTab==='transcript'?'btn-primary':'btn-default'" @click="activeTab='transcript'">导出成绩单</button>
    </div>

    <!-- 成绩列表 -->
    <div v-if="activeTab==='list'">
      <!-- 查询区域 -->
      <div class="card">
        <div class="form-row">
          <input v-model="searchForm.StudentKeyword" placeholder="学号或姓名" @keyup.enter="handleSearch" />
          <input v-model="searchForm.CourseKeyword" placeholder="课程名称" @keyup.enter="handleSearch" />
          <select v-model="searchForm.Term">
            <option value="">全部学期</option>
            <option v-for="t in [...new Set(courses.map(c=>c.Term))]" :key="t" :value="t">{{ t }}</option>
          </select>
          <button class="btn-primary" @click="handleSearch">查询</button>
          <button class="btn-default" @click="handleReset">重置</button>
        </div>
      </div>

      <!-- 录入区域 -->
      <div class="card">
        <div class="card-title">录入成绩</div>
        <div class="form-row">
          <select v-model.number="form.StudentID">
            <option value="0">选择学生</option>
            <option v-for="s in students" :key="s.ID" :value="s.ID">{{ s.Name }} ({{ s.StudentID }})</option>
          </select>
          <select v-model.number="form.CourseID">
            <option value="0">选择课程</option>
            <option v-for="c in courses" :key="c.ID" :value="c.ID">{{ c.CourseName }} ({{ c.Term }})</option>
          </select>
          <input v-model.number="form.Score" type="number" placeholder="分数 0-100" min="0" max="100" style="width:120px;" />
          <button class="btn-primary" @click="handleCreate">录入成绩</button>
        </div>
      </div>

      <!-- 成绩列表 -->
      <div class="card">
        <div class="card-title">成绩列表（共 {{ grades.length }} 条）</div>
        <table class="data-table">
          <thead>
            <tr>
              <th>学号</th><th>姓名</th><th>课程</th><th>学期</th><th class="col-center">学分</th><th class="col-center">分数</th><th class="col-center">绩点</th><th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="g in grades" :key="g.ID">
              <td>{{ g.Student?.StudentID }}</td>
              <td>{{ g.Student?.Name }}</td>
              <td>{{ g.Course?.CourseName }}</td>
              <td><span class="tag tag-green">{{ g.Course?.Term }}</span></td>
              <td class="col-center">{{ g.Course?.Credit }}</td>
              <td class="col-center">
                <template v-if="editId === g.ID">
                  <input v-model.number="editScore" type="number" min="0" max="100" style="width:70px;" />
                  <button class="btn-primary btn-sm" @click="handleUpdate" style="margin-left:4px;">保存</button>
                  <button class="btn-default btn-sm" @click="cancelEdit" style="margin-left:2px;">取消</button>
                </template>
                <template v-else>
                  <span :class="g.Score >= 60 ? 'msg-success' : 'msg-error'" style="font-weight:600;">{{ g.Score }}</span>
                </template>
              </td>
              <td class="col-center"><span class="tag" :class="g.GradePoint >= 2.0 ? 'tag-blue' : 'tag-red'">{{ g.GradePoint?.toFixed(1) }}</span></td>
              <td>
                <button v-if="editId !== g.ID && canModify(g)" class="btn-warning btn-sm" @click="openEdit(g)">修改</button>
                <button v-if="canModify(g)" class="btn-danger btn-sm" @click="handleDelete(g.ID)" style="margin-left:4px;">删除</button>
              </td>
            </tr>
            <tr v-if="grades.length === 0">
              <td colspan="8" style="text-align:center;color:#999;padding:24px;">暂无数据</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 批量导入 -->
    <div v-if="activeTab==='import'" class="card">
      <div class="card-title">批量导入成绩</div>
      <p style="color:#888;font-size:13px;margin-bottom:8px;">
        每行一条，格式：<code>学号,课程代码,分数</code>（用逗号、Tab 或空格分隔）
      </p>
      <div style="display:flex;gap:24px;margin-bottom:12px;">
        <div style="flex:1;">
          <span style="font-size:12px;color:#999;">已有学号：</span>
          <span v-for="s in students" :key="s.ID" class="tag tag-blue" style="margin:2px;">{{ s.StudentID }}</span>
        </div>
        <div style="flex:1;">
          <span style="font-size:12px;color:#999;">已有课程代码：</span>
          <span v-for="c in courses" :key="c.ID" class="tag tag-green" style="margin:2px;">{{ c.CourseCode }}</span>
        </div>
      </div>
      <textarea v-model="importText" rows="8"
        style="width:100%;font-family:monospace;font-size:13px;box-sizing:border-box;"
        placeholder="2024001,10871,87&#10;2024002,10871,92&#10;2024001,10872,78"></textarea>
      <div style="margin-top:12px;">
        <button class="btn-primary" @click="handleBatchImport">开始导入</button>
      </div>
      <div v-if="importResults.length > 0" style="margin-top:16px;">
        <div class="card-title">导入结果</div>
        <div v-for="(r, i) in importResults" :key="i"
          :class="r.startsWith('✓') ? 'msg-success' : 'msg-error'"
          style="font-size:13px;margin:4px 0;">{{ r }}</div>
      </div>
    </div>

    <!-- 批量调整 -->
    <div v-if="activeTab==='adjust'" class="card">
      <div class="card-title">按课程 + 分数段批量加减分</div>
      <div class="form-row">
        <select v-model.number="batchForm.CourseID">
          <option value="0">选择课程 *</option>
          <option v-for="c in courses" :key="c.ID" :value="c.ID">{{ c.CourseName }}（{{ c.Term }}）</option>
        </select>
        <span>分数范围：</span>
        <input v-model.number="batchForm.MinScore" type="number" placeholder="最低分" style="width:100px;" />
        <span>—</span>
        <input v-model.number="batchForm.MaxScore" type="number" placeholder="最高分" style="width:100px;" />
        <span>调整：</span>
        <input v-model.number="batchForm.Delta" type="number" placeholder="+5 或 -3" style="width:100px;" />
        <button class="btn-warning" @click="handleBatchAdjust">确认调整</button>
      </div>
      <div v-if="batchResults" style="margin-top:16px;">
        <p class="msg-success">共调整 {{ batchResults.AffectedCount }} 条记录</p>
        <div v-for="(d, i) in batchResults.Details" :key="i" style="font-size:13px;color:#555;margin:2px 0;">{{ d }}</div>
      </div>
    </div>

    <!-- 数据汇总 -->
    <div v-if="activeTab==='aggregate'" class="card">
      <div class="card-title">跨课程 / 跨学期成绩汇总</div>
      <div class="form-row">
        <select v-model="aggTerm">
          <option value="">全部学期</option>
          <option v-for="t in [...new Set(courses.map(c=>c.Term))]" :key="t" :value="t">{{ t }}</option>
        </select>
        <input v-model="aggKeyword" placeholder="课程关键词（可选）" />
        <button class="btn-primary" @click="handleAggregate">查询汇总</button>
      </div>
      <table v-if="aggData.length > 0" class="data-table" style="margin-top:16px;">
        <thead>
          <tr><th>学号</th><th>姓名</th><th>课程</th><th>学期</th><th class="col-center">学分</th><th class="col-center">分数</th><th class="col-center">绩点</th></tr>
        </thead>
        <tbody>
          <tr v-for="(a, i) in aggData" :key="i">
            <td>{{ a.StudentID }}</td><td>{{ a.StudentName }}</td><td>{{ a.CourseName }}</td>
            <td>{{ a.Term }}</td><td class="col-center">{{ a.Credit }}</td><td class="col-center">{{ a.Score }}</td><td class="col-center">{{ a.GradePoint?.toFixed(1) }}</td>
          </tr>
        </tbody>
      </table>
      <p v-if="aggTerm && aggData.length === 0" style="color:#999;margin-top:12px;">暂无匹配数据，请点击"查询汇总"</p>
    </div>

    <!-- 导出成绩单 -->
    <div v-if="activeTab==='transcript'" class="card">
      <div class="card-title">导出标准化成绩单</div>
      <p style="color:#888;font-size:13px;margin-bottom:16px;">
        导出包含学生详细信息、各科成绩及绩点汇总的标准 Excel 成绩单。
      </p>
      <div class="form-row">
        <select v-model="transcriptTerm">
          <option value="">全部学期</option>
          <option v-for="t in [...new Set(courses.map(c=>c.Term))]" :key="t" :value="t">{{ t }}</option>
        </select>
        <button class="btn-primary" @click="handleExport">导出 Excel</button>
      </div>
    </div>
  </div>
</template>
