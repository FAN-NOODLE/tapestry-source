import VueRouter, { Route } from "vue-router";
import Common from "./modules/common";
import Base from "./modules/base";
import StoreVisit from "./modules/store-visit";
import DisplayTask from './modules/display-task'
import User from "./modules/user";
import Store from "../store/index";
import { loginFnc } from "../utils/login";
import wx from "../utils/wx";
import { Router } from "@/types/router";
import { Toast } from "vant";
import sensors from "sa-sdk-javascript";
import { getDeviceInfo } from "@/utils/tools";
import { i18n } from '@/locales'
const store: any = Store;
const hash = window.location.pathname.split("/")[2];
const IDENTIFIER = window.location.pathname.split("/")[1] || "";
Store.commit("common/setEnterpriseHash", hash);
const routes: Router[] = [...Common, ...Base, ...DisplayTask, ...StoreVisit, ...User];
const router = new VueRouter({
  base: IDENTIFIER + `/${hash}/`,
  mode: "history",
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition && to.meta.keepAlive) {
      return savedPosition;
    }
    return { x: 0, y: 0 };
  },
});
// 全局前置守卫
router.beforeEach(async (to: Route, from: Route, next) => {
  console.log("beforeEach-to", to, from);
  
  Toast.loading({
    duration: 0, // 持续展示 toast
    message: String(i18n.t("common.toast.22")),
    forbidClick: true,
  });
  console.log("设置前", document.title);
  
  // 设置标题
  if (to.meta.title) {
    document.title = String(i18n.t(to.meta.title)) ;
  }
  console.log("设置后", document.title);

  // 判断是否需要微信授权
  const needWxSdk =
    to.meta.need_wx_sdk &&
    process.env.NODE_ENV !== "test" &&
    process.env.NODE_ENV !== "development";

  // 不需要登陆
  console.log("no_need_login", to.meta);
  if (to.meta.no_need_login) {
    
    // 已经登陆了, 还需要微信授权
    if(needWxSdk) {
      const host = window.location.origin
      const hash = window.location.pathname.split("/")[2];
      const IDENTIFIER = window.location.pathname.split("/")[1] || "";
      const url = `${host}/${IDENTIFIER}/${hash}${to.fullPath}`
      const { platform } = getDeviceInfo()
      console.log(11111, window.location.href);
      console.log(11111, url);
      
      await wx.wxConfig(platform === 'android' ? url : decodeURIComponent(window.location.href).replaceAll(' ', '+'))
    }
    next();
    return false;
  }
  // 判断登陆
  console.log("IS_LOGINED", store.state.common.IS_LOGINED);
  
  // 判断是否登陆
  if (!store.state.common.IS_LOGINED) {

    // 需要微信授权
    Promise.all([
      needWxSdk ? wx.wxConfig() : Promise.resolve(),
      loginFnc(to.query, router),
    ])
      .then(function () {
        console.log("允许进入, 芝麻开门 then");
        window.wx.onNetworkStatusChange(function(res) {
          const isConnected = res.isConnected; // true 为网络通畅 false 为断网
          const networkType = res.networkType; // 返回网络类型2g，3g，4g，wifi
          console.log(isConnected, "isConnected");
          console.log(networkType, "networkType");
          // 断网提示
          if (!isConnected || networkType === 'none') {
            Toast("你没网络了哦，请检查下网络设置~");
          }
          // 弱网提示
          if (networkType === '2g' || networkType === '3g') {
            Toast("你的网络不太好哦~");
          }
        })
        // 判断是否允许外部进入
        if (to.meta.allow_outer_enter) {
          next();
          return false;
        }
        // 判断是否有用户信息
        if (!store.state.user.USERINFO.hasOwnProperty("user_id")) {
          console.log("允许进入, 芝麻开门 no allow");
          router.replace("/no-allow");
          return false;
        }
        next();
      })
      .catch(function (e) {
        console.log(e, "允许进入, 芝麻开门 catch");
        next();
      });
    return false;
  }
  // 能走到这, 不容易,
  // 表示你的身份允许登陆
  console.log("允许进入, 芝麻开门 success");
  next();
});
// todo: 全局后置钩子
router.afterEach((to: Route, from: Route) => {
  Toast.clear();
  console.log('我是给上报准备的标题', to.meta.title);
  
  if(to.meta.title) {
    store.commit("storeVisit/setTitle", to.meta.title);
  }
  // sensors.quick('autoTrack')
});
export default router;
