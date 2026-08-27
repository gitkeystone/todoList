import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'

import TodoItem from '@/features/todo/TodoItem.vue'
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

beforeEach(() => {
  setActivePinia(createPinia())
})

describe('TodoItem', () => {
  it('渲染标题与时间', () => {
    const wrapper = mount(TodoItem, {
      props: { todo: makeTodo({ title: '学习 Go 语言' }) },
    })
    expect(wrapper.text()).toContain('学习 Go 语言')
    expect(wrapper.find('time').exists()).toBe(true)
  })

  it('完成态：添加 done 类、勾选状态与划线样式', () => {
    const wrapper = mount(TodoItem, {
      props: { todo: makeTodo({ completed: true, completedAt: '2025-01-01T12:00:00Z' }) },
    })
    expect(wrapper.classes()).toContain('done')
    expect(wrapper.get('.checkbox').attributes('aria-pressed')).toBe('true')
    expect(wrapper.get('.checkbox').attributes('aria-label')).toBe('标记为未完成')
  })

  it('未完成态：aria 语义正确', () => {
    const wrapper = mount(TodoItem, {
      props: { todo: makeTodo({ completed: false }) },
    })
    expect(wrapper.classes()).not.toContain('done')
    expect(wrapper.get('.checkbox').attributes('aria-pressed')).toBe('false')
    expect(wrapper.get('.checkbox').attributes('aria-label')).toBe('标记为已完成')
  })
})
