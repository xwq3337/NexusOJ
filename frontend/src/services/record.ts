import Request from "@/services/api";
import { ApiResponse } from "@/types/api";
import { Record, GetRecordListParams, GetRecordListResponse, GetRecordDetailResponse } from "@/types/record";

export const recordApi = {
    getRecordList: (params : Partial<GetRecordListParams>): Promise<ApiResponse<{
        data : GetRecordListResponse[],
        total : number
    }>> => {
        return Request.get('/record/list', { params })
    },
    getRecordDetail: (id : string): Promise<ApiResponse<GetRecordDetailResponse>> => {
        return Request.get(`/record/${id}`)
    }
}