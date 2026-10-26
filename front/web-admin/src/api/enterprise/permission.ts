import http from "@/utils/axios";

export default {
  // 管理者权限列表
  getManagerAuthListApi(data: any) {
    return http("/dashboard/manager/auth-list", {
      method: "POST",
      data: data
    });
  },
  // 管理者权限详情
  getManagerAuthDetailApi(data: any) {
    return http("/dashboard/manager/auth-detail", {
      method: "POST",
      data: data
    });
  },
  // 管理者权限新增/更新
  saveManagerAuthApi(data: any) {
    return http("/dashboard/manager/auth-save", {
      method: "POST",
      data: data
    });
  },
  // 管理者权限删除
  deleteManagerAuthApi(data: any) {
    return http("/dashboard/manager/auth-delete", {
      method: "POST",
      data: data
    });
  },
  // 已配置的职位memberId集合
  getMemberIdsApi(data: any) {
    return http("/dashboard/manager/all-member-ids", {
      method: "POST",
      data: data
    });
  },
}