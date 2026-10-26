import http from "@/utils/axios";

export default {
  // 企微配置
  settingsApi(data?: any) {
    return http("settings", {
      method: "POST",
      data,
    });
  },

  // 获取会话存档配置
  onGetChatSetting(data?: any) {
    return http("chat-settings", {
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

  onUpdateChatSetting(data: any) {
    return http("session-archiving-update", {
      method: "POST",
      data,
    });
  },

  //自建应用更新
  workPluginUpdateApi(data: any) {
    return http("work-plugin-update", {
      method: "POST",
      data: data,
    });
  },

  // DL配置详情
  getDLConfig(data: any) {
    return http("data-lake/get-config", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },

  // 编辑DL配置
  editDLConfig(data: any) {
    return http("data-lake/edit-config", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },

  // 修改token
  editApiToken(data: any) {
    return http("set-wxqy-token", {
      method: "POST",
      data,
    });
  },
};
