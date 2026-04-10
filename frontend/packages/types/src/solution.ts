export interface Solution {
  id: number
  user_id: number
  problem_id: number
  title: string
  excerpt: string
  context: string
  tags: string[]
  like: number
  collection: number
  view: number
  status: string
  created_at: string
  updated_at: string
}

export interface SolutionWithAuthor extends Solution {
  username: string
  avatar: string | null
  problem_title: string
}

export interface CreateSolutionParam {
  problem_id: number
  title: string
  excerpt?: string
  context: string
  tags?: string[]
  status?: string
}
