import http from "@/utils/axios";

/**
 * 根据活动码获取导购二维码或者门店列表
 * @param data 
 * @returns 
 */
export const getQrcodeByUnionidApi = (data: any) => {
  return http("activity-qrcode/get-qrcode-by-unionid", {
    method: "POST",
    data,
  });
}

/**
 * 手动获取门店列表
 * @param data 
 * @returns 
 */
export const getAllStoreListApi = (data: any) => {
  return http("activity-qrcode/get-all-store-list", {
    method: "POST",
    data,
  });
}

/**
 * 获取活动详情
 * @param data 
 * @returns 
 */
export const getActivityDetailApi = (data: any) => {
  return http("activity-qrcode/activity-detail", {
    method: "POST",
    data,
  });
}