import AuthSettings from "./authorization-settings/authorization-settings";
import OperationalLog from "./operational-log/operational-log";
import accountsManagerModule from "./accounts/accounts";
import merchantModule from "./merchant-info/merchant-info";
// todo: 下载中心
import downloadCenterModule from "./download-center/download-center";
/**
 * @description 模块配置
 */
export default {
  icon: "setting-icon", // 一级菜单所显示的 icon
  path: "",
  cname: "设置", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "setting", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    AuthSettings,
    OperationalLog,
    accountsManagerModule,
    merchantModule,
    downloadCenterModule
  ] //子模块
};
