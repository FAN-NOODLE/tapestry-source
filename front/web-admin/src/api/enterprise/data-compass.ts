import http from "@/utils/axios";

export default {
  // 社群数据统计
  getCommunityIndexApi(data: any) {
    return http("community-statistics/index", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "brand" } }
    });
  },
  // 	社群明细列表--按社群
  getCommunityListApi(data: any) {
    return http("community-statistics/list", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "brand" } }
    });
  },
  // 	社群明细列表--按导购
  getCommunityListByMemberApi(data: any) {
    return http("community-statistics/list-by-member", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "brand" } }
    });
  },
  // 社群明细列表-导出
  exportCommunityApi(data: any) {
    return http("community-statistics/export", {
      method: "POST",
      responseType: "blob",
      data: { ...data, ...{ apiVersion: "brand" } }
    });
  },
};
