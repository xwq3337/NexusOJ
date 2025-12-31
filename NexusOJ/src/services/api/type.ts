import type {
  AxiosRequestConfig,
  InternalAxiosRequestConfig,
  AxiosResponse,
  AxiosInstance,
  AxiosError
} from 'axios'

export interface AxiosOptions extends AxiosRequestConfig {
  directlyGetdata?: boolean
  interceptors?: RequestInterceptors
  abortRepetitiveRequest?: boolean
  retryConfig?: {
    count: number
    waitTime: number
  }
}

export abstract class RequestInterceptors {
  abstract requestInterceptors?: (config: InternalAxiosRequestConfig) => InternalAxiosRequestConfig
  abstract requestInterceptorsCatch: (err: Error) => Error
  abstract responseInterceptors?: (
    res: AxiosResponse<any, any>
  ) => AxiosResponse<any, any> | Promise<any>
  abstract responseInterceptorsCatch?: (axiosInstance: AxiosInstance, error: AxiosError) => void
}

export interface Respones<T = any> {
  code: number
  result: T
}
