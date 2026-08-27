<script setup lang="ts">
import { computed } from 'vue'

import { useTodoStore } from '@/stores/todo'
import type { FilterStatus } from '@/types/todo'

const store = useTodoStore()

interface Tab {
  key: FilterStatus
  label: string
  count: number
}

const tabs = computed<Tab[]>(() => [
  { key: 'all', label: '全部', count: store.allCount },
  { key: 'active', label: '进行中', count: store.activeCount },
  { key: 'completed', label: '已完成', count: store.completedCount },
])
</script>

<template>
  <div class="segmented" role="tablist" aria-label="筛选待办">
    <button
      v-for="tab in tabs"
      :key="tab.key"
      type="button"
      role="tab"
      :aria-selected="store.status === tab.key"
      :class="{ active: store.status === tab.key }"
      @click="store.status = tab.key"
    >
      <span class="label">{{ tab.label }}</span>
      <span class="count">{{ tab.count }}</span>
    </button>
  </div>
</template>

<style scoped>
.segmented {
  display: flex;
  gap: 3px;
  margin-top: 20px;
  padding: 3px;
  border-radius: 12px;
  background: var(--bg-elevated);
  border: 1px solid var(--line);
}

.segmented button {
  flex: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 7px 10px;
  border: none;
  border-radius: 9px;
  background: transparent;
  color: var(--text-secondary);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    color 0.2s ease,
    box-shadow 0.2s ease;
}

.segmented button.active {
  background: var(--bg);
  color: var(--text-primary);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.segmented button:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: -2px;
}

.count {
  font-size: 12px;
  color: var(--text-tertiary);
  font-variant-numeric: tabular-nums;
}

.segmented button.active .count {
  color: var(--text-secondary);
}
</style>
