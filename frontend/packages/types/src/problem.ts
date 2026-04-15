export interface Problem {
  id: number
  user_id : number
  title: string
  context: string
  input_description: string
  output_description: string
  judge_case : JudgeCase[]
  judge_config : JudgeConfig
  judge_sample : JudgeCase[]
  tips?: string
  difficulty: number
  tags: string[]
  accept: number
  submission: number
  collection?: number
  created_at?: string
  updated_at?: string
  deleted_at?: string
}

export interface JudgeCase {
  input: string
  expected: string
}
export interface JudgeConfig {
  time_limit: number
  memory_limit: number
}


export interface ProblemListDTO extends Problem {
  status: "unattempted" | "attempted" | "solved"
}

export interface RecommendedProblem {
  problem_id: number
  title: string
  difficulty: number
  tags: string[]
  score: number
  reason: string
}