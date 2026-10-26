import http from "@/utils/axios";
// 裂变活动列表接口入参类型
// 裂变活动列表
export const getFissionActivityListApi = (data: IFissionActivityList) =>
  http("fission_activity/list", {
    method: "POST",
    data,
  });
