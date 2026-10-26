import http from "@/utils/axios";

export default {
  // 邀请函列表
  invitationLetterList(data: any) {
    return http("invite/template/list", {
      method: "POST",
      data
    });
  },

  /** 关联预约服务列表 */
  associateServiceList(data: any) {
    return http("reserve/service/list", {
      method: "POST",
      data
    });
  },

  /** 邀请数据列表 */
  invitationDataList(data: any) {
    return http("invite/record-order-list", {
      method: "POST",
      data
    });
  },

  /** 邀请数据导出 */
  invitationDataExport(data: any) {
    return http("invite/record-order-list-export", {
      method: "POST",
      data: data,
      responseType: "arraybuffer"
    });
  },

  /** 删除 */
  invitationDataDelete(data: any) {
    return http("invite/template/delete", {
      method: "POST",
      data
    });
  },

  /** 详情 */
  invitationDataInfo(data: any) {
    return http("invite/template/info", {
      method: "POST",
      data
    });
  },

  /** 新增 */
  invitationDataAdd(data: any) {
    return http("invite/template/save", {
      method: "POST",
      data
    });
  },

  /** 服务详情 */
  serviceInfo(data: any) {
    return http("reserve/service/service-info", {
      method: "POST",
      data
    });
  },

  /** 预览 */
  invitationPreview(data: any) {
    return http("invite/create-qrcode", {
      method: "POST",
      data
    });
  }
};
