<template>
  <div>
    <button
      v-if="!isOpen"
      @click="isOpen = true"
      class="fixed bottom-6 right-6 w-14 h-14 bg-blue-600 hover:bg-blue-700 rounded-full shadow-lg flex items-center justify-center text-white transition-all hover:scale-110 z-50 group"
    >
      <Bot :size="28" />
      <span
        class="absolute right-full mr-4 text-white px-3 py-1 rounded text-sm whitespace-nowrap opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none"
        :style="{ backgroundColor: 'var(--surface-primary)' }"
        >Ask NexusAI</span
      >
    </button>

    <div
      v-else
      :class="[
        'fixed bottom-6 right-6 shadow-2xl rounded-2xl flex flex-col z-50 transition-all duration-300',
        isMinimized ? 'w-72 h-14' : 'w-[380px] h-[600px] max-h-[80vh]'
      ]"
      :style="{
        backgroundColor: 'var(--surface-primary)',
        borderColor: 'var(--border-color)',
        borderWidth: '1px',
        borderStyle: 'solid'
      }"
    >
      <div
        class="p-4 border-b flex items-center justify-between rounded-t-2xl cursor-pointer"
        @click="!isMinimized && (isMinimized = !isMinimized)"
        :style="{
          backgroundColor: 'var(--surface-secondary)',
          borderColor: 'var(--border-color)',
          borderWidth: '1px',
          borderStyle: 'solid'
        }"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-8 h-8 bg-blue-600/20 text-blue-400 rounded-lg flex items-center justify-center"
          >
            <Bot :size="20" />
          </div>
          <div>
            <h3 class="font-bold text-white text-sm">NexusAI Assistant</h3>
            <p class="text-xs text-green-400 flex items-center gap-1">
              <span class="w-1.5 h-1.5 bg-green-400 rounded-full animate-pulse"></span>Online
            </p>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <button
            @click.stop="isMinimized = !isMinimized"
            class="text-gray-400 hover:text-white p-1"
          >
            <component :is="isMinimized ? Maximize2 : Minimize2" :size="16" />
          </button>
          <button @click="isOpen = false" class="text-gray-400 hover:text-white p-1">
            <X :size="18" />
          </button>
        </div>
      </div>

      <div v-if="!isMinimized">
        <div
          class="flex-1 overflow-y-auto p-4 space-y-4"
          :style="{ backgroundColor: 'var(--surface-tertiary)' }"
        >
          <div
            v-for="msg in messages"
            :key="msg.id"
            :class="['flex', msg.role === 'user' ? 'justify-end' : 'justify-start']"
          >
            <div
              :class="[
                'max-w-[85%] rounded-2xl px-4 py-3 text-sm leading-relaxed',
                msg.role === 'user'
                  ? 'bg-blue-600 text-white rounded-br-none'
                  : 'text-gray-200 rounded-bl-none border'
              ]"
              :style="
                msg.role === 'user'
                  ? {}
                  : {
                      backgroundColor: 'var(--surface-tertiary)',
                      borderColor: 'var(--border-light)',
                      borderWidth: '1px',
                      borderStyle: 'solid'
                    }
              "
            >
              {{ msg.text }}
            </div>
          </div>
          <div v-if="isLoading" class="flex justify-start">
            <div
              class="text-gray-400 rounded-2xl rounded-bl-none px-4 py-3 flex items-center gap-2"
              :style="{
                backgroundColor: 'var(--surface-tertiary)',
                borderColor: 'var(--border-light)',
                borderWidth: '1px',
                borderStyle: 'solid'
              }"
            >
              <Loader2 class="animate-spin" :size="14" />
              <span class="text-xs">NexusAI is thinking...</span>
            </div>
          </div>
          <div ref="messagesEndRef" />
        </div>

        <div
          class="p-4 border-t rounded-b-2xl"
          :style="{
            backgroundColor: 'var(--surface-primary)',
            borderColor: 'var(--border-color)',
            borderWidth: '1px',
            borderStyle: 'solid'
          }"
        >
          <div class="relative">
            <input
              v-model="input"
              @keydown.enter="handleSend"
              type="text"
              placeholder="Ask about an algorithm..."
              class="w-full text-white rounded-xl pl-4 pr-12 py-3 text-sm focus:outline-none focus:border-blue-500 placeholder-gray-500"
              :style="{
                backgroundColor: 'var(--input-bg)',
                borderColor: 'var(--input-border)',
                borderWidth: '1px',
                borderStyle: 'solid'
              }"
            />
            <button
              @click="handleSend"
              :disabled="isLoading || !input.trim()"
              class="absolute right-2 top-1/2 -translate-y-1/2 p-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <Send :size="16" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Bot, X, Send, Loader2, Minimize2, Maximize2 } from 'lucide-vue-next'
import { getGeminiResponse } from '@/services/gemini'
import type { ChatMessage } from '@/types'

const isOpen = ref(false)
const isMinimized = ref(false)
const input = ref('')
const messages = ref<ChatMessage[]>([
  {
    id: '0',
    role: 'model',
    text: 'Hello! I am NexusAI. Ask me anything about algorithms or debugging.',
    timestamp: Date.now()
  }
])
const isLoading = ref(false)
const messagesEndRef = ref<HTMLElement | null>(null)

const scrollToBottom = () => {
  messagesEndRef.value?.scrollIntoView({ behavior: 'smooth' })
}

watch(messages, () => scrollToBottom())

const handleSend = async () => {
  if (!input.value.trim() || isLoading.value) return

  const userMsg: ChatMessage = {
    id: Date.now().toString(),
    role: 'user',
    text: input.value,
    timestamp: Date.now()
  }

  messages.value.push(userMsg)
  input.value = ''
  isLoading.value = true

  try {
    const responseText = await getGeminiResponse(userMsg.text)
    const botMsg: ChatMessage = {
      id: (Date.now() + 1).toString(),
      role: 'model',
      text: responseText,
      timestamp: Date.now()
    }
    messages.value.push(botMsg)
  } catch (err) {
    console.error(err)
  } finally {
    isLoading.value = false
  }
}
</script>
