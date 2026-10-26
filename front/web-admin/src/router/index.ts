import Vue from "vue";
import VueRouter, { RouteConfig, Route } from "vue-router";
import { IKeysItem } from "@/types";
import Routes from "./Register";
import NProgress from "nprogress";
import "nprogress/nprogress.css";
import MenuModule from "@/store/modules/menu";
import BreadCrumbModule from "@/store/modules/bread-crumb";
import { getKeyByBCNodeItem } from "@/utils/keys-map";
import settingApis from "@/api/setting/setting";
import PerModule from "@/store/modules/permission";
import LoginAboat from "@/utils/LoginManager";

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [...Routes];
console.log("scrm-routes", routes);
const router = new VueRouter({
  mode: "history",
  base: "/web-admin",
  routes: routes,
});

/**
 * @description 路由跳转时设置菜单选中的根目录和叶子目录，用于菜单选中高亮
 * @param {Route} to 即将前往的路由
 * */
const setMenuActive: (to: Route) => void = function(to) {
  const root = `myapp_${to.meta.root || to.meta.ename}`;
  const path_list: Array<string> = to.path.split("/");
  const path_length: number = path_list.length;
  let leaf = "";
  MenuModule.setRoot(root);
  // 根据情况存储选中的叶子结点
  if (path_length === 4) {
    leaf = path_list.splice(1, 3).join("_");
  } else if (path_length === 5) {
    if (to.meta.hide) {
      leaf = path_list.splice(1, 3).join("_");
    } else {
      leaf = path_list.splice(1, 4).join("_");
    }
  } else if (path_length > 5) {
    leaf = path_list.splice(1, 4).join("_");
  }
  MenuModule.setLeaf(leaf);
  const url = to.path.replace(/\//g, "_");
  const keyArray = url.split("_myapp");
  const key = "myapp" + keyArray[1];
  const keys = JSON.parse(localStorage.getItem("KEYS") || "{}"); //权限集合
  const route = getKeyByBCNodeItem(keys, key);
  if (route) {
    BreadCrumbModule.setBreadCrumb({ nodeItem: route, to });
  }
};

/**
 * @description 获取权限key集合
 * @returns {Array<KeysItem>}
 * */
const getKeys: () => Array<IKeysItem> = function() {
  return JSON.parse(localStorage.getItem("KEYS") || "[]");
};

/**
 * @description 路由跳转时校验是否拥有权限
 * @param {Route} to 即将前往的路由
 * */
const checkAuth: (to: Route) => Promise<number> = function(to) {
  return new Promise((resolve) => {
    const auths = getKeys();
    const key = to.meta.key;
    const f_keys = key
      .split("_")
      .slice(0, key.split("_").length - 1)
      .join("_");
    if (f_keys === "myapp") {
      resolve(1);
    }
    if (auths.some((item) => item.name === to.meta.title) || auths.some((item) => item.key === f_keys)) {
      resolve(1);
    } else {
      resolve(2);
    }
  });
};

NProgress.configure({ showSpinner: false });

async function getEpConfig() {
  if(!LoginAboat.getToken()) return
  const res: any = await settingApis.getEpConfigApi();
  if (res.code !== 200) return 

  // 企业名称：coach;ks;sw
  const { name } = res.data;
  // 缓存企业名称，给详情页面用
  PerModule.setEqName(name);
}

/**
 * @description 全局路由前置守卫
 * */
router.beforeEach(async (to: Route, from: Route, next: any) => {
  console.log('beforeEach', to.path);
  
  NProgress.start();
  if(window.pluginWebUpdateNotice_){
    window.pluginWebUpdateNotice_.checkUpdate()
  }
  if(to.path !== '/callback/sso'){
    await getEpConfig();
  }
    
  //如果meta中有key，说明需要进行权限校验
  if (to.meta && to.meta.key) {
    const number = await checkAuth(to);
    console.log("beforeEach-number", number)
    switch (number) {
      case 1:
        next();
        break;
      case 2:
        next("/402");
        break;
    }
  } else {
    next();
  }
});

/**
 * @description 全局路由后置守卫
 * */
router.afterEach((to: Route) => {
  setMenuActive(to);
  NProgress.done();
});

export default router;
