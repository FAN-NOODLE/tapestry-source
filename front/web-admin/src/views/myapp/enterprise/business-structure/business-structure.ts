import unassignedListModule from "./unassigned-list/unassigned-list";
import list from "./list/list";
import organizeDetail from "./staff-detail/staff-detail";
import storeList from "./store-list/store-list";
import syncList from "./sync-list/sync-list";
import orgStructure from "./org-structure/org-structure";
import permission from "./permission/permission";
/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "业务架构", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "business-structure", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "业务架构模块", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "enterprise", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    unassignedListModule, 
    list, 
    storeList, 
    organizeDetail, 
    syncList, 
    orgStructure, 
    permission
  ]
};
