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
      // 以下为纯格式化规则，与 Prettier 冲突，交由 Prettier 统一处理
      'vue/max-attributes-per-line': 'off',
      'vue/html-self-closing': 'off',
      '@typescript-eslint/no-explicit-any': 'warn',
    },
  },
)
