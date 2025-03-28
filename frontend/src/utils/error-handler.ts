import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'

// 错误类型
export enum ErrorType {
  NETWORK = 'NETWORK',
  AUTH = 'AUTH',
  VALIDATION = 'VALIDATION',
  BUSINESS = 'BUSINESS',
  UNKNOWN = 'UNKNOWN'
}

// 错误信息映射
const errorMessages = {
  [ErrorType.NETWORK]: '网络连接失败，请检查网络设置',
  [ErrorType.AUTH]: '登录已过期，请重新登录',
  [ErrorType.VALIDATION]: '输入数据验证失败',
  [ErrorType.BUSINESS]: '操作失败',
  [ErrorType.UNKNOWN]: '系统错误，请稍后重试'
}

// 处理错误
export const handleError = (error: any, type: ErrorType = ErrorType.UNKNOWN) => {
  console.error('错误详情:', error)

  // 获取错误信息
  const message = error.response?.data?.message || error.message || errorMessages[type]

  // 显示错误提示
  ElMessage.error(message)

  // 处理认证错误
  if (type === ErrorType.AUTH) {
    const router = useRouter()
    // 清除用户信息
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    // 跳转到登录页
    router.push('/login')
  }

  return Promise.reject(error)
}

// 处理业务错误
export const handleBusinessError = (error: any) => {
  return handleError(error, ErrorType.BUSINESS)
}

// 处理网络错误
export const handleNetworkError = (error: any) => {
  return handleError(error, ErrorType.NETWORK)
}

// 处理认证错误
export const handleAuthError = (error: any) => {
  return handleError(error, ErrorType.AUTH)
}

// 处理验证错误
export const handleValidationError = (error: any) => {
  return handleError(error, ErrorType.VALIDATION)
} 