<!-- 绩点规则页面 -->
<script setup>
import { ref, onMounted } from 'vue'
import { GetGpaRules, SaveGpaRules, ResetGpaRules, RecalculateAllGPA } from '../../wailsjs/go/api/GpaAPI'
import { useNotify } from '../composables/useNotify'

const notify = useNotify()

const formulaText = ref('')
const editMode = ref(false)
const editedFormula = ref('')
const recalcResult = ref(-1)
const recalcMsg = ref('')

async function loadRules() {
  try {
    formulaText.value = await GetGpaRules()
  } catch (_) {
    formulaText.value = ''
  }
}

function startEdit() {
  editedFormula.value = formulaText.value
  editMode.value = true
}
function cancelEdit() {
  editMode.value = false
}
async function handleSave() {
  try {
    await SaveGpaRules(editedFormula.value)
    formulaText.value = editedFormula.value
    editMode.value = false
    await notify.success('保存成功')
  } catch (error) { await notify.error(String(error)) }
}
async function handleReset() {
  if (!await notify.confirm('确认重置为默认绩点换算公式？')) return
  try {
    await ResetGpaRules()
    await loadRules()
    await notify.success('已重置')
  } catch (error) { await notify.error(String(error)) }
}
async function handleRecalculate() {
  if (!await notify.confirm('确认重新计算所有学生的绩点？此操作将覆盖所有已有绩点数据。')) return
  recalcMsg.value = '正在计算...'
  try {
    const count = await RecalculateAllGPA()
    recalcResult.value = count
    recalcMsg.value = `计算完成：共更新 ${count} 条绩点记录`
  } catch (error) {
    recalcMsg.value = ''
    await notify.error(String(error))
  }
}

onMounted(() => { loadRules() })
</script>

<template>
  <div>
    <div class="card">
      <div class="card-title">绩点换算规则</div>
      <p style="color:#888;font-size:13px;margin-bottom:16px;">
        绩点计算公式存储在 <code>data/gpa_rules.txt</code>，可在此查看和编辑。
        实际计算始终使用 <code>utils/gpa.go</code> 中的公式：score ≥ 60 → gp = round(score/10 − 5, 1)；score &lt; 60 → gp = 0
      </p>

      <div v-if="!editMode">
        <div class="report-text">{{ formulaText || '（未找到规则文件）' }}</div>
        <div style="margin-top:12px;display:flex;gap:8px;">
          <button class="btn-primary" @click="startEdit">编辑规则</button>
          <button class="btn-default" @click="handleReset">重置为默认</button>
        </div>
      </div>

      <div v-else>
        <textarea v-model="editedFormula" rows="8"
          style="width:100%;font-family:monospace;font-size:13px;box-sizing:border-box;"></textarea>
        <div style="margin-top:12px;display:flex;gap:8px;">
          <button class="btn-primary" @click="handleSave">保存</button>
          <button class="btn-default" @click="cancelEdit">取消</button>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="card-title">批量重新计算绩点</div>
      <p style="color:#888;font-size:13px;">
        使用当前公式重新计算所有成绩记录的绩点，并更新成绩文件。
        仅当绩点值发生变化时才会更新对应记录。
      </p>
      <button class="btn-warning" @click="handleRecalculate">开始重新计算</button>
      <p v-if="recalcMsg" :class="recalcResult >= 0 ? 'msg-success' : 'msg-error'" style="margin-top:12px;font-weight:500;">
        {{ recalcMsg }}
      </p>
    </div>
  </div>
</template>
