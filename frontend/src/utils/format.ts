import dayjs from 'dayjs';

/**
 * 格式化日期
 * @param date 日期字符串或Date对象
 * @param format 格式化模板，默认为 'YYYY-MM-DD HH:mm:ss'
 * @returns 格式化后的日期字符串
 */
export const formatDate = (date: string | Date | undefined, format: string = 'YYYY-MM-DD HH:mm:ss'): string => {
  if (!date) return '';
  
  const d = new Date(date);
  if (isNaN(d.getTime())) return '';

  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  const hours = String(d.getHours()).padStart(2, '0');
  const minutes = String(d.getMinutes()).padStart(2, '0');
  const seconds = String(d.getSeconds()).padStart(2, '0');

  return format
    .replace('YYYY', String(year))
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds);
};

// 格式化日期时间
export function formatDateTime(date: string | Date | undefined): string {
  if (!date) return '';
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss');
}

// 格式化金额
export function formatAmount(amount: number): string {
  return amount.toFixed(2);
}

// 格式化百分比
export function formatPercent(value: number): string {
  return `${(value * 100).toFixed(2)}%`;
} 