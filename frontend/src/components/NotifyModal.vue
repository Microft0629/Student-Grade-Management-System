<!-- 全局通知弹窗 -->
<script setup>
import { useNotify } from '../composables/useNotify'
const { state, handleConfirm, handleCancel } = useNotify()
</script>

<template>
  <div v-if="state.visible" class="modal-mask" @click.self="state.type === 'confirm' ? handleCancel() : handleConfirm()">
    <div class="notify-box" :class="'notify-' + state.type">
      <div class="notify-icon">
        <span v-if="state.type === 'success'">✓</span>
        <span v-else-if="state.type === 'error'">✗</span>
        <span v-else-if="state.type === 'confirm'">?</span>
        <span v-else>!</span>
      </div>
      <div class="notify-msg">{{ state.message }}</div>
      <div class="notify-actions">
        <button v-if="state.type === 'confirm'" class="btn-default" @click="handleCancel">取消</button>
        <button
          :class="state.type === 'error' ? 'btn-danger' : state.type === 'confirm' ? 'btn-danger' : 'btn-primary'"
          @click="handleConfirm"
        >{{ state.type === 'confirm' ? '确认' : '确定' }}</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.notify-box {
  background: #fff;
  border-radius: 16px;
  padding: 36px 40px 28px;
  text-align: center;
  min-width: 380px;
  max-width: 480px;
  box-shadow: 0 12px 40px rgba(0,0,0,0.18);
  animation: notifyIn 0.2s ease;
}
@keyframes notifyIn {
  from { opacity: 0; transform: scale(0.9) translateY(-10px); }
  to   { opacity: 1; transform: scale(1) translateY(0); }
}
.notify-icon {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26px;
  font-weight: 700;
  margin: 0 auto 16px;
}
.notify-info .notify-icon    { background: #ecf3fc; color: #4a90d9; }
.notify-success .notify-icon { background: #edf9e8; color: #52c41a; }
.notify-error .notify-icon   { background: #fff1f0; color: #ff4d4f; }
.notify-confirm .notify-icon { background: #fffbe6; color: #faad14; }

.notify-msg {
  font-size: 15px;
  color: #333;
  line-height: 1.6;
  margin-bottom: 24px;
  white-space: pre-wrap;
  word-break: break-word;
}
.notify-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}
.notify-actions button {
  min-width: 100px;
  padding: 10px 24px;
  font-size: 14px;
}
</style>
