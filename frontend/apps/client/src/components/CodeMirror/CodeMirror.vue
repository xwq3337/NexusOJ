<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { EditorView, basicSetup } from 'codemirror'
import { EditorState, Compartment } from '@codemirror/state'
import { oneDark } from '@codemirror/theme-one-dark'
import { StreamLanguage } from '@codemirror/language'
import { cpp } from '@codemirror/lang-cpp'
import { python } from '@codemirror/lang-python'
import { java } from '@codemirror/lang-java'
import { javascript } from '@codemirror/lang-javascript'
import { useLocalStorage } from '@vueuse/core'
import {
  EDITOR_THEME_OPTIONS,
  ACE_MODE_OPTIONS,
  type ACE_MODE,
  type EDITOR_THEHE
} from '@/constants'

const props = defineProps({
  height: { type: [String, Number], default: '100%' },
  width: { type: [String, Number], default: '100%' },
  language: { type: String as () => ACE_MODE, default: 'c_cpp' },
  theme: { type: String as () => EDITOR_THEHE, default: 'chrome' }
})

const code = defineModel<string>('value', { required: true })

const emit = defineEmits(['change', 'update:languageOptions', 'update:themeOptions', 'onMounted'])

const editorContainer = ref<HTMLElement>()
let editorView: EditorView | null = null

// 创建配置隔间用于动态切换语言和主题
const languageConf = new Compartment()
const themeConf = new Compartment()

