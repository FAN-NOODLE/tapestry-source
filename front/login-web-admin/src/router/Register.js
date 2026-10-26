// 基础路由
import Login from "@/views/base/common/login";
import OldLogin from "@/views/base/common/old-login";
import Forget from "@/views/base/common/forget";
import NotFound from "@/views/base/common/404";
import SsoLogin from "@/views/base/common/sso-login";

let routes = [
  {
    path: "/account-login",
    name: "account-login",
    component: Login
  },
  {
    path: "/login",
    name: "login",
    component: SsoLogin
  },
  {
    path: "/zijidengluyongde",
    name: "oldlogin",
    component: OldLogin
  },
  {
    path: "/",
    name: "login",
    component: SsoLogin
  },
  {
    path: "/login/forget",
    name: "Forget",
    component: Forget
  },
  {
    path: "/404",
    name: "404",
    component: NotFound
  }
];
export default routes;
