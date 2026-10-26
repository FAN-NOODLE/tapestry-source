import ListModule from "./list/list"
import AddMarketingPlan from "./add-marketing-plan/add-marketing-plan"
import editMarketingPlan from "./edit-marketing-plan/edit-marketing-plan"
import copyMarketingPlan from "./copy-marketing-plan/copy-marketing-plan"
import MarketingPlanDetail from "./marketing-plan-detail/marketing-plan-detail"

/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "营销日历", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "marketing-calendar", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "wisdom", //模块所在根结点
  type: 4, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    ListModule, 
    AddMarketingPlan,
    editMarketingPlan,
    copyMarketingPlan,
    MarketingPlanDetail
  ] //子模块
};