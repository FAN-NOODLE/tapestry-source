import { Meta } from "@/types";
import http from "@/utils/axios";

export interface Credential {
  bucket: string;
  region: string;
  secret_id: string;
  secret_key: string;
  session_token: string;
  start_time: number;
  expired_time: number;
}

/**
 * 敏感词列表入参的类型
 */
export interface ReqListParams extends Meta {
  sensitive_text?: string;
}

/**
 * 配置接口的保存参数
 */
export interface ReqSaveConfigParams {
  id?: string;
  overtime?: number;
  receive_start?: string;
  receive_end?: string;
  notice_type?: string;
  notice_status?: string;
  stat_txt?: string;
}

/**
 * 敏感词详细列表接口的请求参数
 */
export interface ReqDetailListParams extends Meta {
  id: number;
}

/**
 * 获取聊天历史的请求参数类型
 */
export interface ReqMsgListParams {
  size: number;
  from_user: string;
  to_user: string;
  sort: string;
  msg_id: number;
  is_need_self: boolean;
}

/**
 * 请求敏感词列表
 * @param {ReqListParams} data
 */
export const sensitiveListApi = (data: ReqListParams) =>
  http("sensitive/list", {
    method: "POST",
    data
  });

/**
 * 请求敏感词详情
 * @param id
 */
export const sensitiveDetailApi = (id: number) =>
  http("sensitive/detail", {
    method: "POST",
    data: { id }
  });

/**
 * 请求敏感词详情列表
 * @param {ReqDetailListParams} data
 */
export const sensitiveDetailListApi = (data: ReqDetailListParams) =>
  http("sensitive/detail-list", {
    method: "POST",
    data
  });

/**
 * 添加敏感词
 * @param {ReqSaveParams} data
 */
export const sensitiveSaveApi = (data: any) =>
  http("sensitive/save", {
    method: "POST",
    data
  });
/**
 * 删除敏感词
 * @param id
 */
export const sensitiveDeleteApi = (id: number) =>
  http("sensitive/delete", {
    method: "POST",
    data: { id }
  });

/**
 * 获取提醒设置信息
 */
export const configInfoApi = (notice_type: number) =>
  http("remind-config/config-info", {
    method: "POST",
    data: { notice_type }
  });

/**
 *保存提醒设置信息
 * @param {ReqSaveConfigParams} data
 */
export const configSaveApi = (data: ReqSaveConfigParams) =>
  http("remind-config/config-save", {
    method: "POST",
    data
  });

export const getChatlistByMsgIdApi = (data: ReqMsgListParams) =>
  http("chat/msg-list-by-id", {
    method: "POST",
    data
  });

export const getFileCredentialApi = () =>
  http("file/credential", { method: "POST" });

export const getPrivateByKeyApi = (data) =>
  http("file/private-by-key", { method: "POST", data });

export const getPrivateFileUrlApi = (data) =>
  http("file/private-url", { method: "POST", data });
