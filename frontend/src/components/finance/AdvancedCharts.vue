<template>
  <div class="advanced-charts">
    <!-- 同比环比分析 -->
    <el-card class="mb-4">
      <template #header>
        <div class="card-header">
          <span>同比环比分析</span>
          <el-date-picker
            v-model="analysisDate"
            type="month"
            placeholder="选择月份"
            format="YYYY-MM"
            value-format="YYYY-MM"
            @change="updateComparisonData"
          />
        </div>
      </template>
      <div class="comparison-analysis">
        <el-row :gutter="20">
          <el-col :span="12">
            <h4>同比分析（与去年同期相比）</h4>
            <div ref="yearComparisonChartRef" class="chart-container"></div>
          </el-col>
          <el-col :span="12">
            <h4>环比分析（与上月相比）</h4>
            <div ref="monthComparisonChartRef" class="chart-container"></div>
          </el-col>
        </el-row>
      </div>
    </el-card>

    <!-- 交易趋势分析 -->
    <el-card class="mb-4">
      <template #header>
        <div class="card-header">
          <span>交易趋势分析</span>
          <div class="chart-controls">
            <el-select v-model="trendType" placeholder="图表类型" style="width: 120px">
              <el-option label="折线图" value="line" />
              <el-option label="柱状图" value="bar" />
              <el-option label="散点图" value="scatter" />
            </el-select>
            <el-select v-model="trendGroupBy" placeholder="统计周期" style="width: 120px">
              <el-option label="按日" value="day" />
              <el-option label="按周" value="week" />
              <el-option label="按月" value="month" />
            </el-select>
            <el-date-picker
              v-model="trendDateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              @change="updateTrendData"
            />
          </div>
        </div>
      </template>
      <div ref="trendChartRef" class="chart-container"></div>
    </el-card>

    <!-- 支付方式趋势 -->
    <el-card>
      <template #header>
        <div class="card-header">
          <span>支付方式趋势</span>
          <div class="chart-controls">
            <el-select v-model="paymentGroupBy" placeholder="统计周期" style="width: 120px">
              <el-option label="按日" value="day" />
              <el-option label="按周" value="week" />
              <el-option label="按月" value="month" />
            </el-select>
            <el-date-picker
              v-model="paymentDateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              @change="updatePaymentData"
            />
          </div>
        </div>
      </template>
      <div ref="paymentChartRef" class="chart-container"></div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue';
import * as echarts from 'echarts';
import { getComparisonAnalysis, getTransactionTrends, getPaymentMethodTrends } from '@/api/finance';
import type { ComparisonAnalysis, TransactionTrend, PaymentMethodTrend } from '@/types/finance';

// 图表引用
const yearComparisonChartRef = ref<HTMLElement>();
const monthComparisonChartRef = ref<HTMLElement>();
const trendChartRef = ref<HTMLElement>();
const paymentChartRef = ref<HTMLElement>();

// 图表实例
let yearComparisonChart: echarts.ECharts | null = null;
let monthComparisonChart: echarts.ECharts | null = null;
let trendChart: echarts.ECharts | null = null;
let paymentChart: echarts.ECharts | null = null;

// 控制参数
const analysisDate = ref<string>(new Date().toISOString().slice(0, 7));
const trendType = ref<'line' | 'bar' | 'scatter'>('line');
const trendGroupBy = ref<'day' | 'week' | 'month'>('month');
const trendDateRange = ref<[string, string] | null>(null);
const paymentGroupBy = ref<'day' | 'week' | 'month'>('month');
const paymentDateRange = ref<[string, string] | null>(null);

