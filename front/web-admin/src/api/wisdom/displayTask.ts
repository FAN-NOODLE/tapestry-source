import http from "@/utils/axios";

export default {
  /** 新建任务 */
  addDisplayTaskApi(data: any) {
    return http("display/task/add", {
      method: "POST",
      data
    });
  },
  /** 获取任务列表 */
  getDisplayTaskListApi(data: any) {
    return http("display/task/list", {
      method: "POST",
      data
    });
  },
  /** 获取任务详情 */
  getDisplayTaskDetailApi(data: any) {
    return http("display/task/detail", {
      method: "POST",
      data
    });
  },
  /** 复制任务 */
  copyDisplayTaskApi(data: any) {
    return http("display/task/copy", {
      method: "POST",
      data
    });
  },
  /** 更新任务 */
  updateDisplayTaskApi(data: any) {
    return http("display/task/edit", {
      method: "POST",
      data
    });
  },
  /** 更新任务状态 */
  updateDisplayTaskStatusApi(data: any) {
    return http("display/task/update-show", {
      method: "POST",
      data
    });
  },
  /** 获取陈列数据 */
  getAllotListApi(data: any) {
    return http("display/task/allot-list", {
      method: "POST",
      data
    });
  },
  /** 获取陈列数据统计 */
  getAllotNumApi(data: any) {
    return http("display/task/allot-num", {
      method: "POST",
      data
    });
  },
  /** 导出数据 */
  exportAllotListApi(data: any) {
    return http("display/task/allot-list-export", {
      method: "POST",
      data
    });
  },
  /** 设置提醒 */
  setNoticeApi(data: any) {
    return http("display/task/notice-setting", {
      method: "POST",
      data
    });
  },
  /** 查看陈列数据 */
  getAllotInfoApi(data: any) {
    return http("display/task/allot-info", {
      method: "POST",
      data
    });
  },
  /** 审核 */
  saveAllotApi(data: any) {
    return http("display/task/allot-check", {
      method: "POST",
      data
    });
  },
  /** 获取提醒配置 */
  getNoticeSettingApi(data: any) {
    return http("display/task/notice-setting-info", {
      method: "POST",
      data
    });
  },
  /** 一键提醒 */
  sendNoticeApi(data: any) {
    return http("display/task/send-notice", {
      method: "POST",
      data
    });
  },
  /** 一键提醒 */
  getNoticeNumApi(data: any) {
    return http("display/task/notice-num", {
      method: "POST",
      data
    });
  },
  /** 获取VM列表 */
  getVmListApi() {
    return http("display/task/vm-list", {
      method: "POST",
    });
  },
  /** cos url */
  getCosUrlApi(data) {
    return http("display/task/cos-avatar", {
      method: "POST",
      data
    });
  },
  /** 批量导入 */
  importListApi(data) {
    return http("display/task/review-file-check", {
      method: "POST",
      data
    });
  },
  /** 下载分配数据 */
  reviewExportApi(data) {
    return http("display/task/review-export", {
      method: "POST",
      data
    });
  },
}