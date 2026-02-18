import { useWebSocket as useVueUseWebSocket } from '@vueuse/core'
import { ref, computed, onUnmounted } from 'vue'
import type { Ref } from 'vue'
import type { ConnectionStatus } from '@/constants'

/**
 * WebSocket 连接选项
 */
export interface WebSocketOptions {
  /** 是否立即连接 */
  immediate?: boolean
  /** WebSocket 协议 (ws 或 wss) */
  protocol?: 'ws' | 'wss'
  /** 自动重连配置 */
  autoReconnect?: {
    /** 重连次数 */
    retries: number
    /** 重连延迟 */
    delay: number | ((retries: number) => number)
    /** 重连失败回调 */
    onFailed?: () => void
  }
  /** 心跳配置 */
  heartbeat?: {
    /** 心跳消息内容 */
    message: string | ArrayBuffer | Blob
    /** 心跳间隔 */
    interval: number
    /** pong 超时时间 */
    pongTimeout: number
  }
  /** 连接成功回调 */
  onConnected?: (ws: WebSocket) => void
  /** 连接关闭回调 */
  onDisconnected?: (ws: WebSocket, event: CloseEvent) => void
  /** 连接错误回调 */
  onError?: (ws: WebSocket, event: Event) => void
  /** 接收消息回调 */
  onMessage?: (ws: WebSocket, event: MessageEvent) => void
}

/**
 * WebSocket 组合式函数返回值
 */
export interface WebSocketReturn {
  /** 连接状态 */
  status: Ref<ConnectionStatus>
  /** 最后接收到的数据 */
  data: Ref<string | null>
  /** 错误信息 */
  error: Ref<Event | null>
  /** 是否正在连接 */
  connecting: Ref<boolean>
  /** 是否已连接 */
  connected: Ref<boolean>
  /** 发送数据 */
  send: (data: string | ArrayBuffer | Blob) => void
  /** 打开连接 */
  open: () => void
  /** 关闭连接 */
  close: () => void
  /** 原始 WebSocket 实例 */
  ws: Ref<WebSocket | null>
}

/**
 * 通用的 WebSocket 组合式函数
 *
 * @param url - WebSocket 服务端点 URL
 * @param options - 配置选项
 * @returns WebSocket 控制对象
 *
 * @example
 * ```ts
 * const { status, data, send, close } = useWebSocket('ws://localhost:8080/ws', {
 *   onMessage: (_, event) => console.log(event.data),
 *   immediate: true
 * })
 *
 * // 发送消息
 * send('Hello Server')
 * ```
 */
export function useWebSocket(
  url: string,
  options: WebSocketOptions = {}
): WebSocketReturn {
  // 解构选项，设置默认值
  const {
    immediate = true,
    autoReconnect,
    heartbeat,
    onConnected,
    onDisconnected,
    onError,
    onMessage
  } = options

  // 使用 VueUse 的 useWebSocket 作为基础
  const {
    status: vueUseStatus,
    data,
    send,
    open,
    close,
    ws
  } = useVueUseWebSocket(url, {
    immediate,
    autoReconnect: autoReconnect
      ? {
          retries: autoReconnect.retries,
          delay: autoReconnect.delay,
          onFailed: autoReconnect.onFailed
        }
      : undefined,
    heartbeat,
    onConnected: (ws) => {
      onConnected?.(ws)
    },
    onDisconnected: (ws, event) => {
      onDisconnected?.(ws, event)
    },
    onError: (ws, event) => {
      onError?.(ws, event)
    },
    onMessage: (ws, event) => {
      onMessage?.(ws, event)
    }
  })

  // 映射状态到统一的 ConnectionStatus
  const status = computed<ConnectionStatus>(() => {
    const s = vueUseStatus.value
    if (s === 'OPEN') return 'connected'
    if (s === 'CONNECTING') return 'connecting'
    return 'disconnected'
  })

  // 派生状态
  const connecting = computed(() => status.value === 'connecting')
  const connected = computed(() => status.value === 'connected')

  // 错误状态 (VueUse 没有直接导出 error，用 ws 实例的 error 事件处理)
  const error = ref<Event | null>(null)

  // 监听错误事件
  if (ws.value) {
    ws.value.addEventListener('error', (e) => {
      error.value = e
    })
  }

  // 组件卸载时自动关闭连接
  onUnmounted(() => {
    if (connected.value) {
      close()
    }
  })

  return {
    status,
    data,
    error,
    connecting,
    connected,
    send,
    open,
    close,
    ws
  }
}

/**
 * 便捷函数：创建带心跳的 WebSocket 连接
 *
 * @param url - WebSocket 服务端点 URL
 * @param options - 配置选项
 */
export function useWebSocketWithHeartbeat(
  url: string,
  options: Omit<WebSocketOptions, 'heartbeat'> & {
    /** 心跳消息 */
    heartbeatMessage?: string
    /** 心跳间隔 */
    heartbeatInterval?: number
    /** pong 超时时间 */
    pongTimeout?: number
  } = {}
): WebSocketReturn {
  const {
    heartbeatMessage = 'ping',
    heartbeatInterval = 30000,
    pongTimeout = 10000,
    ...restOptions
  } = options

  return useWebSocket(url, {
    ...restOptions,
    heartbeat: {
      message: heartbeatMessage,
      interval: heartbeatInterval,
      pongTimeout
    }
  })
}

/**
 * 便捷函数：创建带自动重连的 WebSocket 连接
 *
 * @param url - WebSocket 服务端点 URL
 * @param options - 配置选项
 */
export function useWebSocketWithReconnect(
  url: string,
  options: Omit<WebSocketOptions, 'autoReconnect'> & {
    /** 重连次数 */
    retries?: number
    /** 重连延迟基数 */
    retryDelay?: number
    /** 重连失败回调 */
    onRetryFailed?: () => void
  } = {}
): WebSocketReturn {
  const {
    retries = 3,
    retryDelay = 1000,
    onRetryFailed,
    ...restOptions
  } = options

  return useWebSocket(url, {
    ...restOptions,
    autoReconnect: {
      retries,
      delay: (attempt) => Math.min(retryDelay * 2 ** (attempt - 1), 30000),
      onFailed: onRetryFailed
    }
  })
}
