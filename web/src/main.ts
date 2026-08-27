import { createPinia } from 'pinia'
import { createApp } from 'vue'

// vue-sonner 内置样式（toast 卡片/定位/动画，需显式导入）
import 'vue-sonner/style.css'
import './style.css'
import App from './App.vue'

createApp(App).use(createPinia()).mount('#app')
