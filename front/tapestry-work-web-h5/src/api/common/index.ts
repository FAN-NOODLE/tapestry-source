import http from "@/utils/axios";
// 授权登录
export const loginApi = (data: ILoginParams) =>
  http("auth/login", {
    method: "POST",
    data,
  });
// 授权地址
export const redirectUrlApi = (data: IRedirectUrlParams) =>
  http("auth/redirect_url", {
    method: "POST",
    data,
  });
// 应用签名(获取企业微信授权信息)
export const getAgentConfigApi = (data: IConfigParams) =>
  http("auth/agent-config", {
    method: "POST",
    data,
  });
// 企业签名(获取微信授权信息)
export const getConfigApi = (data: IConfigParams) =>
  http("auth/config", {
    method: "POST",
    data,
  });
// 模拟登录
export const mockApi = (data: any) =>
  http("auth/mock", {
    method: "GET",
    params: data,
  });
// 获取导购信息
export const getBaseInfoApi = () =>
  http("member/base-info", {
    method: "POST",
  });
// 获取getWelearnTokenApi
export const getWelearnTokenApi = (data:any) =>
  http("display/api/welearn-token", {
    method: "POST",
    data
  });
