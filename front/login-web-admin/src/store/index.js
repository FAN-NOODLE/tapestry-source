import Vue from "vue";
import Vuex from "vuex";
import User from "./modules/UserInfo";
import Menu from "./modules/MenuList";
import App from "./modules/AppKey";
import Temp from "./modules/Template.js";
import Permission from "./modules/Permission";
Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    USER: User,
    MENU: Menu,
    TEMP: Temp,
    APP: App,
    PERMISSION: Permission
  }
});
