import axios from 'axios'

import type { ApiEnvelope, PageMeta, Todo, TodoQuery } from '@/types/todo'

/** 业务错误：携带后端信封中的错误码与可直接展示的中文信息（PRD §6.2） */
export class ApiError extends Error {
  code: number
  status?: number

  constructor(message: string, code: number, status?: number) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.status = status
  }
}

const http = axios.create({
  baseURL: '/api/v1',
  timeout: 10_000,
})

// 响应拦截器：解包信封；非 2xx 统一转为 ApiError（PRD §8.3）
http.interceptors.response.use(
  (res) => res.data,
  (error) => {
    const env = error.response?.data as ApiEnvelope<unknown> | undefined
    const message = env?.message ?? '网络请求失败，请稍后重试'
    const code = env?.code ?? -1
    return Promise.reject(new ApiError(message, code, error.response?.status))
  },
)

/** 导出实例供测试（axios-mock-adapter 挂载） */
export { http }

export const todoApi = {
  list: (params: TodoQuery = {}) => http.get('/todos', { params }) as Promise<ApiEnvelope<Todo[]>>,

  create: (title: string) => http.post('/todos', { title }) as Promise<ApiEnvelope<Todo>>,

  get: (id: number) => http.get(`/todos/${id}`) as Promise<ApiEnvelope<Todo>>,

  patch: (id: number, payload: { title?: string; completed?: boolean }) =>
    http.patch(`/todos/${id}`, payload) as Promise<ApiEnvelope<Todo>>,

  remove: (id: number) => http.delete(`/todos/${id}`) as Promise<ApiEnvelope<null>>,

  clearCompleted: () => http.delete('/todos/completed') as Promise<ApiEnvelope<null>>,
}

export type { PageMeta }
