import { createRouter, createWebHashHistory } from 'vue-router'
import MainLayout from '../layout/MainLayout.vue'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import StudentView from '../views/StudentView.vue'
import CourseView from '../views/CourseView.vue'
import GradeView from '../views/GradeView.vue'
import StatisticsView from '../views/StatisticsView.vue'
import GpaView from '../views/GpaView.vue'
import DataManagementView from '../views/DataManagementView.vue'
import OperationLogView from '../views/OperationLogView.vue'
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
    next()
})

export default router
