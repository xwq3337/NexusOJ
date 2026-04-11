export type ContestType = 'ACM' | 'OI'
export type ContestStatus = 'Upcoming' | 'Live' | 'Ended'
export type ParticipantStatus = 'registered' | 'disqualified'

export interface Contest {
  id: string
  title: string
  introduction?: string
  begin_at: string
  end_at: string
  duration: number
  contest_type: ContestType
  status: ContestStatus
  is_private: boolean
  like: number
  collection: number
  participants: number
  submission: number
  accept: number
  seal_rank: boolean
  userID: number
  created_at: string
  updated_at: string
}

export interface ContestDetail {
  contest: Contest
  is_registered: boolean
}

export interface ContestProblemItem {
  id: number
  contest_id: string
  label: string
  score: number
  title: string
  context: string
  input_description?: string
  output_description?: string
  tips?: string
  difficulty: number
  judge_case: { input: string; expected: string }[]
  judge_config: { time_limit: number; memory_limit: number }
  judge_sample: { input: string; expected: string }[]
  tags: string[]
  submission: number
  accept: number
  source_problem_id?: number
  my_status?: 'accepted' | 'wrong' | 'pending' | 'unattempted'
  my_score?: number
}

export interface ContestRecordItem {
  id: number
  contest_id: string
  user_id: number
  problem_label: string
  language: string
  verdict: string
  max_time: number
  max_memory: number
  created_at: string
  username?: string
  problem_title?: string
}

export interface ContestRankItem {
  rank: number
  user_id: number
  username: string
  avatar?: string
  solved: number
  total_penalty: number
  score: number
  problems: Record<string, ContestProblemRankDetail>
}

export interface ContestProblemRankDetail {
  attempts: number
  accepted: boolean
  time?: string
  score?: number
  penalty?: number
}

export interface RegisterContestParams {
  contest_id: string
  password?: string
}

export interface ContestSubmitParams {
  problem_label: string
  code: string
  language: string
}

export interface ContestReport {
  contest_id: string
  title: string
  type: ContestType
  generated_at: string
  total_participants: number
  total_submissions: number
  problem_stats: ContestProblemStat[]
  top_participants: ContestTopParticipant[]
}

export interface ContestProblemStat {
  label: string
  attempts: number
  accepted: number
  accept_rate: number
}

export interface ContestTopParticipant {
  user_id: number
  username: string
  rank: number
  score: number
  problems_solved: number
  total_penalty: number
}
