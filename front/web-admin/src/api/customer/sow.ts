import http from "@/utils/axios";

export default {
  // 播种模版列表
  getSowListApi(data: any) {
    return http("customer/sow-template/list", {
      method: "POST",
      data
    });
  },
  // 播种模版详情
  getSowDetailApi(data: any) {
    return http("customer/sow-template/detail", {
      method: "POST",
      data
    });
  },
  // 播种模版记录列表
  getSowRecordListApi(data: any) {
    return http("customer/sow-template/all-task-list", {
      method: "POST",
      data
    });
  },
  // 播种模版创建/更新
  saveSowApi(data: any) {
    return http("customer/sow-template/save", {
      method: "POST",
      data
    });
  },
  // 播种模版创建/更新
  updateSowStatusApi(data: any) {
    return http("customer/sow-template/update-status", {
      method: "POST",
      data
    });
  },
  // 播种模版删除
  removeSowApi(data: any) {
    return http("customer/sow-template/delete", {
      method: "POST",
      data
    });
  },
  // 播种模版到处
  exportSowListApi(data: any) {
    return http("customer/sow-template/template-list-export", {
      method: "POST",
      responseType: "blob",
      data
    });
  },
  // 播种模版到处
  exportSowRecordListApi(data: any) {
    return http("customer/sow-template/download-sow-task-record", {
      method: "POST",
      data
    });
  },

}