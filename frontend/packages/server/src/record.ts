import Request from "./api";
import type { ApiResponse } from "@nexusoj/type";
import type { Record, GetRecordListParams, GetRecordListResponse, GetRecordDetailResponse } from "@nexusoj/type";

export const recordApi = {
    getRecordList: (params : Partial<GetRecordListParams>): Promise<ApiResponse<{
        data : GetRecordListResponse[],
        total : number
    }>> => {
        return Request.get('/record/list', { params })
    },
    getRecordDetail: (id : string): Promise<ApiResponse<GetRecordDetailResponse>> => {
        return Request.get(`/record/${id}`)
    },
    getDailyActivity: (): Promise<ApiResponse<{ date: string, count: number }[]>> => {
        return Request.get('/daily-activity')
    }
}