export interface User {
  id: Number
  username: string
  password: string
  email: string
  nickname: string
  introduction: string
  rating: number
  school: string
  avatar: string
  user_role: 'user' | 'admin' | 'super_admin'
  gender: string
  submission: number
  accept: number
  codeforces: string
  birthday: string
  status: number // // 0 正常 1 封禁
  created_at: string
  updated_at: string
  banned_to: string
  balance: number // "DouDou" / Beans
}
export interface FriendShip {
  id: string
  user_id: Number
  friend_id: Number
  friend_username: string
  friend_nickname: string
  unread_count : number
  latest_message: string
  friend_avatar: string
  remark: string
  created_at: string
}
export interface FriendShipRequest {
  id: string
  user_id: Number
  friend_id: Number
  friend_username: string
  friend_nickname: string
  friend_avatar: string
  message: string
  status: string // 0 待处理 1 已接受 2 已拒绝
  created_at: string
}
