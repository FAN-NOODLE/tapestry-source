import changepasswordModule from "./change-password/change-password";
import accountListModule from "./list/list";
import accountEditModule from "./accounts-edit/accounts-edit";
import accountDetailModule from "./accounts-detail/accounts-detail";

/**
 * @description 模块配置
 */
export default {
  icon: "", // 菜单以及权限验证所显示的 icon
  cname: "账号管理", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "accounts-news", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "账号管理", // 模块描述
  is_authorization: true, // 是否要进行权限验证
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  root: "setting", //模块所在根结点
  children: [
    changepasswordModule, accountListModule, accountEditModule, accountDetailModule] // 子模块
};
