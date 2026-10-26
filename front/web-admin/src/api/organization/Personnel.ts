import http from "@/utils/axios";

export default {
  // 部门列表
  depListApi(data: any) {
    return http("department/list", {
      method: "GET",
      params: data,
    });
  },
  // 部门列表
  depListTwoApi(data: any) {
    return http("department/list", {
      method: "GET",
      params: data,
    });
  },
  // 部门分页列表
  depPageList(data: any) {
    return http("department/page-list", {
      method: "GET",
      params: data,
    });
  },
  // 部门排序
  depListSort(data: any) {
    return http("department/sort", {
      method: "POST",
      params: data,
    });
  },
  // 获取部门详情
  depDetailApi(data: any) {
    return http("department/detail", {
      method: "GET",
      params: data,
    });
  },
  // 部门新增
  depAddApi(data: any) {
    return http("department/add", {
      method: "POST",
      data: data,
    });
  },
  // 部门编辑
  depEditApi(data: any) {
    return http("department/update", {
      method: "POST",
      data: data,
    });
  },
  // 删除部门
  depDelete(data: any) {
    return http("department/delete", {
      method: "POST",
      data: data,
    });
  },
  // 获取员工列表
  saffListApi(data: any) {
    return http("user/list", {
      method: "GET",
      params: data,
    });
  },
  // 离职员工 user/leave
  saffLeaveApi(data: any) {
    return http("user/leave", {
      method: "POST",
      data: data,
    });
  },
  // 新增员工
  saffAddApi(data: any) {
    return http("user/add", {
      method: "POST",
      data: data,
    });
  },
  // 编辑员工
  saffEditApi(data: any) {
    return http("user/update", {
      method: "POST",
      data: data,
    });
  },
  // 员工详情 user/detail
  saffDetailApi(data: any) {
    return http("user/detail", {
      method: "GET",
      params: data,
    });
  },
  // 待分配客户列表
  handoverListApi(data: any) {
    return http("user/handover-list", {
      method: "GET",
      params: data,
    });
  },
  // 分配客户
  handoverCustomerApi(data: any) {
    return http("user/handover", {
      method: "POST",
      data: data,
    });
  },
  // 分配列表
  getGuideListApi(data: any) {
    return http("user/guide_list", {
      method: "POST",
      data,
    });
  },
  // 更改企业分配客户状态
  updateEnterpriseStatusApi(data: any) {
    return http("user/change_enterprise_allocation_status", {
      method: "POST",
      data,
    });
  },
  // 更改员工分配客户状态
  updateAllocationStatusApi(data: any) {
    return http("user/change_allocation_status", {
      method: "POST",
      data,
    });
  },
  // 待添加客户列表
  getAdditionListApi(data: any) {
    return http("customer/addition_list", {
      method: "GET",
      params: data,
    });
  },
  // 标记是否复制手机号
  setMarkUserApi(data: any) {
    return http("customer/mark_user", {
      method: "GET",
      params: data,
    });
  },
  // 导购客户列表
  getGuideUserListApi(data: any) {
    return http("user/guide_users", {
      method: "POST",
      data,
    });
  },
  // 转交导购列表
  getOtherGuidesListApi(data: any) {
    return http("user/other_guides", {
      method: "POST",
      data,
    });
  },
  // 转交客户
  setCustomerTransferApi(data: any) {
    return http("user/customer_transfer", {
      method: "POST",
      data,
    });
  },
  // 联系我二维码列表
  getContactListApi(data: any) {
    return http("contact/list", {
      method: "GET",
      params: data,
    });
  },
  // 联系我二维码列表(下拉框选择)
  getContactSelectionListApi(data: any) {
    return http("contact/selection_list", {
      method: "GET",
      params: data,
    });
  },
  // 删除二维码
  delContactApi(data: any) {
    return http("contact/delete_qr_code", {
      method: "POST",
      data,
    });
  },
  // 新增二维码
  addContactApi(data: any) {
    return http("contact/create_qr_code", {
      method: "POST",
      data,
    });
  },
  // 编辑二维码
  editContactApi(data: any) {
    return http("contact/update_qr_code", {
      method: "POST",
      data,
    });
  },
  // 二维码详情
  getContactDateilApi(data: any) {
    return http("contact/select_qr_code", {
      method: "POST",
      data,
    });
  },
  // 选择员工列表
  getContactTreeListApi(data: any) {
    return http("contact/member_tree", {
      method: "POST",
      data,
    });
  },
  // 欢迎语列表
  getWelcomeSpeechListApi(data: any) {
    return http("welcome_speech/list", {
      method: "GET",
      params: data,
    });
  },
  // 欢迎语删除
  delWelcomeSpeechApi(data: any) {
    return http("welcome_speech/delete", {
      method: "POST",
      data,
    });
  },
  // 新增欢迎语
  addWelcomeSpeechApi(data: any) {
    return http("welcome_speech/add", {
      method: "POST",
      data,
    });
  },
  // 编辑欢迎语
  editWelcomeSpeechApi(data: any) {
    return http("welcome_speech/edit", {
      method: "POST",
      data,
    });
  },
  // 欢迎语详情
  getWelcomeSpeechDetailApi(data: any) {
    return http("welcome_speech/detail", {
      method: "GET",
      params: data,
    });
  },
  // 上下架红包
  setRedPacketSaleApi(data: any) {
    return http("red_packet/sale", {
      method: "POST",
      data,
    });
  },
  // 删除红包
  delRedPacketApi(data: any) {
    return http("red_packet/delete", {
      method: "POST",
      data,
    });
  },
  // 红包发送人列表
  getRedPacketMembersListApi(data: any) {
    return http("red_packet/members", {
      method: "POST",
      data,
    });
  },
  // 新增红包
  addRedPacketApi(data: any) {
    return http("red_packet/create", {
      method: "POST",
      data,
    });
  },
  // 红包详情
  getRedPacketDetailApi(data: any) {
    return http("red_packet/detail", {
      method: "POST",
      data,
    });
  },
  // 编辑红包
  editRedPacketApi(data: any) {
    return http("red_packet/update", {
      method: "POST",
      data,
    });
  },
  // 对账列表
  getReconciliationListApi(data: any) {
    return http("red_packet/reconciliation_list", {
      method: "GET",
      params: data,
    });
  },
  // 红包列表
  getRedPacketListApi(data: any) {
    return http("red_packet/list", {
      method: "GET",
      params: data,
    });
  },
  // 上次同步时间
  syncLastTime(data: any) {
    return http("contact/sync_last_time", {
      method: "GET",
      params: data,
    });
  },
  // 同步组织结构
  syncOrg(data: any) {
    return http("contact/sync_org", {
      method: "POST",
      data: { apiVersion: "old", ...data },
    });
  },
  // 获取同步状态
  syncStatus(data: any) {
    return http("contact/sync_status", {
      method: "GET",
      params: data,
    });
  },
  // 获取同步状态
  TreePeople(data: any) {
    return http("user/tree-list", {
      method: "POST",
      data: data,
    });
  },
  // 组织管理页面 获取部门下员工接口
  GuidePeople(data: any) {
    return http("user/guide-v-user-list", {
      method: "POST",
      data: data,
    });
  },
  // 批量添加导购
  GuideAdd(data: any) {
    return http("guide-manage/add", {
      method: "POST",
      data: data,
    });
  },
  // 删除导购
  GuideDelete(data: any) {
    return http("guide-manage/delete", {
      method: "POST",
      data: data,
    });
  },
  // 获取联系我明细
  ContactMeDetail(data: any) {
    return http("contact/statistics-detail", {
      method: "POST",
      data: data,
    });
  },
  // 查询是否有通讯录
  checkMember(data: any) {
    return http("enterprise/check-member", {
      method: "POST",
      data,
    });
  },
};
