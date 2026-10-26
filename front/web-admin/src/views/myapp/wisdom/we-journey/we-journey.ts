import indexModule from "./index/index";
import gradeSetListModule from "./grade-set-list/grade-set-list";
import gradeStatisticsModule from "./grade-statistics/grade-statistics";
import integralExchangeModule from "./integral-exchange/integral-exchange";
import integralExplainModule from "./integral-explain/integral-explain";
import integralSetAddModule from "./integral-set-add/integral-set-add";
import integralSetEditModule from "./integral-set-edit/integral-set-edit";
import integralSetListModule from "./integral-set-list/integral-set-list";
import integralSetLookModule from "./integral-set-look/integral-set-look";
import medalSetAdd from "./medal-set-add/medal-set-add";
import medalSetEdit from "./medal-set-edit/medal-set-edit";
import medalSetList from "./medal-set-list/medal-set-list";
import medalSetLook from "./medal-set-look/medal-set-look";
import medalStatistics from "./medal-statistics/medal-statistics";
import coinsIndex from "./coins-index/coins-index";
/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "We Journey", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "we-journey", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "wisdom", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    indexModule,
    gradeSetListModule,
    gradeStatisticsModule,
    integralExchangeModule,
    integralExplainModule,
    integralSetAddModule,
    integralSetEditModule,
    integralSetListModule,
    integralSetLookModule,
    medalSetAdd,
    medalSetEdit,
    medalSetList,
    medalSetLook,
    medalStatistics,
    coinsIndex
  ] //子模块
};
