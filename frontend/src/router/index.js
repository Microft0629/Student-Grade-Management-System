import { createRouter, createWebHashHistory } from 'vue-router'
import MainLayout from '../layout/MainLayout.vue'
import LoginView from '../views/LoginView.vue'
import ChangePasswordView from '../views/ChangePasswordView.vue'
import DashboardView from '../views/DashboardView.vue'
import StudentView from '../views/StudentView.vue'
import CourseView from '../views/CourseView.vue'
import GradeView from '../views/GradeView.vue'
import StatisticsView from '../views/StatisticsView.vue'
import GpaView from '../views/GpaView.vue'
import DataManagementView from '../views/DataManagementView.vue'
import OperationLogView from '../views/OperationLogView.vue'
import UserManagementView from '../views/UserManagementView.vue'
import { useAuthStore } from '../store/auth'

const routes = [
    {
        path: '/',
        component: LoginView,
    },
    {
        path: '/main',
        component: MainLayout,
        children: [
            { path: 'dashboard',  component: DashboardView },
            { path: 'students',   component: StudentView },
            { path: 'courses',    component: CourseView },
            { path: 'grades',     component: GradeView },
            { path: 'statistics', component: StatisticsView },
            { path: 'gpa',        component: GpaView },
            { path: 'datamgmt',   component: DataManagementView },
            { path: 'logs',       component: OperationLogView },
            { path: 'users',      component: UserManagementView },
            { path: 'password',   component: ChangePasswordView },
        ],
    },
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()
    authStore.loadUser()
    if (to.path.startsWith('/main') && !authStore.isLogin) {
        next('/')
        return
    }
    // 首次登录强制改密
    if (authStore.isLogin && authStore.user?.MustChangePassword && to.path !== '/main/password') {
        next('/main/password')
        return
    }
    // 管理员专属路由：数据管理、操作日志、用户管理
    const adminRoutes = ['/main/datamgmt', '/main/logs', '/main/users']
    if (adminRoutes.includes(to.path) && !authStore.isAdmin()) {
        next('/main/dashboard')
        return
    }
    next()
})

export default router
