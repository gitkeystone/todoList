import { defineStore } from 'pinia'
import { toast } from 'sonner'

import { ApiError, todoApi } from '@/api/todo'
import type { FilterStatus, Todo } from '@/types/todo'

/** 单用户本地应用：一次拉取全部（上限 100 条），视图/统计在客户端即时过滤（PRD §4.4 I-10 数据刷新） */
const PAGE_SIZE = 100

interface TodoState {
  items: Todo[]
  status: FilterStatus
  keyword: string
  loading: boolean
  clearing: boolean
  loaded: boolean
}

function notifyError(e: unknown) {
  const message = e instanceof ApiError ? e.message : '操作失败，请稍后重试'
  toast.error(message)
}

export const useTodoStore = defineStore('todo', {
  state: (): TodoState => ({
    items: [],
    status: 'all',
    keyword: '',
    loading: false,
    clearing: false,
    loaded: false,
  }),

  getters: {
    /** 当前视图（状态筛选 + 关键词）的待办，按创建时间倒序 */
    visibleItems(state): Todo[] {
      const kw = state.keyword.trim().toLowerCase()
      return state.items
        .filter((t) =>
          state.status === 'all'
            ? true
            : state.status === 'completed'
              ? t.completed
              : !t.completed,
        )
        .filter((t) => (kw ? t.title.toLowerCase().includes(kw) : true))
        .sort((a, b) => Date.parse(b.createdAt) - Date.parse(a.createdAt))
    },
    activeCount: (s) => s.items.filter((t) => !t.completed).length,
    completedCount: (s) => s.items.filter((t) => t.completed).length,
    allCount: (s) => s.items.length,
  },

  actions: {
    async fetchList() {
      this.loading = true
      try {
        const res = await todoApi.list({ status: 'all', page: 1, pageSize: PAGE_SIZE })
        this.items = res.data
        this.loaded = true
      } catch (e) {
        notifyError(e)
      } finally {
        this.loading = false
      }
    },

    async create(title: string) {
      const res = await todoApi.create(title)
      this.items.unshift(res.data)
      return res.data
    },

    /** 完成/取消完成：乐观更新，失败回滚（PRD §8.2） */
    async toggle(todo: Todo) {
      const next = !todo.completed
      const prevCompleted = todo.completed
      const prevCompletedAt = todo.completedAt
      todo.completed = next
      todo.completedAt = next ? new Date().toISOString() : null
      try {
        const res = await todoApi.patch(todo.id, { completed: next })
        Object.assign(todo, res.data)
      } catch (e) {
        todo.completed = prevCompleted
        todo.completedAt = prevCompletedAt
        notifyError(e)
        throw e
      }
    },

    /** 编辑标题：先本地生效，失败回滚 */
    async updateTitle(todo: Todo, title: string) {
      const prev = todo.title
      todo.title = title
      try {
        const res = await todoApi.patch(todo.id, { title })
        Object.assign(todo, res.data)
      } catch (e) {
        todo.title = prev
        notifyError(e)
        throw e
      }
    },

    /** 删除：乐观移除，失败恢复原位 */
    async remove(id: number) {
      const idx = this.items.findIndex((t) => t.id === id)
      const [removed] = this.items.splice(idx, 1)
      try {
        await todoApi.remove(id)
      } catch (e) {
        if (removed) this.items.splice(Math.min(idx, this.items.length), 0, removed)
        notifyError(e)
        throw e
      }
    },

    async clearCompleted() {
      this.clearing = true
      try {
        await todoApi.clearCompleted()
        this.items = this.items.filter((t) => !t.completed)
      } catch (e) {
        notifyError(e)
        throw e
      } finally {
        this.clearing = false
      }
    },
  },
})
