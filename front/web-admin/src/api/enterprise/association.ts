import { Meta } from "@/types";
import http from "@/utils/axios";

/**
 * 敏感词详细列表接口的请求参数
 */
export interface IGroupList extends Meta {
  name: string;
  status: number;
  member_id: number;
  start_at: number;
  end_at: number;
}
// 明细
export interface ITaskDetail extends Meta {
  id: number;
  member_id: number;
}
export interface IGroupSendAdd {
  id?: number;
  name: string;
  type: number;
  send_time?: number;
  content: any;
  member_ids: [];
}
export interface IChatList {
  chat_name: string;
  chat_id: string;
  channel_id: number;
  owner: string;
  operator: string;
  department_id: number;
  chat_type_id: number;
  sts?: number | string;
  ets?: number | string;
}

export interface IChatType {
  opt: number;
  id?: number;
  name?: string;
}
export interface IchatSet {
  chat_id: string;
  channel_id?: number;
  chat_type_id?: number | string;
  operator_id?: string;
  chat_target?: string;
}

export interface IchatInfo {
  enterprise_id?: number;
  chat_id: string;
}

export interface IchatMember {
  chat_id: string;
  member_type: number;
  sts?: number;
  ets?: number;
}

export interface IchatData {
  chat_id: string;
  sts?: number;
  ets?: number;
}
/**
 * 群发消息列表
 */
export const groupListApi = (data: IGroupList) =>
  http("group-msg/list", {
    method: "POST",
    data
  });
/**
 * 群发消息详情
 */
export const taskDetailApi = (data: ITaskDetail) =>
  http("group-msg/task-detail", {
    method: "POST",
    data
  });
/**
 * 群发消息明细
 */
export const taskInfoApi = (id: number) =>
  http("group-msg/task-info", {
    method: "POST",
    data: { id }
  });
/**
 * 执行群主下拉列表
 */
export const memberlist = (id?: number) =>
  http("group-msg/member-list", {
    method: "POST",
    data: { id }
  });
/**
 * 新增群发消息
 */
export const groupSendAdd = (data: IGroupSendAdd) =>
  http("group-msg/group-send-add", {
    method: "POST",
    data
  });
/**
 * 编辑群发消息
 */
export const groupSendEdit = (data: IGroupSendAdd) =>
  http("group-msg/group-send-edit", {
    method: "POST",
    data
  });
/**
 * 社群列表
 */
export const chatListApi = (data: IChatList) =>
  http("chat-group/list", {
    method: "POST",
    data
    // data: { ...data, ...{ apiVersion: "yang" } }
  });
/**
 * 社群类型管理
 */
export const chatTypeManage = (data: IChatType) =>
  http("chat-group/type-manage", {
    method: "POST",
    data
    // data: { ...data, ...{ apiVersion: "yang" } }
  });

/**
 * 社群类型列表
 */
export const chatTypeList = () =>
  http("chat-group/types", {
    method: "POST"
    // data: { ...{ apiVersion: "yang" } }
  });

/**
 * 社群设置
 */
export const chatSettingApi = (data: IchatSet) =>
  http("chat-group/setting", {
    method: "POST",
    data
    // data: { ...data, ...{ apiVersion: "yang" } }
  });
/**
 * 社群信息
 */
export const chatInfoApi = (data: IchatInfo) =>
  http("chat-group/info", {
    method: "POST",
    data
    // data: { ...data, ...{ apiVersion: "yang" } }
  });
/**
 * 群成员明细
 */
export const chatMemberApi = (data: IchatMember) =>
  http("chat-group/members", {
    method: "POST",
    data
    // data: { ...data, ...{ apiVersion: "yang" } }
  });
/**
 * 群数据指标、饼图
 */
export const chatDataApi = (data: IchatData) =>
  http("chat-group/data-grid", {
    method: "POST",
    data
    // data: { ...data, ...{ apiVersion: "yang" } }
  });

/**
 * 渠道列表
 */
export const chatChannelApi = (data: any) =>
  http("chat-group/filters", {
    method: "POST",
    data
    // data: { ...data, ...{ apiVersion: "yang" } }
  });
/**
 * 渠道列表
 */
export const chatAllChannelApi = () =>
  http("drainage/channels-simple-list", {
    method: "POST"
    // data
    // data: { ...{ apiVersion: "yang" } }
  });

/**
 * 获取数据更新时间
 */
export const getRefreshTime = (data: any) =>
  http("chat-group/data-update-time", {
    method: "POST",
    data
    // data: { ...data, ...{ apiVersion: "yang" } }
  });

/**
 * 同步群信息
 */
export const updateGroup = (data: any) =>
  http("chat-group/sync", {
    method: "POST",
    data
    // data: { ...data, ...{ apiVersion: "yang" } }
  });
/**
 * 同步群信息
 */
export const batchSetting = (data: any) =>
  http("chat-group/batch-setting", {
    method: "POST",
    data
    // data: { ...data, ...{ apiVersion: "yang" } }
  });

/**
 * 同步群信息
 */
export const downLoadExcel = (data: any) =>
  http("group-chat/download", {
    method: "GET",
    params: data
    // data: { ...data, ...{ apiVersion: "yang" } }
  });

/**
 * 员工明细
 */
export const taskMemberDetailApi = (data: any) =>
  http("group-msg/task-member-detail", {
    method: "POST",
    data
  });
/**
 * 员工详情
 */
export const taskSendApi = (data: any, responseType) =>
  http("group-msg/task-member-send-list", {
    method: "POST",
    ...responseType,
    data
  });
