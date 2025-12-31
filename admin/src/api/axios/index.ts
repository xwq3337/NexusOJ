import AbortAxios from '@/api/axios/AbortAxios'
import { _RequstInterceptors } from '@/api/axios/Interceptors'
import { staticAxiosConfig } from '@/api/axios/staticAxiosConfig'
import type {
	AxiosError,
	AxiosInstance,
	AxiosRequestConfig,
	AxiosResponse,
	CreateAxiosDefaults,
	InternalAxiosRequestConfig
} from 'axios'
import axios from 'axios'
import type { AxiosOptions, Respones, RequestInterceptors } from '@/api/axios/type'

class Axios {
	#axiosInstance: AxiosInstance
	#options: AxiosOptions
	#interceptors

	constructor(options) {
		this.#axiosInstance = axios.create(options)
		this.#options = options
		this.#interceptors = options.interceptors
		this.setInterceptors()
	}

	request<T = any>(config: AxiosRequestConfig): Promise<T> {
		return new Promise((resolve, reject) => {
			this.#axiosInstance
				.request<any, AxiosResponse<Respones>>(config)
				.then((res) => {
					return resolve(res as unknown as Promise<T>)
				})
				.catch((err) => {
					return reject(err)
				})
		})
	}

	get<T = any>(config: AxiosRequestConfig): Promise<T> {
		return this.request<T>({ ...config, method: 'GET' })
	}

	post<T = any>(config: AxiosRequestConfig): Promise<T> {
		return this.request<T>({ ...config, method: 'POST' })
	}

	put<T = any>(config: AxiosRequestConfig): Promise<T> {
		return this.request<T>({ ...config, method: 'PUT' })
	}

	delete<T = any>(config: AxiosRequestConfig): Promise<T> {
		return this.request<T>({ ...config, method: 'DELETE' })
	}

	setInterceptors() {
		if (!this.#interceptors) return
		const {
			requestInterceptors,
			requestInterceptorsCatch,
			responseInterceptors,
			responseInterceptorsCatch
		} = this.#interceptors
		const abortAxios = new AbortAxios()

		this.#axiosInstance.interceptors.request.use((config: InternalAxiosRequestConfig) => {
			const abortRepetitiveRequest =
				(config as unknown as any)?.abortRepetitiveRequest ??
				this.#options.abortRepetitiveRequest
			if (abortRepetitiveRequest) {
				abortAxios.addPending(config)
			}
			if (requestInterceptors) config = requestInterceptors(config)
			return config
		}, requestInterceptorsCatch ?? undefined)

		this.#axiosInstance.interceptors.response.use(
			async (response: AxiosResponse) => {
				response && abortAxios.removePending(response.config)
				if (responseInterceptors) response = await responseInterceptors(response)
				if (this.#options.directlyGetdata && response.data) response = response.data
				return response
			},
			async (err: AxiosError) => {
				if (responseInterceptorsCatch) {
					return await responseInterceptorsCatch(this.#axiosInstance, err)
				}
				return err
			}
		)
	}
}

const Request = new Axios({
	directlyGetdata: true,
	baseURL: staticAxiosConfig.baseURL,
	timeout: 10 * 1000,
	interceptors: _RequstInterceptors,
	headers: { 'Content-Type': 'application/json' },
	abortRepetitiveRequest: false,
	retryConfig: {
		count: 5,
		waitTime: 500
	}
})

export default Request
