import Vue from "vue";
import "./plugins/Axios";
import App from "./App.vue";
import Router from "./router";
import Store from "./store/index";
import "./plugins/Element.js";
import Global from "./plugins/Global";
import md5 from "js-md5";
import SlideVerify from "vue-monoplasty-slide-verify";
Vue.prototype.$md5 = md5;
Vue.prototype.$global = Global;
Vue.config.productionTip = false;
Vue.use(SlideVerify);
new Vue({
  router: Router,
  store: Store,
  render: h => h(App)
}).$mount("#app");
