import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: null,
        isLogin: false,
        role: '',
    }),

    actions: {
        setUser(user) {
            this.user = user
            this.isLogin = true
            this.role = user.Role || 'teacher'
            // 只持久化非敏感信息，避免把密码哈希等字段写入 localStorage
            localStorage.setItem('user', JSON.stringify({
                Username: user.Username,
                Role: user.Role,
            }))
        },

        logout() {
            this.user = null
            this.isLogin = false
            this.role = ''
            localStorage.removeItem('user')
        },

        loadUser() {
            try {
                const raw = localStorage.getItem('user')
                if (!raw) return
                const saved = JSON.parse(raw)
                if (!saved || !saved.Username) {
                    localStorage.removeItem('user')
                    return
                }
                this.user = saved
                this.isLogin = true
                this.role = saved.Role || 'teacher'
            } catch (_) {
                localStorage.removeItem('user')
            }
        },

        isAdmin() {
            return this.role === 'admin'
        },
    },
})
