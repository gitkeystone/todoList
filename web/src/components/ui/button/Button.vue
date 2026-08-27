<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import { cn } from '@/lib/utils'
import { buttonVariants, type ButtonVariants } from './button'

interface Props {
  variant?: ButtonVariants['variant']
  size?: ButtonVariants['size']
  class?: HTMLAttributes['class']
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'default',
  class: '',
})
</script>

<template>
  <button :class="cn(buttonVariants({ variant, size }), props.class)">
    <slot />
  </button>
</template>
