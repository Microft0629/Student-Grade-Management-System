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
            localStorage.setItem('user', JSON.stringify(user))
        },

        logout() {
            this.user = null
            this.isLogin = false
            this.role = ''
            localStorage.removeItem('user')
        },

        loadUser() {
            const user = localStorage.getItem('user')
            if (user) {
                this.user = JSON.parse(user)
                this.isLogin = true
                this.role = this.user.Role || 'teacher'
            }
        },

        isAdmin() {
            return this.role === 'admin'
        },
    },
})
