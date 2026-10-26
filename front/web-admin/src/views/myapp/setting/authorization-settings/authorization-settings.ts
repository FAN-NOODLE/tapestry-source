import dlconfigModule from "./dl-config/dl-config";
          import StoreSet from "./store-set/store-set";
import EnterpriseConfigModule from "./enterprise-config/enterprise-config";
import OfficialAccount from "./official-account/official-account";
import WisdomSet from "./wisdom-set/wisdom-set";
import WxPay from "./wx-pay/wx-pay";

/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "授权配置", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "authorization-settings", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "授权配置", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "setting", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
      dlconfigModule,
    StoreSet,
    OfficialAccount,
    EnterpriseConfigModule,
    WisdomSet,
    WxPay
  ] //子模块
};