// 语言映射
const languageExtensions: Record<string, any> = {
  c_cpp: cpp(),
  python: python(),
  java: java(),
  javascript: javascript(),
  go: StreamLanguage.define({
    token(stream: any) {
      // Go 基础语法高亮
      if (stream.eatSpace()) return null
      if (stream.match(/\/\/.*/)) return 'lineComment'
      if (stream.match(/\/\*/)) {
        while (!stream.match(/.*?\*\//) && !stream.eol()) stream.next()
        return 'comment'
      }
      if (stream.match(/"(?:[^"\\]|\\.)*"/)) return 'string'
      if (stream.match(/'(?:[^'\\]|\\.)*'/)) return 'string'
      if (stream.match(/\d+(\.\d+)?/)) return 'number'
      if (stream.match(/[a-zA-Z_]\w*/)) {
        const word = stream.current()
        if (/\b(package|import|func|var|const|type|struct|if|else|for|range|return|break|continue|go|select|case|default|switch|defer|fallthrough|goto|interface|map|chan)\b/.test(word)) {
          return 'keyword'
        }
        return 'variableName'
      }
      stream.next()
      return null
    }
  }),
  rust: StreamLanguage.define({
    token(stream: any) {
      // Rust 基础语法高亮
      if (stream.eatSpace()) return null
      if (stream.match(/\/\/.*/)) return 'lineComment'
      if (stream.match(/\/\*/)) {
        while (!stream.match(/.*?\*\//) && !stream.eol()) stream.next()
        return 'comment'
      }
      if (stream.match(/"(?:[^"\\]|\\.)*"/)) return 'string'
      if (stream.match(/'(?:[^'\\]|\\.)*'/)) return 'string'
      if (stream.match(/\d+(\.\d+)?/)) return 'number'
      if (stream.match(/[a-zA-Z_]\w*/)) {
        const word = stream.current()
        if (/\b(fn|let|mut|const|if|else|for|while|loop|match|return|break|continue|struct|enum|impl|trait|type|use|mod|crate|pub|unsafe|where|move|ref|static)\b/.test(word)) {
          return 'keyword'
        }
        return 'variableName'
      }
      stream.next()
      return null
    }
  }),
  csharp: StreamLanguage.define({
    token(stream: any) {
      // C# 基础语法高亮
      if (stream.eatSpace()) return null
      if (stream.match(/\/\/.*/)) return 'lineComment'
      if (stream.match(/\/\*/)) {
        while (!stream.match(/.*?\*\//) && !stream.eol()) stream.next()
        return 'comment'
      }
      if (stream.match(/"(?:[^"\\]|\\.)*"/)) return 'string'
      if (stream.match(/'(?:[^'\\]|\\.)*'/)) return 'string'
      if (stream.match(/\d+(\.\d+)?/)) return 'number'
      if (stream.match(/[a-zA-Z_]\w*/)) {
        const word = stream.current()
        if (/\b(using|namespace|class|struct|interface|enum|delegate|event|where|if|else|switch|case|default|for|foreach|while|do|return|break|continue|goto|try|catch|finally|throw|new|this|base|null|true|false|static|public|private|protected|internal|virtual|override|abstract|sealed|readonly|const|void|int|string|bool|double|float|long|short|byte|char|object|var)\b/.test(word)) {
          return 'keyword'
        }
        return 'variableName'
      }
      stream.next()
      return null
    }
  }),
  text: StreamLanguage.define({
    token(stream: any) {
      stream.next()
      return null
    }
  })
}

// 主题映射
const getThemeExtension = (theme: string) => {
  const darkThemes = ['monokai', 'dracula', 'terminal']
  return darkThemes.includes(theme) ? oneDark : []
}

const computedHeight = computed(() => {
  if (typeof props.height === 'number') return `${props.height}px`
  if (props.height.includes('%')) {
    return props.height
  }
  return `${props.height}px`
})

const computedWidth = computed(() => {
  if (typeof props.width === 'number') return `${props.width}px`
  if (props.width.includes('%')) {
    return props.width
  }
  return `${props.width}px`
})

const fontSize = useLocalStorage('editor_font_size', String(14))

onMounted(() => {
  emit('update:languageOptions', ACE_MODE_OPTIONS)
  emit('update:themeOptions', EDITOR_THEME_OPTIONS)

  if (!editorContainer.value) return

  const extensions = [
    basicSetup,
    languageConf.of(languageExtensions[props.language] || languageExtensions.c_cpp),
    themeConf.of(getThemeExtension(props.theme)),
    EditorView.updateListener.of((update: any) => {
      if (update.docChanged) {
        const newValue = update.state.doc.toString()
        code.value = newValue
        emit('change', newValue)
      }
    }),
    EditorView.theme({
      '&': {
        fontSize: `${fontSize.value}px`,
        height: computedHeight.value,
        width: computedWidth.value
      },
      '.cm-scroller': {
        overflow: 'auto',
        height: '100%'
      },
      '.cm-content': {
        padding: '10px 0'
      },
      '.cm-line': {
        padding: '0 10px'
      }
    })
  ]

  const state = EditorState.create({
    doc: code.value,
    extensions
  })

  editorView = new EditorView({
    state,
    parent: editorContainer.value
  })

  emit('onMounted', editorView)
})

onUnmounted(() => {
  if (editorView) {
    editorView.destroy()
    editorView = null
  }
})

// 监听语言变化
watch(
  () => props.language,
  (newLanguage) => {
    if (editorView) {
      const languageExtension = languageExtensions[newLanguage] || languageExtensions.c_cpp
      editorView.dispatch({
        effects: languageConf.reconfigure(languageExtension)
      })
    }
  }
)

// 监听主题变化
watch(
  () => props.theme,
  (newTheme) => {
    if (editorView) {
      editorView.dispatch({
        effects: themeConf.reconfigure(getThemeExtension(newTheme))
      })
    }
  }
)

// 监听外部代码变化
watch(
  code,
  (newCode) => {
    if (editorView && editorView.state.doc.toString() !== newCode) {
      const transaction = editorView.state.update({
        changes: {
          from: 0,
          to: editorView.state.doc.length,
          insert: newCode
        }
      })
      editorView.dispatch(transaction)
    }
  }
)

// 暴露编辑器实例
defineExpose({
  getEditor: () => editorView
})
</script>

<template>
  <div ref="editorContainer" class="codemirror-container"></div>
</template>

<style scoped>
.codemirror-container {
  height: v-bind(computedHeight);
  width: v-bind(computedWidth);
  border: 1px solid #d1d5db;
  border-radius: 4px;
  overflow: hidden;
}

:deep(.cm-editor) {
  height: 100%;
}

:deep(.cm-scroller) {
  height: 100%;
  overflow: auto;
}

:deep(.cm-content) {
  padding: 10px 0;
}

:deep(.cm-focused) {
  outline: none;
}

:deep(.cm-line) {
  padding: 0 10px;
}
</style>
