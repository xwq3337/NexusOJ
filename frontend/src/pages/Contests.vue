<template>
  <div class="animate-fade-in max-w-4xl mx-auto">
    <h1 class="text-3xl font-bold mb-6" :style="{ color: 'var(--text-primary)' }">比赛</h1>
    <div class="grid gap-4">
      <div
        v-for="c in MOCK_CONTESTS"
        :key="c.id"
        class="p-4 rounded-lg flex items-center justify-between"
        :style="{
          backgroundColor: 'var(--surface-primary)',
          borderColor: 'var(--border-color)',
          borderWidth: '1px',
          borderStyle: 'solid'
        }"
      >
        <div>
          <h3 class="font-medium" :style="{ color: 'var(--text-primary)' }">
            {{ c.title }}
          </h3>
          <p class="text-xs text-gray-400">
            {{ c.startTime }} • {{ c.duration }} • {{ c.registered }} 人参加
          </p>
        </div>
        <div class="flex items-center gap-2">
          <button
            class="text-xs px-3 py-1.5 rounded"
            :style="{
              color: 'var(--text-primary)',
              backgroundColor: 'var(--btn-secondary)',
              '--tw-bg-opacity': '1'
            }"
            :class="['hover:bg-gray-700']"
          >
            报名
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { MOCK_CONTESTS } from '@/constants'
</script>
