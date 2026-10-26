export default [
  {
    path: "/no-allow",
    name: "no-allow",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/common/no-auth.vue")),
        "no-auth"
      ),
    meta: {
      title: "授权失败",
      no_need_login: true, // 不需要 login
      allow_outer_enter: true, // 外部人不可以进入
    },
  },
  {
    path: "/link",
    name: "link",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/online-order/info.vue")),
        "online-order-info"
      ),
    meta: {
      title: "订单详情",
      keepAlive: false,
      allow_outer_enter: true, // 是否允许外部人访问
      no_need_login: true, // 需要sdk
    },
  },
];
