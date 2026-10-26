/**
 * @description 模块配置
 */
import newcreateModule from "./newcreate/newcreate";
import neweditModule from "./newedit/newedit";
import list from "./list/list";
import create from "./create/create";
import edit from "./edit/edit";
import channelcreate from "./channel-create/channel-create";
import channeledit from "./channel-edit/channel-edit";
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "客户欢迎语", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "welcome", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "客户欢迎语", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "enterprise", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    newcreateModule,
    neweditModule,
    list,
    create,
    edit,
    channelcreate,
    channeledit
  ]
};
