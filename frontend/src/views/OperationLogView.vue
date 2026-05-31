<!-- 操作日志页面 -->
<script setup>
import { ref, onMounted } from 'vue'
import {
  GetAllOperationLogs, GetOperationLogsByTerm,
  GetOperationLogsByStudent, SearchOperationLogs,
  ReadErrorLogs,
} from '../../wailsjs/go/api/LogAPI'
import { ExportOperationLogs } from '../../wailsjs/go/api/ExcelAPI'
import { useNotify } from '../composables/useNotify'

const notify = useNotify()

const logs = ref([])
const errorLogs = ref([])
const activeTab = ref('operation')
const searchForm = ref({ Action: '', StudentID: '', Course: '', Term: '', StartTime: '', EndTime: '' })
const termFilter = ref('')
const studentFilter = ref('')

async function loadLogs() {
  try { logs.value = await GetAllOperationLogs() || [] } catch (_) { logs.value = [] }
}
async function loadErrorLogs() {
  try { errorLogs.value = await ReadErrorLogs() || [] } catch (_) { errorLogs.value = [] }
}
async function handleSearch() {
  logs.value = await SearchOperationLogs(
    searchForm.value.Action, searchForm.value.StudentID,
    searchForm.value.Course, searchForm.value.Term,
    searchForm.value.StartTime, searchForm.value.EndTime,
  )
}
async function handleSearchByTerm() {
  if (!termFilter.value) { await loadLogs(); return }
  logs.value = await GetOperationLogsByTerm(termFilter.value)
}
async function handleSearchByStudent() {
  if (!studentFilter.value) { await loadLogs(); return }
  logs.value = await GetOperationLogsByStudent(studentFilter.value)
}
async function handleExport() {
  try {
    const path = await ExportOperationLogs()
    await notify.success('日志已导出：' + path)
  } catch (error) { await notify.error(String(error)) }
}

function handleReset() {
  searchForm.value = { Action: '', StudentID: '', Course: '', Term: '', StartTime: '', EndTime: '' }
  loadLogs()
}

function formatTime(t) {
  if (!t) return ''
  const d = new Date(t)
  return d.toLocaleString('zh-CN')
}

onMounted(() => { loadLogs(); loadErrorLogs() })
</script>

<template>
  <div>
    <div style="display:flex;gap:8px;margin-bottom:16px;">
      <button :class="activeTab==='operation'?'btn-primary':'btn-default'" @click="activeTab='operation';loadLogs()">操作日志</button>
      <button :class="activeTab==='error'?'btn-primary':'btn-default'" @click="activeTab='error';loadErrorLogs()">错误日志</button>
    </div>

    <!-- 操作日志 -->
    <div v-if="activeTab==='operation'">
      <!-- 快捷筛选 -->
      <div class="card">
        <div class="form-row">
          <input v-model="studentFilter" placeholder="按学号快速筛选" @keyup.enter="handleSearchByStudent" />
          <button class="btn-primary btn-sm" @click="handleSearchByStudent">筛选</button>
          <select v-model="termFilter" @change="handleSearchByTerm">
            <option value="">按学期筛选</option>
            <option value="2024-2025-1">2024-2025-1</option>
            <option value="2024-2025-2">2024-2025-2</option>
          </select>
          <button class="btn-default btn-sm" @click="termFilter='';studentFilter='';loadLogs()">清除</button>
        </div>
      </div>

      <!-- 高级搜索 -->
      <div class="card">
        <div class="card-title">高级追溯</div>
        <div class="form-row">
          <select v-model="searchForm.Action">
            <option value="">全部操作</option>
            <option value="新增">新增</option>
            <option value="修改">修改</option>
            <option value="删除">删除</option>
            <option value="批量调整">批量调整</option>
            <option value="导入">导入</option>
          </select>
          <input v-model="searchForm.StudentID" placeholder="学号" />
          <input v-model="searchForm.Course" placeholder="课程名称" />
          <input v-model="searchForm.Term" placeholder="学期" />
          <input v-model="searchForm.StartTime" placeholder="开始日期" style="width:130px;" />
          <input v-model="searchForm.EndTime" placeholder="结束日期" style="width:130px;" />
          <button class="btn-primary" @click="handleSearch">搜索</button>
          <button class="btn-default" @click="handleReset">重置</button>
        </div>
      </div>

      <!-- 日志列表 -->
      <div class="card">
        <div style="display:flex;justify-content:space-between;align-items:center;">
          <span class="card-title" style="margin:0;padding:0;border:none;">操作记录（共 {{ logs.length }} 条，按时间倒序）</span>
          <button class="btn-primary btn-sm" @click="handleExport">导出 Excel</button>
        </div>
        <table class="data-table">
          <thead>
            <tr><th>时间</th><th>操作人</th><th>操作</th><th>学号</th><th>学生</th><th>课程</th><th>学期</th><th class="col-center">旧分</th><th class="col-center">新分</th><th>详情</th></tr>
          </thead>
          <tbody>
            <tr v-for="l in logs" :key="l.ID">
              <td style="font-size:12px;">{{ formatTime(l.Time) }}</td>
              <td>{{ l.Operator }}</td>
              <td>
                <span class="tag" :class="
                  l.Action==='新增'?'tag-green':
                  l.Action==='修改'?'tag-blue':
                  l.Action==='删除'?'tag-red':
                  l.Action==='批量调整'?'tag-orange':'tag-blue'
                ">{{ l.Action }}</span>
              </td>
              <td>{{ l.StudentID }}</td>
              <td>{{ l.Student }}</td>
              <td>{{ l.Course }}</td>
              <td>{{ l.Term }}</td>
              <td class="col-center">{{ l.OldScore }}</td>
              <td class="col-center">{{ l.NewScore }}</td>
              <td style="font-size:12px;max-width:200px;overflow:hidden;text-overflow:ellipsis;">{{ l.Detail }}</td>
            </tr>
            <tr v-if="logs.length === 0">
              <td colspan="10" style="text-align:center;color:#999;padding:24px;">暂无操作记录</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 错误日志 -->
    <div v-if="activeTab==='error'" class="card">
      <div class="card-title">校验错误日志</div>
      <p style="color:#888;font-size:13px;">来源：<code>data/error.log</code></p>
      <table v-if="errorLogs.length > 0" class="data-table">
        <thead><tr><th>时间</th><th>学生</th><th>课程</th><th>分数</th><th>原因</th></tr></thead>
        <tbody>
          <tr v-for="(e, i) in errorLogs" :key="i">
            <td style="font-size:12px;">{{ e.Time }}</td>
            <td>{{ e.Student }}</td>
            <td>{{ e.Course }}</td>
            <td>{{ e.Score }}</td>
            <td class="msg-error">{{ e.Reason }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else style="color:#999;">暂无错误日志</p>
    </div>
  </div>
</template>
