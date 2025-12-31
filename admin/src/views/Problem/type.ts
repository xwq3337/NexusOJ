export interface Problem {
  title: string
  context: string
  user_id : string
  input_description: string
  output_description: string
  tip: string
  difficulty: string
  judge_case: JudgeCase[]
  judge_config: JudgeConfig
  judge_sample : JudgeSample[]
  tags: string[]
}
interface JudgeConfig {
  time_limit: number
  memory_limit: number
}
interface JudgeCase {
  input: string
  expected: string
}
interface JudgeSample {
  input: string
  expected: string
}
