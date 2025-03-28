// 交易类型
export const TransactionType = {
  INCOME: 'income',  // 收入
  EXPENSE: 'expense' // 支出
} as const;

// 交易状态
export const TransactionStatus = {
  PENDING: 'pending',   // 待处理
  COMPLETED: 'completed', // 已完成
  FAILED: 'failed'      // 失败
} as const;

// 交易记录
export interface Transaction {
  id: number;
  orderId?: number;
  type: typeof TransactionType[keyof typeof TransactionType];
  amount: number;
  status: typeof TransactionStatus[keyof typeof TransactionStatus];
  paymentMethod: string;
  remark: string;
  createdAt: string;
  updatedAt: string;
  order?: Order;
}

// 创建交易记录参数
export interface CreateTransactionParams {
  orderId?: number;
  type: typeof TransactionType[keyof typeof TransactionType];
  amount: number;
  paymentMethod: string;
  remark?: string;
}

// 更新交易状态参数
export interface UpdateTransactionStatusParams {
  status: typeof TransactionStatus[keyof typeof TransactionStatus];
}

// 交易记录查询参数
export interface TransactionQueryParams {
  page?: number;
  pageSize?: number;
  type?: typeof TransactionType[keyof typeof TransactionType];
  status?: typeof TransactionStatus[keyof typeof TransactionStatus];
  startTime?: string;
  endTime?: string;
}

// 交易记录列表响应
export interface TransactionListResponse {
  list: Transaction[];
  total: number;
  page: number;
  pageSize: number;
}

// 财务汇总
export interface FinancialSummary {
  totalIncome: number;
  totalExpense: number;
  netIncome: number;
  period: string;
}

// 财务趋势
export interface FinancialTrends {
  dates: string[];
  income: number[];
  expense: number[];
}

// 支付方式分布
export interface PaymentDistribution {
  name: string;
  value: number;
  percentage: number;
}

// 交易状态分布
export interface StatusDistribution {
  name: string;
  value: number;
  percentage: number;
}

// 财务报表
export interface FinancialReport {
  id: number;
  startDate: string;
  endDate: string;
  type: string;
  data: string;
  createdAt: string;
  updatedAt: string;
}

// 生成财务报表参数
export interface GenerateReportParams {
  startDate: string;
  endDate: string;
  type: 'daily' | 'monthly';
}

// 订单类型
export interface Order {
  id: number;
  orderNo: string;
  memberId: number;
  storeId: number;
  totalAmount: number;
  status: string;
  createdAt: string;
  updatedAt: string;
}

// 同比环比分析数据
export interface ComparisonAnalysis {
  current: {
    period: string;
    income: number;
    expense: number;
    netIncome: number;
  };
  previous: {
    period: string;
    income: number;
    expense: number;
    netIncome: number;
  };
  yearOnYear: {
    income: number;
    expense: number;
    netIncome: number;
  };
  monthOnMonth: {
    income: number;
    expense: number;
    netIncome: number;
  };
}

// 交易趋势分析
export interface TransactionTrend {
  date: string;
  income: number;
  expense: number;
  netIncome: number;
  avgAmount: number;
  transactionCount: number;
}

// 支付方式趋势
export interface PaymentMethodTrend {
  date: string;
  method: string;
  amount: number;
  count: number;
}

// 导出模板类型
export const ExportTemplateType = {
  BASIC: 'basic',           // 基础模板
  DETAILED: 'detailed',     // 详细模板
  ANALYSIS: 'analysis'      // 分析模板
} as const;

// 导出格式
export const ExportFormat = {
  EXCEL: 'excel',
  CSV: 'csv',
  PDF: 'pdf'
} as const;

// 导出报表参数
export interface ExportReportParams extends GenerateReportParams {
  format: typeof ExportFormat[keyof typeof ExportFormat];
  template: typeof ExportTemplateType[keyof typeof ExportTemplateType];
}

// 导入模板类型
export enum ImportTemplateType {
  BASIC = 'basic',
  DETAILED = 'detailed'
}

export interface ImportPreviewResponse {
  fileId: string;
  headers: string[];
  mapping: Record<string, string>;
  preview: Array<Record<string, any>>;
  errors?: Array<{
    row: number;
    message: string;
  }>;
}

export interface ImportConfirmResponse {
  success: number;
  failed: number;
}

export interface ImportConfirmParams {
  fileId: string;
  mapping: Record<string, string>;
  skipErrors: boolean;
}

// 导入结果
export interface ImportResult {
  success: number;
  failed: number;
  total: number;
  errors: Array<{
    row: number;
    message: string;
    data?: Record<string, any>;
  }>;
  successRecords: Transaction[];
} 