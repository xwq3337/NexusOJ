import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore(
  'user',
  () => {
    const id = ref(0)
    const username = ref('')
    const nickname = ref('')
    const gender = ref()
    const avatar = ref('')
    const rating = ref()
    function initStore({ Id, Username, Nickname, Gender, Avatar, Rating }) {
      id.value = Id
      username.value = Username
      nickname.value = Nickname
      gender.value = Gender
      avatar.value = Avatar
      rating.value = Rating
    }

    const resetStore = () => {
      id.value = 0
      username.value = ''
      nickname.value = ''
      gender.value = ''
      avatar.value = ''
      rating.value = ''
    }
    return {
      id,
      username,
      nickname,
      gender,
      avatar,
      rating,
      resetStore,
      initStore
    }
  },
  {
    persist: {
      key: 'User',
      storage: localStorage
    }
  }
)
