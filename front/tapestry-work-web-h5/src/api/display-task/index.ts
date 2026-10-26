import http from "@/utils/axios";
/**
 * 导购获取任务列表
 * @param data 
 * @returns 
 */
export const getTaskListApi = (data: any) => {
  return http("display/api/task-list", {
    method: "POST",
    data,
  });
}

/**
 * 导购获取任务详情
 */
export const getTaskInfoApi = (data: any) => {
  return http("display/api/task-info", {
    method: "POST",
    data,
  });
}
/**
 * 获取门店任务详情
 */
export const getStoreTaskInfoApi = (data: any) => {
  return http("display/api/store-task-info", {
    method: "POST",
    data,
  });
}

/**
 * 导购操作门店VM任务数据，存草稿和提交
 */
export const saveTaskApi = (data: any) => {
  return http("display/api/save", {
    method: "POST",
    data,
  });
}

/**
 * 查看VM的备注上报
 */
export const lookTaskApi = (data: any) => {
  return http("display/api/look", {
    method: "POST",
    data,
  });
}

/**
 * 获取当前身份是否为VM
 */
export const getAuthApi = () => {
  return http("display/api/auth", {
    method: "POST",
    data: {},
  });
}

/**
 * 获取VM成员列表
 */
export const getVmListApi = () => {
  return http("display/api/vm-list", {
    method: "POST",
  });
}

/**
 * 获取陈列其他门店
 */
export const getStoreListApi = (data: any) => {
  return http("display/api/store-list", {
    method: "POST",
    data,
  });
}

/**
 * 获取各任务状态下的任务数
 */
export const getListTipsApi = (data: any) => {
  return http("display/api/list-tips", {
    method: "POST",
    data,
  });
}

/**
 * VM审核列表接口
 */
export const getCheckListApi = (data: any) => {
  return http("display/api/check-list", {
    method: "POST",
    data,
  });
}

/**
 * VM审核详情门店列表接口
 */
export const getCheckStoreListApi = (data: any) => {
  return http("display/api/check-store-list", {
    method: "POST",
    data,
  });
}

/**
 * 一键提醒
 */
export const sendNoticeApi = (data: any) => {
  return http("display/api/send-notice", {
    method: "POST",
    data,
  });
}

/**
 * 审核
 */
export const saveVmCheckApi = (data: any) => {
  return http("display/api/vm-check", {
    method: "POST",
    data,
  });
}

/**
 * 保存图片
 */
export const saveImageApi = (data: any) => {
  return http("display/api/save-image", {
    method: "POST",
    data,
  });
}
/**
 * VM看板任务列表
 */
export const getAllTaskListApi = (data?: any) => {
  return http("display/api/all-list", {
    method: "POST",
    data,
  });
}
/**
 * VM看板 dashboard
 */
export const getDashboardInfoApi = (data?: any) => {
  return http("display/api/vm-dashboard", {
    method: "POST",
    data,
  });
}
/**
 * VM看板审核未通过列表
 */
export const getFailListApi = (data?: any) => {
  return http("display/api/fail-list", {
    method: "POST",
    data,
  });
}
/**
 * VM看板审核未通过列表
 */
export const getMemberAuthApi = (data?: any) => {
  return http("display/dashboard/member-auth", {
    method: "POST",
    data,
  });
}

/**
 * VM看板 门店树
 */
export const getStoreTreeApi = (data?: any) => {
  return http("display/dashboard/tree", {
    method: "POST",
    data,
  });
}

/**
 * VM看板 层级穿透
 */
export const getDashboardListApi = (data?: any) => {
  return http("display/api/dashboard-list", {
    method: "POST",
    data,
  });
}
/**
 * 删除图片
 */
export const deleteImageApi = (data: any) => {
  return http("display/api/del-image", {
    method: "POST",
    data,
  });
}
/**
 * 保存备注
 */
export const saveCommentApi = (data: any) => {
  return http("display/api/save-comment", {
    method: "POST",
    data,
  });
}
/**
 * 删除备注
 */
export const deleteCommentApi = (data: any) => {
  return http("display/api/del-comment", {
    method: "POST",
    data,
  });
}
/**
 * VM看板 层级穿透
 */
export const getStoreTaskListApi = (data?: any) => {
  return http("display/api/store-task-list", {
    method: "POST",
    data,
  });
}
/**
 *  special 保存草稿/提交
 */
export const specialSaveApi = (data?: any) => {
  return http("display/api/special-save", {
    method: "POST",
    data,
  });
}