<script setup lang="ts">
import { computed, onMounted, watch, ref } from 'vue'
import { VAceEditor } from 'vue3-ace-editor'
import './ace.config'
import { useLocalStorage } from '@vueuse/core'
const props = defineProps({
  height: { type: [String, Number], default: '100%' },
  width: { type: [String, Number], default: '100%' },
  language: { type: String as () => LanguageValue, default: 'c_cpp' },
  theme: { type: String as () => EDITOR_THEHE, default: 'chrome' }
})
const code = defineModel<string>('value', { required: true })

const emit = defineEmits(['change', 'update:languageOptions', 'update:themeOptions', 'onMounted'])
import { EDITOR_THEME_OPTIONS, ACE_MODE_OPTIONS, type ACE_MODE, type EDITOR_THEHE, LanguageValue, LANGUAGE_CONFIG } from '@/constants'
import { ideApi } from '@nexusoj/server'

// 编辑器实例引用
const editorInstance = ref<any>(null)
const formatterLanguageList  = ref<string[]>([])
onMounted(() => {
  // ideApi.FormatLanguageList().then((res) => {
  //   formatterLanguageList.value = res.info
  // })
  emit('update:languageOptions', ACE_MODE_OPTIONS)
  emit('update:themeOptions', EDITOR_THEME_OPTIONS)
})
watch(code, (newVal) => {
  emit('change', newVal)
}
)
const editorInit = (editor: any) => {
  editorInstance.value = editor
  editor.setOptions({
    useWorker: true,
    enableBasicAutocompletion: true,
    enableSnippets: true,
    enableLiveAutocompletion: true
  })

  // 添加格式化命令
  editor.commands.addCommand({
    name: 'formatCode',
    bindKey: { win: 'Shift-Alt-F', mac: 'Shift-Control-F' },
    exec: async () => {
      console.log('格式化代码')
      formateCode(code.value, props.language)
    }
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


const formateCode = async (sourceCode : string, language: string ) => {
    return ideApi.FormatCode(sourceCode, language ).then((res) => {
      const {code: statusCode, info} = res
      if (statusCode === 200) {
        code.value = info
      }
    })
}
</script>
<template>
  <v-ace-editor v-model:value="code" :lang="LANGUAGE_CONFIG[language].aceMode" :theme="theme" :style="{
    height: `${computed_height}`,
    width: `${computed_width}`,
    fontSize: `${fontSize}px`
  }" @init="editorInit" />
</template>
