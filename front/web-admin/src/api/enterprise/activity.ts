
import http from "@/utils/axios";

/**
 * 保存活动
 * @param data 
 * @returns 
 */
export function saveActivityApi(data: any) {
  return http("activity-qrcode/activity-save", {
    method: "POST",
    data: data
  });
}

/**
 * 获取活动列表
 * @param data 
 * @returns 
 */
export function getActivityListApi(data: any) {
  return http("activity-qrcode/activity-list", {
    method: "POST",
    data: data,
  });
}

/**
 * 删除活动
 * @param data 
 * @returns 
 */
export function deleteActivityApi(data: any) {
  return http("activity-qrcode/activity-delete", {
    method: "POST",
    data: data
  });
}

/**
 * 获取活动详情
 * @param data 
 * @returns 
 */
export function getActivityDetailApi(data: any) {
  return http("activity-qrcode/activity-detail", {
    method: "POST",
    data: data
  });
}

/**
 * 活动二维码投放分类列表（投放渠道)
 * @param data 
 * @returns 
 */
export function getCategoryListApi(data: any) {
  return http("activity-qrcode/activity-category-list", {
    method: "POST",
    data: data
  });
}

/**
 * 获取导购列表通过门店ID
 * @param data 
 * @returns 
 */
export function getMemberByThirdsApi(data: any) {
  return http("members/list-by-third-ids", {
    method: "POST",
    data: data
  });
}

/**
 * 根据活动码获取导购二维码或者门店列表
 * @param data 
 * @returns 
 */
export function getQrcodeByUnionidApi(data: any) {
  return http("activity-qrcode/get-qrcode-by-unionid", {
    method: "POST",
    data: data
  });
}

/**
 * 活动二维码管理后台统计
 * @param data 
 * @returns 
 */
export function getStatisticsApi(data: any) {
  return http("activity-qrcode/statistics", {
    method: "POST",
    data: data
  });
}