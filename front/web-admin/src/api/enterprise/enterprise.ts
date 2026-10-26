import http from "@/utils/axios";

export default {
  // 重复客户
  repeatCustomerApi(data: any) {
    return http("customer/repeat-customer", {
      method: "POST",
      data,
    });
  },
  // 获取树形数据
  getTreeList(data) {
    return http("/department/tree-list", {
      method: "POST",
      data: { ...data, apiVersion: "old" },
    });
  },
  // 重复客户员工列表
  repeatFindMemberApi(data: any) {
    return http("customer/repeat-find-member", {
      method: "POST",
      data,
    });
  },
  // 组织列表
  depatListApi(data: any) {
    return http("departments/list", {
      method: "POST",
      data,
    });
  },
  // 员工列表
  memberListApi(data: any) {
    return http("members/list", {
      method: "POST",
      data,
    });
  },
  // 员工详情
  memberDetailApi(data: any) {
    return http("members/detail", {
      method: "POST",
      data,
    });
  },
  // 外部联系人secret更新
  authContacUpdateApi(data: any) {
    return http("auth-contact-update", {
      method: "POST",
      data,
    });
  },
  // 企微配置
  settingsApi() {
    return http("settings", {
      method: "POST",
    });
  },
  //自建应用更新
  workPluginUpdateApi(data: any) {
    return http("work-plugin-update", {
      method: "POST",
      data,
    });
  },
  // 获取会话存档配置
  onGetChatSetting() {
    return http("chat-settings", {
      method: "POST",
    });
  },
  onUpdateChatSetting(data: any) {
    return http("session-archiving-update", {
      method: "POST",
      data,
    });
  },
  // 添加/更新会话存档配置
  //上传校验文件
  uploadCheckFileApi(data: any) {
    const form_data = new FormData();
    form_data.append("file", data.file);
    return http("upload-check-file", {
      method: "POST",
      data: form_data,
    });
  },
  // 获取组织更新员工信息
  updateRecordApi() {
    return http("members/sync", {
      method: "POST",
    });
  },
  // 获取组织更新员工信息
  updateRocordApi() {
    return http("departments/update-record", {
      method: "POST",
    });
  },
  // 消息配置列表更新
  updateMessageApi(data: any) {
    return http("customer/update-message", {
      method: "POST",
      data,
    });
  },
  // 消息配置列表
  messageListApi(data: any) {
    return http("customer/message-list", {
      method: "POST",
      data,
    });
  },
  // 流失客户列表
  losscustomersListApi(data: any) {
    return http("/customer/loss-customers-list", {
      method: "POST",
      data,
    });
  },
  // 门店数据
  storeDataApi(data: any) {
    return http("/department/list", {
      method: "GET",
      data: { apiVersion: "old" },
      params: data,
    });
  },

  // 导购列表
  shoppingGuideListApi(data: any) {
    return http("/wisdom-guide/list", {
      method: "POST",
      data: { ...data },
    });
  },

  // 获取组织架构配置的应用消息配置人员
  departmentsUserLostConfig(data: any) {
    return http("departments/user-lost-info", {
      method: "POST",
      data,
    });
  },

  // 导购丢失客户推送消息配置
  userLostConfig(data: any) {
    return http("departments/user-lost-config", {
      method: "POST",
      data,
    });
  },
};
