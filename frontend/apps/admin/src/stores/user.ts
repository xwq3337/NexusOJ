import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore(
  'user',
  () => {
    const id = ref('')
    const nickname = ref('')
    const username = ref('')
    const gender = ref('')
    const avatar = ref('')
    const rating = ref('')
    const resetStore = () => {
      id.value = ''
      nickname.value = ''
      username.value = ''
      gender.value = ''
      avatar.value = ''
      rating.value = ''
    }
    return {
      id,
      nickname,
      username,
      gender,
      avatar,
      rating,
      resetStore,
    }
  },
  {
    persist: {
      key: 'User',
      storage: localStorage,
    },
  },
)
