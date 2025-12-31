<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { VAceEditor } from 'vue3-ace-editor'
import './ace.config'
import { useLocalStorage } from '@vueuse/core'
const props = defineProps({
  height: { type: [String, Number], default: '100%' },
  width: { type: [String, Number], default: '100%' },
  language: { type: String as () => LanguageOption, default: 'c_cpp' },
  theme: { type: String as () => ThemeOption, default: 'chrome' }
})
const code = defineModel<string>('value', { required: true })

const emit = defineEmits(['change', 'update:languageOptions', 'update:themeOptions', 'onMounted'])

const languageOptions = ['c_cpp', 'javascript', 'sql', 'text', 'python', 'java', 'csharp'] as const
type LanguageOption = (typeof languageOptions)[number]
const themeOptions = [
  'chrome',
  'github',
  'monokai',
  // 'solarized_dark',
  // 'solarized_light',
  // 'tomorrow_night',
  'xcode',
  'dracula',
  'clouds',
  'terminal'
] as const

type ThemeOption = (typeof themeOptions)[number]
onMounted(() => {
  emit('update:languageOptions', languageOptions)
  emit('update:themeOptions', themeOptions)
})
watch(
  () => code.value,
  (newVal) => {
    emit('change', newVal)
  }
)
const editorInit = (editor: any) => {
  editor.setOptions({
    useWorker: true,
    enableBasicAutocompletion: true,
    enableSnippets: true,
    enableLiveAutocompletion: true
  })
  emit('onMounted', editor)
}
const computed_height = computed(() => {
  if (typeof props.height === 'number') return `${props.height}px`
  if (props.height.includes('%')) {
    return props.height
  }
  return `${props.height}px`
})
const computed_width = computed(() => {
  if (typeof props.width === 'number') return `${props.width}px`
  if (props.width.includes('%')) {
    return props.width
  }
  return `${props.width}px`
})
const fontSize = useLocalStorage('editor_font_size', String(14))
</script>
<template>
  <v-ace-editor
    v-model:value="code"
    :lang="language"
    :theme="theme"
    :style="{
      height: `${computed_height}`,
      width: `${computed_width}`,
      fontSize: `${fontSize}px`
    }"
    @init="editorInit"
  />
</template>
