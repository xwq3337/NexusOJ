<template>
  <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
    <div class="p-6 rounded-xl transition-all duration-300 cyber-glow-card"
      :style="{
        backgroundColor: 'var(--surface-primary)',
        border: '1px solid var(--border-color)'
      }">
      <div class="flex items-center gap-4 mb-2">
        <div class="p-3 rounded-lg" :style="{
          backgroundColor: 'rgba(14, 165, 233, 0.1)',
          color: 'var(--accent-color)'
        }">
          <Code2 :size="24" />
        </div>
      </div>
      <h3 class="text-3xl font-bold mb-1 font-terminal" :style="{ color: 'var(--text-primary)' }">2,400+</h3>
      <p class="text-sm" :style="{ color: 'var(--text-tertiary)' }">Problems Available</p>
    </div>

    <div class="p-6 rounded-xl transition-all duration-300 cyber-glow-card"
      :style="{
        backgroundColor: 'var(--surface-primary)',
        border: '1px solid var(--border-color)'
      }">
      <div class="flex items-center gap-4 mb-2">
        <div class="p-3 rounded-lg" :style="{
          backgroundColor: 'rgba(20, 184, 166, 0.1)',
          color: 'var(--secondary)'
        }">
          <Users :size="24" />
        </div>
      </div>
      <h3 class="text-3xl font-bold mb-1 font-terminal" :style="{ color: 'var(--text-primary)' }">150k+</h3>
      <p class="text-sm" :style="{ color: 'var(--text-tertiary)' }">Active Developers</p>
    </div>

    <div class="p-6 rounded-xl transition-all duration-300 cyber-glow-card"
      :style="{
        backgroundColor: 'var(--surface-primary)',
        border: '1px solid var(--border-color)'
      }">
      <div class="flex items-center gap-4 mb-2">
        <div class="p-3 rounded-lg" :style="{
          backgroundColor: 'rgba(74, 222, 128, 0.1)',
          color: 'var(--success-color)'
        }">
          <Zap :size="24" />
        </div>
      </div>
      <h3 class="text-3xl font-bold mb-1 font-terminal" :style="{ color: 'var(--text-primary)' }">12ms</h3>
      <p class="text-sm" :style="{ color: 'var(--text-tertiary)' }">Avg. Judge Time</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Code2, Users, Zap } from 'lucide-vue-next'
</script>
