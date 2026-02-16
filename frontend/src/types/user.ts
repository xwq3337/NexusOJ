export interface User {
  id: string
  username: string
  password: string
  email: string
  nickname: string
  introduction: string
  rating: string
  school: string
  avatar: string
  user_role: string
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
