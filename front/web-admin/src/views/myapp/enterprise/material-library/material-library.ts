import materialList from "./material-list/material-list";
import productCategory from "./product-category/product-category";
import createFriend from "./create-friend/create-friend";
import createImage from "./create-image/create-image";
import createLink from "./create-link/create-link";
import createVideo from "./create-video/create-video";

/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "素材库", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "material-library", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "素材库模块", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "enterprise", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [materialList, productCategory,createFriend, createImage, createLink, createVideo]
};
