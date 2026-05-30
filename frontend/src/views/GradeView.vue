<!-- 成绩管理页面 -->
<script setup>
import { ref, onMounted, computed } from 'vue'
import {
  CreateGrade, UpdateGrade, GetAllGrades, DeleteGrade, SearchGrades,
  BatchImportGrades, BatchAdjustScores, AggregateGrades, ExportTranscript,
} from '../../wailsjs/go/api/GradeAPI'
import { GetAllStudents } from '../../wailsjs/go/api/StudentAPI'
import { GetAllCourses } from '../../wailsjs/go/api/CourseAPI'

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
const batchForm = ref({ MinScore: 0, MaxScore: 100, Delta: 0 })
const batchResults = ref(null)

// 成绩汇总
const aggTerm = ref('')
const aggKeyword = ref('')
const aggData = ref([])

// 成绩单
const transcriptTerm = ref('')
const transcriptText = ref('')
const showTranscript = ref(false)

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
  if (form.value.StudentID === 0) { alert('请选择学生'); return }
  if (form.value.CourseID === 0) { alert('请选择课程'); return }
  if (form.value.Score < 0 || form.value.Score > 100) { alert('成绩必须在0-100之间'); return }
  try {
    await CreateGrade(form.value)
    form.value = { StudentID: 0, CourseID: 0, Score: 0 }
    await loadData()
  } catch (error) { alert(error) }
}

function openEdit(grade) {
  editId.value = grade.ID
  editScore.value = grade.Score
}
async function handleUpdate() {
  if (editScore.value < 0 || editScore.value > 100) { alert('成绩必须在0-100之间'); return }
  try {
    await UpdateGrade(editId.value, editScore.value)
    editId.value = 0
    await loadData()
  } catch (error) { alert(error) }
}
function cancelEdit() { editId.value = 0 }

async function handleDelete(id) {
  if (!confirm('确认删除该成绩记录吗？')) return
  await DeleteGrade(id)
  await loadData()
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
  if (!importText.value.trim()) { alert('请粘贴成绩数据'); return }
  const lines = importText.value.trim().split('\n')
  const gradeList = []
  for (const line of lines) {
    const parts = line.split(/[,\t]/)
    if (parts.length < 3) continue
    const sid = students.value.find(s => s.StudentID === parts[0].trim())
    const cid = courses.value.find(c => c.CourseCode === parts[1].trim())
    if (!sid || !cid) continue
    gradeList.push({
      StudentID: sid.ID,
      CourseID: cid.ID,
      Score: parseFloat(parts[2]),
      GradePoint: 0,
      Student: sid,
      Course: cid,
    })
  }
  if (gradeList.length === 0) { alert('没有解析到有效数据，请检查格式'); return }
  try {
    const [successCount, errors] = await BatchImportGrades(gradeList)
    importResults.value = errors || []
    await loadData()
    alert(`导入完成：成功 ${successCount} 条，失败 ${(errors || []).length} 条`)
  } catch (error) { alert(error) }
}

// 批量调整
async function handleBatchAdjust() {
  const d = batchForm.value.Delta
  if (d === 0) { alert('调整分值不能为0'); return }
  if (!confirm(`确认将分数 [${batchForm.value.MinScore}, ${batchForm.value.MaxScore}] 范围的所有成绩 ${d > 0 ? '+' + d : d} 分？`)) return
  try {
    batchResults.value = await BatchAdjustScores(batchForm.value.MinScore, batchForm.value.MaxScore, d)
    await loadData()
  } catch (error) { alert(error) }
}

// 汇总
async function handleAggregate() {
  aggData.value = await AggregateGrades(aggTerm.value, aggKeyword.value)
}

