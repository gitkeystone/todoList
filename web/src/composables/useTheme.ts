import { computed, onMounted, onUnmounted, ref } from 'vue'

type ThemeMode = 'light' | 'dark'

const STORAGE_KEY = 'todolist.theme'

/** 读取持久化主题；未设置时跟随系统 */
function resolveInitial(): ThemeMode {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored === 'light' || stored === 'dark') return stored
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

function apply(mode: ThemeMode) {
  document.documentElement.classList.toggle('dark', mode === 'dark')
}

/**
 * 深浅色主题 composable（PRD FR-09）
 * 初始值来自 index.html 的防 FOUC 脚本；未手动设置时跟随系统变化。
 */
export function useTheme() {
  const mode = ref<ThemeMode>(resolveInitial())
  const isDark = computed(() => mode.value === 'dark')
  const mql = window.matchMedia('(prefers-color-scheme: dark)')

  function toggle() {
    mode.value = mode.value === 'dark' ? 'light' : 'dark'
    localStorage.setItem(STORAGE_KEY, mode.value)
    apply(mode.value)
  }

  function onSystemChange(e: MediaQueryListEvent) {
    // 用户未手动选择过主题时才跟随系统
    if (!localStorage.getItem(STORAGE_KEY)) apply(e.matches ? 'dark' : 'light')
  }

  onMounted(() => {
    apply(mode.value)
    mql.addEventListener('change', onSystemChange)
  })
  onUnmounted(() => mql.removeEventListener('change', onSystemChange))

  return { isDark, toggle }
}
