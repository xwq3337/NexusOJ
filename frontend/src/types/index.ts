export interface ApiResponse<T = any> {
  code: number // 200, 500等
  msg: string // "success", "error"或其他有用信息
  info?: T // 返回的数据, 若code是错误码，则info字段不存在
}
