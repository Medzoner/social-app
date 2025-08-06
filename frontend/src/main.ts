import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './assets/main.css'
import router from './router'
import { useAuthStore } from './stores/auth'
import { initNotificationSocket } from './notifications/ws'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

const auth = useAuthStore()
auth.restoreSession()

initNotificationSocket()

app.mount('#app')
