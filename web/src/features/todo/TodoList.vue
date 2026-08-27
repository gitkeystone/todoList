<script setup lang="ts">
import { AlertTriangle, Inbox, Sparkles } from 'lucide-vue-next'

import { Skeleton } from '@/components/ui/skeleton'
import { useTodoStore } from '@/stores/todo'
import TodoItem from './TodoItem.vue'

const store = useTodoStore()
</script>

<template>
  <div class="todo-list">
    <!-- 加载骨架 -->
    <div v-if="store.loading" class="skeletons" aria-label="加载中">
      <Skeleton v-for="i in 4" :key="i" class="skeleton" />
    </div>

    <!-- 加载失败 + 重试（PRD §4.5 错误态） -->
    <div v-else-if="store.error" class="empty">
      <AlertTriangle :size="36" class="empty-icon" aria-hidden="true" />
      <p class="empty-title">加载失败</p>
      <p class="empty-desc">{{ store.error }}</p>
      <button class="retry" type="button" @click="store.fetchList()">重试</button>
    </div>

    <template v-else>
      <!-- 列表：入场错峰 40ms + FLIP 位移（PRD §4.2.4） -->
      <TransitionGroup name="list" tag="ul">
        <TodoItem
          v-for="(todo, i) in store.visibleItems"
          :key="todo.id"
          :todo="todo"
          :style="{ '--enter-delay': `${i * 40}ms` }"
        />
      </TransitionGroup>

      <!-- 空状态（PRD §4.5） -->
      <div v-if="store.visibleItems.length === 0" class="empty">
        <template v-if="store.allCount === 0">
          <Sparkles :size="36" class="empty-icon" aria-hidden="true" />
          <p class="empty-title">欢迎使用 Todo List</p>
          <p class="empty-desc">在上方输入一条待办，回车即可创建</p>
        </template>
        <template v-else>
          <Inbox :size="36" class="empty-icon" aria-hidden="true" />
          <p class="empty-title">这里空空如也</p>
          <p class="empty-desc">
            {{ store.status === 'completed' ? '还没有已完成的待办' : '当前筛选下没有待办' }}
          </p>
        </template>
      </div>
    </template>
  </div>
</template>

<style scoped>
.todo-list {
  position: relative;
  margin-top: 20px;
}

ul {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin: 0;
  padding: 0;
}

/* TransitionGroup 动效（PRD §4.2.4） */
.list-enter-active {
  transition:
    opacity 0.3s cubic-bezier(0.32, 0.72, 0, 1),
    transform 0.3s cubic-bezier(0.32, 0.72, 0, 1);
  /* 入场错峰：由 TodoItem 上的 --enter-delay 控制 */
  transition-delay: var(--enter-delay, 0ms);
}

.list-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

.list-leave-active {
  position: absolute;
  width: 100%;
  transition:
    opacity 0.25s cubic-bezier(0.4, 0, 0.2, 1),
    transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  transition-delay: 0ms;
}

.list-leave-to {
  opacity: 0;
  transform: translateY(-8px) scale(0.98);
}

.list-move {
  transition: transform 0.3s cubic-bezier(0.32, 0.72, 0, 1);
}

.skeletons {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.skeleton {
  height: 60px;
  border-radius: 14px;
}

.empty {
  padding: 56px 0;
  text-align: center;
  animation: fade-in-up 0.3s cubic-bezier(0.32, 0.72, 0, 1) both;
}

.empty-icon {
  margin: 0 auto 12px;
  display: block;
  color: var(--text-tertiary);
}

.empty-title {
  margin: 0 0 4px;
  font-size: 17px;
  font-weight: 600;
  color: var(--text-primary);
}

.empty-desc {
  margin: 0;
  font-size: 14px;
  color: var(--text-secondary);
}

.retry {
  margin-top: 16px;
  padding: 8px 20px;
  border: none;
  border-radius: 9999px;
  background: var(--accent);
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: opacity 0.15s ease;
}

.retry:hover {
  opacity: 0.85;
}

.retry:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}
</style>
