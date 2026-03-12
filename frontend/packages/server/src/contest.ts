import Request from "./api"
import type { ApiResponse } from '@nexusoj/type'
import type { Contest } from '@nexusoj/type'


export const contestApi = {
    getContestList: (): Promise<ApiResponse<Contest[]>> => {
    return Request.get('/contest/list')
  },
}