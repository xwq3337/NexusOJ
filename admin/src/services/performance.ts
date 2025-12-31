import { onLCP, onFID, onCLS, onTTFB } from 'web-vitals'
import { usePerformanceStore } from '../stores/performance'
import type { UserBehavior } from '../types/performance'
import * as rrweb from 'rrweb'
import { init as initApm } from '@elastic/apm-rum'

class PerformanceMonitor {
  private store = usePerformanceStore()
  private apm: any

  constructor() {
    this.initWebVitals()
    this.initUserBehaviorTracking()
    this.initApm()
  }

  private initWebVitals() {
    onLCP(metric => {
      this.store.updatePerformanceMetrics({ lcp: metric.value })
    })

    onFID(metric => {
      this.store.updatePerformanceMetrics({ fid: metric.value })
    })

    onCLS(metric => {
      this.store.updatePerformanceMetrics({ cls: metric.value })
    })

    onTTFB(metric => {
      this.store.updatePerformanceMetrics({ ttfb: metric.value })
    })
  }

  private initUserBehaviorTracking() {
    // 记录用户点击
    document.addEventListener('click', (e) => {
      const target = e.target as HTMLElement
      const behavior: UserBehavior = {
        timestamp: Date.now(),
        type: 'click',
        target: target.tagName.toLowerCase(),
        details: {
          path: this.getElementPath(target)
        }
      }
      this.store.addUserBehavior(behavior)
    })

    // 记录页面导航
    window.addEventListener('popstate', () => {
      const behavior: UserBehavior = {
        timestamp: Date.now(),
        type: 'navigation',
        target: window.location.pathname,
        details: {
          path: window.location.pathname
        }
      }
      this.store.addUserBehavior(behavior)
    })

    // 记录错误
    window.addEventListener('error', (e) => {
      const behavior: UserBehavior = {
        timestamp: Date.now(),
        type: 'error',
        target: 'window',
        details: {
          error: e.message
        }
      }
      this.store.addUserBehavior(behavior)
    })

    // 初始化 rrweb 录制
    rrweb.record({
      emit: (event) => {
        // 这里可以将事件发送到后端存储
        console.log('rrweb event', event)
      }
    })
  }

  private initApm() {
    this.apm = initApm({
      serviceName: 'oj-admin',
      serverUrl: 'http://your-apm-server-url',
      environment: import.meta.env.MODE
    })
  }

  private getElementPath(element: HTMLElement): string {
    const path: string[] = []
    let current: HTMLElement | null = element

    while (current) {
      let selector = current.tagName.toLowerCase()
      if (current.id) {
        selector += `#${current.id}`
      } else if (current.className) {
        selector += `.${current.className.split(' ').join('.')}`
      }
      path.unshift(selector)
      current = current.parentElement
    }

    return path.join(' > ')
  }
}

export const performanceMonitor = new PerformanceMonitor()
