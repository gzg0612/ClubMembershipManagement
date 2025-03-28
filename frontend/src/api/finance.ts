import request from '@/utils/request';
import type {
  Transaction,
  CreateTransactionParams,
  UpdateTransactionStatusParams,
  TransactionQueryParams,
  TransactionListResponse,
  FinancialSummary,
  FinancialTrends,
  PaymentDistribution,
  StatusDistribution,
  FinancialReport,
  GenerateReportParams,
  ComparisonAnalysis,
  TransactionTrend,
  PaymentMethodTrend,
  ImportTemplateType,
  ExportReportParams,
  ImportResult
} from '@/types/finance';

// 获取交易记录列表
export const getTransactionList = (params: TransactionQueryParams) => {
  return request<TransactionListResponse>({
    url: '/api/v1/transactions',
    method: 'get',
    params
  });
};

// 创建交易记录
export const createTransaction = (data: CreateTransactionParams) => {
  return request<Transaction>({
    url: '/api/v1/transactions',
    method: 'post',
    data
  });
};

// 更新交易状态
export const updateTransactionStatus = (id: number, data: UpdateTransactionStatusParams) => {
  return request<void>({
    url: `/api/v1/transactions/${id}/status`,
    method: 'put',
    data
  });
};

// 获取财务汇总
export const getFinancialSummary = (params: { period: string; startTime?: string; endTime?: string }) => {
  return request<FinancialSummary>({
    url: '/api/v1/finance/summary',
    method: 'get',
    params
  });
};

// 获取财务趋势
export const getFinancialTrends = (params: { period: 'week' | 'month' | 'year' }) => {
  return request<FinancialTrends>({
    url: '/api/v1/finance/trends',
    method: 'get',
    params
  });
};

// 获取支付方式分布
export const getPaymentDistribution = () => {
  return request<{
    income: PaymentDistribution[];
    expense: PaymentDistribution[];
  }>({
    url: '/api/v1/finance/payment-distribution',
    method: 'get'
  });
};

// 获取交易状态分布
export const getTransactionStatusDistribution = () => {
  return request<StatusDistribution[]>({
    url: '/api/v1/finance/status-distribution',
    method: 'get'
  });
};

// 生成财务报表
export const generateReport = (data: GenerateReportParams) => {
  return request<FinancialReport>({
    url: '/api/v1/finance/reports',
    method: 'post',
    data
  });
};

// 获取报表列表
export const getReportList = (params: { page: number; pageSize: number }) => {
  return request<{
    list: FinancialReport[];
    total: number;
    page: number;
    pageSize: number;
  }>({
    url: '/api/v1/finance/reports',
    method: 'get',
    params
  });
};

// 下载报表
export const downloadReport = (id: number, format: 'excel' | 'csv' | 'pdf' = 'excel') => {
  return request<Blob>({
    url: `/api/v1/finance/reports/${id}/download`,
    method: 'get',
    params: { format },
    responseType: 'blob'
  });
};

// 批量导入交易记录
export const importTransactions = (file: File) => {
  const formData = new FormData();
  formData.append('file', file);
  
  return request<{
    success: number;
    failed: number;
    errors: Array<{
      row: number;
      message: string;
    }>;
  }>({
    url: '/api/v1/transactions/import',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
};

// 获取同比环比分析数据
export const getComparisonAnalysis = (params: { period: string; date: string }) => {
  return request<ComparisonAnalysis>({
    url: '/api/v1/finance/comparison',
    method: 'get',
    params
  });
};

// 获取交易趋势分析
export const getTransactionTrends = (params: { 
  startDate: string;
  endDate: string;
  type?: 'line' | 'bar' | 'scatter';
  groupBy: 'day' | 'week' | 'month';
}) => {
  return request<TransactionTrend[]>({
    url: '/api/v1/finance/transaction-trends',
    method: 'get',
    params
  });
};

// 获取支付方式趋势
export const getPaymentMethodTrends = (params: {
  startDate: string;
  endDate: string;
  groupBy: 'day' | 'week' | 'month';
}) => {
  return request<PaymentMethodTrend[]>({
    url: '/api/v1/finance/payment-trends',
    method: 'get',
    params
  });
};

// 下载导入模板
export const downloadImportTemplate = (type: typeof ImportTemplateType[keyof typeof ImportTemplateType]) => {
  return request<Blob>({
    url: '/api/v1/finance/import-template',
    method: 'get',
    params: { type },
    responseType: 'blob'
  });
};

// 导出报表（支持多种格式和模板）
export const exportReport = (params: ExportReportParams) => {
  return request<Blob>({
    url: '/api/v1/finance/export',
    method: 'post',
    data: params,
    responseType: 'blob'
  });
};

// 批量导入交易记录（支持预览）
export const previewImport = (file: File) => {
  const formData = new FormData();
  formData.append('file', file);
  
  return request<{
    headers: string[];
    preview: Array<Record<string, any>>;
    mapping: Record<string, string>;
  }>({
    url: '/api/v1/finance/import/preview',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
};

// 确认导入
export const confirmImport = (params: {
  fileId: string;
  mapping: Record<string, string>;
  skipErrors?: boolean;
}) => {
  return request<ImportResult>({
    url: '/api/v1/finance/import/confirm',
    method: 'post',
    data: params
  });
}; 