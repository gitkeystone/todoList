<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import type { Component } from 'vue'
import { Check, Monitor, Moon, Sun } from 'lucide-vue-next'

import { useTheme, type ThemeMode } from '@/composables/useTheme'

const { mode, isDark, setMode } = useTheme()

const open = ref(false)
const menuEl = ref<HTMLDivElement | null>(null)

const options: { key: ThemeMode; label: string; icon: Component }[] = [
  { key: 'light', label: '浅色', icon: Sun },
  { key: 'dark', label: '深色', icon: Moon },
  { key: 'system', label: '跟随系统', icon: Monitor },
]

const themeIcon = computed(() => (isDark.value ? Moon : Sun))

function onDocumentClick(e: MouseEvent) {
  if (menuEl.value && !menuEl.value.contains(e.target as Node)) open.value = false
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') open.value = false
}

onMounted(() => {
  document.addEventListener('click', onDocumentClick)
  document.addEventListener('keydown', onKeydown)
})
onUnmounted(() => {
  document.removeEventListener('click', onDocumentClick)
  document.removeEventListener('keydown', onKeydown)
})

function select(next: ThemeMode) {
  setMode(next)
  open.value = false
}
</script>

<template>
  <header class="site-header">
    <div class="inner">
      <span class="logo" aria-hidden="true">◐</span>
      <span class="brand">Todo List</span>
      <div ref="menuEl" class="theme">
        <button
          class="theme-toggle"
          type="button"
          :aria-label="
            '切换主题，当前：' + (mode === 'system' ? '跟随系统' : isDark ? '深色' : '浅色')
          "
          :aria-expanded="open"
          @click="open = !open"
        >
          <component :is="themeIcon" :size="17" />
        </button>
        <div v-if="open" class="menu" role="menu" aria-label="选择主题">
          <button
            v-for="opt in options"
            :key="opt.key"
            type="button"
            role="menuitemradio"
            :aria-checked="mode === opt.key"
            :class="{ active: mode === opt.key }"
            @click="select(opt.key)"
          >
            <component :is="opt.icon" :size="15" class="opt-icon" aria-hidden="true" />
            <span>{{ opt.label }}</span>
            <Check v-if="mode === opt.key" :size="15" class="check" aria-hidden="true" />
          </button>
        </div>
      </div>
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

.theme {
  position: relative;
  margin-left: auto;
}

.theme-toggle {
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

.theme-toggle:hover,
.theme-toggle[aria-expanded='true'] {
  background: var(--bg-elevated);
  color: var(--text-primary);
}

.theme-toggle:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

.menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 160px;
  padding: 5px;
  border-radius: 14px;
  background: var(--bg);
  border: 1px solid var(--line);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
  animation: fade-in-up 0.2s cubic-bezier(0.32, 0.72, 0, 1) both;
}

.menu button {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 8px 10px;
  border: none;
  border-radius: 9px;
  background: transparent;
  color: var(--text-primary);
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.menu button:hover {
  background: var(--bg-elevated);
}

.menu button:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: -2px;
}

.opt-icon {
  color: var(--text-secondary);
}

.check {
  margin-left: auto;
  color: var(--accent);
}
</style>
