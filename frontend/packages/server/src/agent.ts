import { fetchEventSource } from '@microsoft/fetch-event-source'
import axios from 'axios'

interface Message {
  id: number
  role: 'user' | 'model'
  content: string
  timestamp: number
}

interface StreamCallbacks {
  onMessage: (chunk: string) => void
  onError: (error: Error) => void
  onClose: () => void
  onStatus?: (status: string, tool?: string) => void
}

/** 从 localStorage 读取 access_token */
function getAuthToken(): string {
  const raw = localStorage.getItem('access_token')
  if (!raw) return ''
  // useLocalStorage 存的值可能是 JSON 字符串也可能是原始字符串
  try {
    const parsed = JSON.parse(raw)
    return typeof parsed === 'string' ? parsed : ''
  } catch {
    return raw
  }
}

/** 构建 Authorization headers */
function authHeaders(): Record<string, string> {
  const token = getAuthToken()
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  return headers
}

// ==================== SSE 通用处理器 ====================

async function streamSSE(
  url: string,
  body: Record<string, unknown>,
  callbacks: StreamCallbacks,
  abortController: AbortController,
  timeoutMs = 120_000,
) {
  const timeoutId = setTimeout(() => abortController.abort(), timeoutMs)

  await fetchEventSource(url, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify(body),
    signal: abortController.signal,

    onmessage(msg) {
      try {
        const data = JSON.parse(msg.data)

        if (data.done) {
          clearTimeout(timeoutId)
          callbacks.onClose()
          return
        }

        if (data.error) {
          callbacks.onError(new Error(data.error))
          return
        }

        if (data.status) {
          callbacks.onStatus?.(data.status, data.tool)
          return
        }

        if (data.text) {
          callbacks.onMessage(data.text)
        }
      } catch {
        if (msg.data) {
          callbacks.onMessage(msg.data)
        }
      }
    },

    onerror(error) {
      clearTimeout(timeoutId)
      callbacks.onError(new Error('SSE connection error: ' + (error as Error).message))
      throw error
    },

    onclose() {
      clearTimeout(timeoutId)
      callbacks.onClose()
    },
  })
}

// ==================== 对外 API ====================

/** 流式对话（主聊天端点） */
export const streamChat = async (
  messages: Message[],
  callbacks: StreamCallbacks,
  abortController: AbortController,
) => {
  return streamSSE('/ai/chat', { messages }, callbacks, abortController)
}

/** 流式代码分析 */
export const streamCodeAnalysis = async (
  code: string,
  language: string,
  analysisType: string,
  callbacks: StreamCallbacks,
  abortController: AbortController,
) => {
  return streamSSE('/ai/analyze-code', { code, language, analysis_type: analysisType }, callbacks, abortController)
}

/** 流式个性化指导 */
export const streamGuidance = async (
  question: string,
  messages: Message[],
  callbacks: StreamCallbacks,
  abortController: AbortController,
) => {
  return streamSSE('/ai/personalized-guidance', { question, messages }, callbacks, abortController)
}

/** 测试用例生成（非流式，返回 JSON） */
export const generateTestCases = async (
  problemId: number,
  userCode: string = '',
  count: number = 5,
) => {
  const resp = await axios.post('/ai/generate-tests', {
    problem_id: problemId,
    user_code: userCode,
    count,
  }, { headers: authHeaders() })
  return resp.data
}

/**
 * @deprecated 使用 streamChat 替代
 */
export const streamResponse = async (
  messages: Message[],
  onMessage: (chunk: string) => void,
  onError: (error: Error) => void,
  onClose: () => void,
  abortController: AbortController,
) => {
  return streamChat(messages, { onMessage, onError, onClose }, abortController)
}
