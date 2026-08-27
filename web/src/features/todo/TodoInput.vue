<script setup lang="ts">
import { nextTick, onMounted, ref } from 'vue'
import { Plus } from 'lucide-vue-next'

import { useTodoStore } from '@/stores/todo'

const store = useTodoStore()
const title = ref('')
const inputEl = ref<HTMLInputElement | null>(null)
const shaking = ref(false)
let composing = false

onMounted(() => {
  // 桌面端自动聚焦（移动端不弹键盘，PRD §4.4 I-1）
  if (window.matchMedia('(hover: hover)').matches) inputEl.value?.focus()
})

function onCompositionStart() {
  composing = true
}

function onCompositionEnd() {
  composing = false
}

function onKeydown(e: KeyboardEvent) {
  // 中文输入法组词中不触发提交（PRD §12 R-3）
  if (e.key === 'Enter' && !composing && !e.isComposing) void submit()
}

async function submit() {
  const value = title.value.trim()
  if (!value) {
    shake()
    return
  }
  try {
    await store.create(value)
    title.value = ''
    await nextTick()
    inputEl.value?.focus()
  } catch {
    /* 错误已由 store 统一 toast（PRD §4.4 I-9） */
  }
}

async function shake() {
  shaking.value = false
  await nextTick()
  shaking.value = true
}

/** 供外部（⌘K 快捷键）聚焦输入框 */
defineExpose({
  focus: () => inputEl.value?.focus(),
})
</script>

<template>
  <div class="todo-input" :class="{ shaking }">
    <Plus :size="20" class="icon" aria-hidden="true" />
    <input
      ref="inputEl"
      v-model="title"
      type="text"
      placeholder="添加一项待办…"
      maxlength="200"
      aria-label="新待办标题"
      @keydown="onKeydown"
      @compositionstart="onCompositionStart"
      @compositionend="onCompositionEnd"
    />
  </div>
</template>

<style scoped>
.todo-input {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 18px;
  height: 56px;
  border-radius: 9999px;
  background: var(--bg-elevated);
  border: 1px solid var(--line);
  transition:
    box-shadow 0.2s ease,
    border-color 0.2s ease,
    background-color 0.2s ease;
}

.todo-input:focus-within {
  background: var(--bg);
  border-color: var(--accent);
  box-shadow: 0 0 0 4px color-mix(in srgb, var(--accent) 18%, transparent);
}

.icon {
  color: var(--text-tertiary);
  flex-shrink: 0;
}

input {
  flex: 1;
  min-width: 0;
  border: none;
  outline: none;
  background: transparent;
  font-size: 17px;
  color: var(--text-primary);
}

input::placeholder {
  color: var(--text-tertiary);
}

.shaking {
  animation: shake 0.3s ease;
}
</style>
