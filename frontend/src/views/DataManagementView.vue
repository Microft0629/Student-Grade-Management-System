<!-- 数据管理页面 -->
<script setup>
import { ref, onMounted } from 'vue'
import { BackupByTerm, BackupByCourse, ListBackups, RestoreFromBackup } from '../../wailsjs/go/api/BackupAPI'
import { GetAllCourses } from '../../wailsjs/go/api/CourseAPI'

const courses = ref([])
const backups = ref([])
const backupMsg = ref('')
const activeTab = ref('backup')
const termForBackup = ref('')
const courseForBackup = ref({ Term: '', CourseCode: '' })

async function loadBackups() {
  try { backups.value = await ListBackups() || [] } catch (_) { backups.value = [] }
}

async function handleBackupByTerm() {
  if (!termForBackup.value) { alert('请选择一个学期'); return }
  try {
    const dir = await BackupByTerm(termForBackup.value)
    backupMsg.value = '备份成功：' + dir
    await loadBackups()
  } catch (error) { alert(error) }
}

async function handleBackupByCourse() {
  if (!courseForBackup.value.Term || !courseForBackup.value.CourseCode) {
    alert('请选择学期和课程'); return
  }
  try {
    const dir = await BackupByCourse(courseForBackup.value.Term, courseForBackup.value.CourseCode)
    backupMsg.value = '备份成功：' + dir
    await loadBackups()
  } catch (error) { alert(error) }
}

async function handleRestore(backupName) {
  if (!confirm(`确认从备份 [${backupName}] 恢复数据？当前数据将被覆盖。`)) return
  try {
    await RestoreFromBackup(backupName)
    alert('恢复成功，数据已重新加载')
    backupMsg.value = '已从备份恢复：' + backupName
  } catch (error) { alert(error) }
}

const terms = ['2024-2025-1', '2024-2025-2', '2025-2026-1', '2025-2026-2']

onMounted(async () => {
  courses.value = await GetAllCourses()
  await loadBackups()
})
</script>

<template>
  <div>
    <div style="display:flex;gap:8px;margin-bottom:16px;">
      <button :class="activeTab==='backup'?'btn-primary':'btn-default'" @click="activeTab='backup'">数据备份</button>
      <button :class="activeTab==='restore'?'btn-primary':'btn-default'" @click="activeTab='restore'">数据恢复</button>
    </div>

    <!-- 备份 -->
    <div v-if="activeTab==='backup'">
      <div class="card">
        <div class="card-title">按学期备份</div>
        <div class="form-row">
          <select v-model="termForBackup">
            <option value="">选择学期</option>
            <option v-for="t in terms" :key="t" :value="t">{{ t }}</option>
          </select>
          <button class="btn-primary" @click="handleBackupByTerm">开始备份</button>
        </div>
      </div>

      <div class="card">
        <div class="card-title">按课程备份</div>
        <div class="form-row">
          <select v-model="courseForBackup.Term">
            <option value="">选择学期</option>
            <option v-for="t in terms" :key="t" :value="t">{{ t }}</option>
          </select>
          <select v-model="courseForBackup.CourseCode">
            <option value="">选择课程</option>
            <option v-for="c in courses.filter(c => !courseForBackup.Term || c.Term === courseForBackup.Term)" :key="c.CourseCode" :value="c.CourseCode">
              {{ c.CourseName }} ({{ c.CourseCode }})
            </option>
          </select>
          <button class="btn-primary" @click="handleBackupByCourse">开始备份</button>
        </div>
      </div>

      <p v-if="backupMsg" class="msg-success" style="font-weight:500;">{{ backupMsg }}</p>
    </div>

    <!-- 恢复 -->
    <div v-if="activeTab==='restore'" class="card">
      <div class="card-title">从备份恢复</div>
      <p style="color:#888;font-size:13px;">选择备份目录，恢复其中的成绩数据。恢复后会自动重新加载到数据库。</p>

      <table v-if="backups.length > 0" class="data-table">
        <thead><tr><th>备份名称</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="b in backups" :key="b">
            <td><span class="tag tag-blue">{{ b }}</span></td>
            <td><button class="btn-warning btn-sm" @click="handleRestore(b)">恢复</button></td>
          </tr>
        </tbody>
      </table>
      <p v-else style="color:#999;">暂无备份记录</p>

      <p v-if="backupMsg" class="msg-success" style="font-weight:500;margin-top:12px;">{{ backupMsg }}</p>
    </div>
  </div>
</template>
