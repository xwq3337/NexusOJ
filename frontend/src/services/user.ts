import { ApiResponse } from '@/types/api';
import { User } from '@/types/user'
import Request from './api'

export const userApi = {
  getInfoById: (id: string): Promise<ApiResponse<User>> => {
    return Request.get(`/user/${id}`)
  }
}
