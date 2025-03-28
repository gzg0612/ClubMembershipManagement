<template>
  <div class="error-page">
    <div class="error-code">401</div>
    <div class="error-desc">抱歉，请先登录</div>
    <el-button type="primary" @click="goLogin">去登录</el-button>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'

export default {
  name: 'Unauthorized',
  setup() {
    const router = useRouter()

    const goLogin = () => {
      router.push('/login')
    }

    return {
      goLogin
    }
  }
}
</script>

<style scoped>
.error-page {
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #f0f2f5;
}

.error-code {
  font-size: 120px;
  font-weight: bold;
  color: #909399;
  line-height: 1;
  margin-bottom: 20px;
}

.error-desc {
  font-size: 20px;
  color: #666;
  margin-bottom: 30px;
}
</style> 