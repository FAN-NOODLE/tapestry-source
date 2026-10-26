import http from "@/utils/axios";

export default {
  // 组织列表
  depatListApi(data: any) {
    return http("departments/list", {
      method: "POST",
      data: data,
    });
  },
  // 员工列表
  memberListApi(data: any) {
    return http("members/list", {
      method: "POST",
      data: data,
    });
  },
  // 员工详情
  memberDetailApi(data: any) {
    return http("members/detail", {
      method: "POST",
      data: data,
    });
  },
  // 外部联系人secret更新
  authContacUpdateApi(data: any) {
    return http("auth-contact-update", {
      method: "POST",
      data: data,
    });
  },
  // 企微配置
  settingsApi() {
    return http("settings", {
      method: "POST",
    });
  },
  //自建应用更新
  workPluginUpdateApi(data: any) {
    return http("work-plugin-update", {
      method: "POST",
      data: data,
    });
  },
  //上传校验文件
  uploadCheckFileApi(data: any) {
    const form_data = new FormData();
    form_data.append("file", data.file);
    return http("upload-check-file", {
      method: "POST",
      data: form_data,
    });
  },
  // 获取组织更新员工信息
  updateRecordApi() {
    return http("members/sync", {
      method: "POST",
    });
  },
  // 获取组织更新员工信息
  updateRocordApi() {
    return http("departments/update-record", {
      method: "POST",
    });
  },

  // 获取公众号列表
  getOfficialAccount(data: any) {
    return http("auth/official-auth-list", {
      method: "POST",
      data,
    });
  },
  //解绑公众号
  unbindOfficialAccount(data: any) {
    return http("auth/official-auth-unbind", {
      method: "POST",
      data,
    });
  },
  // 保存支付配置
  setPayApi(data: any) {
    const form_data = new FormData();
    form_data.append("private_key", data.private_key);
    form_data.append("public_key", data.public_key);
    form_data.append("merchant_id", data.merchant_id);
    form_data.append("merchant_key", data.merchant_key);
    form_data.append("serial_no", data.serial_no);
    form_data.append("api_v3_key", data.api_v3_key);
    form_data.append("apiVersion", "old");
    return http("payment/save", {
      method: "POST",
      data: form_data,
      // data: form_data
    });
  },
  // 获取支付配置
  getPayApi(data?: any) {
    return http("payment/detail", {
      method: "get",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 安装第三方应用
  installThirdApp(data: any) {
    return http("auth/install-third-app", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // shopAuthDetail(data)
  // 获取申请接口
  registerApply(data: any) {
    return http("auth/register-apply", {
      method: "post",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 获取详情接口
  registerDetail(data: any) {
    return http("/auth/register-detail", {
      method: "post",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 获取更新详情接口
  registerUpdataInfo(data: any) {
    return http("/auth/register-update-info", {
      method: "post",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 获取商城小程序授权详情接口
  shopAuthDetail(data: any) {
    return http("/auth/shop-auth-detail", {
      method: "post",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  //保存客户信息
  EditShopApi(data: any) {
    return http("/auth/shop-save", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 刷新token
  refreshToken() {
    return http("/auth/update-refresh-token", {
      method: "POST",
      data: { apiVersion: "old" },
    });
  },
  // 企业配置
  getEpConfigApi() {
    return http("enterprise/config", {
      method: "POST",
    });
  },
  // 企业配置
  getEpConfig() {
    return http("enterprise/config", {
      method: "POST",
    });
  },
};
