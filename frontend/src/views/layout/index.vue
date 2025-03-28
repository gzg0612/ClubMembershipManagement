<template>
  <div class="app-wrapper">
    <div class="sidebar-container">
      <div class="logo">会员管理系统</div>
      <el-menu
        :default-active="activeMenu"
        class="el-menu-vertical"
        background-color="#304156"
        text-color="#fff"
        active-text-color="#409EFF"
        router
      >
        <el-menu-item index="/dashboard">
          <i class="el-icon-s-home"></i>
          <span slot="title">控制台</span>
        </el-menu-item>
        
        <el-submenu index="/member">
          <template slot="title">
            <i class="el-icon-user"></i>
            <span>会员管理</span>
          </template>
          <el-menu-item index="/member/list">会员列表</el-menu-item>
          <el-menu-item index="/member/edit">添加会员</el-menu-item>
        </el-submenu>
        
        <el-submenu index="/system">
          <template slot="title">
            <i class="el-icon-setting"></i>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/system/user">用户管理</el-menu-item>
          <el-menu-item index="/system/role">角色管理</el-menu-item>
          <el-menu-item index="/system/permission">权限管理</el-menu-item>
        </el-submenu>
      </el-menu>
    </div>
    
    <div class="main-container">
      <div class="navbar">
        <div class="right-menu">
          <el-dropdown trigger="click">
            <span class="el-dropdown-link">
              管理员 <i class="el-icon-arrow-down"></i>
            </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>个人中心</el-dropdown-item>
              <el-dropdown-item divided @click.native="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
      </div>
      
      <div class="app-main">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'Layout',
  computed: {
    ...mapGetters([
      'name',
      'avatar',
      'roles'
    ]),
    activeMenu() {
      const route = this.$route
      const { meta, path } = route
      if (meta.activeMenu) {
        return meta.activeMenu
      }
      return path
    }
  },
  methods: {
    async logout() {
      try {
        await this.$store.dispatch('user/logout')
        this.$router.push(`/login?redirect=${this.$route.fullPath}`)
      } catch (error) {
        console.error(error)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.app-wrapper {
  position: relative;
  height: 100%;
  width: 100%;
  
  .sidebar-container {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    width: 210px;
    background-color: #304156;
    overflow-y: auto;
    
    .logo {
      height: 50px;
      line-height: 50px;
      background: #2b3649;
      text-align: center;
      color: #fff;
      font-size: 18px;
      font-weight: bold;
    }
    
    .el-menu {
      border-right: none;
    }
  }
  
  .main-container {
    margin-left: 210px;
    height: 100%;
    
    .navbar {
      height: 50px;
      border-bottom: 1px solid #dcdfe6;
      display: flex;
      align-items: center;
      justify-content: flex-end;
      padding: 0 15px;
      
      .right-menu {
        margin-right: 20px;
        
        .el-dropdown-link {
          cursor: pointer;
          color: #606266;
        }
      }
    }
    
    .app-main {
      padding: 15px;
      height: calc(100% - 50px);
      overflow-y: auto;
    }
  }
}
</style> 