import { fetchEventSource, EventSourceMessage } from '@microsoft/fetch-event-source'
import { ref, onUnmounted } from 'vue'
import type { Ref } from 'vue'
import type { ConnectionStatus } from '@/constants'

/**
 * EventSource 连接选项
 */
export interface EventSourceOptions {
  /** 请求方法 */
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  /** 请求头 */
  headers?: Record<string, string>
  /** 请求体 (仅 POST/PUT/PATCH) */
  body?: string | FormData | URLSearchParams
  /** 是否立即连接 */
  immediate?: boolean
  /** 连接成功回调 */
  onConnected?: (response: Response) => void | Promise<void>
  /** 连接关闭回调 */
  onDisconnected?: () => void
  /** 连接错误回调 */
  onError?: (error: Error) => void
  /** 接收消息回调 */
  onMessage?: (message: string, event: EventSourceMessage) => void
  /** 自动重连配置 */
  autoReconnect?: {
    /** 重连次数 */
    retries: number
    /** 重连延迟 */
    delay: number | ((retries: number) => number)
    /** 重连失败回调 */
    onFailed?: () => void
  }
}

/**
 * EventSource 组合式函数返回值
 */
export interface EventSourceReturn {
  /** 连接状态 */
  status: Ref<ConnectionStatus>
  /** 最后接收到的数据 */
  data: Ref<string | null>
  /** 错误信息 */
  error: Ref<Error | null>
  /** 是否正在连接 */
  connecting: Ref<boolean>
  /** 是否已连接 */
  connected: Ref<boolean>
  /** 手动连接 */
  connect: () => void
  /** 手动断开 */
  disconnect: () => void
  /** 发送数据 (EventSource 不支持，仅保留接口一致性) */
  send: (data: string) => void
}

/**
 * 通用的 EventSource (SSE) 组合式函数
 *
 * @param url - SSE 服务端点 URL
 * @param options - 配置选项
 * @returns EventSource 控制对象
 *
 * @example
 * ```ts
 * const { status, data, connect, disconnect } = useEventSource('/api/sse', {
 *   onMessage: (msg) => console.log(msg),
 *   immediate: true
 * })
 * ```
 */
export function useEventSource(
  url: string,
  options: EventSourceOptions = {}
): EventSourceReturn {
  // 解构选项，设置默认值
  const {
    method = 'GET',
    headers = {},
    body,
    immediate = true,
    onConnected,
    onDisconnected,
    onError,
    onMessage,
    autoReconnect
  } = options

  // 状态管理
  const status = ref<ConnectionStatus>('disconnected')
  const data = ref<string | null>(null)
  const error = ref<Error | null>(null)
  const connecting = ref(false)
  const connected = ref(false)

  // 用于中止连接的 AbortController
  let abortController: AbortController | null = null
  // 重连计数
  let reconnectAttempts = 0
  const maxReconnectAttempts = autoReconnect?.retries ?? 0

  /**
   * 建立 EventSource 连接
   */
  const connect = () => {
    // 如果已经连接或正在连接，则不重复连接
    if (connected.value || connecting.value) {
      return
    }

    // 创建新的 AbortController
    abortController = new AbortController()
    connecting.value = true
    status.value = 'connecting'
    error.value = null

    // 计算重连延迟
    const getDelay = (): number => {
      if (!autoReconnect) return 0
      const { delay } = autoReconnect
      return typeof delay === 'function' ? delay(reconnectAttempts) : delay
    }

    fetchEventSource(url, {
      method,
      headers: body ? headers : { 'Content-Type': 'text/event-stream', ...headers },
      body,
      signal: abortController.signal,
      async onopen(response) {
        connecting.value = false
        connected.value = true
        status.value = 'connected'
        reconnectAttempts = 0
        await onConnected?.(response)
      },
      onmessage(event) {
        data.value = event.data
        onMessage?.(event.data, event)
      },
      onclose() {
        connecting.value = false
        connected.value = false
        status.value = 'disconnected'
        abortController = null
        onDisconnected?.()
      },
      onerror(err) {
        connecting.value = false
        connected.value = false
        error.value = err instanceof Error ? err : new Error(String(err))

        // 检查是否需要重连
        if (autoReconnect && reconnectAttempts < maxReconnectAttempts) {
          reconnectAttempts++
          status.value = 'connecting'
          onError?.(error.value)
          // 延迟后重连
          setTimeout(() => {
            connect()
          }, getDelay())
        } else {
          status.value = 'error'
          if (reconnectAttempts >= maxReconnectAttempts) {
            autoReconnect.onFailed?.()
          }
          onError?.(error.value)
          throw err
        }
      }
    })
  }

  /**
   * 断开 EventSource 连接
   */
  const disconnect = () => {
    if (abortController) {
      abortController.abort()
      abortController = null
    }
    connecting.value = false
    connected.value = false
    status.value = 'disconnected'
    onDisconnected?.()
  }

  /**
   * EventSource 不支持发送数据，保留此函数仅为接口一致性
   */
  const send = (_data: string) => {
    console.warn('EventSource does not support sending data to the server')
  }

  // 自动连接
  if (immediate) {
    connect()
  }

  // 组件卸载时自动断开
  onUnmounted(() => {
    disconnect()
  })

  return {
    status,
    data,
    error,
    connecting,
    connected,
    connect,
    disconnect,
    send
  }
}

/**
 * 便捷函数：创建带认证的 EventSource 连接
 *
 * @param url - SSE 服务端点 URL
 * @param token - 认证令牌
 * @param onMessage - 消息回调
 * @param options - 其他配置选项
 */
export function useAuthenticatedEventSource(
  url: string,
  token: string,
  onMessage: (message: string) => void,
  options: EventSourceOptions = {}
): EventSourceReturn {
  return useEventSource(url, {
    ...options,
    headers: {
      ...options.headers,
      Authorization: `Bearer ${token}`
    },
    onMessage
  })
}
