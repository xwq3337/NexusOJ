import { forEach, pick } from 'lodash'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore(
  'user',
  () => {
    const id = ref(0)
    const username = ref('')
    const nickname = ref('')
    const gender = ref('')
    const avatar = ref('')
    const rating = ref(0)
    function initStore( obj : { id : Number, username: string, nickname : string, gender : string, avatar: string, rating: number } ) {
      const pickedObj = pick(obj, ['id', 'username', 'nickname', 'gender', 'avatar', 'rating'])
      forEach(pickedObj, (val, key) => {
        eval(key).value = val
      })
    }

    const resetStore = () => {
      id.value = 0
      username.value = ''
      nickname.value = ''
      gender.value = ''
      avatar.value = ''
      rating.value = 0
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
      key: 'user_cache',
      storage: localStorage
    }
  }
)
