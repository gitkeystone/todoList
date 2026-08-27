import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import Sonner from '@/components/ui/sonner/Sonner.vue'

// 回归测试：曾误装 React 版 sonner 导致白屏（Vue 渲染 React hooks 崩溃）。
// 现在使用 vue-sonner（Vue 移植版），必须可正常渲染。
describe('Sonner（vue-sonner 封装）', () => {
  it('可正常渲染且不抛错', () => {
    expect(() => mount(Sonner, { props: { position: 'top-center' } })).not.toThrow()
  })

  it('渲染出 toaster 根节点', () => {
    const wrapper = mount(Sonner, { props: { position: 'top-center' } })
    expect(wrapper.find('.toaster').exists()).toBe(true)
  })
})
