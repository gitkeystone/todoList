import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import Sonner from '@/components/ui/sonner/Sonner.vue'

// 回归测试：曾误装 React 版 sonner 导致白屏；后又发现未导入 vue-sonner 样式
// 导致 toast 无卡片外观。现在必须可渲染出 [data-sonner-toaster] 根节点。
describe('Sonner（vue-sonner 封装）', () => {
  it('可正常渲染且不抛错', () => {
    expect(() => mount(Sonner)).not.toThrow()
  })

  it('渲染出 toaster 根节点', () => {
    const wrapper = mount(Sonner)
    expect(wrapper.find('[data-sonner-toaster]').exists()).toBe(true)
  })
})
