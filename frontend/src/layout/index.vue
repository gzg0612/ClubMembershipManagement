<template>
  <div class="app-wrapper">
    <!-- 侧边栏 -->
    <div class="sidebar-container">
      <div class="logo">
        <img src="@/assets/logo.png" alt="logo">
        <span>会员管理系统</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="el-menu-vertical"
        :collapse="isCollapse"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><Monitor /></el-icon>
          <template #title>控制台</template>
        </el-menu-item>
        <el-menu-item index="/member">
          <el-icon><User /></el-icon>
          <template #title>会员管理</template>
        </el-menu-item>
        <el-menu-item index="/store">
          <el-icon><Shop /></el-icon>
          <template #title>店铺管理</template>
        </el-menu-item>
        <el-menu-item index="/tag">
          <el-icon><Collection /></el-icon>
          <template #title>标签管理</template>
        </el-menu-item>
        <el-menu-item index="/activity">
          <el-icon><Calendar /></el-icon>
          <template #title>活动管理</template>
        </el-menu-item>
        <el-menu-item index="/product">
          <el-icon><Goods /></el-icon>
          <template #title>商品管理</template>
        </el-menu-item>
        <el-menu-item index="/order">
          <el-icon><List /></el-icon>
          <template #title>订单管理</template>
        </el-menu-item>
        <el-menu-item index="/setting">
          <el-icon><Setting /></el-icon>
          <template #title>系统设置</template>
        </el-menu-item>
      </el-menu>
    </div>

    <!-- 主容器 -->
    <div class="main-container">
      <!-- 头部 -->
      <div class="navbar">
        <div class="left">
          <el-icon
            class="hamburger"
            @click="toggleSideBar"
          >
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </el-icon>
          <breadcrumb />
        </div>
        <div class="right">
          <el-dropdown trigger="click">
            <span class="user-info">
              <el-avatar :size="32" :src="userAvatar" />
              <span class="username">{{ username }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleProfile">个人信息</el-dropdown-item>
                <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>

      <!-- 内容区 -->
      <div class="app-main">
        <router-view v-slot="{ Component }">
          <transition name="fade-transform" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Monitor,
  User,
  Shop,
  Collection,
  Calendar,
  Goods,
  List,
  Setting,
  Fold,
  Expand
} from '@element-plus/icons-vue'
import Breadcrumb from './components/Breadcrumb.vue'

export default {
  name: 'Layout',
  components: {
    Monitor,
    User,
    Shop,
    Collection,
    Calendar,
    Goods,
    List,
    Setting,
    Fold,
    Expand,
    Breadcrumb
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const isCollapse = ref(false)
    const username = ref('管理员')
    const userAvatar = ref('https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png')

    // 当前激活的菜单
    const activeMenu = computed(() => {
      return route.path
    })

    // 切换侧边栏
    const toggleSideBar = () => {
      isCollapse.value = !isCollapse.value
    }

    // 处理个人信息
    const handleProfile = () => {
      router.push('/profile')
    }

    // 处理退出登录
    const handleLogout = () => {
      localStorage.removeItem('token')
      router.push('/login')
    }

    return {
      isCollapse,
      activeMenu,
      username,
      userAvatar,
      toggleSideBar,
      handleProfile,
      handleLogout
    }
  }
}
</script>

<style scoped>
.app-wrapper {
  display: flex;
  height: 100vh;
}

.sidebar-container {
  width: 210px;
  height: 100%;
  background-color: #304156;
  transition: width 0.3s;
  overflow-x: hidden;
}

.sidebar-container.collapse {
  width: 64px;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  color: #fff;
  font-size: 16px;
  font-weight: bold;
}

.logo img {
  width: 32px;
  height: 32px;
  margin-right: 12px;
}

.el-menu-vertical {
  border-right: none;
}

.el-menu-vertical:not(.el-menu--collapse) {
  width: 210px;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.navbar {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
}

.left {
  display: flex;
  align-items: center;
}

.hamburger {
  font-size: 20px;
  cursor: pointer;
  margin-right: 20px;
}

.right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.username {
  margin-left: 8px;
  font-size: 14px;
}

.app-main {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #f0f2f5;
}

.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style> 