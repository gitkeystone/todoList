import { vi } from 'vitest'

// jsdom 未实现 window.matchMedia，注入可控 mock（useTheme / TodoInput 自动聚焦依赖）
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query: string): MediaQueryList => {
    const listeners: ((e: MediaQueryListEvent) => void)[] = []
    const mql = {
      matches: false,
      media: query,
      onchange: null,
      addEventListener: (_: string, cb: (e: MediaQueryListEvent) => void) => listeners.push(cb),
      removeEventListener: () => undefined,
      addListener: (cb: (e: MediaQueryListEvent) => void) => listeners.push(cb),
      removeListener: () => undefined,
      dispatchEvent: () => true,
    }
    Object.defineProperty(mql, 'listeners', { value: listeners })
    return mql as unknown as MediaQueryList
  }),
})
