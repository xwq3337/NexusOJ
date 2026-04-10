import type { ApiResponse, SolutionWithAuthor, CreateSolutionParam } from '@nexusoj/type'
import Request from './api'

export const solutionApi = {
  getSolutions: (params: {
    problem_id?: number
    tag?: string
    keyword?: string
    status?: string
    user_id?: number
    page?: number
    page_size?: number
  }): Promise<ApiResponse<{ solutions: SolutionWithAuthor[]; total: number }>> => {
    return Request.get('/solution/list', { params })
  },

  getSolutionDetail: (id: number): Promise<ApiResponse<SolutionWithAuthor>> => {
    return Request.get(`/solution/${id}`)
  },

  createSolution: (data: CreateSolutionParam): Promise<ApiResponse<SolutionWithAuthor>> => {
    return Request.post('/solution/create', data)
  },
}
