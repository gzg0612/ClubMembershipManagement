<template>
  <div class="error-page">
    <div class="error-code">403</div>
    <div class="error-desc">抱歉，您没有权限访问此页面</div>
    <el-button type="primary" @click="goHome">返回首页</el-button>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'

export default {
  name: 'Forbidden',
  setup() {
    const router = useRouter()

    const goHome = () => {
      router.push('/')
    }

    return {
      goHome
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
  color: #E6A23C;
  line-height: 1;
  margin-bottom: 20px;
}

.error-desc {
  font-size: 20px;
  color: #666;
  margin-bottom: 30px;
}
</style> 