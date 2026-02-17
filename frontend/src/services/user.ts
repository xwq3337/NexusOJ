import { ApiResponse } from '@/types/api'
import { User, FriendShip } from '@/types/user'
import Request from './api'

export const userApi = {
  getInfoById: (id: string): Promise<ApiResponse<User>> => {
    return Request.get(`/user/${id}`)
  },
  updateAvatar: (file: File): Promise<ApiResponse<{ url: string }>> => {
    const formData = new FormData()
    formData.append('avatar', file)
    return Request.post('/user/update_avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },
  updateUser: (
    data: Partial<Omit<User, 'id' | 'created_at' | 'password' | 'avatar'>>
  ): Promise<ApiResponse<User>> => {
    return Request.post('/user/update', data)
  },
  // 获取聊天好友列表，TODO: 可能需要分页
  getFriendList: (): Promise<ApiResponse<User[]>> => {
    return Request.get(`/user/friend_list`)
  },
  // 获取新的好友请求列表
  getFriendRequestList: (): Promise<
    ApiResponse<
      (FriendShip & {
        username: string
        nickname: string
        avatar: string
        rating: number
      })[]
    >
  > => {
    return Request.get(`/user/friend_request_list`)
  },
  // 处理好友请求，friendship_id是好友请求的ID，status表示是否接受
  HandleFriendRequest: (friendship_id: string, accept: boolean): Promise<ApiResponse<null>> => {
    return Request.post(`/user/handle_friend_request`, { friendship_id, accept })
  },
  // 发送好友请求，verification是验证消息，可以为空
  FirendRequest: (friend_id: string, verification?: string): Promise<ApiResponse<null>> => {
    return Request.post(`/user/friend_request`, { friend_id, verification })
  },
  //  搜索用户，根据ID，用户名，昵称进行模糊搜索
  searchUser: (keyword: string): Promise<ApiResponse<User[]>> => {
    return Request.get(`/user/search`, { params: { keyword } })
  }
}
