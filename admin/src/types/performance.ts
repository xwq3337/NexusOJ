export interface PerformanceMetrics {
  fcp: number // First Contentful Paint
  lcp: number // Largest Contentful Paint
  fid: number // First Input Delay
  cls: number // Cumulative Layout Shift
  ttfb: number // Time to First Byte
}

export interface UserBehavior {
  timestamp: number
  type: 'click' | 'input' | 'navigation' | 'error'
  target: string
  details?: {
    path?: string
    value?: string | number
    error?: string
    duration?: number
  }
}

export interface ApiMetrics {
  totalRequests: number
  successRequests: number
  failedRequests: number
  averageResponseTime: number
}
