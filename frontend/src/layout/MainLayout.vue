<!-- 主界面布局 -->
<script setup>
    import { useRouter, useRoute } from 'vue-router'
    import { useAuthStore } from '../store/auth'
    import { computed, onMounted } from 'vue'
    import { WindowGetSize, WindowSetSize, WindowGetPosition, WindowSetPosition } from '../../wailsjs/runtime/runtime'

    const router = useRouter()
    const route = useRoute()
    const authStore = useAuthStore()

    // 窗口大小记忆
    let saveTimer = null
    function saveWindowState() {
        WindowGetSize().then(size => {
            localStorage.setItem('win_w', size.w)
            localStorage.setItem('win_h', size.h)
        })
        WindowGetPosition().then(pos => {
            localStorage.setItem('win_x', pos.x)
            localStorage.setItem('win_y', pos.y)
        })
    }
    function restoreWindowState() {
        const w = parseInt(localStorage.getItem('win_w')) || 1024
        const h = parseInt(localStorage.getItem('win_h')) || 768
        const x = parseInt(localStorage.getItem('win_x'))
        const y = parseInt(localStorage.getItem('win_y'))
        if (w > 0 && h > 0) {
            WindowSetSize(w, h)
        }
        if (!isNaN(x) && !isNaN(y)) {
            WindowSetPosition(x, y)
        }
    }

    onMounted(() => {
        restoreWindowState()
        // 监听窗口大小变化，延迟保存避免频繁写入
        window.addEventListener('resize', () => {
            clearTimeout(saveTimer)
            saveTimer = setTimeout(saveWindowState, 500)
        })
    })

    const menuItems = computed(() => {
        const items = [
            { path: '/main/dashboard',   label: '首页概览', icon: '📊' },
            { path: '/main/students',     label: '学生管理', icon: '👤' },
            { path: '/main/courses',      label: '课程管理', icon: '📚' },
            { path: '/main/grades',       label: '成绩管理', icon: '📝' },
            { path: '/main/statistics',   label: '统计分析', icon: '📈' },
            { path: '/main/gpa',          label: '绩点规则', icon: '🎯' },
            { path: '/main/datamgmt',     label: '数据管理', icon: '💾' },
        ]
        if (authStore.isAdmin()) {
            items.push({ path: '/main/users',  label: '用户管理', icon: '⚙️' })
            items.push({ path: '/main/logs',   label: '操作日志', icon: '📋' })
        }
        return items
    })

    function isActive(path) {
        return route.path === path
    }

    function handleLogout() {
        authStore.logout()
        router.push('/')
    }
</script>

<template>
    <div class="layout">
        <div class="sidebar">
            <div class="sidebar-brand">
                <span class="brand-icon">🎓</span>
                <span class="brand-text">成绩管理系统</span>
            </div>

            <div class="sidebar-menu">
                <div
                    v-for="item in menuItems"
                    :key="item.path"
                    class="menu-item"
                    :class="{ active: isActive(item.path) }"
                    @click="router.push(item.path)"
                >
                    <span class="menu-icon">{{ item.icon }}</span>
                    <span class="menu-label">{{ item.label }}</span>
                </div>
            </div>

            <div class="sidebar-footer">
                <div class="user-info">
                    <span class="user-avatar">👤</span>
                    <span class="user-name">{{ authStore.user?.Username }}</span>
                </div>
                <button class="logout-btn" @click="handleLogout">退出登录</button>
            </div>
        </div>

        <div class="main">
            <div class="header">
                <div class="header-title">
                    {{ menuItems.find(m => isActive(m.path))?.label || '' }}
                </div>
            </div>
            <div class="content">
                <router-view />
            </div>
        </div>
    </div>
</template>

<style scoped>
.layout {
    display: flex;
    height: 100vh;
    background: #f0f2f5;
}

/* 侧边栏 */
.sidebar {
    width: 230px;
    background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
    color: white;
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
}
.sidebar-brand {
    padding: 24px 20px 20px;
    display: flex;
    align-items: center;
    gap: 10px;
    border-bottom: 1px solid rgba(255,255,255,0.08);
}
.brand-icon { font-size: 24px; }
.brand-text { font-size: 16px; font-weight: 700; letter-spacing: 0.5px; }

.sidebar-menu {
    flex: 1;
    padding: 12px 0;
    overflow-y: auto;
}
.menu-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 20px;
    cursor: pointer;
    transition: all 0.2s;
    color: rgba(255,255,255,0.65);
    font-size: 14px;
    border-left: 3px solid transparent;
}
.menu-item:hover {
    background: rgba(255,255,255,0.06);
    color: #fff;
}
.menu-item.active {
    background: rgba(74,144,217,0.2);
    color: #fff;
    border-left-color: #4a90d9;
    font-weight: 600;
}
.menu-icon { font-size: 16px; width: 22px; text-align: center; }

.sidebar-footer {
    padding: 16px 20px;
    border-top: 1px solid rgba(255,255,255,0.08);
}
.user-info {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 10px;
    color: rgba(255,255,255,0.75);
    font-size: 13px;
}
.user-avatar { font-size: 18px; }
.logout-btn {
    width: 100%;
    padding: 8px;
    background: rgba(255,255,255,0.1);
    color: rgba(255,255,255,0.7);
    border: 1px solid rgba(255,255,255,0.12);
    border-radius: 6px;
    font-size: 13px;
    cursor: pointer;
    transition: all 0.2s;
}
.logout-btn:hover {
    background: rgba(255,77,79,0.25);
    color: #ff6b6b;
    border-color: rgba(255,77,79,0.3);
}

/* 主内容区 */
.main {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}
.header {
    height: 56px;
    background: #fff;
    display: flex;
    align-items: center;
    padding: 0 28px;
    box-shadow: 0 1px 4px rgba(0,0,0,0.06);
    flex-shrink: 0;
}
.header-title {
    font-size: 16px;
    font-weight: 600;
    color: #1a1a1a;
}
.content {
    flex: 1;
    padding: 24px 28px;
    overflow-y: auto;
}
</style>
