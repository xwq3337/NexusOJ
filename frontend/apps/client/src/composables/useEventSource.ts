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
export function useNexusEventSource(url: string, options?: Partial<Options>) {
  const ab = new AbortController()
  onMounted(() => {
    fetchEventSource(url, {
      signal: ab.signal,
      headers: {
        'Content-Type': 'text/event-stream',
        Authorization: `Bearer ${AccessToken.value}`
      },
      async onopen() {
        options?.onOpen()
      },
      onmessage(msg) {
        options?.onMessage(msg)
      },
      onerror(err) {
        options?.onError(err)
      },
      onclose() {
        options?.onClose()
      }
    })
  })
  onUnmounted(() => {
    ab.abort()
  })
  const close = ab.abort
  return {
    close
  }
}
