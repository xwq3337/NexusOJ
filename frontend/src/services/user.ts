import { User } from '@/types/user'
import Request from './api'

export const userApi = {
  getInfoById: (id: string): Promise<User> => {
    return Request.get(`/user/${id}`)
  }
}
