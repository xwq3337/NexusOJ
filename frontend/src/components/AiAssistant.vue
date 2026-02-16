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
      class="fixed bottom-6 right-6 w-[calc(100%-3rem)] md:w-96 h-125 rounded-xl shadow-2xl flex flex-col z-50 overflow-hidden"
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
            <span class="text-xs" :style="{ color: 'var(--text-secondary)' }">输入中...</span>
          </div>
        </div>
        <button @click="toggleChat" class="flex items-center justify-center w-8 h-8 rounded-lg transition-colors cursor-pointer"
          :style="{ color: 'var(--text-secondary)' }" :class="`hover:bg-(--surface-secondary)`">
          <X :size="18" />
        </button>
      </div>

      <!-- Messages Container -->
      <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-4">
        <!-- Welcome Message -->
        <div v-if="messages.length === 0"
          class="flex flex-col items-center justify-center h-full text-center space-y-4">
          <div class="w-16 h-16 rounded-full flex items-center justify-center"
            :style="{ backgroundColor: 'var(--accent-color)' }">
            <Bot :size="32" style="color: white" />
          </div>
          <div>
            <h3 class="text-lg font-semibold mb-1" :style="{ color: 'var(--text-primary)' }">
              你好！我是 NexusAI
            </h3>
            <p class="text-sm" :style="{ color: 'var(--text-secondary)' }">
              我可以帮助你理解算法、调试代码和提升编程技能
            </p>
          </div>
          <div class="grid grid-cols-2 gap-2 w-full">
            <button v-for="prompt in quickPrompts" :key="prompt.text" @click="sendQuickPrompt(prompt.text)"
              class="px-3 py-2 rounded-lg text-xs text-left transition-colors" :style="{
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
            <div class="text-sm whitespace-pre-wrap leading-relaxed message-content"
              v-html="renderMessage(message.text)"></div>
            <div class="text-xs mt-1 opacity-70"
              :style="{ color: message.role === 'user' ? 'white' : 'var(--text-tertiary)' }">
              {{ formatTime(message.timestamp) }}
            </div>
          </div>
        </div>

        <!-- Loading Indicator -->
        <div v-if="isLoading" class="flex gap-3">
          <div class="w-8 h-8 rounded-full flex items-center justify-center"
            :style="{ backgroundColor: 'var(--surface-secondary)' }">
            <Bot :size="16" :style="{ color: 'var(--accent-color)' }" />
          </div>
          <div class="flex items-center gap-1 px-4 py-2 rounded-2xl"
            :style="{ backgroundColor: 'var(--surface-secondary)' }">
            <div class="flex gap-1">
              <div class="w-2 h-2 rounded-full"
                :style="{ backgroundColor: 'var(--accent-color)', animation: 'bounce 1.4s infinite ease-in-out both' }">
              </div>
              <div class="w-2 h-2 rounded-full"
                :style="{ backgroundColor: 'var(--accent-color)', animation: 'bounce 1.4s infinite ease-in-out both 0.16s' }">
              </div>
              <div class="w-2 h-2 rounded-full"
                :style="{ backgroundColor: 'var(--accent-color)', animation: 'bounce 1.4s infinite ease-in-out both 0.32s' }">
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Input Area -->
      <div class="p-4 border-t" :style="{ borderColor: 'var(--border-color)' }">
        <form @submit.prevent="sendMessage" class="flex gap-2">
          <input v-model="inputMessage" ref="inputRef" type="text" placeholder="输入你的问题..." :disabled="isLoading"
            class="flex-1 px-4 py-2 rounded-lg outline-none transition-colors text-sm" :style="{
              backgroundColor: 'var(--input-bg)',
              border: '1px solid var(--input-border)',
              color: 'var(--text-primary)'
            }" :class="`focus:border-(--input-focus)`" @input="handleInput" />
          <button type="submit" :disabled="!inputMessage.trim() || isLoading"
            class="px-4 py-2 rounded-lg flex items-center justify-center transition-colors disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
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
            你确定要清空对话记录吗
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
import { ref, computed, nextTick, onMounted, onUnmounted, watch } from 'vue'
import { Bot, X, User, Send, Loader2 } from 'lucide-vue-next'
import { useLocalStorage } from '@vueuse/core'
import { streamGeminiResponse } from '@/services/agent'
import type { ChatMessage } from '@/types/chat'
import MarkdownIt from 'markdown-it'
import { NPopconfirm, NButton } from "naive-ui"
import { useMessage } from 'naive-ui'
const message = useMessage()
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
const messages = ref<ChatMessage[]>([])
const isLoading = ref(false)
const isTyping = ref(false)
const messagesContainer = ref<HTMLElement>()
const inputRef = ref<HTMLInputElement>()
const isMobile = ref(false)

