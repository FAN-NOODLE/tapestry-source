import http from "@/utils/axios";
export default {
  // 待分配人员列表
  allotExternalList(data: any) {
    return http("customer/inherit/allot_external_list", {
      method: "POST",
      data: { ...data },
    });
  },
  // 已分配导购列表
  takeoverList(data: any) {
    return http("customer/inherit/takeover_list", {
      method: "POST",
      data: { ...data },
    });
  },
  // 离职员工列表
  inheritList(data: any) {
    return http("customer/inherit/list", {
      method: "POST",
      data: { ...data },
    });
  },
  // 分配离职成员的客户
  transferCustomer(data: any) {
    return http("customer/inherit/transfer_customer", {
      method: "POST",
      data: { ...data },
    });
  },
  // 分配离职成员的客户群
  transferCustomerGroup(data: any) {
    return http("customer/inherit/transfer_customer_group", {
      method: "POST",
      data: { ...data },
    });
  },
  // 已分配外部联系人列表(已接替)
  takeoverExternalList(data: any) {
    return http("customer/inherit/takeover_external_list", {
      method: "POST",
      data: { ...data },
    });
  },
  // 已分配群列表(已接替)
  takeoverGroupList(data: any) {
    return http("customer/inherit/takeover_group_list", {
      method: "POST",
      data: { ...data },
    });
  },
  // 待分配群列表
  allotGroupList(data: any) {
    return http("customer/inherit/allot_group_list", {
      method: "POST",
      data: { ...data },
    });
  },

  //  部门列表接口
  departmentTree() {
    return http("department/tree-list", {
      method: "POST",
      data: { apiVersion: "old" },
    });
  },
  // 数据更新
  syncCustomerData(data: any) {
    return http("customer/inherit/sync_customer_data", {
      method: "POST",
      data: { ...data },
    });
  },
};
