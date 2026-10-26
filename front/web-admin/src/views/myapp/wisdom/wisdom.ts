import marketingactivitiesModule from "./marketing-activities/marketing-activities";
import researchactivitiesModule from "./research-activities/research-activities";
import taskcenterModule from "./task-center/task-center";
import subdistributionModule from "./sub-distribution/sub-distribution";
import DataCompassModule from "./data-compass/data-compass";
import GuideManageModule from "./guide-manage/guide-manage";
import StaffDistributionModule from "./staff-distribution/staff-distribution";
import MarketingCalendarModule from "./marketing-calendar/marketing-calendar";
import QuitInheritModule from "./quit-inherit/quit-inherit";
import weJourneyModule from "./we-journey/we-journey";
import DisplayTaskModule from "./display-task/display-task";
import CommunicationModule from "./communication/communication";
import StoreVisitModule from "./store-visit/store-visit";

/**
 * @description 模块配置
 */
export default {
  icon: "wisdom-icon", // 一级菜单所显示的 icon
  path: "",
  cname: "智慧导购", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "wisdom", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    marketingactivitiesModule,
    researchactivitiesModule,
    taskcenterModule,
    subdistributionModule,
    DataCompassModule,
    MarketingCalendarModule,
    GuideManageModule,
    StaffDistributionModule,
    QuitInheritModule,
    weJourneyModule,
    DisplayTaskModule,
    CommunicationModule,
    StoreVisitModule
  ] //子模块
};
