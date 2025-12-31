// 继承了我们在最开始实现的抽象类RequstInterceptors，主要关心responseInterceptorsCatch内容
import Request from '@/api/axios'
import axios, { HttpStatusCode, type AxiosError } from 'axios'
import { checkErrorStatus } from '@/api/axios/checkErrorStatus'
import type { RequestInterceptors } from '@/api/axios/type'
import { useLocalStorage } from '@vueuse/core'
import type { AxiosInstance, AxiosResponse, AxiosError } from 'axios'
import type { CustomAxiosRequestConfig } from './type'
import { usePerformanceStore } from '@/stores/performance'

export class Interceptors {
  private instance: AxiosInstance
  private performanceStore = usePerformanceStore()
  private requestStartTimes = new Map<string, number>()

  constructor(instance: AxiosInstance) {
    this.instance = instance
    this.initInterceptors()
  }

  private initInterceptors() {
    // 请求拦截器
    this.instance.interceptors.request.use(
      (config: CustomAxiosRequestConfig) => {
        const requestId = Math.random().toString(36).substring(7)
        this.requestStartTimes.set(requestId, Date.now())
        config.requestId = requestId
        return config
      },
      (error: AxiosError) => {
        this.updateMetrics(false, 0)
        return Promise.reject(error)
      }
    )

    // 响应拦截器
    this.instance.interceptors.response.use(
      (response: AxiosResponse) => {
        const requestId = (response.config as CustomAxiosRequestConfig).requestId
        const startTime = this.requestStartTimes.get(requestId)
        const duration = startTime ? Date.now() - startTime : 0
        this.requestStartTimes.delete(requestId)
        this.updateMetrics(true, duration)
        return response
      },
      (error: AxiosError) => {
        const requestId = (error.config as CustomAxiosRequestConfig)?.requestId
        const startTime = requestId ? this.requestStartTimes.get(requestId) : null
        const duration = startTime ? Date.now() - startTime : 0
        this.requestStartTimes.delete(requestId)
        this.updateMetrics(false, duration)
        return Promise.reject(error)
      }
    )
  }

  private updateMetrics(isSuccess: boolean, duration: number) {
    const currentMetrics = this.performanceStore.apiMetrics
    const totalRequests = currentMetrics.totalRequests + 1
    const successRequests = isSuccess ? currentMetrics.successRequests + 1 : currentMetrics.successRequests
    const failedRequests = !isSuccess ? currentMetrics.failedRequests + 1 : currentMetrics.failedRequests
    const totalTime = currentMetrics.averageResponseTime * currentMetrics.totalRequests + duration
    const averageResponseTime = totalTime / totalRequests

    this.performanceStore.updateApiMetrics({
      totalRequests,
      successRequests,
      failedRequests,
      averageResponseTime
    })
  }
}

const AccessToken = useLocalStorage('access_token', null)
const RefreshToken = useLocalStorage('refresh_token', null)
const userName = useLocalStorage('userName', null)
const password = useLocalStorage('password', null)
export const _RequstInterceptors: RequestInterceptors = {
  requestInterceptors(config) {
    config.headers['Authorization'] ||= AccessToken.value
    return config
  },
  requestInterceptorsCatch(err) {
    return err
  },
  responseInterceptors(res) {
    console.log(res.data)
    if (res.data.code === HttpStatusCode.Unauthorized) {
      return refreshToken().then(async () => {
        const newConfig = {
          ...res.config,
          headers: {
            ...res.config.headers,
            'Authorization': `${AccessToken.value}`,
          }
        };
        return await Request.request(newConfig);
      }).catch(error => {
        console.error('Failed to retry request after refreshing token:', error);
        return res;
      });
    }
    return res
  },
  responseInterceptorsCatch(axiosInstance, err: AxiosError) {
    const message = err.code === 'ECONNABORTED' ? '请求超时' : undefined
    console.log('catch', err)
    if (axios.isCancel(err)) {
      return Promise.reject(err)
    }
    checkErrorStatus((err as AxiosError).response?.status, message, (message) =>
      console.log(message)
    )
    // return retry(axiosInstance, err as AxiosError)
  }
}

const refreshToken = () => {
  return Request.post({
    url: '/user/refresh',
    data: {
      username: userName.value,
      password: password.value
    },
    headers: {
      Authorization: RefreshToken.value
    }
  })
    .then((response) => {
      if (response.code == HttpStatusCode.Ok) {
        [AccessToken.value, RefreshToken.value] = [response.info[0], response.info[1]]
        return true
      } else {
        [AccessToken.value, RefreshToken.value] = [null, null]
        console.error('请重新登录')
        return false
      }
    })
    .catch((err) => {
      console.log(err)
      return false
    })
}
