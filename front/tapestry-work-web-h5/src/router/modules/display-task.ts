export default [

  {
    path: "/display-task",
    name: "display-task",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/display-task/list.vue")),
        "display-task"
      ),
    meta: {
      title: "WeStore VM",
      keepAlive: true,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/display-task-detail",
    name: "display-task-detail",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/display-task/detail.vue")),
        "display-task-detail"
      ),
    meta: {
      title: "任务详情",
      keepAlive: true,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/display-task-detail-vm",
    name: "display-task-detail-vm",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/display-task/detail-vm.vue")),
        "display-task-detail-vm"
      ),
    meta: {
      title: "任务详情",
      keepAlive: true,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/display-task-comment",
    name: "display-task-comment",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/display-task/comment.vue")),
        "display-task-comment"
      ),
    meta: {
      title: "添加备注",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/display-task-preview",
    name: "display-task-preview",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/display-task/preview.vue")),
        "display-task-preview"
      ),
    meta: {
      title: "预览长图",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/display-task-audit",
    name: "display-task-audit",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/display-task/audit.vue")),
        "display-task-audit"
      ),
    meta: {
      title: "审核",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/display-store-list",
    name: "display-store-list",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/display-task/store-list.vue")),
        "display-store-list"
      ),
    meta: {
      title: "",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: false, // 需要sdk
    },
  },
]