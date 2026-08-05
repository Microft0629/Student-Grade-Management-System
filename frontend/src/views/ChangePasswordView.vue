<!-- 修改密码页面（首次登录强制修改） -->
<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ChangePassword } from '../../wailsjs/go/api/UserAPI'
import { useAuthStore } from '../store/auth'
import { useNotify } from '../composables/useNotify'

const notify = useNotify()
const router = useRouter()
const authStore = useAuthStore()

const oldPwd = ref('')
const newPwd = ref('')
const confirmPwd = ref('')
const loading = ref(false)

async function handleSubmit() {
  if (!oldPwd.value) { await notify.info('请输入旧密码'); return }
  if (newPwd.value.length < 8 || newPwd.value.length > 12) { await notify.info('新密码长度须为8-12位'); return }
  if (newPwd.value !== confirmPwd.value) { await notify.info('两次输入的新密码不一致'); return }
  loading.value = true
  try {
    await ChangePassword(oldPwd.value, newPwd.value)
    authStore.markPasswordChanged()
    oldPwd.value = ''
    newPwd.value = ''
    confirmPwd.value = ''
    await notify.success('密码修改成功')
    await router.push('/main/dashboard')
  } catch (error) {
    await notify.error(String(error))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="card" style="max-width:460px;margin:40px auto;">
    <div class="card-title">修改密码</div>
    <p style="color:#888;font-size:13px;">密码长度须为 8-12 位。</p>
    <div style="display:flex;flex-direction:column;gap:12px;margin-top:12px;">
      <input v-model="oldPwd" type="password" placeholder="旧密码" @keyup.enter="handleSubmit" />
      <input v-model="newPwd" type="password" placeholder="新密码" @keyup.enter="handleSubmit" />
      <input v-model="confirmPwd" type="password" placeholder="确认新密码" @keyup.enter="handleSubmit" />
      <button class="btn-primary" :disabled="loading" @click="handleSubmit">
        {{ loading ? '提交中...' : '确认修改' }}
      </button>
    </div>
  </div>
</template>
