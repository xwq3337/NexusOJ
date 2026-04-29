import Request from './api'

export interface JudgePerLanguage {
  avg_max_time_ms: number
  avg_total_test_case_time_ms: number
  count: number
  total_test_cases: number
}

export interface JudgeStatsResponse {
  global: {
    avg_max_time_ms: number
    avg_total_test_case_time_ms: number
  }
  per_language: Record<string, JudgePerLanguage>
  total_submissions: number
  total_test_cases: number
  verdict_counts: Record<string, number>
}

export const judgeApi = {
  GetStats: (): Promise<JudgeStatsResponse> => {
    return Request.get('/judge/stats')
  },
}
