import { fetchEventSource } from '@microsoft/fetch-event-source'
import { useLocalStorage } from '@vueuse/core'
import { onMounted, onUnmounted } from 'vue'
const AccessToken = useLocalStorage('access_token', null)
interface Options {
  onMessage: (msg?: any) => void
  onOpen: () => void
  onClose: () => void
  onError: (err?: any) => void
}
type Response = {
  data: string
  event: string
}
/**
 * @warning 如果事件监听器被命中，则不会调用 onMessage 回调函数
 * @warning 组件卸载时会自动关闭连接，组件挂载会自动连接，不需要放在 OnMounted 中
 * @param url sse 的 url 地址
 * @param options  可选项，包含 onMessage, onOpen, onClose, onError 四个回调函数
 * @returns close: 关闭连接的函数, addListener: 添加事件监听的函数, removeListener: 移除事件监听的函数
 * @description 这是一个封装了 fetchEventSource 的 Vue 组合式函数，提供了更方便的接口来使用 SSE。它支持动态添加和移除事件监听器，并且在组件卸载时自动关闭连接。
 */
export function useNexusEventSource(url: string, options?: Partial<Options>) {
  const ab = new AbortController()
  const ListenerMap = new Map<string, (msg?: any) => void>()
  const addListener = (event: string, callback: (msg?: any) => void) => {
    ListenerMap.set(event, callback)
  }
  const removeListener = (event: string) => {
    ListenerMap.delete(event)
  }

  onMounted(() => {
    fetchEventSource(url, {
      signal: ab.signal,
      headers: {
        'Content-Type': 'text/event-stream',
        Authorization: `Bearer ${AccessToken.value}`
      },
      async onopen() {
        if (options?.onOpen) options.onOpen()
      },
      onmessage(msg: Response) {
        const callback = ListenerMap.get(msg.event)
        if (callback) {
          callback(msg)
        } else if (options?.onMessage) options.onMessage(msg)
      },
      onerror(err) {
        if (options?.onError) options.onError(err)
      },
      onclose() {
        if (options?.onClose) options.onClose()
      }
    })
  })

  const close = () => {
    ab.abort()
    ListenerMap.clear()
  }
  onUnmounted(close)
  return {
    close,
    addListener,
    removeListener
  }
}
