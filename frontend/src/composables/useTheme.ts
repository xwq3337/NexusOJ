import { ref, watch, onMounted, Ref } from 'vue'
import { darkTheme, lightTheme, useOsTheme } from 'naive-ui'
import { RemovableRef, useLocalStorage } from '@vueuse/core'
export function useTheme() {
  const theme = useLocalStorage('theme', useOsTheme().value) as RemovableRef<'light' | 'dark'>
  const naiveTheme = ref(theme.value === 'dark' ? darkTheme : lightTheme)

  // 更新 body class 的函数
  const updateBodyClass = (newTheme: 'light' | 'dark') => {
    // 移除旧的 class
    document.body.classList.remove('light', 'dark')
    document.documentElement.classList.remove('light', 'dark')

    // 添加新的 class
    document.body.classList.add(newTheme)
    document.documentElement.classList.add(newTheme)

    // 更新 naive-ui 主题
    naiveTheme.value = newTheme === 'dark' ? darkTheme : lightTheme
  }

  // 初始化主题
  onMounted(() => {
    updateBodyClass(theme.value)
  })

  // 监听主题变化
  watch(theme, (newTheme) => {
    localStorage.setItem('theme', newTheme)
    // // 刷新页面
    // location.reload()
    // history.go(0)
    updateBodyClass(newTheme)
  })

  // 切换主题函数
  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }

  // 设置主题函数
  const setTheme = (newTheme: 'light' | 'dark') => {
    theme.value = newTheme
  }

  return {
    theme,
    naiveTheme,
    toggleTheme,
    setTheme
  }
}
