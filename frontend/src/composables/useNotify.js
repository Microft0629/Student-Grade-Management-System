// 全局通知弹窗
import { reactive } from 'vue'

const state = reactive({
  visible: false,
  message: '',
  type: 'info',   // 'info' | 'confirm' | 'success' | 'error'
  resolve: null,
})

export function useNotify() {
  // 信息提示（仅确认按钮）
  function info(msg) {
    return new Promise((resolve) => {
      state.message = msg
      state.type = 'info'
      state.visible = true
      state.resolve = resolve
    })
  }

  // 成功提示（仅确认按钮，绿色）
  function success(msg) {
    return new Promise((resolve) => {
      state.message = msg
      state.type = 'success'
      state.visible = true
      state.resolve = resolve
    })
  }

  // 错误提示（仅确认按钮，红色）
  function error(msg) {
    return new Promise((resolve) => {
      state.message = msg
      state.type = 'error'
      state.visible = true
      state.resolve = resolve
    })
  }

  // 二次确认（取消 + 确认按钮）
  function confirm(msg) {
    return new Promise((resolve) => {
      state.message = msg
      state.type = 'confirm'
      state.visible = true
      state.resolve = resolve
    })
  }

  function handleConfirm() {
    state.visible = false
    if (state.resolve) {
      state.resolve(state.type === 'confirm' ? true : undefined)
      state.resolve = null
    }
  }

  function handleCancel() {
    state.visible = false
    if (state.resolve) {
      state.resolve(false)
      state.resolve = null
    }
  }

  return { state, info, success, error, confirm, handleConfirm, handleCancel }
}
