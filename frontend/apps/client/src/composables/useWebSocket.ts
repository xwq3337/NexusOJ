import { onMounted, onUnmounted } from 'vue'
import { useWebSocket } from '@vueuse/core'
import { ChatMessage } from '@nexusoj/type'

interface Options {
  onMessage: (msg?: MessageEvent<ChatMessage | "pong">) => void
  onOpen: () => void
  onClose: () => void
  onError: (err?: any) => void
}

export function useNexusWebSocket(url: string, options?: Partial<Options>) {
  onMounted(() => {})
  onUnmounted(() => {})
  const { send, status, open, close } = useWebSocket(url, {
    autoReconnect: {
      retries: 3, // 最多重试3次
      delay: (retries) => Math.min(1000 * 2 ** (retries - 1), 30000)
    },
    heartbeat: {
      message: 'ping',
      interval: 5000, // 每5秒发送一次ping
      pongTimeout: 10000 // 10秒内未收到pong则认为连接断开
    },
    immediate: true,
    onConnected: () => {
      options?.onOpen()
    },
    onMessage: (_, msg : MessageEvent<ChatMessage>) => {
      options?.onMessage(msg)
    },
    onDisconnected: () => {
      options?.onClose()
    },
    onError: (_, err) => {
      options?.onError(err)
    }
  })
  return { send, status, open , close }
}