import Request from './api'
import type { ApiResponse } from '@nexusoj/type'
import type { Problem, ProblemListDTO } from '@nexusoj/type'

export const problemApi = {
  getProblemList: (): Promise<ApiResponse<ProblemListDTO[]>> => {
    return Request.get('/problem/list')
  },
  getProblemDetail: (id: string): Promise<ApiResponse<{
    problem : Problem,
    my_status : "unattempted" | "attempted" | "accepted",
  }>> => {
    return Request.get(`/problem/${id}`)
  },
  Count: (): Promise<ApiResponse<number>> => {
    return Request.get(`/problem/count`)
  },
  SubmitCode: (options: SubmitCodeOptions): Promise<ApiResponse<any>> => {
    return Request.post('/problem/submit', { ...options })
  },
  // TODO:标明类型
  updateProblem: (id: string, options: any): Promise<ApiResponse<void>> => {
    return Request.post('/problem/update', { id, ...options })
  },
  createProblem: (options: any): Promise<ApiResponse<void>> => {
    return Request.post('/problem/create', { ...options })
  }
}
interface SubmitCodeOptions {
  problem_id: string
  user_id: Number
  code: string
  language: string
}


