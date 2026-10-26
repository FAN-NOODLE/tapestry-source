import http from "@/utils/axios";

export default {
  // 结算记录列表
  getSettlementListApi(data: any) {
    return http("/partner-distribution/list", {
      method: "POST",
      data
    });
  },

  // 结算记录明细
  getSettlementDetailApi(data: any) {
    return http("/partner-distribution/detail", {
      method: "POST",
      data
    });
  },

  // 结算列表导出
  exportSettlementListApi(data: any) {
    return http("/partner-distribution/export", {
      method: "POST",
      params: data,
      responseType: "blob"
    });
  }
};
