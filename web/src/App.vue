<script setup lang="ts">
import { onMounted } from 'vue'

import { Sonner } from '@/components/ui/sonner'
import Hero from '@/features/todo/Hero.vue'
import SiteHeader from '@/features/todo/SiteHeader.vue'
import TodoFooter from '@/features/todo/TodoFooter.vue'
import TodoInput from '@/features/todo/TodoInput.vue'
import TodoList from '@/features/todo/TodoList.vue'
import TodoSegmented from '@/features/todo/TodoSegmented.vue'
import { useTodoStore } from '@/stores/todo'

const store = useTodoStore()

onMounted(() => {
  void store.fetchList()
})
</script>

<template>
  <div class="app">
    <SiteHeader />
    <main class="content">
      <Hero />
      <TodoInput />
      <TodoSegmented />
      <TodoList />
      <TodoFooter />
    </main>
    <Sonner position="top-center" />
  </div>
</template>

<style scoped>
.content {
  max-width: 672px;
  margin: 0 auto;
  padding: 0 20px 48px;
}
</style>
