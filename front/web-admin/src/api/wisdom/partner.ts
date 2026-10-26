import http from "@/utils/axios";

export default {
  // 合伙人列表
  partnerListApi(data: any) {
    return http("/partner/list", {
      method: "POST",
      data
    });
  },
  // 合伙人列表导出
  partnerListExportApi(data: any) {
    return http("/partner/export_list", {
      method: "POST",
      params: data,
      responseType: "blob"
    });
  },
  // 合伙人详情列表
  partnerDetailApi(data: any) {
    return http("/partner/detail", {
      method: "POST",
      data
    });
  },
  // 合伙人详情导出
  partnerDetailExportApi(data: any) {
    return http("/partner/export_detail", {
      method: "POST",
      params: data,
      responseType: "blob"
    });
  },
  // 合伙人启用禁用开关
  partnerSwitchApi(data: any) {
    return http("/partner/switch", {
      method: "POST",
      data
    });
  }
};
