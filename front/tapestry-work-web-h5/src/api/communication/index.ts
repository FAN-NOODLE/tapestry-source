import http from "@/utils/axios";
/**
 * 传讯 - 获取RO或FT身份
 * @param data 
 * @returns 
 */
export const getManagerInfoApi = (data: any) => {
  return http("notification/manager", {
    method: "POST",
    data,
  });
}
/**
 * 传讯 - 获取未读数量
 * @param data 
 * @returns 
 */
export const getReadNumApi = (data: any) => {
  return http("notification/list-read-num", {
    method: "POST",
    data,
  });
}
/**
 * 传讯 - 获取传讯列表
 * @param data 
 * @returns 
 */
export const getContentListApi = (data: any) => {
  return http("notification/info-list", {
    method: "POST",
    data,
  });
}
/**
 * 传讯 - 获取传讯详情
 * @param data 
 * @returns 
 */
export const getContentInfoApi = (data: any) => {
  return http("notification/info", {
    method: "POST",
    data,
  });
}
/**
 * 传讯 - 撤回
 * @param data 
 * @returns 
 */
export const recallApi = (data: any) => {
  return http("notification/recall", {
    method: "POST",
    data,
  });
}
/**
 * 传讯 - 审批
 * @param data 
 * @returns 
 */
export const checkApi = (data: any) => {
  return http("notification/check", {
    method: "POST",
    data,
  });
}
/**
 * 传讯 - 作废
 * @param data 
 * @returns 
 */
export const cancelApi = (data: any) => {
  return http("notification/cancel", {
    method: "POST",
    data,
  });
}
/**
 * 传讯 - 传讯附件通过邮件发送
 * @param data 
 * @returns 
 */
export const sendFileByEmailApi = (data: any) => {
  return http("notification/send-file-by-email", {
    method: "POST",
    data,
  });
}
/**
 * 传讯 - RO整合传讯列表
 * @param data 
 * @returns 
 */
export const getMergeListApi = (data: any) => {
  return http("notification/merge-list", {
    method: "POST",
    data,
  });
}
/**
 * 传讯 - RO整合传讯详情
 * @param data 
 * @returns 
 */
export const getMergeDetailApi = (data: any) => {
  return http("notification/merge-detail", {
    method: "POST",
    data,
  });
}
/**
 * 传讯 - 保存pdf
 * @param data 
 * @returns 
 */
export const saveMergePdfApi = (data: any) => {
  return http("notification/save-pdf", {
    method: "POST",
    data,
  });
}