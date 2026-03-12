import type { Log } from '@nexusoj/type';
import type { ApiResponse } from "@nexusoj/type";
import Request from "./api";

export const logApi = {
  getDate: (): Promise<ApiResponse<{
    dates : string[];
  }>> => {
    return Request.get("/log/date");
  },
  /**
   * 
   * @param date 2006-01-02
   * @param page 1 
   * @param pageSize 10
   * @returns 
   */
  getLogList: (date: string): Promise<string> => {
    return Request.get("/log/list", { params: { date } });
  },
};

