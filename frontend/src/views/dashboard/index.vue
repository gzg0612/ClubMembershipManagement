<template>
  <div class="dashboard">
    <h1 class="text-2xl font-bold mb-6">欢迎使用会员管理系统</h1>
    
    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
      <el-card shadow="hover" v-loading="loading.stats">
        <template #header>
          <div class="flex items-center">
            <el-icon class="mr-2"><User /></el-icon>
            <span>会员总数</span>
          </div>
        </template>
        <div class="text-3xl font-bold text-primary">{{ stats.memberCount }}</div>
      </el-card>
      
      <el-card shadow="hover" v-loading="loading.stats">
        <template #header>
          <div class="flex items-center">
            <el-icon class="mr-2"><Shop /></el-icon>
            <span>店铺总数</span>
          </div>
        </template>
        <div class="text-3xl font-bold text-primary">{{ stats.storeCount }}</div>
      </el-card>
      
      <el-card shadow="hover" v-loading="loading.stats">
        <template #header>
          <div class="flex items-center">
            <el-icon class="mr-2"><Money /></el-icon>
            <span>今日营业额</span>
          </div>
        </template>
        <div class="text-3xl font-bold text-primary">¥{{ stats.todayRevenue }}</div>
      </el-card>
      
      <el-card shadow="hover" v-loading="loading.stats">
        <template #header>
          <div class="flex items-center">
            <el-icon class="mr-2"><Calendar /></el-icon>
            <span>今日新增会员</span>
          </div>
        </template>
        <div class="text-3xl font-bold text-primary">{{ stats.todayNewMembers }}</div>
      </el-card>
    </div>

    <!-- 图表区域 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <el-card shadow="hover" v-loading="loading.memberChart">
        <template #header>
          <div class="flex items-center justify-between">
            <span>会员增长趋势</span>
            <el-radio-group v-model="memberChartType" size="small" @change="fetchMemberTrend">
              <el-radio-button label="week">周</el-radio-button>
              <el-radio-button label="month">月</el-radio-button>
              <el-radio-button label="year">年</el-radio-button>
            </el-radio-group>
          </div>
        </template>
        <div ref="memberChartRef" class="h-80"></div>
      </el-card>

      <el-card shadow="hover" v-loading="loading.revenueChart">
        <template #header>
          <div class="flex items-center justify-between">
            <span>营业额统计</span>
            <el-radio-group v-model="revenueChartType" size="small" @change="fetchRevenueTrend">
              <el-radio-button label="week">周</el-radio-button>
              <el-radio-button label="month">月</el-radio-button>
              <el-radio-button label="year">年</el-radio-button>
            </el-radio-group>
          </div>
        </template>
        <div ref="revenueChartRef" class="h-80"></div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { User, Shop, Money, Calendar } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { dashboardApi } from '@/api/dashboard'
import type { DashboardStats, ChartData } from '@/api/dashboard'

// 统计数据
const stats = reactive<DashboardStats>({
  memberCount: 0,
  storeCount: 0,
  todayRevenue: 0,
  todayNewMembers: 0
})

// 加载状态
const loading = reactive({
  stats: false,
  memberChart: false,
  revenueChart: false
})

// 图表类型
const memberChartType = ref<'week' | 'month' | 'year'>('week')
const revenueChartType = ref<'week' | 'month' | 'year'>('week')

// 图表实例
const memberChartRef = ref<HTMLElement>()
const revenueChartRef = ref<HTMLElement>()
let memberChart: echarts.ECharts | null = null
let revenueChart: echarts.ECharts | null = null

// 初始化图表
const initCharts = () => {
  if (memberChartRef.value) {
    memberChart = echarts.init(memberChartRef.value)
  }
  if (revenueChartRef.value) {
    revenueChart = echarts.init(revenueChartRef.value)
  }
}

// 更新会员趋势图表
const updateMemberChart = (data: ChartData[]) => {
  if (!memberChart) return

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.date)
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '会员数量',
        type: 'line',
        smooth: true,
        data: data.map(item => item.value),
        areaStyle: {
          opacity: 0.1
        }
      }
    ]
  }

  memberChart.setOption(option)
}

// 更新营业额图表
const updateRevenueChart = (data: ChartData[]) => {
  if (!revenueChart) return

  const option = {
    tooltip: {
      trigger: 'axis',
      formatter: '{b}: ¥{c}'
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.date)
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: '¥{value}'
      }
    },
    series: [
      {
        name: '营业额',
        type: 'bar',
        data: data.map(item => item.value),
        itemStyle: {
          color: '#409EFF'
        }
      }
    ]
  }

  revenueChart.setOption(option)
}

// 获取统计数据
const fetchStats = async () => {
  loading.stats = true
  try {
    const { data } = await dashboardApi.getStats()
    Object.assign(stats, data)
  } catch (error) {
    console.error('获取统计数据失败:', error)
  } finally {
    loading.stats = false
  }
}

// 获取会员趋势数据
const fetchMemberTrend = async () => {
  loading.memberChart = true
  try {
    const { data } = await dashboardApi.getMemberTrend(memberChartType.value)
    updateMemberChart(data)
  } catch (error) {
    console.error('获取会员趋势数据失败:', error)
  } finally {
    loading.memberChart = false
  }
}

// 获取营业额趋势数据
const fetchRevenueTrend = async () => {
  loading.revenueChart = true
  try {
    const { data } = await dashboardApi.getRevenueTrend(revenueChartType.value)
    updateRevenueChart(data)
  } catch (error) {
    console.error('获取营业额趋势数据失败:', error)
  } finally {
    loading.revenueChart = false
  }
}

// 监听窗口大小变化
const handleResize = () => {
  memberChart?.resize()
  revenueChart?.resize()
}

// 页面加载时初始化
onMounted(() => {
  initCharts()
  fetchStats()
  fetchMemberTrend()
  fetchRevenueTrend()
  window.addEventListener('resize', handleResize)
})

// 页面卸载时清理
onUnmounted(() => {
  memberChart?.dispose()
  revenueChart?.dispose()
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}
</style> 