import organizeStructure from "./organize-structure/organize-structure";
import enterpriseConfig from "./enterprise-config/enterprise-config";
import friendCircle from "./friend-circle/friend-circle";
import statistics from "./statistics/statistics";
import materialIbrary from "./material-library/material-library";
import drainage from "./drainage/drainage";
import association from "./association/association";
import businessStructure from "./business-structure/business-structure";
import launch from "./launch/launch";
import hkBusinessStructure from "./hk-business-structure/hk-business-structure";
/**
 * @description 模块配置
 */
export default {
  icon: "enterprise-icon", // 一级菜单所显示的 icon
  path: "",
  cname: "企微", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "enterprise", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    statistics,
    organizeStructure,
    enterpriseConfig,
    friendCircle,
    drainage,
    materialIbrary,
    association,
    businessStructure,
    hkBusinessStructure,
    launch
  ] //子模块
};
