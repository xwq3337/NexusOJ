import Request from "@/services/api";
import { ApiResponse } from "@/types/api";
import { Problem } from "@/types/problem";

export const problemApi = {
    getProblemList: (): Promise<ApiResponse<Problem[]>> =>  {
        return Request.get("/problem/list");
    },
    getProblemDetail: (id: string): Promise<ApiResponse<Problem>> => {
        return Request.get(`/problem/${id}`);
    },
    Count:() : Promise<ApiResponse<number>> => {
    return Request.get(`/problem/count`)
  }
}