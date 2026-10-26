import ContentModule from './content/content'
import WhitelistModule from './whitelist/whitelist'
import IntegrationModule from './integration/integration'
import AuditingModule from './auditing/auditing'
/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "传讯管理", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "communication", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "传讯管理", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "wisdom", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    ContentModule,
    AuditingModule,
    IntegrationModule,
    WhitelistModule,
  ] //子模块
};