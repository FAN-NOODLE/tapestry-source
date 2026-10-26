import http from "@/utils/axios";

/**
 * 巡店 - 获取预置的部门和人员
 * @param data 
 * @returns 
 */
export const getDeptMemberListApi = (data: any) => {
  return http("store-visit/dept-member-list", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 获取个人邮箱通讯录
 * @param data 
 * @returns 
 */
export const getEmailContactListApi = (data: any) => {
  return http("store-visit/email-contact", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 个人邮箱通讯录新增/修改
 * @param data 
 * @returns 
 */
export const saveEmailApi = (data: any) => {
  return http("store-visit/save-email", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 个人邮箱通讯录删除
 * @param data 
 * @returns 
 */
export const deleteEmailApi = (data: any) => {
  return http("store-visit/del-email", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 根据姓名或邮箱搜索导购
 * @param data 
 * @returns 
 */
export const searchMemberApi = (data: any) => {
  return http("store-visit/search-member", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 门店搜索
 * @param data 
 * @returns 
 */
export const searchStoreApi = (data: any) => {
  return http("display/dashboard/search-store", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 创建巡店计划
 * @param data 
 * @returns 
 */
export const saveVisitPlanApi = (data: any) => {
  return http("store-visit/save-plan", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 巡店计划列表
 * @param data 
 * @returns 
 */
export const getStoreVisitListApi = (data: any) => {
  return http("store-visit/list", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 当前登录用户信息及权限
 * @param data 
 * @returns 
 */
export const getStoreVisitAuthApi = (data: any) => {
  return http("store-visit/login-user-info", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 巡店记录详情
 * @param data 
 * @returns 
 */
export const getStoreVisitInfoApi = (data: any) => {
  return http("store-visit/info", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 无预约开始巡店/巡店打卡/有预约的开始巡店
 * @param data 
 * @returns 
 */
export const startVisitApi = (data: any) => {
  return http("store-visit/start-visit", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 巡店记录的反馈记录列表
 * @param data 
 * @returns 
 */
export const getVisitRecordListApi = (data: any) => {
  return http("store-visit/show-record-list", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 新增/修改反馈记录
 * @param data 
 * @returns 
 */
export const saveVisitRecordApi = (data: any) => {
  return http("store-visit/save-show-record-content", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 删除反馈记录
 * @param data 
 * @returns 
 */
export const deleteVisitRecordApi = (data: any) => {
  return http("store-visit/del-show-record", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 新增反馈记录素材
 * @param data 
 * @returns 
 */
export const saveVisitRecordMediaApi = (data: any) => {
  return http("store-visit/save-show-record-media", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 删除反馈记录素材
 * @param data 
 * @returns 
 */
export const deleteVisitRecordMediaApi = (data: any) => {
  return http("store-visit/del-show-record-media", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 提交巡店记录（完成巡店
 * @param data 
 * @returns 
 */
export const completeVisitApi = (data: any) => {
  return http("store-visit/complete-visit", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 更新重置概览、整体表现 / 保存巡店报告pdf
 * @param data 
 * @returns 
 */
export const saveVisitInfoApi = (data: any) => {
  return http("store-visit/save-visit-info", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 保存更新 巡店报告查看权限及新增巡店人
 * @param data 
 * @returns 
 */
export const saveVisitReportAuthApi = (data: any) => {
  return http("store-visit/save-report-permission", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 发送巡店报告通知
 * @param data 
 * @returns 
 */
export const sendVisitReportApi = (data: any) => {
  return http("store-visit/send-complete-report", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 新增评论
 * @param data 
 * @returns 
 */
export const saveCommentApi = (data: any) => {
  return http("store-visit/save-comment", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 评论列表（分页）
 * @param data 
 * @returns 
 */
export const getCommentListApi = (data: any) => {
  return http("store-visit/comment-list", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 删除评论
 * @param data 
 * @returns 
 */
export const deleteCommentApi = (data: any) => {
  return http("store-visit/del-comment", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 首次打开
 * @param data 
 * @returns 
 */
export const getFirstOpenStateApi = (data: any) => {
  return http("member/first-open", {
    method: "POST",
    data,
  });
}

/**
 * 巡店 - 巡店计划同步至其他人
 * @param data 
 * @returns 
 */
export const addCopyMemberApi = (data: any) => {
  return http("store-visit/add-copy-member", {
    method: "POST",
    data,
  });
}

/**
 * 
 * 巡店 - 看板统计
 * @param data 
 * @returns 
 */
export const getDashboardTreeApi = (data: any) => {
  return http("store-visit/dashboard-tree", {
    method: "POST",
    data,
  });
}

/**
 * 
 * 巡店 - 看板统计导出
 * @param data 
 * @returns 
 */
export const exportDashboardListApi = (data: any) => {
  return http("store-visit/export-dashboard-list", {
    method: "POST",
    data,
  });
}

/**
 * 
 * 巡店 - 下载pdf
 * @param data 
 * @returns 
 */
export const sendPdfToEmailApi = (data: any) => {
  return http("store-visit/send-pdf-to-email", {
    method: "POST",
    data,
  });
}

/**
 * 
 * 巡店 - 获取语言配置文件
 * @param data 
 * @returns 
 */
export const getConfigApi = (data: any) => {
  return http("store-visit/get-config", {
    method: "POST",
    data,
  });
}

/**
 * 
 * 巡店 - 保存语言配置文件
 * @param data 
 * @returns 
 */
export const saveConfigApi = (data: any) => {
  return http("store-visit/save-language-config", {
    method: "POST",
    data,
  });
}


/**
 *  获取门店管理员列表
 */
export const getStoreManagerApi = (data?: any) => {
  return http("store-visit/get-store-manager", {
    method: "POST",
    data,
  });
}


/**
 *  获取门店管理员列表
 */
export const getStoreMMApi = (data?: any) => {
  return http("store-visit/get-store-mm", {
    method: "POST",
    data,
  });
}