// 导出成绩单
async function handleExport() {
  transcriptText.value = await ExportTranscript(transcriptTerm.value)
  showTranscript.value = true
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
              <th>ID</th><th>学号</th><th>姓名</th><th>课程</th><th>学期</th><th>学分</th><th>分数</th><th>绩点</th><th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="g in grades" :key="g.ID">
              <td>{{ g.ID }}</td>
              <td>{{ g.Student?.StudentID }}</td>
              <td>{{ g.Student?.Name }}</td>
              <td>{{ g.Course?.CourseName }}</td>
              <td><span class="tag tag-green">{{ g.Course?.Term }}</span></td>
              <td>{{ g.Course?.Credit }}</td>
              <td>
                <template v-if="editId === g.ID">
                  <input v-model.number="editScore" type="number" min="0" max="100" style="width:70px;" />
                  <button class="btn-primary btn-sm" @click="handleUpdate" style="margin-left:4px;">保存</button>
                  <button class="btn-default btn-sm" @click="cancelEdit" style="margin-left:2px;">取消</button>
                </template>
                <template v-else>
                  <span :class="g.Score >= 60 ? 'msg-success' : 'msg-error'" style="font-weight:600;">{{ g.Score }}</span>
                </template>
              </td>
              <td><span class="tag" :class="g.GradePoint >= 2.0 ? 'tag-blue' : 'tag-red'">{{ g.GradePoint?.toFixed(1) }}</span></td>
              <td>
                <button class="btn-warning btn-sm" @click="openEdit(g)" v-if="editId !== g.ID">修改</button>
                <button class="btn-danger btn-sm" @click="handleDelete(g.ID)" style="margin-left:4px;">删除</button>
              </td>
            </tr>
            <tr v-if="grades.length === 0">
              <td colspan="9" style="text-align:center;color:#999;padding:24px;">暂无数据</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 批量导入 -->
    <div v-if="activeTab==='import'" class="card">
      <div class="card-title">批量导入成绩</div>
      <p style="color:#888;font-size:13px;">每行一条，格式：<code>学号,课程代码,分数</code>（用逗号或Tab分隔）</p>
      <textarea v-model="importText" rows="8"
        style="width:100%;font-family:monospace;font-size:13px;box-sizing:border-box;"
        placeholder="2024001,10871,87&#10;2024002,10871,92&#10;2024001,10872,78"></textarea>
      <div style="margin-top:12px;">
        <button class="btn-primary" @click="handleBatchImport">开始导入</button>
      </div>
      <div v-if="importResults.length > 0" style="margin-top:16px;">
        <div class="card-title">导入结果</div>
        <div v-for="(r, i) in importResults" :key="i" class="msg-error" style="font-size:13px;margin:4px 0;">{{ r }}</div>
      </div>
    </div>

    <!-- 批量调整 -->
    <div v-if="activeTab==='adjust'" class="card">
      <div class="card-title">按分数段批量加减分</div>
      <div class="form-row">
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
          <tr><th>学号</th><th>姓名</th><th>课程</th><th>学期</th><th>学分</th><th>分数</th><th>绩点</th></tr>
        </thead>
        <tbody>
          <tr v-for="(a, i) in aggData" :key="i">
            <td>{{ a.StudentID }}</td><td>{{ a.StudentName }}</td><td>{{ a.CourseName }}</td>
            <td>{{ a.Term }}</td><td>{{ a.Credit }}</td><td>{{ a.Score }}</td><td>{{ a.GradePoint?.toFixed(1) }}</td>
          </tr>
        </tbody>
      </table>
      <p v-if="aggTerm && aggData.length === 0" style="color:#999;margin-top:12px;">暂无匹配数据，请点击"查询汇总"</p>
    </div>

    <!-- 导出成绩单 -->
    <div v-if="activeTab==='transcript'" class="card">
      <div class="card-title">导出标准化成绩单</div>
      <div class="form-row">
        <select v-model="transcriptTerm">
          <option value="">全部学期</option>
          <option v-for="t in [...new Set(courses.map(c=>c.Term))]" :key="t" :value="t">{{ t }}</option>
        </select>
        <button class="btn-primary" @click="handleExport">生成成绩单</button>
      </div>
      <div v-if="showTranscript" style="margin-top:16px;">
        <div class="report-text">{{ transcriptText }}</div>
      </div>
    </div>
  </div>
</template>
