/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "积分规则说明", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "integral-explain", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "", // 模块描述备注
  hide: true, //是否在菜单隐藏
  root: "wisdom", //模块所在根结点
  type: 2, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [] //子模块
};
