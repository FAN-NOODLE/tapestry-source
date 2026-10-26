import Vue from "vue";
import Vuex from "vuex";
import createPersistedState from 'vuex-persistedstate';
import common from "./modules/common";
import user from "./modules/user";
import customer from "./modules/customer";
import materialLibrary from "./modules/material-library";
import displayTask from "./modules/display-task";
import onlineOrder from "./modules/online-order";
import storeVisit from "./modules/store-visit";
import lang from "./modules/lang";
Vue.use(Vuex);
const modules: any = {};
modules.user = user;
modules.customer = customer;
modules.common = common; // https://webpack.js.org/guides/dependency-management/#requirecontext
modules.materialLibrary = materialLibrary;
modules.displayTask = displayTask;
modules.onlineOrder = onlineOrder;
modules.storeVisit = storeVisit;
modules.lang = lang;
// const modulesFiles = require.context('./module', true, /\.ts$/);
// const modules: any = {};
// modulesFiles.keys().forEach((key: any) => {
//   modules[key.replace(/(\.\/|\.ts)/g, '')] = modulesFiles(key).default;
// });
export default new Vuex.Store({
  modules,
  plugins: [
    createPersistedState({
      storage: window.sessionStorage,
    }),
  ],
});
