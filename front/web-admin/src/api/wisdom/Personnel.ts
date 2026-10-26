import http from "@/utils/axios";

export default {
  // 部门列表
  depListApi(data: any) {
    return http("department/list", {
      method: "GET",
      // loadingHide: true, // 隐藏 loading
      params: data,
      data: { apiVersion: "old" },
    });
  },
  // 部门列表
  depListTwoApi(data: any) {
    return http("department/list", {
      method: "GET",
      params: data,
    });
  },
  // 部门分页列表
  depPageList(data: any) {
    return http("department/page-list", {
      method: "GET",
      params: data,
    });
  },
  // 上次同步时间
  syncLastTime(data: any) {
    return http("contact/sync_last_time", {
      method: "GET",
      params: data,
    });
  },
  // 同步组织结构
  syncOrg(data: any) {
    return http("contact/sync_org", {
      method: "POST",
      data: { apiVersion: "old", ...data },
    });
  },
  // 获取同步状态
  syncStatus(data: any) {
    return http("contact/sync_status", {
      method: "GET",
      params: data,
    });
  },
  // 获取同步状态
  TreePeople(data?: any) {
    return http("user/tree-list", {
      method: "POST",
      data: {
        ...data,
        apiVersion: 'old'
      }
    });
  },
};
