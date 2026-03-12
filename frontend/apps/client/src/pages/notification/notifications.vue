<template>
    <div class="min-h-screen p-6">
        <div class="max-w-7xl mx-auto">
            <div class="rounded-lg shadow-md" :style="{ backgroundColor: 'var(--card-bg)' }">
                <RouterView />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
// No additional logic needed - navigation is now handled via navbar dropdown
</script>
