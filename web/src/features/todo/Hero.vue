<script setup lang="ts">
import { computed } from 'vue'

import { useTodoStore } from '@/stores/todo'

const store = useTodoStore()

const todayText = new Intl.DateTimeFormat('zh-CN', {
  month: 'long',
  day: 'numeric',
  weekday: 'long',
}).format(new Date())

const summary = computed(
  () =>
    `还剩 ${store.activeCount} 项 · 已完成 ${store.completedCount} 项 · 全部 ${store.allCount} 项`,
)
</script>

<template>
  <section class="hero">
    <p class="date">{{ todayText }}</p>
    <h1 class="title">今天</h1>
    <p class="summary">{{ summary }}</p>
  </section>
</template>

<style scoped>
.hero {
  padding: 40px 0 24px;
  text-align: center;
}

.date {
  margin: 0 0 4px;
  font-size: 15px;
  color: var(--accent);
  font-weight: 500;
}

.title {
  margin: 0;
  font-size: 48px;
  line-height: 1.1;
  font-weight: 700;
  letter-spacing: -0.03em;
}

.summary {
  margin: 8px 0 0;
  font-size: 15px;
  color: var(--text-secondary);
}

@media (min-width: 768px) {
  .hero {
    padding: 56px 0 32px;
  }
  .title {
    font-size: 56px;
  }
}
</style>
