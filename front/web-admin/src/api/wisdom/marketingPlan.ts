import http from "@/utils/axios";
import { IPublicParam } from "@/api/public/public-interface" 

export interface IListCreateTime {
  create_start?: number | string;
  create_end?: number | string;
}
export interface IListEditTime {
  edit_start?: number | string;
  edit_end?: number | string;
}

export interface IMarketingPlanListParams {
  page?: number;
  size: number;
  plan_title: string;
  plan_sn: string;
  create_time?: IListCreateTime;
  edit_time?: IListEditTime;
  operator_name?: string;
}

export interface IMarketingPlanDetailParams {
  plan_id: number | string;
}
export interface IObjectArray {
  object_id: string;
  range_id: string;
  object_name: string;
}
export interface IResource {
  type: string;
  resource_id: string;
  resource_url: string;
  applet_path: string;
  applet_card_title: string;
  applet_card_cover_url: string;
}
export interface ITaskArray {
  task_title: string;
  task_id: string;
  start_time: string;
  end_time: string;
  content: string;
  resource?: Array<IResource>;
}
export interface IMarketingPlanCreateParams {
  plan_title?: string;
  plan_id: string;
  object_type?: string | number;
  object_array?: Array<IObjectArray>;
  task_array?: Array<ITaskArray>;
}
export interface IMarketingPlanDeleteParams {
  plan_id: number | string;
}
export interface IMarketingPlanExportParams {
  plan_id: string;
  task_id: string;
}

// 上传图片/视频成功接口返回
export interface IUploadImageSuccess extends IPublicParam{
  request_id?: string;
}
export interface ISunCode {
  plan_id?: string|number;
  applet_path?: string;
}

// 营销计划列表
export function marketingPlanList<T>(data: IMarketingPlanListParams) {
  return http.post<T>("marketing-plan/list", data);
}

// 营销详情
export function marketingPlanDetail<T>(data: IMarketingPlanDetailParams) {
  return http.post<T>("marketing-plan/detail", data);
}

// 新增营销计划
export function marketingPlanCreate<T>(data: IMarketingPlanCreateParams) {
  return http.post<T>("marketing-plan/create", data);
}
// 编辑营销计划
export function marketingPlanUpdate<T>(data: IMarketingPlanCreateParams) {
  return http.post<T>("marketing-plan/update", data);
}

// 删除营销计划
export function marketingPlanDelete<T>(data: IMarketingPlanDeleteParams) {
  return http.post<T>("marketing-plan/delete", data);
}
// 商城小程序信息
export function miniProgram<T>() {
  return http.post<T>("common/mini-program");
}

// 未完成名单下载
export function marketingPlanExport<T>(data: IMarketingPlanExportParams) {
  return http.post<T>("marketing-plan/export", data);
}

// 未完成名单下载
export function uploadImg<T>(data: any) {
  return http.post<T>("common/upload-img", data);
}
// 小程序海报获取太阳码
export function sunCode<T>(data: ISunCode) {
  return http.post<T>("marketing-plan/sun-code", data);
}
