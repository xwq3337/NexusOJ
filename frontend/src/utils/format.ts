// 格式化日期
import dayjs from 'dayjs'
export const formatDate = (dateString: string) => {
  return dayjs(dateString).format('YYYY-MM-DD HH:mm:ss')
}

// 格式化内存
export const formatMemory = (memory: number) => {
  return `${Math.round(memory / 1024 / 1024 * 100) / 100} MB`
}

// 格式化时间
export const formatTime = (time: number) => {
  return `${time} ms`
}

// 格式化日期（1天前，29天前，1月前，12月前，一年前， 。。。）
export const formatRelativeTime = (dateString: string): string => {
  const now = new Date()
  const date = new Date(dateString)
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)

  // 小于1分钟
  if (diffInSeconds < 60) {
    return '刚刚'
  }

  // 分钟
  const minutes = Math.floor(diffInSeconds / 60)
  if (minutes < 60) {
    return `${minutes}分钟前`
  }

  // 小时
  const hours = Math.floor(minutes / 60)
  if (hours < 24) {
    return `${hours}小时前`
  }

  // 天
  const days = Math.floor(hours / 24)
  if (days < 30) {
    return `${days}天前`
  }

  // 月
  const months = Math.floor(days / 30)
  if (months < 12) {
    return `${months}月前`
  }

  // 年
  const years = Math.floor(months / 12)
  return `${years}年前`
}

export const formatAcceptance = (accept: number, total: number) => {
  if (total === 0) return '0.00%'
  return `${((accept / total) * 100).toFixed(2)}%`
}