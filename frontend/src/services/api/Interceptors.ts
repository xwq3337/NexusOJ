import Request from './index'
import axios, { HttpStatusCode, type AxiosError } from 'axios'
import { checkErrorStatus } from './checkErrorStatus'
import type { RequestInterceptors } from './type'
import { useLocalStorage } from '@vueuse/core'
const AccessToken = useLocalStorage('access_token', null)
const RefreshToken = useLocalStorage('refresh_token', null)
const account = useLocalStorage('account', null)
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
    if (res.data.code === HttpStatusCode.Unauthorized) {
      // 尝试刷新令牌
      if (RefreshToken.value) {
        console.error('登录已过期，正在尝试刷新令牌...')
        return refreshToken()
          .then(async () => {
            const newConfig = {
              ...res.config,
              headers: {
                ...res.config.headers,
                Authorization: `${AccessToken.value}`
              }
            }
            return await Request.request(newConfig)
          })
          .catch((error) => {
            console.error('Failed to retry request after refreshing token:', error)
            return res
          })
      } else return Promise.reject(res)
    }
    return res
  },
  responseInterceptorsCatch(axiosInstance, err: AxiosError) {
    const message = err.code === '服务端主动终止连接' ? '请求超时' : undefined
    if (axios.isCancel(err)) {
      return Promise.reject(err)
    }
    // TODO: 处理返回状态
    // checkErrorStatus(
    //     (err as AxiosError).response?.status,
    //     message,
    //     (message) =>
    //         console.log(message)
    // )
    // return retry(axiosInstance, err as AxiosError)
  }
}

const refreshToken = () => {
  console.log('正在刷新令牌，请稍候...')
  return Request.post(
    '/user/refresh',
    {
      /// TODO: 刷新令牌接口，字段变更
      username: account.value,
      password: password.value
    },
    {
      headers: {
        Authorization: `${RefreshToken.value}`
      }
    }
  )
    .then((response) => {
      if (response.code == HttpStatusCode.Ok) {
        ;[AccessToken.value, RefreshToken.value] = [response.info[0], response.info[1]]
        return true
      } else {
        ;[AccessToken.value, RefreshToken.value] = [null, null]
        console.error('登录已过期，请重新登录', response.message)
        return false
      }
    })
    .catch((err) => {
      console.error(err.message || '登录已过期，请重新登录')
      return false
    })
}
