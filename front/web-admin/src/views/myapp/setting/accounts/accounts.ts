import roleModule from "./roles/roles";
import accountsModule from "./accounts-news/accounts-news";
/**
 * @description 模块配置
 */
export default {
  icon: "", // 菜单以及权限验证所显示的 icon
  cname: "账号权限", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "accounts", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "账号权限", // 模块描述
  is_authorization: true, // 是否要进行权限验证
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  root: "setting",
  children: [roleModule, accountsModule] // 子模块
};
