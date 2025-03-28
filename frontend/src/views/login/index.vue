<template>
  <div class="login-container">
    <el-form
      ref="loginFormRef"
      :model="loginForm"
      :rules="loginRules"
      class="login-form"
      autocomplete="on"
      label-position="left"
    >
      <div class="title-container">
        <h3 class="title">会员管理系统</h3>
      </div>

      <el-form-item prop="username">
        <el-input
          v-model="loginForm.username"
          placeholder="用户名"
          type="text"
          tabindex="1"
          autocomplete="on"
        >
          <template #prefix>
            <el-icon><User /></el-icon>
          </template>
        </el-input>
      </el-form-item>

      <el-form-item prop="password">
        <el-input
          v-model="loginForm.password"
          :type="passwordVisible ? 'text' : 'password'"
          placeholder="密码"
          tabindex="2"
          autocomplete="on"
          @keyup.enter="handleLogin"
        >
          <template #prefix>
            <el-icon><Lock /></el-icon>
          </template>
          <template #suffix>
            <el-icon
              class="cursor-pointer"
              @click="passwordVisible = !passwordVisible"
            >
              <View v-if="passwordVisible" />
              <Hide v-else />
            </el-icon>
          </template>
        </el-input>
      </el-form-item>

      <el-button
        :loading="loading"
        type="primary"
        style="width: 100%; margin-bottom: 30px"
        @click="handleLogin"
      >
        登录
      </el-button>
    </el-form>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, View, Hide } from '@element-plus/icons-vue'
import { login } from '@/api/auth'

export default {
  name: 'Login',
  components: {
    User,
    Lock,
    View,
    Hide
  },
  setup() {
    const router = useRouter()
    const loginFormRef = ref(null)
    const loading = ref(false)
    const passwordVisible = ref(false)

    const loginForm = reactive({
      username: '',
      password: ''
    })

    const loginRules = {
      username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
      password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
    }

    const handleLogin = async () => {
      if (!loginFormRef.value) return

      await loginFormRef.value.validate(async (valid) => {
        if (valid) {
          loading.value = true
          try {
            const res = await login(loginForm)
            localStorage.setItem('token', res.token)
            ElMessage.success('登录成功')
            router.push('/')
          } catch (error) {
            ElMessage.error('登录失败')
          } finally {
            loading.value = false
          }
        }
      })
    }

    return {
      loginFormRef,
      loginForm,
      loginRules,
      loading,
      passwordVisible,
      handleLogin
    }
  }
}
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  width: 100%;
  background-color: #2d3a4b;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;

  .login-form {
    position: relative;
    width: 520px;
    max-width: 100%;
    padding: 160px 35px 0;
    margin: 0 auto;
    overflow: hidden;

    :deep(.el-form-item) {
      border: 1px solid rgba(255, 255, 255, 0.1);
      background: rgba(0, 0, 0, 0.1);
      border-radius: 5px;
      color: #454545;

      .el-input {
        display: inline-block;
        height: 47px;
        width: 100%;

        input {
          background: transparent;
          border: 0;
          -webkit-appearance: none;
          border-radius: 0;
          padding: 12px 5px 12px 15px;
          color: #fff;
          height: 47px;
          caret-color: #fff;

          &:-webkit-autofill {
            box-shadow: 0 0 0 1000px #283443 inset !important;
            -webkit-text-fill-color: #fff !important;
          }
        }

        .el-input__prefix {
          color: #889aa4;
          margin-left: 10px;
        }
      }
    }

    .title-container {
      position: relative;

      .title {
        font-size: 26px;
        color: #eee;
        margin: 0 auto 40px auto;
        text-align: center;
        font-weight: bold;
      }
    }
  }
}
</style> 