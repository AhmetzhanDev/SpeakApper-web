import { createApp } from 'vue'
import App from './App.vue'
import './style.css'
import { initAllAnimations } from './animations.js'

createApp(App).mount('#app')

// Инициализация анимаций после монтирования приложения
document.addEventListener('DOMContentLoaded', () => {
  initAllAnimations();
}); 