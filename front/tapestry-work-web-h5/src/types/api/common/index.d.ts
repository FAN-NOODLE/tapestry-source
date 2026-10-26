// 授权登录接口入参类型
interface ILoginParams {
  code: string; // 企业微信授权code
  hash: number; // 企业id
}
// 授权登录接口入参类型
interface IRedirectUrlParams {
  redirect_uri: string; // 授权回调地址
  hash: number; // 企业id
}
// 应用签名、企业签名接口入参类型
interface IConfigParams {
  url: string; // URL
  hash: number; // 企业id
}

interface IUserInfo {
  /** 导购工号 */
  user_id: string;
  /** ??? */
  device_id: string;
  /** 导购头像 */
  avatar: string;
  /** 导购ID */
  member_id: number;
  /** 导购昵称 */
  name: string;
  member_name: string;
  /** 部门名称 */
  dept_name: string;
  /** 是否为VM的管理员，1为是，0为否 */
  is_manager: number;
  /** 应用来源 */
  app_source: number;
  /** 是否为VM，1为是，0为否 */
  is_vm: number;
  /** 是否为特殊用户 */
  is_special: number;
  /** 传讯 FT */
  is_ft: number;
  /** 传讯 RO */
  is_ro: number;
  /** 门店ID */
  third_id: number;
  /** 是否为店长 */
  is_store_manager?: boolean
  /** 是否为巡店管理员 */
  is_store_visit_manager?: boolean
  /** 门店名称 */
  store_name?: string
  /** 门店编吗 */
  store_code?: string
}
