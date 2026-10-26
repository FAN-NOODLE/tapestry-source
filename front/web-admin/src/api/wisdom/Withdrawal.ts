import http from "@/utils/axios";

export default {
  // 提现审批列表
  approvalList(data: any) {
    return http(`withdrawal/list`, {
      method: "POST",
      data
    });
  },
  // 提现审批操作
  operationStatus(data: any) {
    return http(`withdrawal/approval-check`, {
      method: "POST",
      data
    });
  },
  // 提现审批导出
  exportApprovalList(data: any) {
    return http(`withdrawal/export_list`, {
      method: "POST",
      data
    });
  },
  // 获取自动审批详情
  getSetting(data?: any) {
    return http(`withdrawal-setting/auto_detail`, {
      method: "POST",
      data
    });
  },
  // 设置自动审批
  saveAutoSetting(data: any) {
    return http(`withdrawal-setting/auto_add`, {
      method: "POST",
      data
    });
  },
  // 获取提现设置
  getSettingDetail(data?: any) {
    return http("withdrawal-setting/detail", {
      method: "POST",
      data
    });
  },

  // 保存提现设置
  saveSetting(data: any) {
    return http("withdrawal-setting/add", {
      method: "POST",
      data
    });
  }
};
