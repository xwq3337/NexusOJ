// 格式化日期
import dayjs from 'dayjs'

// 格式化内存 字节转MB
export const formatMemory = (memory: number) => {
  return `${Math.round((memory / 1024 / 1024) * 100) / 100} MB`
}

// 格式化时间
export const formatTime = (time: number) => {
  return `${time} ms`
}

/**
 * 格式化日期
 * @param {number} dateString  日期字符串 2024-01-01T00:00:00.000Z
 * @returns {string} 2024-01-01 00:00:00
 */
export const formatDate = (dateString: string) => {
  return dayjs(dateString).format('YYYY-MM-DD HH:mm:ss')
}
/**
 * 格式化时间戳为日期时间字符串
 * @param {number} timestamp - 需要格式化的时间戳（毫秒级）
 * @returns {string} 格式化后的日期时间字符串，格式为 'YYYY-MM-DD HH:mm:ss'
 */
export const formateTimeStamp = (timestamp: number) => {
  return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss') // 使用 dayjs 库将时间戳格式化为指定格式的字符串
}
/***
 * 格式化相对时间
 * @param dateString  日期字符串 2024-01-01T00:00:00.000Z
 * @returns  （1天前，29天前，1月前，12月前，一年前， 。。。）
 */
export const formatRelativeTime = (dateString: string): string => {
  const now = new Date()
  const date = new Date(dateString)
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)
  const minutes = Math.floor(diffInSeconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  const months = Math.floor(days / 30)
  const years = Math.floor(months / 12)
  // 小于1分钟
  if (diffInSeconds < 60) return '刚刚'
  // 分钟
  if (minutes < 60) return `${minutes}分钟前`
  // 小时
  if (hours < 24) return `${hours}小时前`
  // 天
  if (days < 30) return `${days}天前`
  // 月
  if (months < 12) return `${months}月前`
  // 年
  return `${years}年前`
}

export const formatAcceptance = (accept: number, total: number) => {
  if (total === 0) return '0.00%'
  return `${((accept / total) * 100).toFixed(2)}%`
}
