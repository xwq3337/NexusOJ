import type { AxiosError } from 'axios'
export const checkErrorStatus = (error: AxiosError) => {
  const { response } = error
  if (response) {
    const { status } = response
    switch (status) {
      case 400:
        console.error('请求错误')
        break
      case 401:
        console.error('未授权，请登录')
        break
      case 403:
        console.error('拒绝访问')
        break
      case 404:
        console.error('请求地址出错')
        break
      case 408:
        console.error('请求超时')
        break
      case 500:
        console.error('服务器内部错误')
        break
      case 501:
        console.error('服务未实现')
        break
      case 502:
        console.error('网关错误')
        break
      case 503:
        console.error('服务不可用')
        break
      case 504:
        console.error('网关超时')
        break
      case 505:
        console.error('HTTP版本不受支持')
        break
      default:
        console.error('网络错误')
    }
  }
  return Promise.reject(error)
}
