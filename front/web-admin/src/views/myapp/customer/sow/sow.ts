import list from './list/list'
import dataList from './data-list/data-list'
import edit from './edit/edit'
import add from './add/add'
import info from './info/info'
/**
 * @description 模块配置
 */
export default {
  icon: "", // 菜单以及权限验证所显示的 icon
  path: "",
  cname: "播种管理", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "sow", // 权限验证根据 ename 生成 key （*必填）
  root: "customer",
  description: "播种管理", // 模块描述/备注 选填字段
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  hide: false,
  children: [
    list,
    dataList,
    add,
    edit,
    info
  ]
};
