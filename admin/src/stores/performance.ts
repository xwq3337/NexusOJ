import { defineStore } from 'pinia'
import { ApplicationInsights } from '@microsoft/applicationinsights-web'
import type { IConfiguration, IMetricTelemetry } from '@microsoft/applicationinsights-web'
import type { eventWithTime } from '@rrweb/types'

interface PerformanceState {
  webVitals: {
    cls: number
    fid: number
    lcp: number
  }
  apiMetrics: {
    requestCount: number
    averageLatency: number
    errorRate: number
  }
  userBehavior: {
    sessionCount: number
    averageSessionDuration: number
    bounceRate: number
  }
  rrwebEvents: eventWithTime[]
  isMonitoringEnabled: boolean
}

export const usePerformanceStore = defineStore('performance', {
  state: (): PerformanceState => ({
    webVitals: {
      cls: 0,
      fid: 0,
      lcp: 0
    },
    apiMetrics: {
      requestCount: 0,
      averageLatency: 0,
      errorRate: 0
    },
    userBehavior: {
      sessionCount: 0,
      averageSessionDuration: 0,
      bounceRate: 0
    },
    rrwebEvents: [],
    isMonitoringEnabled: false
  }),

  actions: {
    initializeMonitoring() {
      // 检查是否为开发环境
      const isDev = import.meta.env.DEV
      if (isDev) {
        console.warn('Performance monitoring is disabled in development mode')
        return
      }

      // 检查配置是否存在
      const instrumentationKey = import.meta.env.VITE_APP_INSIGHTS_INSTRUMENTATION_KEY
      if (!instrumentationKey) {
        console.warn('Application Insights instrumentation key is not configured')
        return
      }

      try {
        // 初始化 Application Insights
        const appInsights = new ApplicationInsights({
          config: {
            instrumentationKey,
            enableAutoRouteTracking: true,
            enableAjaxPerfTracking: true,
            enableUnhandledPromiseRejectionTracking: true
          } as IConfiguration
        })
        appInsights.loadAppInsights()
        appInsights.trackPageView()
        this.isMonitoringEnabled = true

        // 初始化 rrweb 录制
        if (typeof window !== 'undefined') {
          import('rrweb').then((rrweb) => {
            rrweb.record({
              emit: (event: eventWithTime) => {
                this.rrwebEvents.push(event)
                // 限制存储的事件数量，避免内存溢出
                if (this.rrwebEvents.length > 100) {
                  this.rrwebEvents.shift()
                }
              }
            })
          })
        }

        // 监听性能指标
        if (typeof window !== 'undefined') {
          import('web-vitals').then(({ getCLS, getFID, getLCP }) => {
            getCLS((metric) => {
              this.webVitals.cls = metric.value
              appInsights.trackMetric({
                name: 'CLS',
                average: metric.value
              } as IMetricTelemetry)
            })
            getFID((metric) => {
              this.webVitals.fid = metric.value
              appInsights.trackMetric({
                name: 'FID',
                average: metric.value
              } as IMetricTelemetry)
            })
            getLCP((metric) => {
              this.webVitals.lcp = metric.value
              appInsights.trackMetric({
                name: 'LCP',
                average: metric.value
              } as IMetricTelemetry)
            })
          })
        }
      } catch (error) {
        console.error('Failed to initialize performance monitoring:', error)
        this.isMonitoringEnabled = false
      }
    },

    updateApiMetrics(metrics: Partial<typeof this.apiMetrics>) {
      Object.assign(this.apiMetrics, metrics)
    },

    updateUserBehavior(behavior: Partial<typeof this.userBehavior>) {
      Object.assign(this.userBehavior, behavior)
    },

    clearRrwebEvents() {
      this.rrwebEvents = []
    }
  }
})
