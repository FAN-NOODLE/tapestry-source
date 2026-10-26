import http from "@/utils/axios";

export default {
  // 获取活动列表
  getTypeList(data: any) {
    return http("activity/get-type-list", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 获取活动用户列表
  getUserList(data: any) {
    return http("activity/get-user-list", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 导出数据
  exportUserData(data: any) {
    return http("activity/export-user-data", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
      responseType: "blob",
    });
  },
  // 详情
  getUserInfo(data: any) {
    return http("activity/get-user-info", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 用户审核
  verifyUser(data: any) {
    return http("activity/verify-user", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 签到列表
  getSignInList(data: any) {
    return http("activity/sign-in/list", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 签到详情
  signInDetail(data: any) {
    return http("activity/sign-in/detail", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 签到导出
  exportSignIn(data: any) {
    return http("activity/sign-in/export", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
      responseType: "blob",
    });
  },
};
