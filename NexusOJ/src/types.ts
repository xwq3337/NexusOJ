export interface Problem {
  id: string
  title: string
  difficulty: 'Easy' | 'Medium' | 'Hard'
  acceptance: string
  tags: string[]
  status?: 'solved' | 'attempted' | 'todo'
}

export interface Contest {
  id: string
  title: string
  startTime: string
  duration: string
  registered: number
  type: 'Weekly' | 'Biweekly' | 'Cup'
  status: 'Live' | 'Upcoming' | 'Ended'
}

export interface User {
  username: string
  avatar: string
  rank: number
  solved: number
  coins: number // "DouDou" / Beans
}

export interface ChatMessage {
  id: string
  role: 'user' | 'model'
  text: string
  timestamp: number
}

export interface Blog {
  id: string
  user_id: string
  title: string
  context: string
  tags: string[]
  collection: number
  like: number
  is_private: boolean
  created_at: string
  status: number
  updated_at: string
  DeletedAt: string | null
  username: string
}
