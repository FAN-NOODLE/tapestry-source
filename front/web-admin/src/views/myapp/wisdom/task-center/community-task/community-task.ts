import taskdetailModule from "./task-detail/task-detail";
import taskeditModule from "./task-edit/task-edit";
import taskcopyModule from "./task-copy/task-copy";
import listModule from "./list/list";
import taskaddModule from "./task-add/task-add";
/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "社群任务", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "community-task", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "wisdom", //模块所在根结点
  type: 2, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    taskdetailModule,
    taskeditModule,
    taskcopyModule,
    listModule,
    taskaddModule,
  ], //子模块
};
