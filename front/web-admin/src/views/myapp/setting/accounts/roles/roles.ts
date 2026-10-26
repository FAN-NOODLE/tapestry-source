import roleListModule from "./list/list";
import roleEditModule from "./role-edit/role-edit";
import roleDetailModule from "./role-detail/role-detail";

/**
 * @description 模块配置
 */
export default {
  icon: "", // 菜单以及权限验证所显示的 icon
  cname: "角色管理", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "roles", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "角色管理", // 模块描述
  is_authorization: true, // 是否要进行权限验证
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  root: "setting", //模块所在根结点
  children: [roleListModule, roleEditModule, roleDetailModule] // 子模块
};
