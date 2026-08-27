import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { useTheme } from '@/composables/useTheme'

const STORAGE_KEY = 'todolist.theme'

type ThemeApi = ReturnType<typeof useTheme>

/** 可控 matchMedia：可切换系统偏好并触发 change 监听 */
function mockMatchMedia(initialDark = false) {
  const listeners: ((e: MediaQueryListEvent) => void)[] = []
  let dark = initialDark
  const mql = {
    matches: dark,
    media: '(prefers-color-scheme: dark)',
    onchange: null,
    addEventListener: (_: string, cb: (e: MediaQueryListEvent) => void) => listeners.push(cb),
    removeEventListener: () => undefined,
    addListener: (cb: (e: MediaQueryListEvent) => void) => listeners.push(cb),
    removeListener: () => undefined,
    dispatchEvent: () => true,
  }
  vi.mocked(window.matchMedia).mockReturnValue(mql as unknown as MediaQueryList)
  return {
    setSystemDark(next: boolean) {
      dark = next
      mql.matches = next
      listeners.forEach((cb) => cb({ matches: next } as MediaQueryListEvent))
    },
  }
}

function mountTheme() {
  const wrapper = mount({
    setup() {
      return useTheme()
    },
    template: '<div />',
  })
  return wrapper.vm as unknown as ThemeApi
}

beforeEach(() => {
  localStorage.clear()
  document.documentElement.classList.remove('dark')
})

afterEach(() => {
  vi.restoreAllMocks()
})

describe('useTheme 三态主题', () => {
  it('默认跟随系统：系统深色时应用 dark class', () => {
    mockMatchMedia(true)
    const vm = mountTheme()
    expect(vm.mode).toBe('system')
    expect(document.documentElement.classList.contains('dark')).toBe(true)
  })

  it('setMode(dark) 应用深色并持久化', () => {
    mockMatchMedia(false)
    const vm = mountTheme()
    vm.setMode('dark')
    expect(vm.isDark).toBe(true)
    expect(document.documentElement.classList.contains('dark')).toBe(true)
    expect(localStorage.getItem(STORAGE_KEY)).toBe('dark')
  })

  it('setMode(light) 移除 dark class', () => {
    mockMatchMedia(true)
    const vm = mountTheme()
    vm.setMode('light')
    expect(vm.isDark).toBe(false)
    expect(document.documentElement.classList.contains('dark')).toBe(false)
    expect(localStorage.getItem(STORAGE_KEY)).toBe('light')
  })

  it('system 模式下系统主题变化时实时跟随', () => {
    const { setSystemDark } = mockMatchMedia(false)
    const vm = mountTheme()
    expect(vm.isDark).toBe(false)

    setSystemDark(true)
    expect(vm.isDark).toBe(true)
    expect(document.documentElement.classList.contains('dark')).toBe(true)

    // 切到显式 light 后不再跟随系统
    vm.setMode('light')
    setSystemDark(false)
    expect(vm.isDark).toBe(false)
    expect(localStorage.getItem(STORAGE_KEY)).toBe('light')
  })

  it('读取已持久化的主题', () => {
    localStorage.setItem(STORAGE_KEY, 'dark')
    mockMatchMedia(false)
    const vm = mountTheme()
    expect(vm.mode).toBe('dark')
    expect(vm.isDark).toBe(true)
  })
})
