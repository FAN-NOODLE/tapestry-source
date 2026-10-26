import http from "@/utils/axios";

export default {
  // 角色列表
  roleList(data: any) {
    return http("rbac/role/list", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 新增角色
  addRole(data: any) {
    return http("rbac/role/add", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 删除角色
  deleteRole(data: any) {
    return http("rbac/role/delete", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 编辑角色
  editRole(data: any) {
    return http("rbac/role/edit", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 获取当前账号可分配给角色的管理后台权限树
  configurablePermissions(appCode: string) {
    return http("rbac/role/configurable-permissions", {
      method: "POST",
      data: { app_code: appCode },
    });
  },
  // 当前企业已开通的管理后台
  roleApplications() {
    return http("setting/account/account-system", {
      method: "POST",
      data: { apiVersion: "old" },
    });
  },
  // 所有角色基础权限列表
  rolePower() {
    return http("setting/role/power", {
      method: "GET",
    });
  },
  // 获取角色详情
  roleInfo(data: any) {
    return http("rbac/role/detail", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 获取组织列表
  departmentList(data: any) {
    return http("department/list", {
      method: "GET",
      params: { ...data, ...{ apiVersion: "old" } },
    });
  },
};
