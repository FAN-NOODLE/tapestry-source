/**
 * @description 模块配置
 */
export default {
  icon: "", // 菜单以及权限验证所显示的 icon
  cname: "新增/编辑账号", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "accounts-edit", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "新增/编辑账号", // 模块描述
  hide: true, // 设置其在菜单不显示
  root: "setting", //模块所在根结点
  type: 2, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  is_authorization: false // 是否要进行权限验证
};
