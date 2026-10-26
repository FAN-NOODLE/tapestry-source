import externalchainmanageModule from "./externalchain-manage/externalchain-manage";
import marketStrategy from "./market-strategy/market-strategy";
import insertMarketStrategy from "./insert-market-strategy/insert-market-strategy";
import smsMgt from "./sms-mgt/sms-mgt";

/**
 * @description 模块配置
 */
export default {
  icon: "", // 菜单以及权限验证所显示的 icon
  path: "",
  cname: "自动化营销", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "auto-market", // 权限验证根据 ename 生成 key （*必填）
  root: "customer",
  description: "自动化营销", // 模块描述/备注 选填字段
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  hide: false,
  children: [
    externalchainmanageModule,
    marketStrategy,
    insertMarketStrategy,
    smsMgt
  ]
};
