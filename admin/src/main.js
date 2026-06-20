import { createApp } from 'vue'
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'
import ToastService from 'primevue/toastservice'
import App from './App.vue'
import router from './router'
import './assets/styles/main.css'

// Paksa dark mode pada elemen HTML
document.documentElement.classList.add('dark-mode')

const app = createApp(App)

app.use(PrimeVue, {
    theme: {
        preset: Aura,
        options: {
            darkModeSelector: '.dark-mode'
        }
    }
})
app.use(ToastService)
app.use(router)

app.mount('#app')
