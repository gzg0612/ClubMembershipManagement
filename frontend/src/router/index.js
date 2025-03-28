import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '@/store'
import { getToken } from '@/utils/auth'

Vue.use(VueRouter)

// 路由配置
const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', noAuth: true }
  },
  {
    path: '/',
    component: () => import('@/views/layout/index.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '控制台', icon: 'el-icon-s-home' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人中心', icon: 'el-icon-user' }
      }
    ]
  },
  {
    path: '/member',
    component: () => import('@/views/layout/index.vue'),
    redirect: '/member/list',
    meta: { title: '会员管理', icon: 'el-icon-s-custom' },
    children: [
      {
        path: 'list',
        name: 'MemberList',
        component: () => import('@/views/member/list.vue'),
        meta: { title: '会员列表', icon: 'el-icon-s-order' }
      },
      {
        path: 'add',
        name: 'MemberAdd',
        component: () => import('@/views/member/edit.vue'),
        meta: { title: '添加会员', icon: 'el-icon-plus' }
      },
      {
        path: 'edit/:id',
        name: 'MemberEdit',
        component: () => import('@/views/member/edit.vue'),
        meta: { title: '编辑会员', hidden: true }
      },
      {
        path: 'detail/:id',
        name: 'MemberDetail',
        component: () => import('@/views/member/detail.vue'),
        meta: { title: '会员详情', hidden: true }
      }
    ]
  },
  {
    path: '/system',
    component: () => import('@/views/layout/index.vue'),
    redirect: '/system/user',
    meta: { title: '系统管理', icon: 'el-icon-setting' },
    children: [
      {
        path: 'user',
        name: 'UserList',
        component: () => import('@/views/system/user/index.vue'),
        meta: { title: '用户管理', icon: 'el-icon-user' }
      },
      {
        path: 'role',
        name: 'RoleList',
        component: () => import('@/views/system/role/index.vue'),
        meta: { title: '角色管理', icon: 'el-icon-s-check' }
      },
      {
        path: 'permission',
        name: 'PermissionList',
        component: () => import('@/views/system/permission/index.vue'),
        meta: { title: '权限管理', icon: 'el-icon-lock' }
      }
    ]
  },
  {
    path: '/404',
    component: () => import('@/views/error/404.vue'),
    meta: { noAuth: true, hidden: true }
  },
  {
    path: '*',
    redirect: '/404',
    meta: { hidden: true }
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

// 全局前置守卫
router.beforeEach(async (to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - 会员管理系统` : '会员管理系统'

  // 判断是否需要登录权限
  if (to.meta.noAuth) {
    // 无需登录权限的页面直接放行
    next()
  } else {
    // 判断是否已登录
    const hasToken = getToken()
    if (hasToken) {
      if (to.path === '/login') {
        // 已登录状态访问登录页，重定向到首页
        next({ path: '/' })
      } else {
        // 判断是否已获取用户信息
        const hasUserInfo = store.getters.name
        if (hasUserInfo) {
          next()
        } else {
          try {
            // 获取用户信息
            await store.dispatch('user/getInfo')
            next()
          } catch (error) {
            // 获取用户信息失败，可能是token过期，清除token并跳转登录页
            await store.dispatch('user/resetToken')
            next(`/login?redirect=${to.path}`)
          }
        }
      }
    } else {
      // 未登录，跳转到登录页
      next(`/login?redirect=${to.path}`)
    }
  }
})

export default router 