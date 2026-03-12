import type { ApiResponse, Blog } from '@nexusoj/type'
import Request from './api'
import type { createBlogParam, UpdateBlogParam } from '@nexusoj/type'
export const blogApi = {
  // 获取博客总数(全部，包括审核，未审核，草稿等状态的博客数量)
  Count:() : Promise<ApiResponse<number>> => {
    return Request.get(`/blog/count`)
  },
  // ============ 博客相关接口 ============
  /**
   * 获取博客列表(GET)
   * @param params 查询参数
   * @returns 分页的博客列表
   */
  getAvailableList: (
    keywords?: string,
    page?: number,
    page_size?: number
  ): Promise<ApiResponse<{data : (Blog & { user_id: string; username: string })[], total : number}>> => {
    return Request.get('/blog/available-list', { params: { keywords, page, page_size } })
  },
  // TODO: page, pagesize后端需要更改
  getAllList: (
      page?: number,
      page_size?: number,
  ): Promise<ApiResponse<Blog[]>> => {
    return Request.get('/blog/full-list', { params: { page, page_size } })
  },
  getRecycleList:() : Promise<ApiResponse<Blog[]>> => {
    return Request.get('/blog/recycle-list')
  },
  getVerifyBlogList: (): Promise<ApiResponse<Blog[]>> => {
    return Request.get('/blog/verify-list')
  },
  /**
   * 根据博客ID获取博客信息(GET)
   * @param id 博客ID
   * @returns 博客信息，包含博客作者的用户名,头像等基本信息
   */
  getBlogById: (id: string): Promise<ApiResponse<Blog & { username: string; avatar: string }>> => {
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
  },
  deleteBlog: (id: string): Promise<ApiResponse<void>> => {
    return Request.post(`/blog/delete`, { id })
  },

}
