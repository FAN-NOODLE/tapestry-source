import servicesettingModule from "./service-setting/service-setting";
import datalistModule from "./data-list/data-list";
import commentlistModule from "./comment-list/comment-list";
import subscribeModule from "./subscribe/subscribe";
/**
 * @description 模块配置
 */
export default {
  icon: "", // 菜单以及权限验证所显示的 icon
  path: "",
  cname: "到店邀约", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "store-invite", // 权限验证根据 ename 生成 key （*必填）
  root: "customer",
  description: "到店邀约", // 模块描述/备注 选填字段
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  hide: false,
  children: [
    subscribeModule,
    servicesettingModule,
    datalistModule,
    commentlistModule
  ]
};