// Quick prompts
const quickPrompts = [
  { text: '如何优化时间复杂度？', label: '⚡ 优化算法' },
  { text: '帮我调试这段代码', label: '🐛 调试代码' },
  { text: '解释这道题的思路', label: '💡 解题思路' },
  { text: '推荐一些练习题', label: '📚 练习推荐' }
]

// Load messages from localStorage
onMounted(() => {
  const saved = localStorage.getItem('ai-assistant-messages')
  if (saved) {
    try {
      messages.value = JSON.parse(saved)
    } catch (e) {
      console.error('Failed to load messages:', e)
    }
  }

  // Check if mobile
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
}

// Save messages to localStorage
watch(messages, (newMessages) => {
  localStorage.setItem('ai-assistant-messages', JSON.stringify(newMessages))
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

// Render message with markdown
const renderMessage = (text: string) => {
  if (!text) return ''

  // Render markdown
  let html = md.render(text)

  // Highlight code blocks
  html = html.replace(
    /<pre><code class="language-(\w+)">([\s\S]*?)<\/code><\/pre>/g,
    (match, lang, code) => {
      return `<pre class="bg-[var(--bg-tertiary)] rounded-lg p-3 overflow-x-auto my-2"><code class="text-sm text-[var(--accent-color)]">${escapeHtml(code.trim())}</code></pre>`
    }
  )

  // Highlight inline code
  html = html.replace(
    /<code>([\s\S]*?)<\/code>/g,
    (match, code) => {
      return `<code class="px-1.5 py-0.5 rounded text-sm" style="background-color: var(--bg-tertiary); color: var(--accent-color);">${escapeHtml(code)}</code>`
    }
  )

  return html
}

// Escape HTML
const escapeHtml = (text: string) => {
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML
}

// Format time
const formatTime = (timestamp: number) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  // Less than 1 minute
  if (diff < 60000) {
    return '刚刚'
  }

  // Less than 1 hour
  if (diff < 3600000) {
    const minutes = Math.floor(diff / 60000)
    return `${minutes}分钟前`
  }

  // Today
  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }

  // Otherwise
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}

// Handle input
const handleInput = () => {
  // Auto-resize if needed
}

// Send quick prompt
const sendQuickPrompt = (prompt: string) => {
  inputMessage.value = prompt
  sendMessage()
}

// Send message with streaming support
const sendMessage = async () => {
  const text = inputMessage.value.trim()
  if (!text || isLoading.value) return

  // Add user message
  const userMessage: ChatMessage = {
    id: Date.now().toString(),
    role: 'user',
    text,
    timestamp: Date.now()
  }
  messages.value.push(userMessage)
  inputMessage.value = ''

  // Scroll to bottom
  scrollToBottom()

  // Set loading state
  isLoading.value = true
  isTyping.value = true

  // Create an initial AI message for streaming
  const aiMessageId = (Date.now() + 1).toString()
  const aiMessage: ChatMessage = {
    id: aiMessageId,
    role: 'model',
    text: '',
    timestamp: Date.now()
  }
  messages.value.push(aiMessage)

  // Find the AI message index
  const aiMessageIndex = messages.value.length - 1

  // Create AbortController for this request
  const abortController = new AbortController()

  try {
    await streamGeminiResponse(
      text,
      // onMessage - append each chunk to the AI message
      (chunk: string) => {
        messages.value[aiMessageIndex].text += chunk
        scrollToBottom()
      },
      // onError
      (error: Error) => {
        console.error('SSE Error:', error)
        messages.value[aiMessageIndex].text = '抱歉，我遇到了一些问题。请稍后再试。'
        scrollToBottom()
      },
      // onClose
      () => {
        isLoading.value = false
        isTyping.value = false
      },
      abortController
    )
  } catch (error) {
    console.error('Failed to get AI response:', error)

    // Update message with error
    messages.value[aiMessageIndex].text = '抱歉，我遇到了一些问题。请稍后再试。'
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
  // Ctrl/Cmd + K to toggle chat
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    toggleChat()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})
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

/* Bounce Animation for Loading */
@keyframes bounce {

  0%,
  80%,
  100% {
    transform: scale(0);
  }

  40% {
    transform: scale(1);
  }
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
