<script setup lang="ts">
import { Moon, Sun } from 'lucide-vue-next'

import { useTheme } from '@/composables/useTheme'

const { isDark, toggle } = useTheme()
</script>

<template>
  <header class="site-header">
    <div class="inner">
      <span class="logo" aria-hidden="true">◐</span>
      <span class="brand">Todo List</span>
      <button
        class="theme-toggle"
        type="button"
        :aria-label="isDark ? '切换到浅色模式' : '切换到深色模式'"
        @click="toggle"
      >
        <Sun v-if="isDark" :size="17" />
        <Moon v-else :size="17" />
      </button>
    </div>
  </header>
</template>

<style scoped>
.site-header {
  position: sticky;
  top: 0;
  z-index: 50;
  background: var(--frosted);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--line);
}

.inner {
  max-width: 672px;
  margin: 0 auto;
  padding: 10px 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo {
  font-size: 18px;
  color: var(--accent);
}

.brand {
  font-size: 15px;
  font-weight: 600;
  letter-spacing: -0.01em;
}

.theme-toggle {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border: none;
  border-radius: 9999px;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  transition:
    background-color 0.15s ease,
    color 0.15s ease;
}

.theme-toggle:hover {
  background: var(--bg-elevated);
  color: var(--text-primary);
}

.theme-toggle:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}
</style>
