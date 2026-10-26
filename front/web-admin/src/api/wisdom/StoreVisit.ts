import http from "@/utils/axios";

/**
 * 巡店管理相关API
 */
export default {
  /**
   * 获取部门列表
   * @param data 请求参数
   * @returns Promise
   */
  getDeptNameList(data: any) {
    return http("store-visit/dept_name", {
      method: "POST",
      data
    });
  },
  
  /**
   * 获取巡店记录列表
   * @param data 请求参数
   * @returns Promise
   */
  getStoreVisitList(data: any) {
    return http("store-visit/list", {
      method: "POST",
      data
    });
  },
  
  /**
   * 关闭巡店
   * @param data 请求参数，包含id
   * @returns Promise
   */
  closeStoreVisit(data: any) {
    return http("store-visit/close", {
      method: "POST",
      data
    });
  },
  
  /**
   * 导出巡店列表
   * @param data 请求参数
   * @returns Promise
   */
  exportStoreVisit(data: any) {
    return http("store-visit/export", {
      method: "POST",
      data
    });
  },
  
  /**
   * 获取巡店报表数据
   * @param data 请求参数
   * @returns Promise
   */
  getStoreVisitDashboard(data: any) {
    return http("store-visit/dashboard", {
      method: "POST",
      data
    });
  },
  /**
   * 获取区域列表
   * @param data 请求参数
   * @returns Promise
     */
  getAreaList(data: any) {
    return http("dashboard/area/list", {
      method: "POST",
      data
    });
  },

    /**
   * 导出巡店报表
   * @param data 请求参数
   * @returns Promise
   */
    exportDashboard(data: any) {
      return http("store-visit/dashboard/export", {
        method: "POST",
        data
      });
    },
}
