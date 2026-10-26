/**
 * @description 模块配置
 */
import EditModule from "./edit/edit";
// import DetailModule from "./detail/detail";
export default {
  icon: "", // 菜单以及权限验证所显示的 icon
  cname: "商户信息", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "merchant-info", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "商户信息", // 模块描述
  root: "setting", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  is_authorization: true, // 是否要进行权限验证
  children: [EditModule]
};
