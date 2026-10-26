import childtaskDataModule from "./child-task-data/child-task-data";
import childTaskDetailModule from "./child-task-detail/child-task-detail";
import childTaskCopyModule from "./child-task-copy/child-task-copy";
import childTaskAddModule from "./child-task-add/child-task-add";
import childTaskEditModule from "./child-task-edit/child-task-edit";
import childlistModule from "./child-list/child-list";
import planListModule from "./plan-list/plan-list";
/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "任务管理", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "plan-task", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "wisdom", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    childtaskDataModule,
    childTaskDetailModule,
    childTaskCopyModule,
    childTaskAddModule,
    childTaskEditModule,
    childlistModule,
    planListModule,
  ], //子模块
};
