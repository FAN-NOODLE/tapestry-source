import http from "@/utils/axios";

export default {
  // 账号列表
  accountList(data: any) {
    return http("setting/enterprise-admin/list", {
      method: "GET",
      params: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 操作日志列表
  operationLog(data: any) {
    return http("setting/operation-log/list", {
      method: "GET",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 接口模块列表
  operationModule(data: any) {
    return http("setting/operation-log/module", {
      method: "GET",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 新增账号
  addAccount(data: any) {
    return http("setting/enterprise-admin/add", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 删除账号
  deleteAccount(data: any) {
    return http("setting/enterprise-admin/delete", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 编辑账号
  editAccount(data: any) {
    return http("setting/enterprise-admin/update", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 重置账号密码
  resetAccount(data: any) {
    return http("setting/enterprise-admin/reset-password", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 账号禁用/启用
  statusAccount(data: any) {
    return http("setting/enterprise-admin/status", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // ------------------------new--------------------------------
  // 区域下门店列表
  areaStoreList(data?: any) {
    return http("store-area/areaStore-list", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 角色列表接口
  roleList(data: any) {
    return http("rbac/role/list", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 账号详情
  accountDetail(data: any) {
    return http("setting/account/detail", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 账号状态修改
  accountStatus(data: any) {
    return http("setting/account/status", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 删除账号
  accountDel(data: any) {
    return http("setting/account/delete", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 账号列表
  accountsList(data: any) {
    return http("setting/account/list", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 创建账号可用系统选项
  accountSystem(data?: any) {
    return http("setting/account/account-system", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 添加账号
  accountAdd(data: any) {
    return http("setting/account/add", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 编辑账号
  accountEdit(data: any) {
    return http("setting/account/edit", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 检查子账号数量是否达到上限
  checkCount(data?: any) {
    return http("setting/account/check-count", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 修改密码
  editPassword(data?: any) {
    return http("setting/account/edit-password", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 重置密码
  resetPassword(data?: any) {
    return http("setting/account/reset-password", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  }
};
