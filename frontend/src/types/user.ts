export interface User {
  id: string
  username: string
  password?: string
  email: string
  nickname: string
  introduction: string
  rating: number
  school: string
  avatar: string
  user_role: string
  gender: string
  submission: number
  accept: number
  codeforces: string
  birthday: string
  status: number // // 0 正常 1 封禁
  created_at?: string
  updated_at: string
  banned_to: string
  balance: number // "DouDou" / Beans
}
export interface FriendShip {
  id: string
  user_id: string
  friend_id: string
  status: number // // 0 待处理 1 已接受 2 已拒绝
  created_at: string
  updated_at: string
  deleted_at?: string
}