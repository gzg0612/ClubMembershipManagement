<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          会员管理系统
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          请登录您的账号
        </p>
      </div>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        class="mt-8 space-y-6"
      >
        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <el-form-item prop="username">
              <el-input
                v-model="form.username"
                placeholder="用户名"
                :prefix-icon="User"
              />
            </el-form-item>
          </div>
          <div>
            <el-form-item prop="password">
              <el-input
                v-model="form.password"
                type="password"
                placeholder="密码"
                :prefix-icon="Lock"
                show-password
                @keyup.enter="handleSubmit"
              />
            </el-form-item>
          </div>
        </div>

        <div class="flex items-center justify-between">
          <div class="text-sm">
            <el-button type="primary" link @click="handleForgotPassword">
              忘记密码？
            </el-button>
          </div>
        </div>

        <div>
          <el-button
            :loading="loading"
            type="primary"
            class="w-full"
            @click="handleSubmit"
          >
            登录
          </el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { userApi } from '@/api/user'
import type { LoginParams } from '@/types/user'

const router = useRouter()
const loading = ref(false)
const formRef = ref<FormInstance>()

const form = reactive<LoginParams>({
  username: '',
  password: ''
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const { data } = await userApi.login(form)
        // 保存 token 和用户信息
        localStorage.setItem('token', data.token)
        localStorage.setItem('user', JSON.stringify(data.user))
        ElMessage.success('登录成功')
        router.push('/')
      } catch (error) {
        console.error('登录失败:', error)
        ElMessage.error('登录失败，请检查用户名和密码')
      } finally {
        loading.value = false
      }
    }
  })
}

const handleForgotPassword = () => {
  ElMessage.info('请联系系统管理员重置密码')
}
</script> 