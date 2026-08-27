import pluginVue from 'eslint-plugin-vue'
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript'

// Vue 3 + TypeScript 官方推荐 flat 配置（@vue/eslint-config-typescript 负责 .vue 解析器接线）
export default defineConfigWithVueTs(
  {
    ignores: ['dist/**', 'node_modules/**', 'coverage/**', 'public/**'],
  },

  pluginVue.configs['flat/recommended'],
  vueTsConfigs.recommended,

  {
    rules: {
      'vue/multi-word-component-names': 'off',
      '@typescript-eslint/no-explicit-any': 'warn',
    },
  },
)
