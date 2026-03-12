import { fetchEventSource } from '@microsoft/fetch-event-source'

/**
 * Stream AI response from backend SSE endpoint
 * @param query - User's question or prompt
 * @param onMessage - Callback function for each message chunk
 * @param onError - Callback function for errors
 * @param onClose - Callback function when connection closes
 * @param abortController - Controller to abort the request
 */
export const streamResponse = async (
  query: string,
  onMessage: (chunk: string) => void,
  onError: (error: Error) => void,
  onClose: () => void,
  abortController: AbortController
) => {
  await fetchEventSource('/agent/test', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ query }),
    signal: abortController.signal,

    onmessage(msg) {
      if (msg.data === '[DONE]') {
        onClose()
        return
      }

      try {
        const data = JSON.parse(msg.data)
        if (data.content) {
          onMessage(data.content)
        }
      } catch (e) {
        // If not JSON, treat as plain text
        if (msg.data) {
          onMessage(msg.data)
        }
      }
    },

    onerror(error) {
      onError(new Error('SSE connection error: ' + error.message))
      throw error // Re-throw to stop reconnection
    },

    onclose() {
      onClose()
    }
  })
}
