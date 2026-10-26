import Vue from "vue";

const http = Vue.axios;

let apis = {
  // 登录
  login(data) {
    return http("gw-scrm/login", {
      method: "POST",
      data: data
    });
  },
  // 登录
  whiteLogin(data) {
    return http("gw-scrm/white-login", {
      method: "POST",
      data: data
    });
  },
  // 退出登录
  logout(data) {
    return http("admin/logout", {
      method: "POST",
      data: data
    });
  },
  // 验证码 - 账号密码
  validate(data) {
    return http("gw-scrm/validate-code", {
      method: "POST",
      data: data
    });
  },
  // 验证码登录 - 手机号
  loginCode(data) {
    return http("gw-scrm/get-code", {
      method: "POST",
      data: data
    });
  },
  // 发送验证码(李坚)
  sendCode(data) {
    return http("gw-scrm/get-code", {
      method: "POST",
      data: data
    });
  },
  // 发送短信验证码新(王策)
  sendCodeSms(data) {
    return http("gw-scrm/get-bind-code", {
      method: "POST",
      data: data
    });
  },
  // 获取企业列表
  getEnterpriseToken(data) {
    return http("gw-scrm/get-token", {
      method: "post",
      data: data
    });
  },
  // 企业微信登录接口(php:王策)
  workWechatLogin(data) {
    return http("gw-scrm/work-wechat-login", {
      method: "post",
      data
    });
  },
  // 行业列表
  industryList(data) {
    return http("gw-scrm/get-industry-list", {
      method: "post",
      data
    });
  },
  // 验证邮箱是否已被使用
  checkEmail(data) {
    return http("gw-scrm/check-email", {
      method: "post",
      data
    });
  },
  // 获取安装智慧导购企业微信地址
  getWxInstallUrlApi(data) {
    return http("gw-scrm/install-third-app", {
      method: "post",
      data
    });
  },
  // 自助注册表单提交
  autoAccess(data) {
    return http("gw-scrm/auto-access", {
      method: "POST",
      data
    });
  },
  //修改密码
  editPassWord(data) {
    return http("gw-console/v1/setting/account/edit-password", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } }
    });
  },
  getEnterpriseList(data) {
    return http("gw-console/v1/admin/enterprise-list", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
      loadingHide: true
    });
  },

  checkResetPassword(data) {
    return http("gw-console/v1/admin/check-reset-password", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } }
    });
  },
  resetPassword(data) {
    return http("gw-console/v1/admin/reset-password", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } }
    });
  },
};

export default apis;
