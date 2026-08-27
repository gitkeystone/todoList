<script setup lang="ts">
import { computed } from 'vue'
import { Trash2 } from 'lucide-vue-next'

import { useTodoStore } from '@/stores/todo'
import type { Todo } from '@/types/todo'

const props = defineProps<{ todo: Todo }>()

const store = useTodoStore()

/** 展示相对时间：今天显示 HH:mm，昨天显示"昨天"，更早显示 M月D日（PRD FR-10/§8.5） */
const timeText = computed(() =>
  formatTime(props.todo.completed ? props.todo.completedAt : props.todo.createdAt),
)

function formatTime(iso: string | null): string {
  if (!iso) return ''
  const d = new Date(iso)
  const now = new Date()
  const todayStart = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()
  const dayStart = new Date(d.getFullYear(), d.getMonth(), d.getDate()).getTime()
  const diffDays = Math.floor((todayStart - dayStart) / 86_400_000)
  if (diffDays <= 0) {
    return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  if (diffDays === 1) return '昨天'
  return d.toLocaleDateString('zh-CN', { month: 'numeric', day: 'numeric' })
}

async function onToggle() {
  try {
    await store.toggle(props.todo)
  } catch {
    /* 错误已由 store 统一 toast */
  }
}

async function onDelete() {
  try {
    await store.remove(props.todo.id)
  } catch {
    /* 错误已由 store 统一 toast */
  }
}

/** 键盘操作：焦点在条目内时 Delete/Backspace 删除（PRD §4.4 I-7） */
function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Delete' || e.key === 'Backspace') {
    e.preventDefault()
    void onDelete()
  }
}
</script>

<template>
  <li class="todo-item" :class="{ done: todo.completed }" @keydown="onKeydown">
    <button
      class="checkbox"
      type="button"
      :aria-label="todo.completed ? '标记为未完成' : '标记为已完成'"
      :aria-pressed="todo.completed"
      @click="onToggle"
    >
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path class="check" d="M5 12.5l4.5 4.5L19 7.5" />
      </svg>
    </button>
    <span class="title">{{ todo.title }}</span>
    <time class="time">{{ timeText }}</time>
    <button class="delete" type="button" aria-label="删除待办" @click="onDelete">
      <Trash2 :size="16" />
    </button>
  </li>
</template>

<style scoped>
.todo-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 13px 16px;
  border-radius: 14px;
  background: var(--bg);
  border: 1px solid var(--line);
  list-style: none;
  transition:
    background-color 0.2s ease,
    box-shadow 0.2s ease;
}

.todo-item:hover {
  background: var(--bg-elevated);
}

/* 圆形勾选框（Apple 式） */
.checkbox {
  position: relative;
  width: 24px;
  height: 24px;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  border: 2px solid var(--text-tertiary);
  border-radius: 9999px;
  background: transparent;
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    border-color 0.2s ease,
    transform 0.1s ease;
}

/* 触屏点击目标 ≥ 44px（PRD §4.6） */
.checkbox::after {
  content: '';
  position: absolute;
  inset: -12px;
  border-radius: 9999px;
}

.checkbox:hover {
  border-color: var(--accent);
}

.checkbox:active {
  transform: scale(0.9);
}

.check {
  fill: none;
  stroke: #fff;
  stroke-width: 2.6;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-dasharray: 24;
  stroke-dashoffset: 24;
  transition: stroke-dashoffset 0.25s cubic-bezier(0.32, 0.72, 0, 1);
}

.todo-item.done .checkbox {
  background: var(--accent);
  border-color: var(--accent);
}

.todo-item.done .check {
  stroke-dashoffset: 0;
}

.title {
  flex: 1;
  min-width: 0;
  font-size: 16px;
  line-height: 1.45;
  color: var(--text-primary);
  overflow-wrap: anywhere;
  transition: color 0.2s ease;
}

.todo-item.done .title {
  color: var(--text-secondary);
  text-decoration: line-through;
  text-decoration-color: color-mix(in srgb, var(--text-secondary) 60%, transparent);
}

.time {
  flex-shrink: 0;
  font-size: 12px;
  color: var(--text-tertiary);
  font-variant-numeric: tabular-nums;
}

.delete {
  position: relative;
  width: 30px;
  height: 30px;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 9999px;
  background: transparent;
  color: var(--text-tertiary);
  cursor: pointer;
  opacity: 0;
  visibility: hidden;
  transition:
    opacity 0.15s ease,
    background-color 0.15s ease,
    color 0.15s ease;
}

/* 触屏点击目标 ≥ 44px */
.delete::after {
  content: '';
  position: absolute;
  inset: -7px;
  border-radius: 9999px;
}

.todo-item:hover .delete,
.todo-item:focus-within .delete {
  opacity: 1;
  visibility: visible;
}

.delete:hover {
  background: color-mix(in srgb, var(--danger) 12%, transparent);
  color: var(--danger);
}

/* 触屏设备：删除按钮常显（PRD §4.6） */
@media (hover: none) {
  .delete {
    opacity: 1;
    visibility: visible;
  }
}
</style>
