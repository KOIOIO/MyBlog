import {createApp} from 'vue'
import '@/assets/base.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import {pinia} from "@/stores";
import {i18n} from "@/i18n";
import {config} from "md-editor-v3";
import revealDirective from "@/directives/reveal";

// mermaid 本地化：不走 unpkg CDN（国内加载慢/易失败），固定使用 vendor 内 mermaid@11.3.0
config({editorExtensions: {mermaid: {js: "/vendor/mermaid.min.js"}}});

const app = createApp(App)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.use(pinia).use(router).use(i18n)
app.directive('reveal', revealDirective)

app.mount('#app')
