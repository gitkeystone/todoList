import MockAdapter from 'axios-mock-adapter'
import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it } from 'vitest'

import { http } from '@/api/todo'
import { useTodoStore } from '@/stores/todo'
import type { ApiEnvelope, Todo } from '@/types/todo'

const ok = <T>(data: T, meta: ApiEnvelope<T>['meta'] = null) =>
  ({ code: 0, message: 'ok', data, meta }) as ApiEnvelope<T>

const fail = (code: number, message: string, status: number) => ({
  code,
  message,
  data: null,
  meta: null,
  status,
})

function makeTodo(id: number, title: string, completed = false, minutesAgo = 0): Todo {
  const iso = new Date(Date.now() - minutesAgo * 60_000).toISOString()
  return {
    id,
    title,
    completed,
    createdAt: iso,
    updatedAt: iso,
    completedAt: completed ? iso : null,
  }
}

let mock: MockAdapter

beforeEach(() => {
  setActivePinia(createPinia())
  mock = new MockAdapter(http)
})

afterEach(() => {
  mock.restore()
})

describe('todo store', () => {
  it('fetchList 加载全部待办并计算统计', async () => {
    mock.onGet('/api/v1/todos').reply(200, ok([makeTodo(1, 'A', true), makeTodo(2, 'B')]))
    const store = useTodoStore()
    await store.fetchList()

    expect(store.items).toHaveLength(2)
    expect(store.loaded).toBe(true)
    expect(store.loading).toBe(false)
    expect(store.allCount).toBe(2)
    expect(store.activeCount).toBe(1)
    expect(store.completedCount).toBe(1)
  })

  it('create 新增待办并置顶', async () => {
    mock.onGet('/api/v1/todos').reply(200, ok([]))
    const store = useTodoStore()
    await store.fetchList()

    mock.onPost('/api/v1/todos').reply(201, ok(makeTodo(1, '新任务')))
    await store.create('新任务')

    expect(store.items).toHaveLength(1)
    expect(store.items[0].title).toBe('新任务')
    expect(store.allCount).toBe(1)
  })

  it('toggle 完成/取消完成：乐观更新，失败回滚', async () => {
    const t = makeTodo(1, '任务')
    mock.onGet('/api/v1/todos').reply(200, ok([t]))
    const store = useTodoStore()
    await store.fetchList()
    const item = store.items[0]

    // 成功：完成
    mock
      .onPatch('/api/v1/todos/1')
      .reply(200, ok({ ...t, completed: true, completedAt: new Date().toISOString() }))
    await store.toggle(item)
    expect(item.completed).toBe(true)
    expect(item.completedAt).not.toBeNull()
    expect(store.activeCount).toBe(0)
    expect(store.completedCount).toBe(1)

    // 失败：回滚到完成态
    mock.onPatch('/api/v1/todos/1').reply(500, fail(50000, '服务器内部错误', 500))
    await expect(store.toggle(item)).rejects.toThrow('服务器内部错误')
    expect(item.completed).toBe(true)
    expect(item.completedAt).not.toBeNull()
  })

  it('remove 乐观移除，失败恢复原位', async () => {
    mock.onGet('/api/v1/todos').reply(200, ok([makeTodo(1, 'A'), makeTodo(2, 'B')]))
    const store = useTodoStore()
    await store.fetchList()

    mock.onDelete('/api/v1/todos/1').reply(500, fail(50000, '服务器内部错误', 500))
    await expect(store.remove(1)).rejects.toThrow()
    expect(store.items).toHaveLength(2)
    expect(store.items[0].id).toBe(1)

    mock.onDelete('/api/v1/todos/2').reply(204)
    await store.remove(2)
    expect(store.items).toHaveLength(1)
  })

  it('clearCompleted 仅清除已完成项', async () => {
    mock.onGet('/api/v1/todos').reply(200, ok([makeTodo(1, 'A', true), makeTodo(2, 'B')]))
    const store = useTodoStore()
    await store.fetchList()

    mock.onDelete('/api/v1/todos/completed').reply(204)
    await store.clearCompleted()

    expect(store.items).toHaveLength(1)
    expect(store.items[0].id).toBe(2)
    expect(store.completedCount).toBe(0)
  })

  it('fetchList 失败时记录 error 且不置 loaded', async () => {
    mock.onGet('/api/v1/todos').reply(500, fail(50000, '服务器内部错误', 500))
    const store = useTodoStore()
    await store.fetchList()

    expect(store.loaded).toBe(false)
    expect(store.error).toBe('服务器内部错误')
    expect(store.items).toHaveLength(0)

    // 重试成功后清除 error
    mock.onGet('/api/v1/todos').reply(200, ok([makeTodo(1, 'A')]))
    await store.fetchList()
    expect(store.error).toBeNull()
    expect(store.loaded).toBe(true)
  })

  it('visibleItems 按状态与关键词过滤，按创建时间倒序', async () => {
    mock
      .onGet('/api/v1/todos')
      .reply(
        200,
        ok([
          makeTodo(1, '学习 Go', true, 20),
          makeTodo(2, '写周报', false, 5),
          makeTodo(3, 'Go 复习', false, 10),
        ]),
      )
    const store = useTodoStore()
    await store.fetchList()

    expect(store.visibleItems.map((t) => t.title)).toEqual(['写周报', 'Go 复习', '学习 Go'])

    store.status = 'active'
    expect(store.visibleItems.map((t) => t.title)).toEqual(['写周报', 'Go 复习'])

    store.status = 'completed'
    expect(store.visibleItems.map((t) => t.title)).toEqual(['学习 Go'])

    store.status = 'all'
    store.keyword = 'go'
    expect(store.visibleItems.map((t) => t.title)).toEqual(['Go 复习', '学习 Go'])
  })
})
