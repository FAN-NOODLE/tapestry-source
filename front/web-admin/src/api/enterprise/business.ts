import http from "@/utils/axios";

export default {
  downloadFileApi(data: any) {
    return http("store/export-sftp-error-record-list", {
      method: "POST",
      data: { apiVersion: "old", ...data },
      responseType: "blob"
    });
  },
  // 同步日志列表
  syncListApi(data: any) {
    return http("store/sftp-log-list", {
      method: "POST",
      data: { apiVersion: "old", ...data }
    });
  },
  // 门店列表
  getStoreListApi(data: any) {
    return http("store/department-list", {
      method: "POST",
      data: { apiVersion: "old", ...data }
    });
  },
  // 编辑门店
  editStoreApi(data: any) {
    return http("store/save-department-info", {
      method: "POST",
      data: { apiVersion: "old", ...data }
    });
  },
  // 删除门店
  delStoreApi(data: any) {
    return http("store/delete-department", {
      method: "POST",
      data: { apiVersion: "old", ...data }
    });
  },
  // 同步微信
  syncWechatApi(data: any) {
    return http("store/sync-data-to-wechat", {
      method: "POST",
      data: { apiVersion: "old", ...data }
    });
  },
  // 员工分配门店
  staffMoveStoreApi(data: any) {
    return http("store/move-employee-department", {
      method: "POST",
      data: { apiVersion: "old", ...data }
    });
  },
  // 获取员工列表
  getStaffListApi(data: any) {
    return http("store/employee-list", {
      method: "POST",
      data: { apiVersion: "old", ...data }
    });
  },
  // 企业配置
  getEpConfig() {
    return http("enterprise/config", {
      method: "POST"
    });
  },
  // 更新门店数据 / 获取上次更新时间
  getStoreSyncData(data: any) {
    return http("store/sync-data-by-departments", {
      method: "POST",
      data: { apiVersion: "old", ...data }
    });
  },
  // 下载门店列表批量更新文件模板
  downloadStoreList(data: any) {
    return http("store/download-store-list", {
      method: "POST",
      responseType: "blob",
      data: { ...data, ...{ apiVersion: "old" } }
    });
  },
  // 下载门店列表批量更新文件模板
  storeImportInfo(data: any) {
    return http("store/import-store-info", {
      method: "POST",
      // responseType: "blob",
      data: { ...data, ...{ apiVersion: "old" } }
    });
  },
  // 关闭门店
  closeStoreApi(data: any) {
    return http("store/close-or-open-department", {
      method: "POST",
      // responseType: "blob",
      data: { ...data, ...{ apiVersion: "old" } }
    });
  }
};
