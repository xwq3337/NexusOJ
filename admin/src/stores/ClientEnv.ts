import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useClientEnv = defineStore(
  'clientEnv',
  () => {
    const os = ref<any>(null)
    const ip = ref<any>(null)
    const brower = ref<any>(null)
    const language = ref<any>(null)
    const geolocation = ref<any>(null)
    const innerWidth = ref<any>(null)
    const innerHeight = ref<any>(null)
    const screen = ref<any>(null)
    const updateClientEnv = ({
      Os,
      Ip,
      Brower,
      Language,
      Geolocation,
      InnerWidth,
      InnerHeight,
      Screen
    }) => {
      os.value = Os || os.value
      ip.value = Ip || ip.value
      brower.value = Brower || brower.value
      language.value = Language || language.value
      geolocation.value = Geolocation || geolocation.value
      innerWidth.value = InnerWidth || innerWidth.value
      innerHeight.value = InnerHeight || innerHeight.value
      screen.value = Screen || screen.value
    }
    const clearClientEnv = () => {
      os.value = ip.value = null
      brower.value = language.value = null
    }
    return {
      os,
      ip,
      brower,
      language,
      geolocation,
      innerHeight,
      innerWidth,
      screen,
      updateClientEnv,
      clearClientEnv
    }
  },
  {
    persist: {
      key: 'ClientEnv',
      storage: sessionStorage
    }
  }
)
