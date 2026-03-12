import { ref, watch, onUnmounted } from 'vue'

const originalTitle = ref(document.title)

export function useTitle(newTitle?: string) {
  const title = ref(newTitle ?? document.title)

  const setTitle = (t: string) => {
    title.value = t
    document.title = t
  }

  if (newTitle) {
    setTitle(newTitle)
  }

  watch(title, (newTitle) => {
    document.title = newTitle
  })

  onUnmounted(() => {
    document.title = originalTitle.value
  })

  return {
    title,
    setTitle
  }
}
