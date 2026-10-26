import http from "@/utils/axios";
// 获取金币明细记录
export const getGoldRecordList = (data: any) => {
  return http("excitation/coins/record-list", {
    method: "POST",
    data: {...data },
  });
};
// 金币明细记录导出
export const exportGoldRecordList = (data: any) => {
  return http("excitation/coins/download-record", {
    method: "POST",
    data: {...data },
  });
}
// 导购金币配置信息保存
export const saveGoldConfig = (data: any) => {
  return http("excitation/coins/save-coins-settings", {
    method: "POST",
    data: {...data },
  });
}
// 获取导购金币配置信息
export const getGoldConfig = (data: any) => {
  return http("excitation/coins/coins-settings-info", {
    method: "POST",
    data: {...data },
  });
}
// 下载补发金币模板
export const downloadCoinsTemp = () => {
  return http("excitation/download-manual-coins-template", {
    method: "GET",
    responseType: "blob",
    data: {apiVersion: "old" },
  });
}
// 上传补发金币模板
export const importCoinsTemp = (data: any) => {
  return http("excitation/import-manual-coins", {
    method: "POST",
    data: {apiVersion: "old",...data },
  });
}