<template>
  <div class="min-h-screen bg-gray-100">
    <!-- 顶部导航栏 -->
    <el-menu
      :default-active="activeMenu"
      class="bg-white shadow"
      mode="horizontal"
      router
    >
      <div class="flex items-center px-4">
        <h1 class="text-xl font-bold">会员管理系统</h1>
      </div>
      <div class="flex-1">
        <el-menu-item index="/">首页</el-menu-item>
        <el-menu-item
          v-if="hasPermission([UserRole.SuperAdmin, UserRole.Admin])"
          index="/users"
        >
          用户管理
        </el-menu-item>
        <el-menu-item
          v-if="hasPermission([UserRole.SuperAdmin, UserRole.Admin, UserRole.Manager])"
          index="/stores"
        >
          店铺管理
        </el-menu-item>
        <el-menu-item index="/members">会员管理</el-menu-item>
        <el-menu-item
          v-if="hasPermission([UserRole.SuperAdmin, UserRole.Admin])"
          index="/trash"
        >
          数据回收站
        </el-menu-item>
        <el-menu-item
          v-if="hasPermission([UserRole.SuperAdmin])"
          index="/settings"
        >
          系统设置
        </el-menu-item>
      </div>
      <div class="flex items-center px-4">
        <el-dropdown @command="handleCommand">
          <span class="flex items-center cursor-pointer">
            {{ user?.name }}
            <el-icon class="ml-1"><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="changePassword">修改密码</el-dropdown-item>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-menu>

    <!-- 主要内容区域 -->
    <div class="container mx-auto py-6">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowDown } from '@element-plus/icons-vue'
import { UserRole } from '@/types/user'
import type { User } from '@/types/user'

const router = useRouter()
const route = useRoute()

// 当前用户信息
const user = ref<User | null>(null)
const userStr = localStorage.getItem('user')
if (userStr) {
  user.value = JSON.parse(userStr)
}

// 当前激活的菜单项
const activeMenu = computed(() => route.path)

// 检查权限
const hasPermission = (roles: number[]) => {
  return user.value && roles.includes(user.value.role)
}

// 下拉菜单命令处理
const handleCommand = (command: string) => {
  switch (command) {
    case 'changePassword':
      router.push('/change-password')
      break
    case 'logout':
      // 清除登录信息
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      // 跳转到登录页
      router.push('/login')
      break
  }
}
</script> 