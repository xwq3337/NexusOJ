export interface Problem {
  id: string
  title: string
  difficulty: 'Easy' | 'Medium' | 'Hard'
  acceptance: string
  tags: string[]
  status?: 'solved' | 'attempted' | 'todo'
}
