import { ApiResponse } from '@/types'
import { User, FriendShip, FriendShipRequest } from '@/types/user'
import Request from './api'
import { GetRecordListParams, GetRecordListResponse } from '@/types/record'
import { ChatMessage } from '@/types/chat'

export const userApi = {
  Login: (username: string, password: string): Promise<ApiResponse<User>> => {
    return Request.post('/user/login', { username, password })
  },
  Register: (username: string, password: string, nickname: string, email?: string): Promise<ApiResponse<User>> => {
    return Request.post('/user/register', { username, password, nickname, email })
  },
  ValidateToken: (): Promise<ApiResponse<string>> => {
    return Request.get('/user/validate-token')
  },
  getInfoById: (id: string): Promise<ApiResponse<User>> => {
    return Request.get(`/user/${id}`)
  },
  updateAvatar: (file: File): Promise<ApiResponse<string>> => {
    const formData = new FormData()
    formData.append('avatar', file)
    return Request.post('/user/update-avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },
  updatePassword: (old_password: string, new_password: string): Promise<ApiResponse<null>> => {
    return Request.post('/user/update-password', { old_password, new_password })
  },
  updateUser: (
    data: Partial<Omit<User, 'id' | 'created_at' | 'password' | 'avatar'>>
  ): Promise<ApiResponse<User>> => {
    return Request.post('/user/update', data)
  },
  // 获取聊天好友列表
  getFriendList: (): Promise<ApiResponse<FriendShip[]>> => {
    return Request.get(`/user/friend-list`)
  },
  // 获取新的好友请求列表
  getFriendRequestList: (): Promise<ApiResponse<FriendShipRequest[]>> => {
    return Request.get(`/user/friend-request-list`)
  },
  // 处理好友请求，request_id是好友请求的ID，status表示是否接受
  HandleFriendRequest: (
    request_id: string,
    status: 'accepted' | 'rejected'
  ): Promise<ApiResponse<null>> => {
    return Request.post(`/user/handle-friend-request`, { request_id, status })
  },
  // 发送好友请求，message是验证消息，可以为空
  FriendRequest: (friend_id: string, message?: string): Promise<ApiResponse<null>> => {
    return Request.post(`/user/friend-request`, { friend_id, message })
  },
  //  搜索用户，根据ID，用户名，昵称进行模糊搜索
  searchUser: (keyword: string): Promise<ApiResponse<User[]>> => {
    return Request.get(`/user/search`, { params: { keyword } })
  },
  Count: (): Promise<ApiResponse<number>> => {
    return Request.get(`/user/count`)
  },
  getUserRecordList: (
    id: string,
    params: Partial<GetRecordListParams>
  ): Promise<ApiResponse<GetRecordListResponse[]>> => {
    return Request.get(`/record/user/${id}`, { params: params })
  },
  getChatRecordList: (
    friend_id: string,
    page: number = 1
  ): Promise<ApiResponse<ChatMessage[]>> => {
    return Request.get(`/chat/record`, {
      params: { friend_id, page }
    })
  }
}
