import DistributionSettingModule from "./distribution-setting/distribution-setting";
import SettlementRecordModule from "./settlement-record/settlement-record";
import WithdrawalApprovalModule from "./withdrawal-approval/withdrawal-approval";
import WithdrawalSettingModule from "./withdrawal-setting/withdrawal-setting";
import BalanceManagementModule from "./balance-management/balance-management";
/**
 * @description 模块配置
 */
export default {
  icon: "", // 一级菜单所显示的 icon
  path: "",
  cname: "全员分销", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "staff-distribution", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "wisdom", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    DistributionSettingModule,
    SettlementRecordModule,
    WithdrawalApprovalModule,
    WithdrawalSettingModule,
    BalanceManagementModule
  ] //子模块
};
