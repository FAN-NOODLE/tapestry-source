/**
 * @description 模块配置
 */
import list from "./list/list";
import groupData from "./group-data/group-data";
import groupDetail from "./group-detail/group-detail";

export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "社群", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "customer-group", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "客户社群", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "enterprise", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [list, groupData, groupDetail]
};
