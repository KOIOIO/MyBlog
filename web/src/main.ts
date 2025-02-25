import {createApp} from 'vue'
import '@/assets/base.css'
import App from './App.vue'
import router from './router'
import {pinia} from "@/stores";

const app = createApp(App)

app.use(pinia).use(router)

app.mount('#app')