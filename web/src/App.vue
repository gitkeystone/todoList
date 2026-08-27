<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

import { Sonner } from '@/components/ui/sonner'
import Hero from '@/features/todo/Hero.vue'
import SiteHeader from '@/features/todo/SiteHeader.vue'
import TodoFooter from '@/features/todo/TodoFooter.vue'
import TodoInput from '@/features/todo/TodoInput.vue'
import TodoList from '@/features/todo/TodoList.vue'
import TodoSegmented from '@/features/todo/TodoSegmented.vue'
import { useTodoStore } from '@/stores/todo'

const store = useTodoStore()
const inputRef = ref<InstanceType<typeof TodoInput> | null>(null)

/** 全局快捷键：⌘K / Ctrl+K 聚焦输入框（PRD §4.4 I-6） */
function onGlobalKeydown(e: KeyboardEvent) {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    inputRef.value?.focus()
  }
}

onMounted(() => {
  window.addEventListener('keydown', onGlobalKeydown)
  void store.fetchList()
})

onUnmounted(() => {
  window.removeEventListener('keydown', onGlobalKeydown)
})
</script>

<template>
  <div class="app">
    <SiteHeader />
    <main class="content">
      <Hero />
      <TodoInput ref="inputRef" />
      <TodoSegmented />
      <TodoList />
      <TodoFooter />
    </main>
    <Sonner />
  </div>
</template>

<style scoped>
.content {
  max-width: 672px;
  margin: 0 auto;
  padding: 0 20px 48px;
}
</style>
