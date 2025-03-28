<template>
  <div class="financial-charts">
    <!-- 收入支出趋势图 -->
    <el-card class="mb-4">
      <template #header>
        <div class="card-header">
          <span>收入支出趋势</span>
          <el-radio-group v-model="trendPeriod" size="small">
            <el-radio-button label="week">周</el-radio-button>
            <el-radio-button label="month">月</el-radio-button>
            <el-radio-button label="year">年</el-radio-button>
          </el-radio-group>
        </div>
      </template>
      <div ref="trendChartRef" class="chart-container"></div>
    </el-card>

    <!-- 支付方式分布 -->
    <el-row :gutter="20" class="mb-4">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>收入支付方式分布</span>
            </div>
          </template>
          <div ref="incomePaymentChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>支出支付方式分布</span>
            </div>
          </template>
          <div ref="expensePaymentChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 交易状态分布 -->
    <el-card>
      <template #header>
        <div class="card-header">
          <span>交易状态分布</span>
        </div>
      </template>
      <div ref="statusChartRef" class="chart-container"></div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue';
import * as echarts from 'echarts';
import { getFinancialTrends, getPaymentDistribution, getTransactionStatusDistribution } from '@/api/finance';
import type { FinancialTrends, PaymentDistribution, StatusDistribution } from '@/types/finance';

// 图表引用
const trendChartRef = ref<HTMLElement>();
const incomePaymentChartRef = ref<HTMLElement>();
const expensePaymentChartRef = ref<HTMLElement>();
const statusChartRef = ref<HTMLElement>();

// 图表实例
let trendChart: echarts.ECharts | null = null;
let incomePaymentChart: echarts.ECharts | null = null;
let expensePaymentChart: echarts.ECharts | null = null;
let statusChart: echarts.ECharts | null = null;

// 趋势周期
const trendPeriod = ref<'week' | 'month' | 'year'>('month');

// 初始化趋势图
const initTrendChart = async () => {
  if (!trendChartRef.value) return;
  
  trendChart = echarts.init(trendChartRef.value);
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: ['收入', '支出']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: []
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: '¥{value}'
      }
    },
    series: [
      {
        name: '收入',
        type: 'bar',
        data: [],
        itemStyle: {
          color: '#67c23a'
        }
      },
      {
        name: '支出',
        type: 'bar',
        data: [],
        itemStyle: {
          color: '#f56c6c'
        }
      }
    ]
  };
  
  trendChart.setOption(option);
  await updateTrendData();
};

// 初始化支付方式分布图
const initPaymentCharts = async () => {
  if (!incomePaymentChartRef.value || !expensePaymentChartRef.value) return;
  
  incomePaymentChart = echarts.init(incomePaymentChartRef.value);
  expensePaymentChart = echarts.init(expensePaymentChartRef.value);
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [
      {
        type: 'pie',
        radius: '50%',
        data: [],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  };
  
  incomePaymentChart.setOption(option);
  expensePaymentChart.setOption(option);
  await updatePaymentData();
};

// 初始化状态分布图
const initStatusChart = async () => {
  if (!statusChartRef.value) return;
  
  statusChart = echarts.init(statusChartRef.value);
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [
      {
        type: 'pie',
        radius: '50%',
        data: [],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  };
  
  statusChart.setOption(option);
  await updateStatusData();
};

// 更新趋势数据
const updateTrendData = async () => {
  if (!trendChart) return;
  
  try {
    const res = await getFinancialTrends({ period: trendPeriod.value });
    const { dates, income, expense } = res.data;
    
    trendChart.setOption({
      xAxis: {
        data: dates
      },
      series: [
        {
          name: '收入',
          data: income
        },
        {
          name: '支出',
          data: expense
        }
      ]
    });
  } catch (error) {
    console.error('获取趋势数据失败:', error);
  }
};

// 更新支付方式分布数据
const updatePaymentData = async () => {
  if (!incomePaymentChart || !expensePaymentChart) return;
  
  try {
    const res = await getPaymentDistribution();
    const { income, expense } = res.data;
    
    incomePaymentChart.setOption({
      series: [{
        data: income
      }]
    });
    
    expensePaymentChart.setOption({
      series: [{
        data: expense
      }]
    });
  } catch (error) {
    console.error('获取支付方式分布数据失败:', error);
  }
};

// 更新状态分布数据
const updateStatusData = async () => {
  if (!statusChart) return;
  
  try {
    const res = await getTransactionStatusDistribution();
    statusChart.setOption({
      series: [{
        data: res.data
      }]
    });
  } catch (error) {
    console.error('获取状态分布数据失败:', error);
  }
};

// 监听趋势周期变化
watch(trendPeriod, () => {
  updateTrendData();
});

// 监听窗口大小变化
const handleResize = () => {
  trendChart?.resize();
  incomePaymentChart?.resize();
  expensePaymentChart?.resize();
  statusChart?.resize();
};

// 初始化
onMounted(() => {
  initTrendChart();
  initPaymentCharts();
  initStatusChart();
  window.addEventListener('resize', handleResize);
});

// 清理
onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  trendChart?.dispose();
  incomePaymentChart?.dispose();
  expensePaymentChart?.dispose();
  statusChart?.dispose();
});
</script>

<style scoped>
.financial-charts {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  height: 400px;
}

.mb-4 {
  margin-bottom: 16px;
}
</style> 