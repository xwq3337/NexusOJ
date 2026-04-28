<template>
  <!-- Floating Action Button -->
  <Transition name="fab-fade">
    <button v-if="!isOpen" @click="toggleChat"
      class="fixed bottom-6 right-6 w-14 h-14 rounded-full shadow-lg flex items-center justify-center z-50 hover:scale-110 transition-transform duration-200"
      :style="{ backgroundColor: 'var(--accent-color)', color: 'white' }" :title="isMobile ? 'AI助手' : ''">
      <Bot :size="28" />
    </button>
  </Transition>

  <!-- Chat Window -->
  <Transition name="chat-slide">
    <div v-if="isOpen"
      class="fixed bottom-6 right-6 w-[calc(100%-3rem)] md:w-lg h-125 rounded-xl shadow-2xl flex flex-col z-50 overflow-hidden"
      :style="{ backgroundColor: 'var(--surface-primary)', border: '1px solid var(--border-color)' }">
      <!-- Header -->
      <div class="flex items-center justify-between px-4 py-3 border-b" :style="{ borderColor: 'var(--border-color)' }">
        <div class="flex items-center gap-2">
          <Bot :size="20" :style="{ color: 'var(--accent-color)' }" />
          <span class="font-semibold" :style="{ color: 'var(--text-primary)' }">
            NexusAI
          </span>
          <div v-if="isTyping" class="flex items-center gap-1">
            <div class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></div>
            <span class="text-xs" :style="{ color: 'var(--text-secondary)' }">{{ typingStatus }}</span>
          </div>
        </div>
        <button @click="toggleChat"
          class="flex items-center justify-center w-8 h-8 rounded-lg transition-colors cursor-pointer"
          :style="{ color: 'var(--text-secondary)' }" :class="`hover:bg-(--surface-secondary)`">
          <X :size="18" />
        </button>
      </div>

      <!-- Messages Container -->
      <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-4">
        <!-- Welcome Message -->
        <div v-if="messages.length === 0"
          class="flex flex-col items-center justify-center h-full text-center space-y-4">
          <div class="w-16 h-16 rounded-full flex items-center justify-center bg-(--accent-color)">
            <Bot :size="32" style="color: white" />
          </div>
          <div>
            <h3 class="text-lg font-semibold mb-1" :style="{ color: 'var(--text-primary)' }">
              你好！我是 NexusAI
            </h3>
            <p class="text-sm text-(--text-secondary)">
              我可以帮助你理解算法、调试代码、分析代码质量和生成测试用例
            </p>
          </div>
          <div class="grid grid-cols-2 gap-2 w-full">
            <button v-for="prompt in quickPrompts" :key="prompt.text" @click="sendQuickPrompt(prompt.text, prompt.action)"
              class="px-3 py-2 rounded-lg text-xs text-left transition-colors " :style="{
                backgroundColor: 'var(--surface-secondary)',
                color: 'var(--text-secondary)',
                border: '1px solid var(--border-color)'
              }" :class="`hover:bg-(--surface-tertiary) hover:text-(--text-primary)`">
              {{ prompt.label }}
            </button>
          </div>
        </div>

        <!-- Messages List -->
        <div v-for="message in messages" :key="message.id" class="flex gap-3"
          :class="message.role === 'user' ? 'flex-row-reverse' : 'flex-row'">
          <!-- Avatar -->
          <div class="w-8 h-8 rounded-full shrink-0 flex items-center justify-center" :style="{
            backgroundColor: message.role === 'user' ? 'var(--accent-color)' : 'var(--surface-secondary)'
          }">
            <User v-if="message.role === 'user'" :size="16" style="color: white" />
            <Bot v-else :size="16" :style="{ color: 'var(--accent-color)' }" />
          </div>

          <!-- Message Content -->
          <div class="max-w-[75%] rounded-2xl px-4 py-2" :style="{
            backgroundColor: message.role === 'user' ? 'var(--accent-color)' : 'var(--surface-secondary)',
            color: message.role === 'user' ? 'white' : 'var(--text-primary)'
          }">
            <div class="text-sm  leading-relaxed message-content" v-html="renderMessage(message.content)"></div>
            <div class="text-xs mt-1 opacity-70"
              :style="{ color: message.role === 'user' ? 'white' : 'var(--text-tertiary)' }">
              {{ formateTimeStamp(message.timestamp) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Input Area -->
      <div class="p-4 border-t" :style="{ borderColor: 'var(--border-color)' }">
        <form @submit.prevent="sendMessage" class="flex gap-2 items-end">
          <n-input type="textarea" v-model:value="inputMessage" ref="inputRef" placeholder="输入你的问题..."
            :disabled="isLoading" :autosize="{ minRows: 1, maxRows: 6 }" :show-count="false"
            :input-props="{ style: { fontSize: '0.875rem' } }" class="flex-1" @keydown="handleKeydown" />
          <button type="submit" :disabled="!inputMessage.trim() || isLoading"
            class="px-4 py-2.5 rounded-lg flex items-center justify-center transition-colors disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer shrink-0"
            :style="{
              backgroundColor: inputMessage.trim() && !isLoading ? 'var(--accent-color)' : 'var(--btn-secondary)',
              color: 'white'
            }" :class="`hover:bg-(--accent-hover)`">
            <Send v-if="!isLoading" :size="18" />
            <Loader2 v-else :size="18" class="animate-spin" />
          </button>
        </form>
        <div class="flex items-center justify-between mt-2">
          <p class="text-xs" :style="{ color: 'var(--text-tertiary)' }">
            {{ messages.length }} 条消息
          </p>
          <n-popconfirm v-if="messages.length > 0" @positive-click="clearMessages">
            <template #trigger>
              <n-button size="tiny" type="error" class="text-xs"> 清空对话
              </n-button>
            </template>
            确定要清空对话记录吗
          </n-popconfirm>
        </div>
      </div>
    </div>
  </Transition>

  <!-- Overlay for mobile -->
  <Transition name="overlay-fade">
    <div v-if="isOpen && isMobile" @click="toggleChat" class="fixed inset-0 bg-black/50 z-40"></div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted, watch, computed } from 'vue'
