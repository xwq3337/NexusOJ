<template>
  <div
    class="min-h-screen text-slate-100 flex flex-col"
    :style="{ backgroundColor: 'var(--bg-primary)' }"
  >
    <Navbar />
    <main class="flex-1 container mx-auto px-4 py-3">
      <router-view />
    </main>
    <footer
      class="border-t py-6 mt-8"
      :style="{
        borderColor: 'var(--border-color)',
        borderWidth: '1px',
        borderStyle: 'solid'
      }"
    >
      <div class="container mx-auto px-4 text-center text-gray-600 text-sm">
        &copy; 2025 NexusOJ. Built with Vue & Tailwind.
      </div>
    </footer>
    <AiAssistant />
  </div>
</template>

<script setup lang="ts">
import Navbar from '../components/Navbar.vue'
import AiAssistant from '../components/AiAssistant.vue'
</script>
