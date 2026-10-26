/*
 * @Author: amujoe
 * @Date: 2020-01-02 18:12:23
 * @Description: file content
 */
export default [
  {
    path: "/index",
    name: "index",
    component: (r: (arg0: any) => void) =>
      require.ensure([], () => r(require("@/pages/index.vue")), "index"),
    meta: {
      title: "首页",
      allow_outer_enter: true, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/login",
    name: "login",
    component: (r: (arg0: any) => void) =>
      require.ensure([], () => r(require("@/pages/login.vue")), "login"),
    meta: {
      title: "登录",
      allow_outer_enter: true, // 是否允许外部人访问
      need_wx_sdk: false, // 需要sdk
    },
  },
  {
    path: "/demo",
    name: "demo",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/demo.vue")),
        "demo"
      ),
    meta: {
      title: "demo",
      no_need_login: true, // 不需要 login
      allow_outer_enter: true, // 是否允许外部人访问
      need_wx_sdk: false, // 需要sdk
    },
  },
  {
    path: "/marketing-promotion",
    name: "marketing-promotion",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/marketing-promotion/marketing-promotion.vue")),
        "marketing-promotion"
      ),
    meta: {
      title: "营销推广",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/material-library",
    name: "material-library",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/material-library/material-library.vue")),
        "material-library"
      ),
    meta: {
      title: "素材库",
      keepAlive: true,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/material-library-coach",
    name: "material-library-coach",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/material-library-coach/material-library.vue")),
        "material-library"
      ),
    meta: {
      title: "素材库",
      keepAlive: true,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/my-collect",
    name: "my-collect",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/material-library-coach/my-collect.vue")),
        "my-collect"
      ),
    meta: {
      title: "我的收藏",
      keepAlive: true,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/normal-material",
    name: "normal-material",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/material-library/normal-material.vue")),
        "normal-material"
      ),
    meta: {
      title: "素材库",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/goods-material",
    name: "goods-material",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/material-library/goods-material.vue")),
        "goods-material"
      ),
    meta: {
      title: "素材库",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/material-detail",
    name: "material-detail",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/material-library/goods-detail.vue")),
        "material-detail"
      ),
    meta: {
      title: "素材库",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/text-library",
    name: "text-library",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/material-library/text-library.vue")),
        "text-library"
      ),
    meta: {
      title: "新增文本素材",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/video-library",
    name: "video-library",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/material-library/video-library.vue")),
        "video-library"
      ),
    meta: {
      title: "新增视频素材",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/image-library",
    name: "image-library",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/material-library/image-library.vue")),
        "image-library"
      ),
    meta: {
      title: "新增图片素材",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/material-result",
    name: "material-result",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/material-library/material-result.vue")),
        "material-result"
      ),
    meta: {
      title: "素材库",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/customer-detail",
    name: "customer-detail",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () =>
          r(
            require("@/pages/customer-portrait/portrait-customer/children/customer-detail.vue")
          ),
        "customer-detail"
      ),
    meta: {
      title: "客户详情",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/order-detail",
    name: "order-detail",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () =>
          r(
            require("@/pages/customer-portrait/portrait-customer/children/order-detail.vue")
          ),
        "order-detail"
      ),
  },
  {
    path: "/customer-portrait",
    name: "customer-portrait",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/customer-portrait/customer-portrait.vue")),
        "customer-portrait"
      ),
    meta: {
      title: "客户画像",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/portrait-customer-details",
    name: "portrait-customer-details",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () =>
          r(
            require("@/pages/customer-portrait/portrait-customer-details/portrait-customer-details.vue")
          ),
        "portrait-customer-details"
      ),
    meta: {
      title: "客户画像",
    },
  },
  {
    path: "/public-page",
    name: "public-page",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () =>
          r(
            require("@/pages/customer-portrait/group-portrait/public-page/public-page")
          ),
        "public-page"
      ),
    meta: {
      title: "",
      hide_share_btn: true, // 隐藏系统分享按钮
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/join-activites",
    name: "join-activites",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () =>
          r(
            require("@/pages/customer-portrait/group-portrait/join-activites/join-activites")
          ),
        "join-activites"
      ),
    meta: {
      title: "参加活动",
      hide_share_btn: true, // 隐藏系统分享按钮
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/fission-activities",
    name: "fission-activities",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/fission-activities/fission-activities.vue")),
        "fission-activities"
      ),
    meta: {
      title: "裂变活动",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/single-active",
    name: "single-active",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/marketing-campaign/single-active.vue")),
        "single-active"
      ),
    meta: {
      title: "营销活动",
      // need_wx_open: true, // 必须在微信端打开，否则提示缺省
      hide_share_btn: true, // 隐藏系统分享按钮
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/store-invited",
    name: "store-invited",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/store-invited/store-invited.vue")),
        "store-invited"
      ),
    meta: {
      title: "到店邀约",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/invite-config",
    name: "invite-config",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/invite-config/invite-config.vue")),
        "invite-config"
      ),
    meta: {
      title: "邀请函",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  // 开发需要这个路由调试客户画像
  {
    path: "/new-portrait-customer",
    name: "new-portrait-customer",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/customer-portrait/new-portrait-customer/index.vue")),
        "new-portrait-customer"
      ),
    meta: {
      title: "客户画像",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/communication",
    name: "communication",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/communication/index.vue")),
        "communication"
      ),
    meta: {
      title: "企业传讯",
      keepAlive: true,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/communication-detail/:type",
    name: "communication-detail",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/communication/detail.vue")),
        "communication-detail"
      ),
    meta: {
      title: "传讯详情",
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/corporate-detail",
    name: "corporate-detail",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/communication/corporate-detail.vue")),
        "corporate-detail"
      ),
    meta: {
      title: "",
      keepAlive: true,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/online-order",
    name: "online-order",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/online-order/index.vue")),
        "online-order"
      ),
    meta: {
      title: "网仓订单",
      keepAlive: true,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/online-order-list",
    name: "online-order-list",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/online-order/list.vue")),
        "online-order-list"
      ),
    meta: {
      title: "历史订单",
      keepAlive: true,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/online-order-create",
    name: "online-order-create",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/online-order/create.vue")),
        "online-order-create"
      ),
    meta: {
      title: "新建订单",
      keepAlive: false,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/online-order-detail",
    name: "online-order-detail",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/online-order/detail.vue")),
        "online-order-detail"
      ),
    meta: {
      title: "订单详情",
      keepAlive: false,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/online-order-info",
    name: "online-order-info",
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
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/online-order-select-store",
    name: "select-store",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/online-order/select-store.vue")),
        "select-store"
      ),
    meta: {
      title: "选择指定门店",
      keepAlive: false,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/online-order-select-member",
    name: "select-member",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/online-order/select-member.vue")),
        "select-member"
      ),
    meta: {
      title: "选择指定导购",
      keepAlive: false,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },
  {
    path: "/online-order-select-date",
    name: "select-date",
    component: (r: (arg0: any) => void) =>
      require.ensure(
        [],
        () => r(require("@/pages/online-order/select-date.vue")),
        "select-date"
      ),
    meta: {
      title: "选择时间",
      keepAlive: false,
      allow_outer_enter: false, // 是否允许外部人访问
      need_wx_sdk: true, // 需要sdk
    },
  },

];
