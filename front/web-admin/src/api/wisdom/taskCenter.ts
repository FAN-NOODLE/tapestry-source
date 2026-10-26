import http from "@/utils/axios";

export interface TaskCenterCreateParams {
  task_id?: string | number;
  object_type?: string | number;
  object_array?: Array<IObjectArray>;
  task_array?: Array<ITaskArray>;
}

export interface IObjectArray {
  object_id: string;
  range_id: string;
  object_name: string;
}

export interface ITaskArray {
  task_title: string;
  start_time: number | string;
  end_time: number | string;
  type: number;
  tag_ids?: Array<any>;
  before_start_time?: number;
  before_end_time?: number;
  fixed_time?: number;
  content: string;
  resource?: Array<IResource>;
}

export interface IResource {
  type: string;
  resource_id: string;
  resource_url: string;
  applet_path: string;
  applet_card_title: string;
  applet_card_cover_url: string;
}

// 任务列表
export function taskCenterList<T>(data: any) {
  return http("marketing-plan/list", {
    method: "POST",
    data
  });
}

// 新增任务
export function taskCenterCreate<T>(data: TaskCenterCreateParams) {
  return http("marketing-plan/create", {
    method: "POST",
    data
  });
}

// 编辑任务
export function taskCenterUpdate<T>(data: TaskCenterCreateParams) {
  return http("marketing-plan/update", {
    method: "POST",
    data
  });
}

// 任务详情
export function taskCenterDetail<T>(data: any) {
  return http("marketing-plan/detail", {
    method: "POST",
    data
  });
}

// 删除任务
export function taskCenterDelete<T>(data: TaskCenterCreateParams) {
  return http("marketing-plan/delete", {
    method: "POST",
    data
  });
}

// 未完成名单下载
export function marketingPlanExport<T>(data: any) {
  return http("marketing-plan/export", {
    method: "POST",
    data,
    responseType: "blob"
  });
}