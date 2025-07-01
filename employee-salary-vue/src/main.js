import { createApp } from 'vue'
import App from './App.vue'
import locale from 'element-plus/dist/locale/zh-cn.js'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
// 导入路由
import router from "@/router/index.js";
// 导入 Pinia 和持久化插件
import { createPinia } from 'pinia'
import piniaPersistedState from 'pinia-plugin-persistedstate'

const pinia = createPinia()

// 使用持久化插件
pinia.use(piniaPersistedState)

const app = createApp(App);

// 注册路由、Pinia 和 ElementPlus
app.use(router);
app.use(pinia);
app.use(ElementPlus, { locale });

// 挂载应用
app.mount('#app')