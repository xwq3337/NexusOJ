export interface Problem {
  id: string
  title: string
  context: string
  inputDescription: string
  outputDescription: string
  tip: string
  difficulty: number
  judgeCase: string
  judgeConfig: string
  tags: string
  submission: number
  accept: number
  collection: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
  createdBy: string
}

export interface JudgeConfig {
  空间限制: number
  时间限制: number
}

export interface JudgeCase {
  input: string
  output: string
}
