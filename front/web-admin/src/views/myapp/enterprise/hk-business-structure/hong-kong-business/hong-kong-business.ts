import storeManagement from "./store-management/store-management";
import permissionManagement from "./permission-management/permission-management";
import permissionAdd from "./permission-add/permission-add";
import permissionEdit from "./permission-edit/permission-edit";
import permissionDetail from "./permission-detail/permission-detail";

export default {
  icon: "",
  path: "",
  cname: "港澳业务架构",
  ename: "hong-kong-business",
  description: "港澳业务架构",
  hide: false,
  root: "enterprise",
  type: 4,
  children: [
    storeManagement,
    permissionManagement,
    permissionAdd,
    permissionEdit,
    permissionDetail
  ]
};
