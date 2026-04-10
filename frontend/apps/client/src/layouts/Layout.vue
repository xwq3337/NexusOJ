<template>
  <div
    class="min-h-screen flex flex-col"
    :style="{ backgroundColor: 'var(--bg-primary)' }"
  >
    <CyberGridCanvas />
    <Navbar />
    <main class="flex-1 container mx-auto px-4 py-3 relative" style="z-index: 1">
      <router-view />
    </main>
    <footer
      class="py-6 mt-8 relative"
      style="z-index: 1"
      :style="{
        borderTop: '1px solid var(--border-color)',
        boxShadow: '0 -4px 20px rgba(14, 165, 233, 0.05)'
      }"
    >
      <div class="container mx-auto px-4 text-center text-sm" :style="{ color: 'var(--text-tertiary)' }">
        <span class="font-terminal text-xs" :style="{ color: 'var(--accent-color)' }">&gt;_</span>
        &copy; 2025 NexusOJ. Built with Vue & Tailwind.
      </div>
    </footer>
    <AiAssistant />
  </div>
</template>

<script setup lang="ts">
import Navbar from '@/components/Navbar.vue'
import AiAssistant from '@/components/AiAssistant.vue'
import CyberGridCanvas from '@/components/CyberGridCanvas.vue'
</script>
