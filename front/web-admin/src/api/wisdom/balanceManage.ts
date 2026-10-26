import http from "@/utils/axios";

export default {
  // 余额管理列表
  balanceListListApi(data: any) {
    return http("distribution-balance/list", {
      method: "POST",
      data
    });
  },
  // 余额列表导出
  balanceListLisExporttApi(data: any) {
    return http("distribution-balance/list-export", {
      method: "POST",
      params: data,
      responseType: "blob"
    });
  },
  //   余额明细列表
  balanceDetailApi(data: any) {
    return http("distribution-balance/details", {
      method: "POST",
      data
    });
  },
  //   余额明细导出
  balanceDetailExportApi(data: any) {
    return http("distribution-balance/details-export", {
      method: "POST",
      params: data,
      responseType: "blob"
    });
  }
};
