import store from "@/store"
import { API } from "./env-config";
import { getWelearnTokenApi } from '@/api/common'
import { encode64 } from "./base64";
function stringify(data: Record<string, any>) {
  return Object
  .keys(data)
  .map((key) => `${key}=${data[key]}`)
  .join('&')
}

export function getRedirectParams(redirectUri: string) {

  // 环境变量
  const host = API.welearning_url
  const token = API.welearning_token
  const companyCode = API.welearning_company_code
  // 参数
  const params = {
    // 目标跳转地址
    redirectUri: /^https?:\/\//.test(redirectUri) ? encode64(redirectUri): redirectUri,
    // 标识key
    uniqueKey: 'workNo',
    // 标识value
    uniqueValue: store.getters["user/getUserInfo"].user_id,
    // 企业code
    companyCode: companyCode,
    // 事件
    time: Date.now(),
    // token
    token: token,
  }

  async function uri() {// 常量
    const g = 'Wap'
    const m = 'SingleLogin'
    const a = 'index'
    
    const res =  await getWelearnTokenApi({
      login_str: stringify(params)
    })

    return `${host}/index.php?${stringify({ 
      g, 
      m, 
      a, 
      ...params, 
      loginToken: res.data.login_token.toUpperCase()
    })}`
  }

  return {
    params: stringify(params),
    uri
  }
}