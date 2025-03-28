<template>
  <el-breadcrumb separator="/">
    <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
    <el-breadcrumb-item
      v-for="(item, index) in breadcrumbs"
      :key="index"
      :to="item.path"
    >
      {{ item.title }}
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

export default {
  name: 'Breadcrumb',
  setup() {
    const route = useRoute()
    const breadcrumbs = ref([])

    // 获取面包屑数据
    const getBreadcrumbs = () => {
      const matched = route.matched.filter(item => item.meta && item.meta.title)
      breadcrumbs.value = matched.map(item => ({
        title: item.meta.title,
        path: item.path
      }))
    }

    // 监听路由变化
    watch(
      () => route.path,
      () => {
        getBreadcrumbs()
      },
      { immediate: true }
    )

    return {
      breadcrumbs
    }
  }
}
</script>

<style scoped>
.el-breadcrumb {
  line-height: 60px;
}
</style> 