import { Bot, X, User, Send, Loader2 } from 'lucide-vue-next'
import { useLocalStorage, useEventListener } from '@vueuse/core'
import { streamChat, streamCodeAnalysis, streamGuidance } from '@nexusoj/server'
import MarkdownIt from 'markdown-it'
import { NPopconfirm, NButton, NInput, useMessage } from "naive-ui"
import { formateTimeStamp } from '@/utils/format'
import { useRoute } from 'vue-router'

const message = useMessage()

type ActionType = 'chat' | 'analyze' | 'guidance'

interface AIMessage {
  id: number
  role: 'model' | 'user'
  content: string
  timestamp: number
}

const AiAssistantMessages = useLocalStorage<AIMessage[]>('ai-assistant-messages', [])

// Initialize markdown-it
const md = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true,
  breaks: true
})

// State
const isOpen = ref(false)
const inputMessage = ref('')
const messages = ref<AIMessage[]>([])
const isLoading = ref(false)
const isTyping = ref(false)
const typingStatus = ref('输入中...')
const messagesContainer = ref<HTMLElement>()
const inputRef = ref<HTMLInputElement>()
const isMobile = ref(false)
const currentAction = ref<ActionType>('chat')

// Quick prompts
const quickPrompts = [
  { text: '如何优化时间复杂度？', label: '⚡ 优化算法', action: 'chat' as ActionType },
  { text: '帮我分析这段代码的质量', label: '🔍 代码分析', action: 'analyze' as ActionType },
  { text: '解释这道题的思路', label: '💡 解题思路', action: 'chat' as ActionType },
  { text: '给我一些个性化练习建议', label: '📚 个性指导', action: 'guidance' as ActionType },
]

