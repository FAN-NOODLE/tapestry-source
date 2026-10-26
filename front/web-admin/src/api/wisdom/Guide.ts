import http from "@/utils/axios";

export default {
  // 导购列表
  wisdomGuideList(data: any) {
    return http("wisdom-guide/list", {
      method: "POST",
      data: { apiVersion: "old", ...data },
    });
  },
  // 生产海报
  workerPoster(data: any) {
    return http("wisdom-guide/generate-worker-poster", {
      method: "POST",
      data: { apiVersion: "old", ...data },
    });
  },
  // 导购新增
  wisdomGuideAdd(data: any) {
    return http("wisdom-guide/add", {
      method: "POST",
      data,
    });
  },
  // 编辑企业管理导购的管理范围
  wisdomGuideEdit(data: any) {
    return http("wisdom-guide/save-admin-scope", {
      method: "POST",
      data: { apiVersion: "old", ...data },
    });
  },
  // 企业管理导购下的门店列表
  wisdomGuideStoreList(data: any) {
    return http("wisdom-guide/admin-store-list", {
      method: "POST",
      data: { apiVersion: "old", ...data },
    });
  },
  // 企业管理导购下的部门列表
  wisdomGuideDpList(data: any) {
    return http("wisdom-guide/admin-dept-list", {
      method: "POST",
      data: { apiVersion: "old", ...data },
    });
  },
  // 导购导出
  wisdomGuideExport(data: any) {
    return http("wisdom-guide/export-list", {
      method: "POST",
      data,
    });
  },
  // 转化为企业管理/导购
  upgradeAdmin(data: any) {
    return http("wisdom-guide/upgrade-admin", {
      method: "POST",
      data: { apiVersion: "old", ...data },
    });
  },
  // 手动同步
  batchAddGuide() {
    return http("wisdom-guide/batch-add-guide", {
      method: "POST",
       data: { apiVersion: "old"},
    });
  },
  // 导购设置详情
  achievemenDetail(data: any) {
    return http("guide-achievement-setting/detail", {
      method: "POST",
      data,
    });
  },
  // 导购设置保存
  achievemenSave(data: any) {
    return http(`guide-achievement-setting/save`, {
      method: "POST",
      data,
    });
  },
  // 支付即导购列表
  guidePayList(data: any) {
    return http(`user/guide_list`, {
      method: "POST",
      data,
    });
  },
  // 添加支付即导购
  guideManageAdd(data: any) {
    return http(`guide-manage/add`, {
      method: "POST",
      data,
    });
  },
  // 删除支付即导购
  guideManageDelete(data: any) {
    return http(`guide-manage/delete`, {
      method: "POST",
      data,
    });
  },
  // 开启关闭支付即导购
  changeAllocationStatus(data: any) {
    return http(`user/change_allocation_status`, {
      method: "POST",
      data,
    });
  },
  // 获取门店区域
  getArea(data: any) {
    return http(`store-area/get-area`, {
      method: "POST",
      data,
    });
  },
  // 通过区域ID获取门店列表
  getStoreListApi(data: any) {
    return http("store-area/get-store", {
      method: "POST",
      data: { apiVersion: "old", ...data },
    });
  },

  // 根据门店或部门查询导购列表
  guideListByStore(data: any) {
    return http("user/guide-list-by-store", {
      method: "POST",
      data: { apiVersion: "old", ...data },
    });
  },
  // 获取导购客户列表
  getStaffBindList(data: any) {
    return http("staff-right/staff-bind-list", {
      method: "POST",
      data,
    });
  },
  // 业绩关系继承保存
  saveInherit(data: any) {
    return http("staff-right/inher-staff", {
      method: "POST",
      data,
    });
  },
  // 业绩关系部分继承
  partialInherit(data: any) {
    return http("staff-right/partial-inherit", {
      method: "POST",
      data,
    });
  },
  // 业绩关系全部继承
  allInherit(data: any) {
    return http("staff-right/all-inherit", {
      method: "POST",
      data,
    });
  },
  // BA业绩指标导入模板下载
  baPerformanceDownload(data: any) {
    return http("ba-performance/download", {
      method: "GET",
      data,
      responseType: "blob",
    });
  },
  // BA业绩指标导入模板下载
  baPerformanceimport(data: any) {
    return http("ba-performance/import", {
      method: "POST",
      data,
      responseType: "blob",
    });
  },
  // BA业绩指标导出
  baPerformanceExport(data: any) {
    return http("ba-performance/export", {
      method: "POST",
      data,
      responseType: "blob",
    });
  },
  // BA业绩指标列表
  baPerformanceList(data: any) {
    return http("ba-performance/list", {
      method: "POST",
      data,
    });
  },
  // 导购丢失客户推送消息详情
  userLostInfo(data: any) {
    return http("wisdom-guide/user-lost-info", {
      method: "POST",
      data,
    });
  },
  // 导购丢失客户推送消息配置
  userLostConfig(data: any) {
    return http("wisdom-guide/user-lost-config", {
      method: "POST",
      data,
    });
  },
};
