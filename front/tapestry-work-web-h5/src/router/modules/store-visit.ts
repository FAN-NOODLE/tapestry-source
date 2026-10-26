export default [
  {
    path: "/store-visit",
    name: "store-visit",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/store-visit/index.vue")),
        "store-visit"
      ),
    meta: {
      title: "common.pageTitle.1",
      keepAlive: true,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/store-visit-dashboard",
    name: "store-visit-dashboard",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/store-visit/dashboard.vue")),
        "store-visit-dashboard"
      ),
    meta: {
      title: "common.pageTitle.2",
      keepAlive: true,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/store-visit-create",
    name: "store-visit-create",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/store-visit/create.vue")),
        "store-visit-create"
      ),
    meta: {
      title: "common.pageTitle.3",
      keepAlive: false,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/store-visit-submit",
    name: "store-visit-submit",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/store-visit/submit.vue")),
        "store-visit-submit"
      ),
    meta: {
      title: "common.pageTitle.3",
      keepAlive: false,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/store-visit-start",
    name: "store-visit-start",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/store-visit/start.vue")),
        "store-visit-start"
      ),
    meta: {
      title: "common.pageTitle.4",
      keepAlive: false,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/store-visit-record",
    name: "store-visit-record",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/store-visit/record.vue")),
        "store-visit-record"
      ),
    meta: {
      title: "common.pageTitle.5",
      keepAlive: true,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/store-visit-detail",
    name: "store-visit-detail",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/store-visit/detail.vue")),
        "store-visit-detail"
      ),
    meta: {
      title: "common.pageTitle.7",
      keepAlive: false,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/store-visit-report",
    name: "store-visit-report",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/store-visit/report.vue")),
        "store-visit-report"
      ),
    meta: {
      title: "common.pageTitle.7",
      keepAlive: false,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/store-visit-edit",
    name: "store-visit-edit",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/store-visit/edit.vue")),
        "store-visit-edit"
      ),
    meta: {
      title: "common.pageTitle.4",
      keepAlive: false,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
]