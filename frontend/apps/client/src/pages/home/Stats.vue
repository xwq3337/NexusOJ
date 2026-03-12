<template>
  <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
    <div class="p-6 border-2 border-transparent rounded-xl hover:border-blue-500  transition-colors"
      :style="{ backgroundColor: 'var(--surface-primary)' }">
      <div class="flex items-center gap-4 mb-2">
        <div class="p-3 bg-blue-500/10 rounded-lg text-blue-400">
          <Code2 :size="24" />
        </div>
      </div>
      <h3 class="text-3xl font-bold mb-1" :style="{ color: 'var(--text-primary)' }">2,400+</h3>
      <p class="text-gray-400 text-sm">Problems Available</p>
    </div>

    <div class="p-6 border-2 border-transparent rounded-xl hover:border-purple-500/50 transition-colors"
      :style="{ backgroundColor: 'var(--surface-primary)' }">
      <div class="flex items-center gap-4 mb-2">
        <div class="p-3 bg-purple-500/10 rounded-lg text-purple-400">
          <Users :size="24" />
        </div>
      </div>
      <h3 class="text-3xl font-bold mb-1" :style="{ color: 'var(--text-primary)' }">150k+</h3>
      <p class="text-gray-400 text-sm">Active Developers</p>
    </div>

    <div class="p-6 border-2 border-transparent rounded-xl hover:border-green-500/50 transition-colors"
      :style="{ backgroundColor: 'var(--surface-primary)' }">
      <div class="flex items-center gap-4 mb-2">
        <div class="p-3 bg-green-500/10 rounded-lg text-green-400">
          <Zap :size="24" />
        </div>
      </div>
      <h3 class="text-3xl font-bold mb-1" :style="{ color: 'var(--text-primary)' }">12ms</h3>
      <p class="text-gray-400 text-sm">Avg. Judge Time</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Code2, Users, Zap } from 'lucide-vue-next'
</script>
