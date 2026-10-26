import ListModule from "./list/list";
import CreateModule from "./create/create";
import SelectModule from "./select/select";
import ShowModule from "./show/show";
import EditModule from "./edit/edit";

/**
 * @description 模块配置
 */
export default {
  icon: "", // 菜单以及权限验证所显示的 icon
  path: "",
  cname: "客户分组", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "customer-groups", // 权限验证根据 ename 生成 key （*必填）
  root: "customer",
  description: "客户分组模块", // 模块描述/备注 选填字段
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  hide: false,
  children: [ListModule, CreateModule, SelectModule, ShowModule, EditModule]
};
