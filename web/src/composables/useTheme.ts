import { computed, onMounted, onUnmounted, ref } from 'vue'

/** 主题模式：显式浅色 / 显式深色 / 跟随系统（PRD FR-09 三态） */
export type ThemeMode = 'light' | 'dark' | 'system'

const STORAGE_KEY = 'todolist.theme'

function readStored(): ThemeMode {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored === 'light' || stored === 'dark' || stored === 'system') return stored
  return 'system'
}

function systemPrefersDark(): boolean {
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

/** 解析生效主题（跟随系统时取系统偏好） */
function resolve(mode: ThemeMode): 'light' | 'dark' {
  if (mode === 'system') return systemPrefersDark() ? 'dark' : 'light'
  return mode
}

function apply(resolved: 'light' | 'dark') {
  document.documentElement.classList.toggle('dark', resolved === 'dark')
}

/**
 * 深浅色主题 composable（PRD FR-09）
 * 初始 class 由 index.html 防 FOUC 脚本预置；本 composable 负责三态切换与系统监听。
 */
export function useTheme() {
  const mode = ref<ThemeMode>(readStored())
  const resolved = ref<'light' | 'dark'>(resolve(mode.value))
  const isDark = computed(() => resolved.value === 'dark')
  const mql = window.matchMedia('(prefers-color-scheme: dark)')

  function setMode(next: ThemeMode) {
    mode.value = next
    localStorage.setItem(STORAGE_KEY, next)
    resolved.value = resolve(next)
    apply(resolved.value)
  }

  function onSystemChange(e: MediaQueryListEvent) {
    // 仅"跟随系统"模式响应系统变化
    if (mode.value === 'system') {
      resolved.value = e.matches ? 'dark' : 'light'
      apply(resolved.value)
    }
  }

  onMounted(() => {
    apply(resolved.value)
    mql.addEventListener('change', onSystemChange)
  })
  onUnmounted(() => mql.removeEventListener('change', onSystemChange))

  return { mode, isDark, setMode }
}