// 根据当前路由推断用户意图
const route = useRoute()
const inferContext = computed(() => {
  const path = route.path
  if (path.includes('/problem/')) {
    return { type: 'problem', hint: '当前在题目页面，用户可能需要解题帮助或代码分析。' }
  }
  if (path.includes('/record')) {
    return { type: 'record', hint: '当前在提交记录页面，用户可能需要代码调试帮助。' }
  }
  return null
})

// Load messages from localStorage
onMounted(() => {
  messages.value = AiAssistantMessages.value
  checkMobile()
})

const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
}
useEventListener('resize', checkMobile)

// Save messages to localStorage
watch(messages, (newMessages) => {
  AiAssistantMessages.value = newMessages
}, { deep: true })

// Toggle chat
const toggleChat = () => {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    nextTick(() => {
      scrollToBottom()
      focusInput()
    })
  }
}

// Focus input
const focusInput = () => {
  inputRef.value?.focus()
}

// Scroll to bottom
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

// Escape angle brackets outside code blocks
const escapeAngleBrackets = (text: string) => {
  const protectedText = text.replace(/&[a-zA-Z]+;/g, (match) => {
    return `__HTML_ENTITY_${match}__`
  })

  const codeBlockRegex = /`[^`]*`|```[\s\S]*?```/g
  const parts: string[] = []
  let lastIndex = 0
  let match: RegExpExecArray | null

  while ((match = codeBlockRegex.exec(protectedText)) !== null) {
    const beforeCode = protectedText.substring(lastIndex, match.index)
    parts.push(beforeCode.replace(/</g, '&lt;').replace(/>/g, '&gt;'))
    parts.push(match[0])
    lastIndex = match.index + match[0].length
  }

  const remaining = protectedText.substring(lastIndex)
  parts.push(remaining.replace(/</g, '&lt;').replace(/>/g, '&gt;'))

  return parts.join('').replace(/__HTML_ENTITY_(&[a-zA-Z]+;)__/g, '$1')
}

// Render message with markdown
const renderMessage = (text: string) => {
  if (!text) return ''

  let escaped = escapeAngleBrackets(text)
  let html = md.render(escaped)

  html = html.replace(
    /<pre><code class="language-(\w+)">([\s\S]*?)<\/code><\/pre>/g,
    (match, lang, code) => {
      return `<pre class="bg-[var(--bg-tertiary)] rounded-lg p-3 overflow-x-auto my-2"><code class="text-sm text-[var(--accent-color)]">${code.trim()}</code></pre>`
    }
  )

  html = html.replace(
    /<code>([\s\S]*?)<\/code>/g,
    (match, code) => {
      return `<code class="px-1.5 py-0.5 rounded text-sm" style="background-color: var(--bg-tertiary); color: var(--accent-color);">${code}</code>`
    }
  )

  return html
}

// Send quick prompt
const sendQuickPrompt = (prompt: string, action: ActionType = 'chat') => {
  currentAction.value = action
  inputMessage.value = prompt
  sendMessage()
}

// Send message with streaming support
const sendMessage = async () => {
  const text = inputMessage.value.trim()
  if (!text || isLoading.value) return

  const action = currentAction.value
  currentAction.value = 'chat' // 重置

  // Add user message
  const userMessage: AIMessage = {
    id: Date.now(),
    role: 'user',
    content: text,
    timestamp: Date.now()
  }
  messages.value.push(userMessage)
  inputMessage.value = ''
  scrollToBottom()

  isLoading.value = true
  isTyping.value = true
  typingStatus.value = '思考中...'

  const aiMessage: AIMessage = {
    id: Date.now() + 1,
    role: 'model',
    content: '',
    timestamp: Date.now()
  }
  messages.value.push(aiMessage)
  const aiMessageIndex = messages.value.length - 1

  const abortController = new AbortController()

  try {
    const callbacks = {
      onMessage: (chunk: string) => {
        messages.value[aiMessageIndex].content += chunk
        scrollToBottom()
      },
      onError: (error: Error) => {
        console.error('SSE Error:', error)
        messages.value[aiMessageIndex].content = '抱歉，我遇到了一些问题。请稍后再试。'
        scrollToBottom()
      },
      onClose: () => {
        isLoading.value = false
        isTyping.value = false
      },
      onStatus: (status: string, tool?: string) => {
        if (status === 'thinking') {
          typingStatus.value = '思考中...'
        } else if (status === 'tool_call') {
          const toolNames: Record<string, string> = {
            knowledge_search: '检索知识库',
            analyze_code: '分析代码',
            generate_test_cases: '生成测试用例',
            get_user_ability_profile: '获取用户画像',
            lookup_problem: '查询题目',
          }
          typingStatus.value = toolNames[tool || ''] || '调用工具中...'
        } else if (status === 'tool_result') {
          typingStatus.value = '生成回答...'
        }
      },
    }

    // 根据功能类型选择不同端点
    if (action === 'analyze') {
      await streamCodeAnalysis(text, 'auto', 'all', callbacks, abortController)
    } else if (action === 'guidance') {
      await streamGuidance(text, messages.value.filter(m => m.id !== aiMessage.id).slice(-10), callbacks, abortController)
    } else {
      await streamChat(messages.value.filter(m => m.id !== aiMessage.id), callbacks, abortController)
    }
  } catch (error) {
    console.error('Failed to get AI response:', error)
    messages.value[aiMessageIndex].content = '抱歉，我遇到了一些问题。请稍后再试。'
    scrollToBottom()
    isLoading.value = false
    isTyping.value = false
  }
}

// Clear messages
const clearMessages = () => {
  messages.value = []
}

// Keyboard shortcut to open/close
const handleKeydown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    toggleChat()
  }

  if (e.key === 'Enter' && !e.shiftKey && !isLoading.value && inputMessage.value.trim()) {
    e.preventDefault()
  }
  else if (e.key === 'Enter' && e.shiftKey) {
    e.preventDefault()
    inputMessage.value += '\n'
  }
}
useEventListener('keydown', handleKeydown)
</script>

<style scoped>
.message-content :deep(pre) {
  margin: 0.5rem 0;
  padding: 0.75rem;
  border-radius: 0.5rem;
  overflow-x: auto;
  font-family: 'Courier New', monospace;
  font-size: 0.875rem;
  line-height: 1.5;
}

.message-content :deep(code) {
  font-family: 'Courier New', monospace;
}

.message-content :deep(p) {
  margin-bottom: 0.5rem;
}

.message-content :deep(p:last-child) {
  margin-bottom: 0;
}

.message-content :deep(ul),
.message-content :deep(ol) {
  margin-left: 1.5rem;
  margin-bottom: 0.5rem;
}

.message-content :deep(li) {
  margin-bottom: 0.25rem;
}

.message-content :deep(a) {
  color: var(--accent-color);
  text-decoration: underline;
}

/* FAB Animation */
.fab-fade-enter-active,
.fab-fade-leave-active {
  transition: all 0.3s ease;
}

.fab-fade-enter-from,
.fab-fade-leave-to {
  opacity: 0;
  transform: scale(0.5);
}

/* Chat Window Animation */
.chat-slide-enter-active,
.chat-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.chat-slide-enter-from {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}

.chat-slide-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}

/* Overlay Animation */
.overlay-fade-enter-active,
.overlay-fade-leave-active {
  transition: opacity 0.3s ease;
}

.overlay-fade-enter-from,
.overlay-fade-leave-to {
  opacity: 0;
}

/* Scrollbar Styling */
.messages-container :deep(::-webkit-scrollbar) {
  width: 6px;
}

.messages-container :deep(::-webkit-scrollbar-track) {
  background: transparent;
}

.messages-container :deep(::-webkit-scrollbar-thumb) {
  background: var(--surface-tertiary);
  border-radius: 3px;
}

.messages-container :deep(::-webkit-scrollbar-thumb:hover) {
  background: var(--border-focus);
}
</style>
