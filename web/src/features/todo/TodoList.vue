<script setup lang="ts">
import { Skeleton } from '@/components/ui/skeleton'
import { useTodoStore } from '@/stores/todo'
import TodoItem from './TodoItem.vue'

const store = useTodoStore()
</script>

<template>
  <div class="todo-list">
    <TransitionGroup v-if="!store.loading && store.loaded" name="list" tag="ul">
      <TodoItem v-for="todo in store.visibleItems" :key="todo.id" :todo="todo" />
    </TransitionGroup>

    <div v-else-if="store.loading" class="skeletons" aria-label="加载中">
      <Skeleton v-for="i in 4" :key="i" class="skeleton" />
    </div>

    <div v-else-if="store.visibleItems.length === 0" class="empty">
      <template v-if="store.allCount === 0">
        <p class="empty-icon" aria-hidden="true">✓</p>
        <p class="empty-title">欢迎使用 Todo List</p>
        <p class="empty-desc">在上方输入一条待办，回车即可创建</p>
      </template>
      <template v-else>
        <p class="empty-icon" aria-hidden="true">○</p>
        <p class="empty-title">这里空空如也</p>
        <p class="empty-desc">
          {{ store.status === 'completed' ? '还没有已完成的待办' : '当前筛选下没有待办' }}
        </p>
      </template>
    </div>
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
  margin: 0 0 12px;
  font-size: 40px;
  line-height: 1;
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
</style>
