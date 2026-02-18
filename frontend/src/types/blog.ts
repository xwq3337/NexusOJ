export type BlogStatus = 'Pending' | 'Draft' | 'Normal' | 'Violation' | 'Deleted'
export interface Blog {
  id: string
  user_id: string
  title: string
  context: string
  excerpt: string
  view : number
  tags: string[]
  collection: number
  like: number
  is_private: boolean
  created_at: string
  status: BlogStatus
  updated_at: string
}

export interface createBlogParam extends Pick<
  Blog,
  'title' | 'context' | 'tags' | 'is_private' | 'status'
> {}

export interface UpdateBlogParam extends Partial<
  Pick<Blog, 'title' | 'context' | 'tags' | 'is_private' | 'status'>
> {
  id: string
}
