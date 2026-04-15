import { fetchEventSource } from '@microsoft/fetch-event-source'

interface Message {
  id: number
  role: 'user' | 'model'
  content: string
  timestamp: number
}

/**
 * Stream AI response from backend SSE endpoint
 * @param messages - Chat message history
 * @param onMessage - Callback function for each message chunk
 * @param onError - Callback function for errors
 * @param onClose - Callback function when connection closes
 * @param abortController - Controller to abort the request
 */
export const streamResponse = async (
  messages: Message[],
  onMessage: (chunk: string) => void,
  onError: (error: Error) => void,
  onClose: () => void,
  abortController: AbortController
) => {
  // RAG 检索耗时较长，需要延长超时
  const timeoutId = setTimeout(() => abortController.abort(), 120_000)

  await fetchEventSource('/agent/chat', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ messages }),
    signal: abortController.signal,

    onmessage(msg) {
      try {
        const data = JSON.parse(msg.data)

        if (data.done) {
          clearTimeout(timeoutId)
          onClose()
          return
        }

        if (data.error) {
          onError(new Error(data.error))
          return
        }

        if (data.text) {
          onMessage(data.text)
        }
      } catch {
        // If not JSON, treat as plain text
        if (msg.data) {
          onMessage(msg.data)
        }
      }
    },

    onerror(error) {
      clearTimeout(timeoutId)
      onError(new Error('SSE connection error: ' + error.message))
      throw error // Re-throw to stop reconnection
    },

    onclose() {
      clearTimeout(timeoutId)
      onClose()
    }
  })
}