// 初始化同比环比图表
const initComparisonCharts = () => {
  if (!yearComparisonChartRef.value || !monthComparisonChartRef.value) return;
  
  yearComparisonChart = echarts.init(yearComparisonChartRef.value);
  monthComparisonChart = echarts.init(monthComparisonChartRef.value);
  
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: ['当期', '对比期', '变化率'],
      right: 10,
      top: 0
    },
    grid: {
      top: 50,
      right: 65,
      bottom: 30,
      left: 65
    },
    xAxis: {
      type: 'category',
      data: ['收入', '支出', '净收入']
    },
    yAxis: [
      {
        type: 'value',
        name: '金额',
        axisLabel: {
          formatter: '¥{value}'
        }
      },
      {
        type: 'value',
        name: '变化率',
        axisLabel: {
          formatter: '{value}%'
        }
      }
    ],
    series: [
      {
        name: '当期',
        type: 'bar',
        data: [],
        itemStyle: {
          color: '#409EFF'
        }
      },
      {
        name: '对比期',
        type: 'bar',
        data: [],
        itemStyle: {
          color: '#909399'
        }
      },
      {
        name: '变化率',
        type: 'line',
        yAxisIndex: 1,
        data: [],
        itemStyle: {
          color: '#E6A23C'
        },
        label: {
          show: true,
          formatter: '{c}%'
        }
      }
    ]
  };
  
  yearComparisonChart.setOption(option);
  monthComparisonChart.setOption(option);
  updateComparisonData();
};

// 初始化趋势图表
const initTrendChart = () => {
  if (!trendChartRef.value) return;
  
  trendChart = echarts.init(trendChartRef.value);
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      }
    },
    legend: {
      data: ['收入', '支出', '净收入', '平均金额', '交易笔数'],
      top: 10
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: []
    },
    yAxis: [
      {
        type: 'value',
        name: '金额',
        axisLabel: {
          formatter: '¥{value}'
        }
      },
      {
        type: 'value',
        name: '笔数',
        splitLine: {
          show: false
        }
      }
    ],
    series: [
      {
        name: '收入',
        type: 'line',
        data: [],
        itemStyle: {
          color: '#67C23A'
        }
      },
      {
        name: '支出',
        type: 'line',
        data: [],
        itemStyle: {
          color: '#F56C6C'
        }
      },
      {
        name: '净收入',
        type: 'line',
        data: [],
        itemStyle: {
          color: '#409EFF'
        }
      },
      {
        name: '平均金额',
        type: 'line',
        data: [],
        itemStyle: {
          color: '#E6A23C'
        }
      },
      {
        name: '交易笔数',
        type: 'bar',
        yAxisIndex: 1,
        data: [],
        itemStyle: {
          color: '#909399'
        }
      }
    ]
  };
  
  trendChart.setOption(option);
  updateTrendData();
};

// 初始化支付方式趋势图表
const initPaymentChart = () => {
  if (!paymentChartRef.value) return;
  
  paymentChart = echarts.init(paymentChartRef.value);
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: [],
      top: 10
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
    yAxis: [
      {
        type: 'value',
        name: '金额',
        axisLabel: {
          formatter: '¥{value}'
        }
      },
      {
        type: 'value',
        name: '笔数',
        splitLine: {
          show: false
        }
      }
    ],
    series: []
  };
  
  paymentChart.setOption(option);
  updatePaymentData();
};

// 更新同比环比数据
const updateComparisonData = async () => {
  if (!yearComparisonChart || !monthComparisonChart || !analysisDate.value) return;
  
  try {
    const res = await getComparisonAnalysis({
      period: 'month',
      date: analysisDate.value
    });
    
    // 更新同比数据
    const yearData = {
      current: [res.current.income, res.current.expense, res.current.netIncome],
      previous: [res.previous.income, res.previous.expense, res.previous.netIncome],
      change: [
        ((res.yearOnYear.income * 100).toFixed(2)),
        ((res.yearOnYear.expense * 100).toFixed(2)),
        ((res.yearOnYear.netIncome * 100).toFixed(2))
      ]
    };
    
    yearComparisonChart.setOption({
      series: [
        { data: yearData.current },
        { data: yearData.previous },
        { data: yearData.change }
      ]
    });
    
    // 更新环比数据
    const monthData = {
      current: [res.current.income, res.current.expense, res.current.netIncome],
      previous: [res.previous.income, res.previous.expense, res.previous.netIncome],
      change: [
        ((res.monthOnMonth.income * 100).toFixed(2)),
        ((res.monthOnMonth.expense * 100).toFixed(2)),
        ((res.monthOnMonth.netIncome * 100).toFixed(2))
      ]
    };
    
    monthComparisonChart.setOption({
      series: [
        { data: monthData.current },
        { data: monthData.previous },
        { data: monthData.change }
      ]
    });
  } catch (error) {
    console.error('获取同比环比数据失败:', error);
  }
};

