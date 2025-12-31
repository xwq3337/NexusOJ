import type { AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios'

export interface CustomAxiosRequestConfig extends AxiosRequestConfig {
  requestId?: string
}

export interface RequestInterceptors {
  requestInterceptors?: (config: CustomAxiosRequestConfig) => CustomAxiosRequestConfig
  requestInterceptorsCatch?: (error: Error) => Error
  responseInterceptors?: (response: AxiosResponse) => AxiosResponse
  responseInterceptorsCatch?: (error: AxiosError) => Error | Promise<Error>
}

export interface Respones<T = any> {
    code: number
    result: T
}

export interface UserInfo {
  id: string;
  account: string;
  password: string;
  email: string;
  firstName: string;
  lastName: string;
  introduction: string | null;
  rating: number;
  school: string;
  avatar: string;
  education: string;
  userRole: number;
  gender: string;
  submission: number;
  accept: number;
  codeforces: string | null;
  birthday: string;
  bannedEndTime: string | null;
  isDelete: boolean;
  isBanned: boolean;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  balance: number;
}

export interface UserListResponse {
  total: number;
  list: UserInfo[];
}

export interface UserListParams {
  page: number;
  pageSize: number;
  keyword?: string;
  role?: number;
  status?: 'all' | 'banned' | 'normal';
}
