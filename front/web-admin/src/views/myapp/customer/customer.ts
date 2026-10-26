import customerLabel from "./customer-label/customer-label";
import customerGroups from "./customer-groups/customer-groups";
import customerList from "./customer-list/customer-list";
import autoMarket from "./auto-market/auto-market";
import greetingCard from "./greeting-card/greeting-card";
import storeInvite from "./store-invite/store-invite";
import invitationLetter from "./invitation-letter/invitation-letter";
import sow from "./sow/sow";
import onlineOrder from "./online-order/online-order";

/**
 * @description 模块配置
 */
export default {
  icon: "customer-icon", // 一级菜单所显示的 icon
  path: "",
  cname: "客户", // 菜单以及权限验证所显示的中文名 （*必填）
  ename: "customer", // 权限验证所使用的 key 值，同时还可当作 name 在路由跳转中使用 （*必填）
  description: "", // 模块描述备注
  hide: false, //是否在菜单隐藏
  root: "", //模块所在根结点
  type: 1, // 1 => 菜单 2 => 视图 4 => 菜单+视图
  children: [
    customerLabel,
    customerGroups,
    customerList,
    autoMarket,
    greetingCard,
    storeInvite,
    invitationLetter,
    sow,
    onlineOrder
  ] //子模块
};
