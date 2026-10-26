// 文本原始数据类型
export interface ITextData {
  admin_id: number;
  content: string;
  created_at: number;
  enterprise_id: number;
  group_id: number;
  group_name: string;
  id: number;
  type: number;
}

// 图片原始数据类型
export interface IPictureData {
  id: number;
  enterprise_id: number;
  ext: string;
  name: string;
  url: number;
  size: number;
}
// 链接原始数据类型
export interface ILinkData {
  admin_id: number;
  author: string;
  content: string;
  cover_url: string;
  created_at: number;
  description: string;
  enterprise_id: number;
  group_id: number;
  group_name: string;
  id: number;
  link_type: number;
  link_url: string;
  title: string;
  type: number;
}
// 小程序原始数据类型
export interface IAppData {
  admin_id: number;
  app_id: string;
  app_name: string;
  cover_url: string;
  created_at: number;
  deleted_at: number;
  enterprise_id: number;
  group_id: number;
  id: number;
  logo: string;
  mini_program_id: number;
  origin_id: string;
  path: string;
  title: string;
  type: number;
  updated_at: number;
}
// 文本/表情添加函数参数类型
export type IAllData = ITextData | Array<IPictureData> | ILinkData | IAppData;
