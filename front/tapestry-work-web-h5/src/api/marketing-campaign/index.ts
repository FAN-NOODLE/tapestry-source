import http from "@/utils/axios";
// 活动列表
export const activityListApi = (data: Activitylist) =>
  http("api/h5/activity", {
    method: "GET",
    params: data,
    data: { ...data, apiVersion: "new" },
  });

// 活动统计
export const activityStatisticsApi = (data: any) =>
  http("api/h5/activity/withcount", {
    method: "GET",
    params: data,
    data: { ...data, apiVersion: "new" },
  });

// 活动详情
export const activityDetailApi = (data: any) =>
  http("api/h5/activity/detail/" + data.id, {
    method: "GET",
    data: { ...data, apiVersion: "new" },
  });

// 统计详情
export const statisticsDetailApi = (data: any) =>
  http("api/h5/activity/statistics/" + data.id, {
    method: "POST",
    data: { ...data, apiVersion: "new" },
  });
