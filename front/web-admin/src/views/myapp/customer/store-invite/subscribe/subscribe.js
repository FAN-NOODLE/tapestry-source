import list from "./list/list";
import detail from "./detail/detail";

export default {
  icon: "", // 菜单以及权限验证所显示的 icon
  path: "",
  cname: "预约管理", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "subscribe", // 权限验证根据 ename 生成 key （*必填）
  root: "customer", //根目录结点
  description: "营销策略", // 模块描述/备注 选填字段
  type: 2, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  hide: false,
  children: [list, detail]
};
