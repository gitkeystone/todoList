<script setup lang="ts">
import { useTodoStore } from '@/stores/todo'

const store = useTodoStore()

async function clear() {
  try {
    await store.clearCompleted()
  } catch {
    /* 错误已由 store 统一 toast */
  }
}
</script>

<template>
  <footer class="todo-footer">
    <span class="remaining" aria-live="polite">还剩 {{ store.activeCount }} 项未完成</span>
    <button
      v-if="store.completedCount > 0"
      class="clear"
      type="button"
      :disabled="store.clearing"
      @click="clear"
    >
      {{ store.clearing ? '清除中…' : `清除已完成 (${store.completedCount})` }}
    </button>
  </footer>
</template>

<style scoped>
.todo-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 20px;
  padding: 0 6px;
  font-size: 13px;
}

.remaining {
  color: var(--text-secondary);
}

.clear {
  border: none;
  background: none;
  padding: 4px 8px;
  border-radius: 8px;
  color: var(--accent);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition:
    background-color 0.15s ease,
    opacity 0.15s ease;
}

.clear:hover {
  background: color-mix(in srgb, var(--accent) 10%, transparent);
}

.clear:disabled {
  opacity: 0.5;
  cursor: default;
}

.clear:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}
</style>
