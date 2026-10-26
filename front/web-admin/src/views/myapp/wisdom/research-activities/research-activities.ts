import signinListModule from "./signIn-list/signIn-list";
import signinDetailModule from "./signIn-detail/signIn-detail";
import researchListModule from "./research-list/research-list";
import eventsdatailModule from "./events-detail/events-detail";
import eventslistModule from "./events-list/events-list";
/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "调研活动", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "research-activities", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "wisdom", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    signinListModule,
    signinDetailModule,
    researchListModule,
    eventsdatailModule,
    eventslistModule,
  ], //子模块
};
