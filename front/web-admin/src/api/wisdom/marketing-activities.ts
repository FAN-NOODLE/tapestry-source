import http from "@/utils/axios";

export default {
  // 活动列表
  activityList(data: any) {
    return http("activity", {
      method: "GET",
      params: data,
      data: { ...data, ...{ apiVersion: "rest" } },
    });
  },
  // 活动上下架状态更改
  activityUpdataOnline(data: any) {
    return http(`activity/online/${data.id}`, {
      method: "PUT",
      data: { ...data, ...{ apiVersion: "rest" } },
    });
  },
  // 活动删除
  activityDelete(data: any) {
    return http(`activity/${data.id}`, {
      method: "DELETE",
      data: { ...data, ...{ apiVersion: "rest" } },
    });
  },
  // 活动新增提交
  activityAddSubmit(data: any) {
    return http(`activity/${data.steps}`, {
      method: "POST",
      data: { ...data, ...{ apiVersion: "rest" } },
    });
  },
  // 活动编辑查询
  activitySteps(data: any) {
    return http(`activity/${data.id}/${data.steps}`, {
      method: "GET",
      data: { ...data, ...{ apiVersion: "rest" } },
    });
  },
  // 活动编辑更新
  activityUpdata(data: any) {
    return http(`activity/id/${data.id}/steps/${data.steps}`, {
      method: "PUT",
      data: { ...data, ...{ apiVersion: "rest" } },
    });
  },
  // 活动详情（基本信息）
  activityBaseInfo(data: any) {
    return http(`activity/detail/base-info/${data.id}`, {
      method: "GET",
      data: { ...data, ...{ apiVersion: "rest" } },
    });
  },
  // 活动详情（表单配置）
  activityContent(data: any) {
    return http(`activity/detail/content/${data.id}`, {
      method: "GET",
      data: { ...data, ...{ apiVersion: "rest" } },
    });
  },
  // 活动详情（统计与详情）
  activityStatistics(data: any) {
    return http(`activity/detail/data/${data.id}/${data.type}`, {
      method: "GET",
      params: data,
      data: { ...data, ...{ apiVersion: "rest" } },
    });
  },
  // 列表导出
  exportStatistics(data: any) {
    return http(`activity/detail/export/${data.id}`, {
      method: "GET",
      responseType: "blob",
      params: data,
      data: { ...data, ...{ apiVersion: "rest" } },
    });
  },
  // 用户详情
  userDetail(data: any) {
    return http(`activity/user/detail/${data.id}`, {
      method: "POST",
      data: { ...data, ...{ apiVersion: "rest" } },
    });
  },
  // 用户审核
  auditUser(data: any) {
    return http(`activity/user/audit/${data.id}`, {
      method: "POST",
      data: { ...data, ...{ apiVersion: "rest" } },
    });
  },
};
