export interface BlogInfo {
  id: string;
  userID: string;
  title: string;
  context: string;
  tags: string;
  region: string;
  collection: number;
  like: number;
  createdBy: string;
  isPrivate: boolean;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}

export interface BlogListResponse {
  data: BlogInfo[];
  total: number;
  page: number;
  page_size: number;
}
