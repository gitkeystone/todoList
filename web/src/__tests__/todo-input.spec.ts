import MockAdapter from 'axios-mock-adapter'
import { createPinia, setActivePinia } from 'pinia'
import { flushPromises, mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it } from 'vitest'

import { http } from '@/api/todo'
import TodoInput from '@/features/todo/TodoInput.vue'

const ok = (data: unknown) => ({ code: 0, message: 'ok', data, meta: null })

function makeTodo(id: number, title: string) {
  const iso = new Date().toISOString()
  return { id, title, completed: false, createdAt: iso, updatedAt: iso, completedAt: null }
}

let mock: MockAdapter

beforeEach(() => {
  setActivePinia(createPinia())
  mock = new MockAdapter(http)
  mock.onGet('/api/v1/todos').reply(200, ok([]))
})

afterEach(() => {
  mock.restore()
})

async function mountInput() {
  const wrapper = mount(TodoInput)
  await wrapper.vm.$nextTick()
  return wrapper
}

describe('TodoInput 创建交互', () => {
  it('回车创建并清空输入框', async () => {
    mock.onPost('/api/v1/todos').reply(201, ok(makeTodo(1, '学习 Go')))
    const wrapper = await mountInput()
    const input = wrapper.get('input')

    await input.setValue('学习 Go')
    await input.trigger('keydown', { key: 'Enter' })
    await flushPromises()

    expect(mock.history.post).toHaveLength(1)
    expect(JSON.parse(mock.history.post[0].data).title).toBe('学习 Go')
    expect((input.element as HTMLInputElement).value).toBe('')
  })

  it('空输入不创建，触发抖动提示', async () => {
    const wrapper = await mountInput()
    const input = wrapper.get('input')

    await input.setValue('   ')
    await input.trigger('keydown', { key: 'Enter' })
    await wrapper.vm.$nextTick()

    expect(mock.history.post).toHaveLength(0)
    expect(wrapper.find('.todo-input').classes()).toContain('shaking')
  })

  it('中文输入法组词中回车不提交（composition 保护）', async () => {
    mock.onPost('/api/v1/todos').reply(201, ok(makeTodo(1, '任务')))
    const wrapper = await mountInput()
    const input = wrapper.get('input')

    await input.setValue('任务')
    await input.trigger('compositionstart')
    await input.trigger('keydown', { key: 'Enter' })
    expect(mock.history.post).toHaveLength(0)

    await input.trigger('compositionend')
    await input.trigger('keydown', { key: 'Enter' })
    await flushPromises()
    expect(mock.history.post).toHaveLength(1)
  })
})
