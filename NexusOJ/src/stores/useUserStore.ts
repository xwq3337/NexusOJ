import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore(
  'user',
  () => {
    const id = ref('')
    const account = ref('')
    const nickname = ref('')
    const gender = ref()
    const avatar = ref('')
    const rating = ref()
    function initStore({ Id, Account, Nickname, Gender, Avatar, Rating }) {
      id.value = Id
      account.value = Account
      nickname.value = Nickname
      gender.value = Gender
      avatar.value = Avatar
      rating.value = Rating
    }

    const resetStore = () => {
      id.value = ''
      account.value = ''
      nickname.value = ''
      gender.value = ''
      avatar.value = ''
      rating.value = ''
    }
    return {
      id,
      account,
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
