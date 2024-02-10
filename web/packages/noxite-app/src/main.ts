// @ts-expect-error easy.ui currently throws a ts error, has to be fixed in the future
import { VueComponents, WebComponents } from '@easy/ui'
import messages from '@intlify/unplugin-vue-i18n/messages'
import { createPinia } from 'pinia'
import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'
import router from './router'
import App from './App.vue'
import '@easy/ui/style.css'
import './style.css'
import '@/assets/tokens.css'

WebComponents.EasyUI.install()

export const i18n = createI18n({
  locale: navigator.language || 'en-US',
  fallbackLocale: ['en-US', 'de-DE'],
  messages,
})

const pinia = createPinia()

createApp(App)
  .use(router)
  .use(i18n)
  .use(pinia)
  .component('EzIcon', VueComponents.Icon)
  .mount('#app')
