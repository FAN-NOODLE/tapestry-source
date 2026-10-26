import http from "@/utils/axios";

export default {
  // 【渠道】列表
  channelsListApi(data: any) {
    return http("drainage-contact/channels-list", {
      method: "POST",
      data: data
    });
  },
  // 【渠道】新增渠道
  createChannelsApi(data: any) {
    return http("drainage-contact/channels-create", {
      method: "POST",
      data: data
    });
  },
  // 【渠道】详情
  channelsDetailApi(data: any) {
    return http("drainage-contact/channels-detail", {
      method: "POST",
      data: data
    });
  },
  // 【渠道】编辑渠道
  channelsEditApi(data: any) {
    return http("drainage-contact/channels-edit", {
      method: "POST",
      data: data
    });
  },
  // 【联系我】联系我新增
  contactCreateApi(data: any) {
    return http("drainage/contact/create", {
      method: "POST",
      data: data
    });
  },
  // 【联系我】编辑
  contactUpdateApi(data: any) {
    return http("drainage/contact/update", {
      method: "POST",
      data: data
    });
  },
  // 【联系我】详情
  contactDetailsApi(data: any) {
    return http("drainage/contact/details", {
      method: "POST",
      data: data
    });
  },
  // 【联系我】删除
  contactDeleteApi(data: any) {
    return http("drainage/contact/delete", {
      method: "POST",
      data: data
    });
  },
  // 【联系我】列表
  contactListApi(data: any) {
    return http("drainage/contact/list", {
      method: "POST",
      data: data
    });
  },
  // 【联系我】返回不同大小二维码文件流
  contactDownloadApi(data: any) {
    return http("drainage/contact/download-qr-code", {
      method: "POST",
      data: data
    });
  },
  // 【客户还迎语】默认
  welcomeDefaultListApi(data: any) {
    return http("drainage/welcome-default", {
      method: "POST",
      data: data
    });
  },
  // 【客户还迎语】新建
  welcomeCreateApi(data: any) {
    return http("drainage/welcome-create", {
      method: "POST",
      data: data
    });
  },
  // 【客户还迎语】删除
  welcomeDeleteApi(data: any) {
    return http("drainage/welcome-delete", {
      method: "POST",
      data: data
    });
  },
  // 【客户还迎语】详情
  welcomeDetailApi(data: any) {
    return http("drainage/welcome-detail", {
      method: "POST",
      data: data
    });
  },
  // 【客户还迎语】判断所选员工是否已配置还迎语
  welcomeValidateApi(data: any) {
    return http("drainage/welcome-validate", {
      method: "POST",
      data: data
    });
  },
  // 【客户还迎语】编辑
  welcomeEditApi(data: any) {
    return http("drainage/welcome-edit", {
      method: "POST",
      data: data
    });
  },
  // 【客户还迎语】list
  welcomeListApi(data: any) {
    return http("drainage/welcome-list", {
      method: "POST",
      data: data
    });
  },
  // 渠道统计列表
  channelStatisticsListApi(data: any) {
    return http("drainage-contact/channels-data", {
      method: "POST",
      data: data
    });
  },
  // 渠道统计excel导出
  channelStatisticsCountExportApi(data: any) {
    return http("drainage-contact/channels-count-export", {
      method: "POST",
      data: data,
      responseType: "arraybuffer"
    });
  },
  // 员工选择渠道列表
  welcomeChannelListApi(data: any) {
    return http("drainage/welcome-channel-list", {
      method: "POST",
      data: data
    });
  },
  //渠道数据统计总览
  channelsDataSumApi(data: any) {
    return http("drainage-contact/channels-data-sum", {
      method: "POST",
      data: data
    });
  },
  // 个性化欢迎语列表
  drainageWelcomeAll(data: any) {
    return http("drainage/welcome-all", {
      method: "POST",
      data
    });
  },
  // 删除个性化欢迎语
  drainageDeleteApi(data: any) {
    return http("drainage/welcome-del", {
      method: "POST",
      data
    });
  },
  // 卡等列表
  drainageGetDetail(data: any) {
    return http("drainage/get-detail", {
      method: "POST",
      data
    });
  },
  // 保存个性化欢迎语
  drainageSaveApi(data: any) {
    return http("drainage/welcome-new-edit", {
      method: "POST",
      data
    });
  },
  // 个性化欢迎语详情
  drainageDetailApi(data: any) {
    return http("drainage/welcome-detailnew", {
      method: "POST",
      data
    });
  },
  // 个性欢迎语 drainage/all-type
  getAllTypeListApi(data: any) {
    return http("drainage/all-type", {
      method: "POST",
      data
    });
  },
  // 活动欢迎语
  getActivityWelcomeListApi(data: any) {
    return http("activity-qrcode/welcome-list", {
      method: "POST",
      data
    });
  },
  // 活动欢迎语
  getActivityWelcomeDetailApi(data: any) {
    return http("activity-qrcode/welcome-detail", {
      method: "POST",
      data
    });
  },
  // 活动欢迎语
  saveActivityWelcomeApi(data: any) {
    return http("activity-qrcode/welcome-save", {
      method: "POST",
      data
    });
  },
  // 活动欢迎语
  deleteActivityWelcomeApi(data: any) {
    return http("activity-qrcode/welcome-delete", {
      method: "POST",
      data
    });
  },
};
