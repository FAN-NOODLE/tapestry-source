import { RouteConfig } from "vue-router";
import myapp from "./Collect";
const isProduction = process.env.NODE_ENV === "production";
const Register: Array<RouteConfig> = [
  {
    path: "/",
    redirect: {
      name: isProduction ? "myapp_home" : "login",
    },
  },
  {
    path: process.env.VUE_APP_HASH + "base/404",
    name: "404",
    component: () => import("@/views/base/404.vue"),
  },
  {
    path: "/forget",
    name: "Forget",
    component: () => import("@/views/base/forget.vue"),
  },
  {
    path: "/download",
    name: "Demo",
    component: () => import("@/views/base/demo.vue"),
  },
  {
    path: "/canvas",
    name: "Canvas",
    component: () => import("@/views/base/canvas.vue"),
  },
  {
    path: "/callback/sso",
    name: "Sso",
    component: () => import("@/views/base/sso.vue"),
  },
];
Register.push(myapp);
// 如果不是线上则将本地登录路由注入
if (!isProduction) {
  Register.push({
    path: "/login",
    name: "login",
    component: () => import("@/views/base/login-transit.vue"),
  });
}
export default Register;
