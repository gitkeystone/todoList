<script setup lang="ts">
import { Toaster as Sonner } from 'vue-sonner'

import { useTheme } from '@/composables/useTheme'

// toast 主题跟随应用（深浅色），样式使用 vue-sonner 自带卡片样式（main.ts 已导入其 CSS）
const { isDark } = useTheme()
</script>

<template>
  <Sonner position="top-center" :theme="isDark ? 'dark' : 'light'" :duration="2500" />
</template>
