import { ref, watch, onMounted, computed } from 'vue'
import { darkTheme, lightTheme, useOsTheme, type GlobalThemeOverrides } from 'naive-ui'
import { RemovableRef, useLocalStorage } from '@vueuse/core'

const lightOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#0ea5e9',
    primaryColorHover: '#38bdf8',
    primaryColorPressed: '#0284c7',
    primaryColorSuppl: '#0ea5e9',
    infoColor: '#0ea5e9',
    infoColorHover: '#38bdf8',
    infoColorPressed: '#0284c7',
    successColor: '#10b981',
    successColorHover: '#34d399',
    successColorPressed: '#059669',
    warningColor: '#f59e0b',
    warningColorHover: '#fbbf24',
    warningColorPressed: '#d97706',
    errorColor: '#ef4444',
    errorColorHover: '#f87171',
    errorColorPressed: '#dc2626',
  },
  Button: {
    colorPrimary: '#0ea5e9',
    colorHoverPrimary: '#38bdf8',
    colorPressedPrimary: '#0284c7',
    borderPrimary: '1px solid #0ea5e9',
    borderHoverPrimary: '1px solid #38bdf8',
    borderPressedPrimary: '1px solid #0284c7',
  },
  Input: {
    borderFocus: '1px solid #0ea5e9',
    borderHover: '1px solid #38bdf8',
    boxShadowFocus: '0 0 0 2px rgba(14, 165, 233, 0.15)',
    caretColor: '#0ea5e9',
  },
  Switch: {
    railColorActive: '#0ea5e9',
  },
  Tag: {
    colorPrimary: 'rgba(14, 165, 233, 0.08)',
    textColorPrimary: '#0284c7',
    borderPrimary: '1px solid rgba(14, 165, 233, 0.2)',
  },
}

const darkOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#38bdf8',
    primaryColorHover: '#7dd3fc',
    primaryColorPressed: '#0ea5e9',
    primaryColorSuppl: '#38bdf8',
    infoColor: '#38bdf8',
    infoColorHover: '#7dd3fc',
    infoColorPressed: '#0ea5e9',
    successColor: '#4ade80',
    successColorHover: '#86efac',
    successColorPressed: '#22c55e',
    warningColor: '#fbbf24',
    warningColorHover: '#fcd34d',
    warningColorPressed: '#f59e0b',
    errorColor: '#f87171',
    errorColorHover: '#fca5a5',
    errorColorPressed: '#ef4444',
  },
  Button: {
    colorPrimary: '#0ea5e9',
    colorHoverPrimary: '#38bdf8',
    colorPressedPrimary: '#0284c7',
    borderPrimary: '1px solid #38bdf8',
    borderHoverPrimary: '1px solid #7dd3fc',
    borderPressedPrimary: '1px solid #0ea5e9',
  },
  Input: {
    borderFocus: '1px solid #38bdf8',
    borderHover: '1px solid #7dd3fc',
    boxShadowFocus: '0 0 0 2px rgba(56, 189, 248, 0.15)',
    caretColor: '#38bdf8',
  },
  Switch: {
    railColorActive: '#38bdf8',
  },
  Tag: {
    colorPrimary: 'rgba(56, 189, 248, 0.1)',
    textColorPrimary: '#7dd3fc',
    borderPrimary: '1px solid rgba(56, 189, 248, 0.2)',
  },
}

export function useTheme() {

  const theme = useLocalStorage('theme', useOsTheme().value) as RemovableRef<'light' | 'dark'>
  const naiveTheme = ref(theme.value === 'dark' ? darkTheme : lightTheme)
  const naiveThemeOverrides = computed(() => theme.value === 'dark' ? darkOverrides : lightOverrides)

  const updateBodyClass = (newTheme: 'light' | 'dark') => {
    document.body.classList.remove('light', 'dark')
    document.documentElement.classList.remove('light', 'dark')
    document.body.classList.add(newTheme)
    document.documentElement.classList.add(newTheme)
    naiveTheme.value = newTheme === 'dark' ? darkTheme : lightTheme
  }

  onMounted(() => {
    updateBodyClass(theme.value)
  })

  watch(theme, (newTheme) => {
    localStorage.setItem('theme', newTheme)
    updateBodyClass(newTheme)
  })

  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }

  const setTheme = (newTheme: 'light' | 'dark') => {
    theme.value = newTheme
  }

  return {
    theme,
    naiveTheme,
    naiveThemeOverrides,
    toggleTheme,
    setTheme
  }
}
