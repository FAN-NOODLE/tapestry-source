export default [

  {
    path: "/user/qrcode",
    name: "qrcode",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/user/qrcode.vue")),
        "qrcode"
      ),
    meta: {
      title: "添加专属顾问",
      no_need_login: true, // 不需要 login
      allow_outer_enter: true, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/user/store",
    name: "userstore",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/user/store.vue")),
        "userstore"
      ),
    meta: {
      title: "门店列表",
      no_need_login: true, // 不需要 login
      allow_outer_enter: true, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/user/preview",
    name: "preview",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/user/preview.vue")),
        "preview"
      ),
    meta: {
      title: "添加专属顾问",
      no_need_login: true, // 不需要 login
      allow_outer_enter: true, // 是否允许外部人访问
      need_wx_sdk: false, // 需要sdk
    },
  },
]