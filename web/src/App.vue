<script setup lang="ts">
import { ref } from 'vue'

// M0 脚手架占位页：验证主题切换与前端工程可用（M2 替换为真实 Todo List 界面）
const theme = ref<'light' | 'dark'>(
  document.documentElement.classList.contains('dark') ? 'dark' : 'light',
)

function toggleTheme() {
  const next = theme.value === 'dark' ? 'light' : 'dark'
  theme.value = next
  document.documentElement.classList.toggle('dark', next === 'dark')
  localStorage.setItem('todolist.theme', next)
}
</script>

<template>
  <main class="page">
    <button class="theme-btn" type="button" aria-label="切换深色模式" @click="toggleTheme">
      {{ theme === 'dark' ? '☀️' : '🌙' }}
    </button>
    <h1>Todo List</h1>
    <p>M0 脚手架就绪 · 前端 :5173 · 后端 :8080 <code>/healthz</code></p>
  </main>
</template>

<style scoped>
.page {
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 24px;
  text-align: center;
}

h1 {
  margin: 0;
  font-size: 48px;
  font-weight: 700;
  letter-spacing: -0.03em;
}

p {
  margin: 0;
  font-size: 15px;
  color: var(--text-secondary);
}

code {
  padding: 2px 6px;
  border-radius: 6px;
  background: rgba(0, 0, 0, 0.06);
}

.dark code {
  background: rgba(255, 255, 255, 0.1);
}

.theme-btn {
  position: absolute;
  top: 14px;
  right: 14px;
  padding: 8px 12px;
  border: none;
  border-radius: 9999px;
  background: none;
  font-size: 20px;
  line-height: 1;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.theme-btn:hover {
  background: rgba(0, 0, 0, 0.05);
}

.dark .theme-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}
</style>
