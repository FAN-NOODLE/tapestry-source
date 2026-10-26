import friendCircleList from "./friend-circle-list/friend-circle-list";
import friendCircleDetail from "./friend-circle-detail/friend-circle-detail";
import friendCircleEdit from "./friend-circle-edit/friend-circle-edit";
import friendCircleSetting from "./friend-circle-setting/friend-circle-setting";
/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "朋友圈", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "friend-circle", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "朋友圈模块", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "enterprise", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    friendCircleList,
    friendCircleDetail,
    friendCircleEdit,
    friendCircleSetting
  ]
};
