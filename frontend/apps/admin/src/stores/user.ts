import { defineStore } from 'pinia'
import { ref } from 'vue'
import { pick, forEach } from 'lodash-es'
export const useUserStore = defineStore(
  'user',
  () => {
    const id = ref(0)
    const nickname = ref('')
    const username = ref('')
    const gender = ref('')
    const avatar = ref('')
    const rating = ref(0)

    function initStore( obj : { id : number, username: string, nickname : string, gender : string, avatar: string, rating: number } ){
      const pickedObj = pick(obj, ['id', 'username', 'nickname', 'gender', 'avatar', 'rating'])
      forEach(pickedObj, (val, key) => {
        eval(key).value = val
      })
    }
    const resetStore = () => {
      id.value = 0
      nickname.value = ''
      username.value = ''
      gender.value = ''
      avatar.value = ''
      rating.value = 0
    }
    return {
      id,
      nickname,
      username,
      gender,
      avatar,
      rating,
      initStore,
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
