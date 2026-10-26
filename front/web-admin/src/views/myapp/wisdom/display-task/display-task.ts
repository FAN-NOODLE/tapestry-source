import ListModule from "./list/list";
import DetailModule from "./detail/detail";
import AddModule from "./add/add";
import EditModule from "./edit/edit";
import AuditModule from "./audit/audit";
import AuditDetailModule from "./audit-detail/audit-detail";

/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "VM陈列任务", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "display-task", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "VM陈列任务", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "wisdom", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [ListModule,EditModule,AddModule, DetailModule,AuditModule, AuditDetailModule] //子模块
};
