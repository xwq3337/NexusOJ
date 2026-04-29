import Request from './api'
import type { ApiResponse } from '@nexusoj/type'
import type { Problem, ProblemListDTO } from '@nexusoj/type'

export const problemApi = {
  getProblemList: (page?: number, pageSize?: number, search?: string, tags?: string[]): Promise<ApiResponse<{ problems: ProblemListDTO[], total: number }>> => {
    return Request.get('/problem/list', { params: { page, page_size: pageSize, search, tags: tags?.join(',') } })
  },
  getAllTags: (): Promise<ApiResponse<string[]>> => {
    return Request.get('/problem/tags')
  },
  getProblemDetail: (id: string): Promise<ApiResponse<{
    problem : Problem,
    my_status : "unattempted" | "attempted" | "solved",
  }>> => {
    return Request.get(`/problem/${id}`)
  },
  Count: (): Promise<ApiResponse<number>> => {
    return Request.get(`/problem/count`)
  },
  SubmitCode: (options: SubmitCodeOptions): Promise<ApiResponse<any>> => {
    return Request.post('/problem/submit', { ...options })
  },
  // TODO:标明options类型
  updateProblem: (id: number, options: any): Promise<ApiResponse<void>> => {
    return Request.post('/problem/update', { id, ...options })
  },
  createProblem: (options: any): Promise<ApiResponse<void>> => {
    return Request.post('/problem/create', { ...options })
  }
}
interface SubmitCodeOptions {
  problem_id: number
  user_id: number
  code: string
  language: string
}


