// Todo 领域类型（与 PRD §6.2 契约一致）

/** 待办事项 */
export interface Todo {
  id: number
  title: string
  completed: boolean
  /** RFC3339 UTC */
  createdAt: string
  updatedAt: string
  /** 完成时间；未完成时为 null */
  completedAt: string | null
}

/** 列表筛选状态 */
export type FilterStatus = 'all' | 'active' | 'completed'

/** 列表查询参数（PRD §6.2） */
export interface TodoQuery {
  status?: FilterStatus
  q?: string
  page?: number
  pageSize?: number
  sort?: string
}

/** 分页元信息（PRD §6.2） */
export interface PageMeta {
  page: number
  pageSize: number
  total: number
  totalPages: number
}

/** 统一响应信封（PRD §6.2） */
export interface ApiEnvelope<T> {
  code: number
  message: string
  data: T
  meta: PageMeta | null
}
