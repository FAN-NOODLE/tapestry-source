import Vue from "vue";
import Router from "vue-router";
import Routes from "./router/Register";
import Ciphertext from "./utils/Ciphertext";

Vue.use(Router);

let user = {};

const mRouter = new Router({
  mode: "history",
  base: process.env.BASE_URL,
  routes: Routes
});

let checkauth = to => {
  return new Promise(resolve => {
    user = JSON.parse(localStorage.getItem("USERINFO"));
    let authstr = localStorage.getItem("TEMP"),
      parse = Ciphertext("decode"), // 解密
      key = to.meta.router_key;
    let auths = authstr ? new Set([...parse(authstr).split(",")]) : resolve(3);
    let f_keys = key
      .split("_")
      .slice(0, key.split("_").length - 1)
      .join("_");
    if (f_keys === "myapp") {
      resolve(1);
    }
    if (auths.has(f_keys)) {
      resolve(auths.has(to.meta.router_key) ? 1 : 2);
    } else {
      resolve(2);
    }
  });
};

// 全局路由导航守卫，用于验证模块权限
mRouter.beforeEach((to, from, next) => {
  if (to.name === null) {
    next("/404");
  }
  // 模块需要进行权限验证
  if (to.meta && to.meta.router_key) {
    // 没有用户信息时跳转至登录
    if (!user) {
      next("/login");
    } else {
      // 验证权限
      checkauth(to).then(res => {
        switch (res) {
          case 1:
            next();
            break;
          case 2:
            next("/402");
            break;
          case 3:
            next("/login");
            break;
        }
        // res === 1 ? next() : next("/402");
      });
    }
  } else {
    next();
  }
});

export default mRouter;
