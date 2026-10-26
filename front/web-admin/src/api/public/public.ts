import http from "@/utils/axios";

export default {
  login(data: any) {
    return http("login", {
      method: "POST",
      data,
    });
  },
  ssologin(data: any) {
    return http("sso/login", {
      method: "POST",
      data,
    });
  },
  logout(data: any) {
    return http("logout", {
      method: "POST",
      data,
    });
  },
  // 获取企业列表
  getEnterpriseToken(data?: any) {
    return http("get-token", {
      method: "POST",
      data,
    });
  },
  getEnterprise(data?: any) {
    return http("change-enterprise-list", {
      method: "post",
      data,
    });
  },
  getHelp(data?: any) {
    return http("index/help", {
      method: "post",
      data,
    });
  },
  // 图片验证码 - 账号密码
  validate(data?: any) {
    return http("validate-code", {
      method: "POST",
      data: data,
    });
  },
  // 获取版本信息
  currentVersionInfo(data?: any) {
    return http("current-version-info", {
      method: "POST",
      data,
    });
  },
  // blob上传获取sas
  getSas(data?: any) {
    return http("upload-video/azure-blob-sas", {
      method: "GET",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  getEnterpriseList(data) {
    return http("admin/enterprise-list", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } }
    });
  },
  checkResetPassword(data) {
    return http("admin/check-reset-password", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } }
    });
  },
  resetPassword(data) {
    return http("admin/reset-password", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } }
    });
  },
  // 发送验证码(李坚)
  sendCode(data) {
    return http("get-code", {
      method: "POST",
      data: data
    });
  },
};
