import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import md5 from "js-md5";
import Axios from "axios";
import Components from "@/components/public";
import "@/icon";
import "@/utils/element";
import "./style/global.scss";
import "./element-variables.scss";
import "@/components/crs-cy-components/dist/components.css";
import CrsComponents from "@/components/crs-cy-components";
import { VUE_APP_API_ROOT, VUE_APP_API_VERSION } from "@/utils/env-config";
import { Component } from "vue-property-decorator";
import * as Echarts from "echarts";
import Global from "./plugins/Global";


Component.registerHooks(["beforeRouteEnter"]);


/************************vue相关配置开始***************************/

/**
 * 解决警告
 * 警告示例 chunk-vendors.js:2 [Violation] Added non-passive event listener to a scroll-blocking 'mousewheel' event. Consider marking event handler as 'passive' to make the page more responsive. See https://www.chromestatus.com/feature/5745543795965952
 * */

import "default-passive-events";
import Filter from "./filter";
Vue.config.productionTip = false;
Vue.use(Filter);
Vue.use(Components);
Vue.use(CrsComponents);
Vue.prototype.$md5 = md5;
Vue.prototype.$axios = Axios;
Vue.prototype.$echarts = Echarts;
Vue.prototype.$VUE_APP_API_ROOT = VUE_APP_API_ROOT;
Vue.prototype.$VUE_APP_API_VERSION = VUE_APP_API_VERSION;
Vue.prototype.$global = Global;
/************************vue相关配置结束**************************/
new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount("#app");
