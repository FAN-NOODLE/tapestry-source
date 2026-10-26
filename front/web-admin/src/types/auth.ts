// 用户信息
export interface IUserInfo {
  account: string;
  account_email: string;
  account_expire_at: number;
  account_global_roaming: string;
  account_mobile: string;
  account_name: string;
  account_start_at: number;
  account_status: number;
  address: string;
  admin_available_system: string;
  admin_id: number;
  contact_mobile: string;
  contact_name: string;
  contact_status: number;
  corp_id: string;
  crs_status: number;
  domain: string;
  enterprise_id: number;
  enterprise_role_list: Array<any>;
  hash: string;
  is_split_flow: number;
  is_super: number;
  logo: string;
  name: string;
  range_type: number;
  role: number;
  serve_type: number;
  shop_name: string;
  store_status: number;
  version: number;
  version_tag: string;
  warehouse_type: number;
}

// 权限树的item
export interface IPermissionItem {
  key: string;
  id: number;
  center_shop?: number;
  center_shop_store?: number;
  direct_chain_store?: number;
  icon: string;
  level: number;
  pid: number;
  name: string;
  path: string;
  type: number;
  children: Array<IPermissionItem>;
}
export interface IKeysItem {
  key: string;
  id: number;
  center_shop?: number;
  center_shop_store?: number;
  direct_chain_store?: number;
  icon: string;
  pid: number;
  name: string;
  path: string;
  type: number;
}

// 角色list的item
export interface IRoleItem {
  admin_id: number;
  created_at: number;
  enterprise_id: number;
  id: number;
  mark: string;
  name: string;
  pid: number;
  released_at: number;
  weight: number;
}

export interface INodeItem {
  children?: Array<INodeItem>;
  path: string;
  name: string;
  ename?: string;
  cname?: string;
  icon: string;
  type: number;
  root?: string;
  meta: any;
  key: string;
  file_path: string;
  hide?: boolean;
  show: boolean;
  component: any;
  fid: string;
  kid: string;
  redirect?: string;
}

// 面包屑
export interface IBCNodeItem {
  direct_chain_store: number;
  icon: string;
  id: number;
  key: string;
  name: string;
  path: string;
  pid: number;
  type: number;
  children: Array<IBCNodeItem>;
}
// 帮助中心config
export interface IHelpConfig {
  key: string;
  content: string | Array<any>;
  type: string;
}
export interface IActionConfig {
  isShow: boolean;
  config: IHelpConfig;
}
export interface IConfig {
  config: IActionConfig;
}
