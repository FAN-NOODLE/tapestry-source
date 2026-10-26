import ListModule from "./list/list";
import DetailModule from "./detail/detail";
import EditModule from "./edit/edit";
import AddModule from "./add/add";
/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "传讯内容库", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "content", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "传讯内容库", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "communication", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    ListModule,
    DetailModule,
    EditModule,
    AddModule
  ] //子模块
};

/**
 * TODO 1. 提交有效时间控制 新增/编辑
 * TODO 2. 复制传讯 编辑
 * TODO 3. 所属部门调整
 * TODO 4. 撤回接口
 */