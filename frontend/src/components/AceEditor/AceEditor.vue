<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { VAceEditor } from 'vue3-ace-editor'
import './ace.config'
import { useLocalStorage } from '@vueuse/core'
const props = defineProps({
  height: { type: [String, Number], default: '100%' },
  width: { type: [String, Number], default: '100%' },
  language: { type: String as () => ACE_MODE, default: 'c_cpp' },
  theme: { type: String as () => EDITOR_THEHE, default: 'chrome' }
})
const code = defineModel<string>('value', { required: true })

const emit = defineEmits(['change', 'update:languageOptions', 'update:themeOptions', 'onMounted'])
import { EDITOR_THEME_OPTIONS, ACE_MODE_OPTIONS, type ACE_MODE, type LanguageValue, type EDITOR_THEHE } from '@/constants'

onMounted(() => {
  emit('update:languageOptions', ACE_MODE_OPTIONS)
  emit('update:themeOptions', EDITOR_THEME_OPTIONS)
})
watch(code, (newVal) => {
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
  <v-ace-editor v-model:value="code" :lang="language" :theme="theme" :style="{
    height: `${computed_height}`,
    width: `${computed_width}`,
    fontSize: `${fontSize}px`
  }" @init="editorInit" />
</template>
