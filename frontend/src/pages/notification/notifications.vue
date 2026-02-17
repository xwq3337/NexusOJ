<template>
    <div class="min-h-screen p-6">
        <div class="max-w-7xl mx-auto">
            <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
                <!-- Sidebar -->
                <div class="lg:col-span-1">
                    <div class="rounded-lg shadow-md overflow-hidden" :style="{ backgroundColor: 'var(--card-bg)' }">
                        <div class="p-6">
                            <h2 class="text-xl font-semibold mb-4" :style="{ color: 'var(--text-primary)' }">
                                通知中心
                            </h2>
                            <n-menu :options="menuOptions" :value="activeKey" />
                        </div>
                    </div>
                </div>

                <!-- Main Content -->
                <div class="lg:col-span-3">
                    <div class="rounded-lg shadow-md p-6" :style="{ backgroundColor: 'var(--card-bg)' }">
                        <RouterView />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, h } from 'vue'
import { useRoute } from 'vue-router'
import { Bell, MessageSquare } from 'lucide-vue-next'
import { NMenu, type MenuOption } from 'naive-ui'
import { RouterLink } from 'vue-router'

const route = useRoute()

const menuOptions: MenuOption[] = [
    {
        label: () =>
            h(
                RouterLink,
                { to: { name: 'Message' } },
                { default: () => '私信' }
            ),
        key: 'message',
        icon: () => h(MessageSquare)
    },
    {
        label: () =>
            h(
                RouterLink,
                { to: { name: 'Info' } },
                { default: () => '系统通知' }
            ),
        key: 'info',
        icon: () => h(Bell)
    },

]

const activeKey = computed(() => {
    const name = route.name as string
    if (name === 'Info') return 'info'
    if (name === 'Message') return 'message'
    return 'info'
})
</script>
