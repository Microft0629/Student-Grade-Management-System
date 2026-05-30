import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {

    state: () => ({

        user: null,

        isLogin: false,

    }),

    actions: {

        setUser(user) {  // 登录成功，保存用户

            this.user = user

            this.isLogin = true

            localStorage.setItem(
                'user',
                JSON.stringify(user)
            )
        },

        logout() {  // 退出登录

            this.user = null

            this.isLogin = false

            localStorage.removeItem('user')
        },

        loadUser() {  // 刷新页面后，恢复登录状态

            const user = localStorage.getItem('user')

            if (user) {

                this.user = JSON.parse(user)

                this.isLogin = true
            }
        }

    }

})