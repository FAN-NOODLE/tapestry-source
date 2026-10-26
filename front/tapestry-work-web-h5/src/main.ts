import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import Helper from "./utils/helper";
import { Lazyload } from 'vant';

import "@vant/touch-emulator";
import "@/styles/index.scss";
import "vant/lib/icon/local.css";
import VueRouter from "vue-router";
import FastClick from "fastclick";
import "./plugInUnit";
import { sensorsInit, sensorsTrack } from '@/utils/sensors/index';
import SelectItem from '@/pages/display-task/components/select-item.vue';
import { i18n } from '@/locales'

import { Locale } from 'vant';
// 引入英文语言包
import enUS from 'vant/es/locale/lang/en-US';
if(i18n.locale === 'en-US') {
  Locale.use('en-US', enUS);
}

Vue.component('select-item', SelectItem)
Vue.use(VueRouter);
Vue.config.productionTip = false;
Vue.prototype.$helper = Helper;
// 神策数据上报 -- 初始化
const sensors = sensorsInit()
// 神策数据上报 -- 初始化并挂载全局
Vue.prototype.sensors = sensors
Vue.prototype.$track = sensorsTrack
Vue.use(Lazyload, {
  lazyComponent: true,
});
Vue.directive("title", {
  inserted: function (el) {
    document.title = el.dataset.title as string;
  },
});
Vue.directive('throttle', {
  bind: (el, binding) => {
    let throttleTime = binding.value; // 防抖时间
    if (!throttleTime) { // 用户若不设置防抖时间，则默认2s
      throttleTime = 2000;
    }
    let cbFun;
    el.addEventListener('click', event => {
      if (!cbFun) { // 第一次执行
        cbFun = setTimeout(() => {
          cbFun = null;
        }, throttleTime);
      } else {
        event && event.stopImmediatePropagation();
      }
    }, true);
  },
});
if ("addEventListener" in document) {
  document.addEventListener(
    "DOMContentLoaded",
    () => {
      (FastClick as any).attach(document.body);
    },
    false
  );
}
new Vue({
  render: (h) => h(App),
  router,
  store,
  i18n
}).$mount("#app");
