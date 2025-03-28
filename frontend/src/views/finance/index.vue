<template>
  <div class="finance-container">
    <!-- 财务概览卡片 -->
    <el-row :gutter="20" class="mb-4">
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>总收入</span>
            </div>
          </template>
          <div class="card-content">
            <span class="amount income">¥{{ summary.totalIncome.toFixed(2) }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>总支出</span>
            </div>
          </template>
          <div class="card-content">
            <span class="amount expense">¥{{ summary.totalExpense.toFixed(2) }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>净收入</span>
            </div>
          </template>
          <div class="card-content">
            <span class="amount" :class="{ income: summary.netIncome >= 0, expense: summary.netIncome < 0 }">
              ¥{{ summary.netIncome.toFixed(2) }}
            </span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 财务分析图表 -->
    <FinancialCharts class="mb-4" />

    <!-- 交易记录列表 -->
    <el-card class="mb-4">
      <template #header>
        <div class="card-header">
          <span>交易记录</span>
          <div class="header-actions">
            <el-button type="success" @click="showImportDialog">批量导入</el-button>
            <el-button type="primary" @click="showCreateDialog">新建交易</el-button>
            <el-button type="info" @click="showExportDialog">导出报表</el-button>
          </div>
        </div>
      </template>

      <!-- 搜索表单 -->
      <el-form :inline="true" :model="queryParams" class="mb-4">
        <el-form-item label="交易类型">
          <el-select v-model="queryParams.type" placeholder="请选择" clearable>
            <el-option label="收入" :value="TransactionType.INCOME" />
            <el-option label="支出" :value="TransactionType.EXPENSE" />
          </el-select>
        </el-form-item>
        <el-form-item label="交易状态">
          <el-select v-model="queryParams.status" placeholder="请选择" clearable>
            <el-option label="待处理" :value="TransactionStatus.PENDING" />
            <el-option label="已完成" :value="TransactionStatus.COMPLETED" />
            <el-option label="失败" :value="TransactionStatus.FAILED" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            @change="handleDateRangeChange"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="transactionList" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === TransactionType.INCOME ? 'success' : 'danger'">
              {{ row.type === TransactionType.INCOME ? '收入' : '支出' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            <span :class="{ income: row.type === TransactionType.INCOME, expense: row.type === TransactionType.EXPENSE }">
              {{ row.type === TransactionType.INCOME ? '+' : '-' }}¥{{ row.amount.toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="paymentMethod" label="支付方式" width="120" />
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="showDetail(row)">详情</el-button>
            <el-button
              v-if="row.status === TransactionStatus.PENDING"
              type="success"
              link
              @click="showUpdateStatusDialog(row)"
            >
              更新状态
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 新建交易对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="新建交易"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        label-width="100px"
      >
        <el-form-item label="交易类型" prop="type">
          <el-select v-model="createForm.type" placeholder="请选择">
            <el-option label="收入" :value="TransactionType.INCOME" />
            <el-option label="支出" :value="TransactionType.EXPENSE" />
          </el-select>
        </el-form-item>
        <el-form-item label="金额" prop="amount">
          <el-input-number
            v-model="createForm.amount"
            :precision="2"
            :step="0.01"
            :min="0"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="支付方式" prop="paymentMethod">
          <el-select v-model="createForm.paymentMethod" placeholder="请选择">
            <el-option label="现金" value="cash" />
            <el-option label="微信支付" value="wechat" />
            <el-option label="支付宝" value="alipay" />
            <el-option label="银行卡" value="bank" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="createForm.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleCreate" :loading="creating">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 更新状态对话框 -->
    <el-dialog
      v-model="updateStatusDialogVisible"
      title="更新交易状态"
      width="400px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="updateStatusFormRef"
        :model="updateStatusForm"
        :rules="updateStatusRules"
        label-width="100px"
      >
        <el-form-item label="状态" prop="status">
          <el-select v-model="updateStatusForm.status" placeholder="请选择">
            <el-option label="待处理" :value="TransactionStatus.PENDING" />
            <el-option label="已完成" :value="TransactionStatus.COMPLETED" />
            <el-option label="失败" :value="TransactionStatus.FAILED" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="updateStatusDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleUpdateStatus" :loading="updating">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 交易详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="交易详情"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="交易ID">{{ currentTransaction?.id }}</el-descriptions-item>
        <el-descriptions-item label="交易类型">
          <el-tag :type="currentTransaction?.type === TransactionType.INCOME ? 'success' : 'danger'">
            {{ currentTransaction?.type === TransactionType.INCOME ? '收入' : '支出' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="金额">
          <span :class="{ income: currentTransaction?.type === TransactionType.INCOME, expense: currentTransaction?.type === TransactionType.EXPENSE }">
            {{ currentTransaction?.type === TransactionType.INCOME ? '+' : '-' }}¥{{ currentTransaction?.amount.toFixed(2) }}
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentTransaction?.status)">
            {{ getStatusText(currentTransaction?.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="支付方式">{{ currentTransaction?.paymentMethod }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentTransaction?.createdAt) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentTransaction?.updatedAt) }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ currentTransaction?.remark }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!-- 导入对话框 -->
    <el-dialog
      v-model="importDialogVisible"
      title="批量导入交易记录"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-upload
        class="upload-demo"
        drag
        action="/api/v1/transactions/import"
        :headers="{ Authorization: `Bearer ${token}` }"
        :on-success="handleImportSuccess"
        :on-error="handleImportError"
        :before-upload="beforeImportUpload"
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            只能上传 xlsx/xls 文件，且文件大小不超过 10MB
          </div>
        </template>
      </el-upload>
    </el-dialog>

    <!-- 导出对话框 -->
    <el-dialog
      v-model="exportDialogVisible"
      title="导出财务报表"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="exportFormRef"
        :model="exportForm"
        :rules="exportRules"
        label-width="100px"
      >
        <el-form-item label="报表类型" prop="type">
          <el-select v-model="exportForm.type" placeholder="请选择">
            <el-option label="日报表" value="daily" />
            <el-option label="月报表" value="monthly" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围" prop="dateRange">
          <el-date-picker
            v-model="exportForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="exportDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleExport" :loading="exporting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import type { FormInstance } from 'element-plus';
import { UploadFilled } from '@element-plus/icons-vue';
import { TransactionType, TransactionStatus } from '@/types/finance';
import {
  createTransaction,
  getTransactionList,
  getFinancialSummary,
  updateTransactionStatus,
  generateReport,
  downloadReport
} from '@/api/finance';
import type {
  Transaction,
  CreateTransactionParams,
  UpdateTransactionStatusParams,
  TransactionQueryParams,
  GenerateReportParams
} from '@/types/finance';
import { formatDate } from '@/utils/format';
import FinancialCharts from '@/components/finance/FinancialCharts.vue';
import { useUserStore } from '@/stores/user';

// 状态定义
const loading = ref(false);
const creating = ref(false);
const updating = ref(false);
const exporting = ref(false);
const transactionList = ref<Transaction[]>([]);
const total = ref(0);
const summary = reactive({
  totalIncome: 0,
  totalExpense: 0,
  netIncome: 0
});

// 查询参数
const queryParams = reactive<TransactionQueryParams>({
  page: 1,
  pageSize: 10
});
const dateRange = ref<[string, string] | null>(null);

// 对话框控制
const createDialogVisible = ref(false);
const updateStatusDialogVisible = ref(false);
const detailDialogVisible = ref(false);
const importDialogVisible = ref(false);
const exportDialogVisible = ref(false);
const currentTransaction = ref<Transaction | null>(null);

// 表单引用
const createFormRef = ref<FormInstance>();
const updateStatusFormRef = ref<FormInstance>();
const exportFormRef = ref<FormInstance>();

// 表单数据
const createForm = reactive<CreateTransactionParams>({
  type: TransactionType.INCOME,
  amount: 0,
  paymentMethod: '',
  remark: ''
});

const updateStatusForm = reactive<UpdateTransactionStatusParams>({
  status: TransactionStatus.PENDING
});

const exportForm = reactive({
  type: 'daily' as 'daily' | 'monthly',
  dateRange: null as [string, string] | null
});

// 表单验证规则
const createRules = {
  type: [{ required: true, message: '请选择交易类型', trigger: 'change' }],
  amount: [{ required: true, message: '请输入金额', trigger: 'blur' }],
  paymentMethod: [{ required: true, message: '请选择支付方式', trigger: 'change' }]
};

const updateStatusRules = {
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
};

const exportRules = {
  type: [{ required: true, message: '请选择报表类型', trigger: 'change' }],
  dateRange: [{ required: true, message: '请选择时间范围', trigger: 'change' }]
};

// 用户 store
const userStore = useUserStore();
const token = userStore.token;

// 获取交易列表
const fetchTransactionList = async () => {
  loading.value = true;
  try {
    const res = await getTransactionList(queryParams);
    transactionList.value = res.data.list;
    total.value = res.data.total;
  } catch (error) {
    console.error('获取交易列表失败:', error);
  } finally {
    loading.value = false;
  }
};

// 获取财务汇总
const fetchFinancialSummary = async () => {
  try {
    const res = await getFinancialSummary({
      period: 'month',
      startTime: queryParams.startTime,
      endTime: queryParams.endTime
    });
    Object.assign(summary, res.data);
  } catch (error) {
    console.error('获取财务汇总失败:', error);
  }
};

// 显示创建对话框
const showCreateDialog = () => {
  createDialogVisible.value = true;
  Object.assign(createForm, {
    type: TransactionType.INCOME,
    amount: 0,
    paymentMethod: '',
    remark: ''
  });
};

// 显示更新状态对话框
const showUpdateStatusDialog = (row: Transaction) => {
  currentTransaction.value = row;
  updateStatusForm.status = row.status;
  updateStatusDialogVisible.value = true;
};

// 显示详情对话框
const showDetail = (row: Transaction) => {
  currentTransaction.value = row;
  detailDialogVisible.value = true;
};

// 显示导入对话框
const showImportDialog = () => {
  importDialogVisible.value = true;
};

// 显示导出对话框
const showExportDialog = () => {
  exportDialogVisible.value = true;
  Object.assign(exportForm, {
    type: 'daily',
    dateRange: null
  });
};

// 处理创建交易
const handleCreate = async () => {
  if (!createFormRef.value) return;
  
  await createFormRef.value.validate(async (valid) => {
    if (valid) {
      creating.value = true;
      try {
        await createTransaction(createForm);
        ElMessage.success('创建成功');
        createDialogVisible.value = false;
        fetchTransactionList();
        fetchFinancialSummary();
      } catch (error) {
        console.error('创建交易失败:', error);
      } finally {
        creating.value = false;
      }
    }
  });
};

// 处理更新状态
const handleUpdateStatus = async () => {
  if (!updateStatusFormRef.value || !currentTransaction.value) return;

  await updateStatusFormRef.value.validate(async (valid) => {
    if (valid) {
      updating.value = true;
      try {
        await updateTransactionStatus(currentTransaction.value.id, updateStatusForm);
        ElMessage.success('更新成功');
        updateStatusDialogVisible.value = false;
        fetchTransactionList();
        fetchFinancialSummary();
      } catch (error) {
        console.error('更新状态失败:', error);
      } finally {
        updating.value = false;
      }
    }
  });
};

// 处理导入
const handleImportSuccess = (response: any) => {
  ElMessage.success('导入成功');
  importDialogVisible.value = false;
  fetchTransactionList();
  fetchFinancialSummary();
};

const handleImportError = (error: any) => {
  ElMessage.error('导入失败');
  console.error('导入失败:', error);
};

const beforeImportUpload = (file: File) => {
  const isExcel = file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' || 
                  file.type === 'application/vnd.ms-excel';
  const isLt10M = file.size / 1024 / 1024 < 10;

  if (!isExcel) {
    ElMessage.error('只能上传 Excel 文件!');
    return false;
  }
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB!');
    return false;
  }
  return true;
};

// 处理导出
const handleExport = async () => {
  if (!exportFormRef.value || !exportForm.dateRange) return;

  await exportFormRef.value.validate(async (valid) => {
    if (valid) {
      exporting.value = true;
      try {
        const params: GenerateReportParams = {
          startDate: exportForm.dateRange[0],
          endDate: exportForm.dateRange[1],
          type: exportForm.type
        };
        
        const res = await generateReport(params);
        const blob = await downloadReport(res.data.id);
        
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `财务报表_${exportForm.dateRange[0]}_${exportForm.dateRange[1]}.xlsx`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);
        
        ElMessage.success('导出成功');
        exportDialogVisible.value = false;
      } catch (error) {
        console.error('导出失败:', error);
      } finally {
        exporting.value = false;
      }
    }
  });
};

// 处理搜索
const handleSearch = () => {
  queryParams.page = 1;
  fetchTransactionList();
  fetchFinancialSummary();
};

// 处理重置
const handleReset = () => {
  Object.assign(queryParams, {
    page: 1,
    pageSize: 10,
    type: undefined,
    status: undefined,
    startTime: undefined,
    endTime: undefined
  });
  dateRange.value = null;
  fetchTransactionList();
  fetchFinancialSummary();
};

// 处理日期范围变化
const handleDateRangeChange = (val: [string, string] | null) => {
  if (val) {
    queryParams.startTime = val[0];
    queryParams.endTime = val[1];
  } else {
    queryParams.startTime = undefined;
    queryParams.endTime = undefined;
  }
};

// 处理分页
const handleSizeChange = (val: number) => {
  queryParams.pageSize = val;
  fetchTransactionList();
};

const handleCurrentChange = (val: number) => {
  queryParams.page = val;
  fetchTransactionList();
};

// 获取状态类型
const getStatusType = (status?: string) => {
  switch (status) {
    case TransactionStatus.PENDING:
      return 'warning';
    case TransactionStatus.COMPLETED:
      return 'success';
    case TransactionStatus.FAILED:
      return 'danger';
    default:
      return 'info';
  }
};

// 获取状态文本
const getStatusText = (status?: string) => {
  switch (status) {
    case TransactionStatus.PENDING:
      return '待处理';
    case TransactionStatus.COMPLETED:
      return '已完成';
    case TransactionStatus.FAILED:
      return '失败';
    default:
      return '未知';
  }
};

// 初始化
onMounted(() => {
  fetchTransactionList();
  fetchFinancialSummary();
});
</script>

<style scoped>
.finance-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.card-content {
  text-align: center;
  padding: 20px 0;
}

.amount {
  font-size: 24px;
  font-weight: bold;
}

.income {
  color: #67c23a;
}

.expense {
  color: #f56c6c;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.mb-4 {
  margin-bottom: 16px;
}

.upload-demo {
  text-align: center;
}

.el-upload__tip {
  margin-top: 10px;
  color: #909399;
}
</style> 