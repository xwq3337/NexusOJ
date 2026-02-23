import { VerdictType } from './../constants/index'
import Request from '@/services/api'
import { ApiResponse } from '@/types/api'
import { Problem } from '@/types/problem'

export const problemApi = {
  getProblemList: (): Promise<ApiResponse<Problem[]>> => {
    return Request.get('/problem/list')
  },
  getProblemDetail: (id: string): Promise<ApiResponse<Problem>> => {
    return Request.get(`/problem/${id}`)
  },
  Count: (): Promise<ApiResponse<number>> => {
    return Request.get(`/problem/count`)
  },
  TestCode: (options: TestCodeOptions): Promise<ApiResponse<TestCodeResponse>> => {
    return Request.post('/ide/submit', {
      resources_limits: {
        cpu_time: 100000,
        memory_bytes: 67108864,
        stack_bytes: 10485760,
        output_bytes: 10485760
      },
      timeout: 10000,
      ...options
    })
  },
  SubmitCode: (options: SubmitCodeOptions): Promise<ApiResponse<string>> => {
    return Request.post('/problem/submit', { ...options })
  }
}
interface SubmitCodeOptions {
  problem_id: string
  user_id: string
  code: string
  language: string
}

interface TestCodeOptions {
  submission_id: number
  code: string
  language: string
  test_cases: {
    case_id: number
    stdin: string
    expected: string
  }[]
  message: string
  seccomp_profile: string // seccomp profile
  resources_limits?: {
    cpu_time: number // cpu quota
    memory_bytes: number // memory limit 64MB = 67108864bytes
    stack_bytes: number // stack limit  10MB
    output_bytes: number // output limit  10MB = 10485760bytes
  }
  timeout?: number // 10 seconds = 10000ms
}

interface TestCodeResponse {
  submission_id: number
  language: string
  code: string
  verdict: VerdictType
  max_time: number
  max_memory: number
  result: {
    case_id: string
    stdin: string
    stdout: string
    stderr: string
    status: VerdictType
    time: number
    memory: number
    expected: string
  }[]
}
