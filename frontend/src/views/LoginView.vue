<!-- 登录页面 -->
<script setup>
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { useAuthStore } from '../store/auth'
  import { Login } from '../../wailsjs/go/api/AuthAPI'

  const username = ref('')
  const password = ref('')
  const loading = ref(false)
  const router = useRouter()
  const authStore = useAuthStore()

  async function handleLogin() {
    if (!username.value || !password.value) {
      alert('请输入用户名和密码')
      return
    }
    loading.value = true
    try {
      const user = await Login(username.value, password.value)
      authStore.setUser(user)
      await router.push('/main/dashboard')
    } catch (err) {
      alert(err)
    } finally {
      loading.value = false
    }
  }
</script>

<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <span class="login-icon">🎓</span>
        <h1>成绩管理系统</h1>
        <p>Student Grade Management System</p>
      </div>

      <div class="login-form">
        <div class="input-group">
          <span class="input-prefix">👤</span>
          <input
            v-model="username"
            placeholder="请输入用户名"
            @keyup.enter="handleLogin"
          />
        </div>
        <div class="input-group">
          <span class="input-prefix">🔒</span>
          <input
            v-model="password"
            type="password"
            placeholder="请输入密码"
            @keyup.enter="handleLogin"
          />
        </div>
        <button class="login-btn" :disabled="loading" @click="handleLogin">
          {{ loading ? '登录中...' : '登 录' }}
        </button>
      </div>

      <div class="login-hint">默认账号：admin / 123456</div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
}
.login-card {
  width: 400px;
  background: #fff;
  border-radius: 16px;
  padding: 48px 40px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}
.login-header {
  text-align: center;
  margin-bottom: 36px;
}
.login-icon { font-size: 48px; }
.login-header h1 {
  font-size: 24px;
  color: #1a1a2e;
  margin: 12px 0 4px;
}
.login-header p {
  font-size: 13px;
  color: #999;
  margin: 0;
}
.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.input-group {
  display: flex;
  align-items: center;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  overflow: hidden;
  transition: border-color 0.2s;
}
.input-group:focus-within {
  border-color: #4a90d9;
}
.input-prefix {
  padding: 0 12px;
  font-size: 16px;
  background: #fafafa;
  border-right: 1px solid #e0e0e0;
  height: 44px;
  line-height: 44px;
}
.input-group input {
  flex: 1;
  border: none;
  border-radius: 0;
  padding: 0 14px;
  height: 44px;
  font-size: 15px;
}
.input-group input:focus {
  box-shadow: none;
}
.login-btn {
  height: 46px;
  background: linear-gradient(135deg, #4a90d9, #357abd);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 8px;
}
.login-btn:hover { opacity: 0.9; transform: translateY(-1px); }
.login-btn:disabled { opacity: 0.6; cursor: not-allowed; transform: none; }
.login-hint {
  text-align: center;
  margin-top: 20px;
  font-size: 12px;
  color: #bbb;
}
</style>
