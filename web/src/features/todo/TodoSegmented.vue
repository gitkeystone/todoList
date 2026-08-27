<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'

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

const buttonEls = ref<(HTMLButtonElement | null)[]>([])
const thumb = ref({ width: 0, x: 0, visible: false })

const activeIndex = computed(() => tabs.value.findIndex((t) => t.key === store.status))

/** 滑块跟随激活项（PRD §4.2.4 分段控件切换动画） */
function updateThumb() {
  const btn = buttonEls.value[activeIndex.value]
  if (!btn) return
  thumb.value = { width: btn.offsetWidth, x: btn.offsetLeft, visible: true }
}

onMounted(() => {
  void nextTick(updateThumb)
  window.addEventListener('resize', updateThumb)
})
onUnmounted(() => window.removeEventListener('resize', updateThumb))

// 激活项或计数（宽度）变化时重算滑块
watch(activeIndex, () => void nextTick(updateThumb))
watch(
  () => tabs.value.map((t) => t.count).join(','),
  () => void nextTick(updateThumb),
)
</script>

<template>
  <div class="segmented" role="tablist" aria-label="筛选待办">
    <div
      class="thumb"
      :class="{ visible: thumb.visible }"
      :style="{ width: `${thumb.width}px`, transform: `translateX(${thumb.x}px)` }"
      aria-hidden="true"
    />
    <button
      v-for="tab in tabs"
      :key="tab.key"
      ref="buttonEls"
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
  position: relative;
  display: flex;
  gap: 3px;
  margin-top: 20px;
  padding: 3px;
  border-radius: 12px;
  background: var(--bg-elevated);
  border: 1px solid var(--line);
}

/* 滑动高亮（Apple 分段控件式） */
.thumb {
  position: absolute;
  top: 3px;
  bottom: 3px;
  left: 0;
  border-radius: 9px;
  background: var(--bg);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  opacity: 0;
  transition:
    transform 0.25s cubic-bezier(0.32, 0.72, 0, 1),
    width 0.25s cubic-bezier(0.32, 0.72, 0, 1),
    opacity 0.2s ease;
}

.thumb.visible {
  opacity: 1;
}

.segmented button {
  position: relative;
  z-index: 1;
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
  transition: color 0.2s ease;
}

.segmented button.active {
  color: var(--text-primary);
}

.segmented button:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: -2px;
}

.count {
  font-size: 12px;
  color: var(--text-tertiary);
  font-variant-numeric: tabular-nums;
  transition: color 0.2s ease;
}

.segmented button.active .count {
  color: var(--text-secondary);
}
</style>
