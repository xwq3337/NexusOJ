<script setup lang="ts">
import { Button } from '@nexusoj/ui'
import { ref } from 'vue'

const handleClick = (event: MouseEvent) => {
  console.log('Button clicked!', event)
}

const isDisabled = ref(false)
</script>

<template>
  <div class="flex flex-col items-center justify-center min-h-screen gap-8 p-8">
    <h1 class="text-3xl font-bold">Shared UI Component Demo</h1>

    <div class="flex gap-4 flex-wrap justify-center">
      <!-- Primary Button -->
      <Button @click="handleClick">
        Primary Button
      </Button>

      <!-- Secondary Button -->
      <Button variant="secondary" @click="handleClick">
        Secondary Button
      </Button>

      <!-- Danger Button -->
      <Button variant="danger" @click="handleClick">
        Danger Button
      </Button>

      <!-- Disabled Button -->
      <Button :disabled="true">
        Disabled Button
      </Button>

      <!-- Different Sizes -->
      <Button size="sm">
        Small Button
      </Button>

      <Button size="lg">
        Large Button
      </Button>

      <!-- Toggle Disabled -->
      <div class="flex items-center gap-4">
        <Button :disabled="isDisabled" @click="handleClick">
          Toggle Disabled
        </Button>
        <button
          class="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300"
          @click="isDisabled = !isDisabled"
        >
          Toggle: {{ isDisabled ? 'Enable' : 'Disable' }}
        </button>
      </div>
    </div>
  </div>
</template>
