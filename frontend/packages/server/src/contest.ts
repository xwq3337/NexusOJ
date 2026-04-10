import Request from "./api"
import type { ApiResponse } from '@nexusoj/type'
import type {
  Contest, ContestRankItem, ContestReport, ContestProblemItem, ContestRecordItem,
  RegisterContestParams, ContestSubmitParams
} from '@nexusoj/type'

export const contestApi = {
  // 用户端
  getContestList: (page?: number, pageSize?: number, search?: string): Promise<ApiResponse<{ list: (Contest & { is_registered: boolean })[], total: number }>> => {
    return Request.get('/contest/list', { params: { page, page_size: pageSize, search } })
  },
  getContestDetail: (id: string): Promise<ApiResponse<{ contest: Contest, participant_count: number, is_registered: boolean }>> => {
    return Request.get(`/contest/${id}`)
  },
  registerContest: (params: RegisterContestParams): Promise<ApiResponse<void>> => {
    return Request.post('/contest/register', params)
  },
  getContestProblems: (id: string): Promise<ApiResponse<{
    contest : Contest,
    problems: ContestProblemItem[],
    is_registered: boolean
  }>> => {
    return Request.get(`/contest/${id}/problems`)
  },
  getContestProblemDetail: (contestId: string, label: string): Promise<ApiResponse<{
    contest: Contest,
    problem: ContestProblemItem
  }>> => {
    return Request.get(`/contest/${contestId}/problems/${label}`)
  },
  submitContestProblem: (contestId: string, params: ContestSubmitParams): Promise<ApiResponse<any>> => {
    return Request.post(`/contest/${contestId}/submit`, params)
  },
  getContestSubmissions: (contestId: string, params?: { page?: number; page_size?: number; verdict?: string; language?: string; problem_label?: string }): Promise<ApiResponse<{ list: ContestRecordItem[], total: number }>> => {
    return Request.get(`/contest/${contestId}/submissions`, { params })
  },
  getContestSubmissionDetail: (contestId: string, recordId: string): Promise<ApiResponse<ContestRecordItem>> => {
    return Request.get(`/contest/${contestId}/submissions/${recordId}`)
  },
  getContestRanking: (id: string): Promise<ApiResponse<ContestRankItem[]>> => {
    return Request.get(`/contest/${id}/ranking`)
  },
  getMyContestStatus: (id: string): Promise<ApiResponse<any>> => {
    return Request.get(`/contest/${id}/my-status`)
  },

  // 管理员端
  getAdminContestList: (page?: number, pageSize?: number, search?: string): Promise<ApiResponse<{ list: Contest[], total: number }>> => {
    return Request.get('/admin/contest/list', { params: { page, page_size: pageSize, search } })
  },
  getAdminContestDetail: (id: string): Promise<ApiResponse<any>> => {
    return Request.get(`/admin/contest/${id}`)
  },
  createContest: (data: Partial<Contest> & { begin_at: string; end_at: string; contest_type: string }): Promise<ApiResponse<Contest>> => {
    return Request.post('/admin/contest/create', data)
  },
  updateContest: (data: Partial<Contest> & { id: string }): Promise<ApiResponse<void>> => {
    return Request.post('/admin/contest/update', data)
  },
  deleteContest: (id: string): Promise<ApiResponse<void>> => {
    return Request.post('/admin/contest/delete', { id })
  },
  setContestProblems: (contestId: string, problems: Partial<ContestProblemItem>[]): Promise<ApiResponse<void>> => {
    return Request.post(`/admin/contest/${contestId}/problems`, { problems })
  },
  manageParticipant: (contestId: string, data: { user_id: string; action: string }): Promise<ApiResponse<void>> => {
    return Request.post(`/admin/contest/${contestId}/participants`, data)
  },
  generateReport: (id: string): Promise<ApiResponse<ContestReport>> => {
    return Request.get(`/admin/contest/${id}/report`)
  },
  getImportPreview: (contestId: string): Promise<ApiResponse<any[]>> => {
    return Request.get(`/admin/contest/${contestId}/import-preview`)
  },
  importProblems: (contestId: string, problemIds: number[]): Promise<ApiResponse<any>> => {
    return Request.post(`/admin/contest/${contestId}/import-problems`, { problem_ids: problemIds })
  },
}
