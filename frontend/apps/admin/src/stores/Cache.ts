import { defineStore } from 'pinia'
import { ref } from 'vue'
interface DynamicObject {
  [key: string]: any // 或更具体的类型如 string | number
}

export const useCache = defineStore(
  'cache',
  () => {
    const cache = ref<DynamicObject>({})

    const resetCache = () => {
      cache.value = {}
    }
    const addCache = (key: string, value: any) => {
      cache.value[key] = value
    }
    const getCache = (key: string) => {
      return cache.value[key]
    }
    return {
      cache,
      resetCache,
      addCache,
      getCache,
    }
  },
  {
    persist: {
      key: 'cache',
      storage: localStorage,
    },
  },
)
