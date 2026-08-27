import MockAdapter from 'axios-mock-adapter'
import { createPinia, setActivePinia } from 'pinia'
import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it } from 'vitest'

import { http } from '@/api/todo'
import TodoItem from '@/features/todo/TodoItem.vue'
import { useTodoStore } from '@/stores/todo'
import type { Todo } from '@/types/todo'

function makeTodo(overrides: Partial<Todo> = {}): Todo {
  return {
    id: 1,
    title: '学习 Go 语言',
    completed: false,
    createdAt: '2025-01-01T10:00:00Z',
    updatedAt: '2025-01-01T10:00:00Z',
    completedAt: null,
    ...overrides,
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

/** 预置 store：模拟列表接口返回一条待办 */
async function mountWithStore(todo: Todo) {
  mock.onGet('/api/v1/todos').reply(200, { code: 0, message: 'ok', data: [todo], meta: null })
  const store = useTodoStore()
  await store.fetchList()
  const wrapper = mount(TodoItem, { props: { todo: store.items[0] } })
  return { wrapper, store }
}

describe('TodoItem 渲染', () => {
  it('渲染标题与时间', async () => {
    const { wrapper } = await mountWithStore(makeTodo({ title: '学习 Go 语言' }))
    expect(wrapper.text()).toContain('学习 Go 语言')
    expect(wrapper.find('time').exists()).toBe(true)
  })

  it('完成态：done 类、勾选状态与划线样式', async () => {
    const { wrapper } = await mountWithStore(
      makeTodo({ completed: true, completedAt: '2025-01-01T12:00:00Z' }),
    )
    expect(wrapper.classes()).toContain('done')
    expect(wrapper.get('.checkbox').attributes('aria-pressed')).toBe('true')
    expect(wrapper.get('.checkbox').attributes('aria-label')).toBe('标记为未完成')
  })

  it('未完成态：aria 语义正确', async () => {
    const { wrapper } = await mountWithStore(makeTodo({ completed: false }))
    expect(wrapper.classes()).not.toContain('done')
    expect(wrapper.get('.checkbox').attributes('aria-pressed')).toBe('false')
    expect(wrapper.get('.checkbox').attributes('aria-label')).toBe('标记为已完成')
  })
})

describe('TodoItem 交互', () => {
  it('点击勾选框完成待办（乐观更新）', async () => {
    const { wrapper, store } = await mountWithStore(makeTodo({ id: 7 }))
    mock.onPatch('/api/v1/todos/7').reply(200, {
      code: 0,
      message: 'ok',
      data: { ...store.items[0], completed: true, completedAt: new Date().toISOString() },
      meta: null,
    })

    await wrapper.get('.checkbox').trigger('click')
    await wrapper.vm.$nextTick()

    expect(store.items[0].completed).toBe(true)
    expect(store.items[0].completedAt).not.toBeNull()
    expect(mock.history.patch).toHaveLength(1)
  })

  it('点击删除按钮删除待办', async () => {
    const { wrapper, store } = await mountWithStore(makeTodo({ id: 7 }))
    mock.onDelete('/api/v1/todos/7').reply(204)

    await wrapper.get('.delete').trigger('click')
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    expect(store.items).toHaveLength(0)
    expect(mock.history.delete).toHaveLength(1)
  })

  it('条目聚焦时按 Delete 键删除', async () => {
    const { wrapper, store } = await mountWithStore(makeTodo({ id: 7 }))
    mock.onDelete('/api/v1/todos/7').reply(204)

    await wrapper.trigger('keydown', { key: 'Delete' })
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    expect(store.items).toHaveLength(0)
    expect(mock.history.delete).toHaveLength(1)
  })
})
