import { IPermissionItem, IRoleItem } from "@/types";

export interface IEnterpriseItem {
  admin_id: number;
  enterprise_id: number; //企业id
  hash: string; //企业hash
  version: number; //版本
  version_tag: string; // 版本标签（正式/试用）
  name: string; //企业名
  contact_name: string;
  checked?: boolean;
  has_expire?: any; //旧版的没有找到这个字段，不清楚类型
  account_expire_at: number;
}

export interface ITokenData {
  admin_list?: any;
  enterprise_admin_info?: any;
  hash?: string;
  js_api_config?: string | any;
  token: string;
  total: number;
  version: number;
  permissions_tree: Array<IPermissionItem>;
  permissions_list: Array<IPermissionItem>;
  role_list: Array<IRoleItem>;
}
