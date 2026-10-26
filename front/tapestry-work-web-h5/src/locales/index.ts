
import VueI18n from 'vue-i18n'
import zhCN from './lang/zh-CN'
import enUS from './lang/en-US'
import Vue from 'vue'

Vue.use(VueI18n)


export const i18n = new VueI18n({
  locale: localStorage.getItem('lang') || 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS
  }
})