// 更新趋势数据
const updateTrendData = async () => {
  if (!trendChart || !trendDateRange.value) return;
  
  try {
    const res = await getTransactionTrends({
      startDate: trendDateRange.value[0],
      endDate: trendDateRange.value[1],
      type: trendType.value,
      groupBy: trendGroupBy.value
    });
    
    const dates = res.map(item => item.date);
    const income = res.map(item => item.income);
    const expense = res.map(item => item.expense);
    const netIncome = res.map(item => item.netIncome);
    const avgAmount = res.map(item => item.avgAmount);
    const count = res.map(item => item.transactionCount);
    
    trendChart.setOption({
      xAxis: {
        data: dates
      },
      series: [
        { data: income },
        { data: expense },
        { data: netIncome },
        { data: avgAmount },
        { data: count }
      ]
    });
  } catch (error) {
    console.error('获取趋势数据失败:', error);
  }
};

// 更新支付方式数据
const updatePaymentData = async () => {
  if (!paymentChart || !paymentDateRange.value) return;
  
  try {
    const res = await getPaymentMethodTrends({
      startDate: paymentDateRange.value[0],
      endDate: paymentDateRange.value[1],
      groupBy: paymentGroupBy.value
    });
    
    // 处理数据
    const dates = Array.from(new Set(res.map(item => item.date)));
    const methods = Array.from(new Set(res.map(item => item.method)));
    const series = methods.map(method => ({
      name: method,
      type: 'bar',
      stack: 'amount',
      data: dates.map(date => {
        const item = res.find(i => i.date === date && i.method === method);
        return item ? item.amount : 0;
      })
    }));
    
    const countSeries = methods.map(method => ({
      name: `${method}(笔数)`,
      type: 'line',
      yAxisIndex: 1,
      data: dates.map(date => {
        const item = res.find(i => i.date === date && i.method === method);
        return item ? item.count : 0;
      })
    }));
    
    paymentChart.setOption({
      legend: {
        data: [...methods, ...methods.map(m => `${m}(笔数)`)]
      },
      xAxis: {
        data: dates
      },
      series: [...series, ...countSeries]
    });
  } catch (error) {
    console.error('获取支付方式趋势数据失败:', error);
  }
};

// 监听参数变化
watch([trendType, trendGroupBy], () => {
  if (trendDateRange.value) {
    updateTrendData();
  }
});

watch(paymentGroupBy, () => {
  if (paymentDateRange.value) {
    updatePaymentData();
  }
});

// 监听窗口大小变化
const handleResize = () => {
  yearComparisonChart?.resize();
  monthComparisonChart?.resize();
  trendChart?.resize();
  paymentChart?.resize();
};

// 初始化
onMounted(() => {
  initComparisonCharts();
  initTrendChart();
  initPaymentChart();
  window.addEventListener('resize', handleResize);
});

// 清理
onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  yearComparisonChart?.dispose();
  monthComparisonChart?.dispose();
  trendChart?.dispose();
  paymentChart?.dispose();
});
</script>

<style scoped>
.advanced-charts {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-controls {
  display: flex;
  gap: 10px;
  align-items: center;
}

.chart-container {
  height: 400px;
}

.mb-4 {
  margin-bottom: 16px;
}

h4 {
  margin: 0 0 16px;
  color: #606266;
}
</style> 