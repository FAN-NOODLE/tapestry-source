import http from "@/utils/axios";

export default {
  /** 白名单 - 列表 */
  getWhitelistApi(data: any) {
    return http("notification/crowd/list", {
      method: "POST",
      data
    });
  },
  /** 白名单 - 新增/编辑 */
  saveWhitelistApi(data: any) {
    return http("notification/crowd/save", {
      method: "POST",
      data
    });
  },
  /** 白名单 - 删除 */
  removeWhitelistApi(data: any) {
    return http("notification/crowd/delete", {
      method: "POST",
      data
    });
  },
  /** 白名单 - 详情 */
  getWhitelistDetailApi(data: any) {
    return http("notification/crowd/detail", {
      method: "POST",
      data
    });
  },
  /** 白名单 - 导入 */
  importWhitelistApi(data: any) {
    return http("notification/import-crowd", {
      method: "POST",
      data: {apiVersion: "old",...data },
    });
  },
  /** ============= 传讯内容库 =============== */
  /** 传讯内容库 - 列表 */
  getContentListApi(data: any) {
    return http("notification/list", {
      method: "POST",
      data
    });
  },
  /** 传讯内容库 - 新增/编辑 */
  saveContentApi(data: any) {
    return http("notification/save", {
      method: "POST",
      data
    });
  },
  /** 传讯内容库 - 撤销 */
  revokeContentApi(data: any) {
    return http("notification/recall", {
      method: "POST",
      data
    });
  },
  /** 传讯内容库 - 详情 */
  getContentDetailApi(data: any) {
    return http("notification/info", {
      method: "POST",
      data
    });
  },
  /** 传讯内容库 - 复制 */
  copyContentApi(data: any) {
    return http("notification/copy", {
      method: "POST",
      data
    });
  },
  /** 传讯内容库 - 作废 */
  cancelContentApi(data: any) {
    return http("notification/cancel", {
      method: "POST",
      data
    });
  },
  /** 传讯内容库 - 审核 */
  auditContentApi(data: any) {
    return http("notification/check", {
      method: "POST",
      data
    });
  },
  /** 获取所属部门列表 */
  getDepartmentListApi(data: any) {
    return http("notification/departments", {
      method: "POST",
      data
    });
  },

  /** ============= 企业传讯 =============== */
  /** 企业传讯 - 列表 */
  getMergeListApi(data: any) {
    return http("notification/merge/list", {
      method: "POST",
      data
    });
  },
  /** 企业传讯 - 新增/编辑 */
  saveMergeApi(data: any) {
    return http("notification/merge/save", {
      method: "POST",
      data
    });
  },
  /** 企业传讯 - 详情 */
  getMergeDetailApi(data: any) {
    return http("notification/merge/detail", {
      method: "POST",
      data
    });
  },
  /** 企业传讯 - 发布 */
  releaseMergeApi(data: any) {
    return http("notification/merge/release", {
      method: "POST",
      data
    });
  },
  /** 企业传讯 - 获取提醒配置 */
  getNoticeConfigApi(data: any) {
    return http("notification/notice-config/detail", {
      method: "POST",
      data
    });
  },
  /** 企业传讯 - 更新提醒配置 */
  updateNoticeConfigApi(data: any) {
    return http("notification/notice-config/save", {
      method: "POST",
      data
    });
  },
  /** 企业传讯 - 保存pdf */
  saveMergePdfApi(data: any) {
    return http("notification/merge/save-pdf", {
      method: "POST",
      data
    });
  },
  
  
}