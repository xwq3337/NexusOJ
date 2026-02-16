import { Blog } from '@/types/blog'
import { ApiResponse } from '@/types/api'
import Request from './api'
import { createBlogParam, UpdateBlogParam } from '@/types/blog'
export const blogApi = {
  // ============ 博客相关接口 ============
  /**
   * 获取博客列表(GET)
   * @param params 查询参数
   * @returns 分页的博客列表
   */
  getList: (): Promise<ApiResponse<(Blog & { user_id: string; username: string })[]>> => {
    return Request.get('/blog/list')
  },
  /**
   * 根据博客ID获取博客信息(GET)
   * @param id 博客ID
   * @returns 博客信息，包含博客作者的用户名
   */
  getBlogById: (id: string): Promise<ApiResponse<Blog & { user_id: string }>> => {
    return Request.get(`/blog/${id}`)
  },

  /**
   * 创建博客(POST)
   * @param blog 博客信息
   * @returns 创建成功的博客
   */
  createBlog: (param: createBlogParam): Promise<ApiResponse<Blog>> => {
    return Request.post('/blog/create', param)
  },
  /**
   * 更新博客(POST)
   * @param param
   * @returns
   */
  updateBlog: (param: UpdateBlogParam): Promise<ApiResponse<Blog>> => {
    return Request.post('/blog/update', param)
  }
}
