<script setup lang="ts">
import { ref } from 'vue';
import { Settings } from 'lucide-vue-next';
import { NForm, NFormItem, NSelect } from 'naive-ui';
const language = ref('zh')
const theme = ref(localStorage.getItem('theme') || 'light')
const themeOptions = [
    { label: '浅色', value: 'light' },
    { label: '深色', value: 'dark' }
]

const languageOptions = [
    { label: '中文', value: 'zh' },
    { label: 'English', value: 'en' }
]
</script>

<template>
    <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
        <Settings class="inline mr-2" />
        设置
    </h3>
    <n-form label-placement="left" label-width="auto">
        <n-form-item label="语言">
            <n-select v-model:value="language" :options="languageOptions" />
        </n-form-item>
    </n-form>
</template>