import { onMounted, onUnmounted } from 'vue'
import { useWebSocket } from '@vueuse/core'

interface Options {
  onMessage: (msg?: any) => void
  onOpen: () => void
  onClose: () => void
  onError: (err?: any) => void
}

export function useNexusWebSocket(url: string, options?: Partial<Options>) {
  onMounted(() => {})
  onUnmounted(() => {})
  const { send, status, open, close } = useWebSocket(url, {
    autoReconnect: {
      retries: 5,
      delay: (retries) => Math.min(1000 * 2 ** (retries - 1), 30000)
    },
    heartbeat: {
      message: 'ping',
      pongTimeout: 1000
    },
    immediate: true,
    onConnected: () => {
      options?.onOpen()
    },
    onMessage: (msg) => {
      options?.onMessage(msg)
    },
    onDisconnected: () => {
      options?.onClose()
    },
    onError: (err) => {
      options?.onError(err)
    }
  })
  return { send, status, open , close }
}
