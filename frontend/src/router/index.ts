import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { UserRole } from '@/types/user'

// 扩展 RouteMeta 接口
declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    requiresAuth?: boolean
    roles?: number[]
  }
}

// 路由配置
const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/login.vue'),
    meta: {
      title: '登录',
      requiresAuth: false
    }
  },
  {
    path: '/change-password',
    name: 'ChangePassword',
    component: () => import('@/views/auth/change-password.vue'),
    meta: {
      title: '修改密码',
      requiresAuth: true
    }
  },
  {
    path: '/',
    component: () => import('@/layouts/default.vue'),
    meta: {
      requiresAuth: true
    },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: '首页',
          requiresAuth: true
        }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/user/index.vue'),
        meta: {
          title: '用户管理',
          requiresAuth: true,
          roles: [UserRole.SuperAdmin, UserRole.Admin]
        }
      },
      {
        path: 'stores',
        name: 'Stores',
        component: () => import('@/views/store/index.vue'),
        meta: {
          title: '店铺管理',
          requiresAuth: true,
          roles: [UserRole.SuperAdmin, UserRole.Admin, UserRole.Manager]
        }
      },
      {
        path: 'members',
        name: 'Members',
        component: () => import('@/views/member/index.vue'),
        meta: {
          title: '会员管理',
          requiresAuth: true
        }
      },
      {
        path: 'trash',
        name: 'Trash',
        component: () => import('@/views/trash/index.vue'),
        meta: {
          title: '数据回收站',
          requiresAuth: true,
          roles: [UserRole.SuperAdmin, UserRole.Admin]
        }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/settings/index.vue'),
        meta: {
          title: '系统设置',
          requiresAuth: true,
          roles: [UserRole.SuperAdmin]
        }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: {
      title: '404',
      requiresAuth: false
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title} - 会员管理系统`

  // 获取 token 和用户信息
  const token = localStorage.getItem('token')
  const userStr = localStorage.getItem('user')
  const user = userStr ? JSON.parse(userStr) : null

  // 检查是否需要认证
  if (to.meta.requiresAuth) {
    if (!token || !user) {
      // 未登录，重定向到登录页
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
      return
    }

    // 检查角色权限
    if (to.meta.roles && user && !to.meta.roles.includes(user.role)) {
      // 无权限，重定向到首页
      next({ path: '/' })
      return
    }
  }

  // 已登录用户访问登录页，重定向到首页
  if (to.path === '/login' && token && user) {
    next({ path: '/' })
    return
  }

  next()
})

export default router 