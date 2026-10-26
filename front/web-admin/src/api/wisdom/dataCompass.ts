import http from "@/utils/axios";

export default {
  // 订单数据概览
  orderStatistical(data: any) {
    return http("order/order-statistical", {
      method: "POST",
      data,
    });
  },
  // 订单统计列表按日期查询
  orderStatisticalDate(data: any) {
    return http("order/order-statistical-by-date", {
      method: "POST",
      data,
    });
  },
  // 订单统计列表按导购查询
  orderStatisticalGuide(data: any) {
    return http("order/order-statistical-by-guide", {
      method: "POST",
      data,
    });
  },
  // 订单数据导出
  orderStatisticalExport(data: any) {
    return http("order/order-statistical-export", {
      method: "POST",
      data,
      responseType: "blob",
    });
  },
  // 订单明细导出
  orderDetailExport(data: any) {
    return http("order/order-statistical-detail-export", {
      method: "POST",
      data,
      responseType: "blob",
    });
  },
  // 分销数据概览
  taskStatisticalChart(data: any) {
    return http("distribution-statistics/statistic-console", {
      method: "POST",
      data,
    });
  },
  // 分销统计列表按日期查询
  taskStatisticalDate(data: any) {
    return http("distribution-statistics/statistic-date", {
      method: "POST",
      data,
    });
  },
  // 分销统计列表按导购查询
  taskStatisticalWorker(data: any) {
    return http("distribution-statistics/statistic-guide", {
      method: "POST",
      data,
    });
  },
  // 分销统计列表导出
  taskStatisticalEXPORT(data: any) {
    return http("distribution-statistics/statistic-export", {
      method: "POST",
      data,
      responseType: "blob",
    });
  },
  // 顾客数据概览
  statisticGuideCustomer(data: any) {
    return http("customer/statistic-guide-customer", {
      method: "POST",
      data,
    });
  },
  // 顾客数据按日期数据列表
  guideCustomerDate(data: any) {
    return http("customer/statistic-guide-customer-date", {
      method: "POST",
      data,
    });
  },
  // 顾客数据按导购列表
  guideCustomerList(data: any) {
    return http("customer/statistic-guide-customer-list", {
      method: "POST",
      data,
    });
  },
  // 顾客数据导出
  guideCustomerExport(data: any) {
    return http("customer/statistic-guide-customer-export", {
      method: "POST",
      data,
      responseType: "blob",
    });
  },
  // 获取提现设置
  getSettingDetail(data?: any) {
    return http("withdrawal-setting/detail", {
      method: "POST",
      data,
    });
  },
  // 获取项目名称
  getNameByCode(data: any) {
    return http("withdrawal-setting/check", {
      method: "POST",
      data,
    });
  },
  // 保存提现设置
  saveSetting(data: any) {
    return http("withdrawal-setting/add", {
      method: "POST",
      data,
    });
  },
  // 绑客明细导出
  customerDetailExport(data: any) {
    return http("data_compass/export_binding_customer_detail", {
      method: "POST",
      data,
      responseType: "blob",
    });
  },
  // BA部门导购
  baSourceDepte(data: any) {
    return http("ba_kanban/ba_source_depte", {
      method: "POST",
      data: { ...data },
    });
  },
  // BA业绩看板
  baKanBan(data: any) {
    return http("ba_kanban/v1/console/sales", {
      method: "POST",
      data: { ...data },
    });
  },
  // BA业绩看板-销售数据和运营数据
  baKanBanSalesOrOperation(data: any) {
    return http("ba_kanban/v1/operation", {
      method: "POST",
      data: { ...data },
    });
  },
  // BA业绩看板-销售数据导出
  baKanBanSalesExport(data: any) {
    return http("ba_kanban/sales-export", {
      method: "POST",
      responseType: "blob",
      data: { ...data },
    });
  },
  // BA业绩看板-获取同比增长
  baKanBanGrowthRate(data: any) {
    return http("ba_kanban/growth_rate", {
      method: "POST",
      data: { ...data },
    });
  },
  // BA业绩看板-运营数据导出
  baKanBanOperationExport(data: any) {
    return http("ba_kanban/ba/sale_export", {
      method: "POST",
      responseType: "blob",
      data: { ...data },
    });
  },
  // 社群、社群活动、好友数据汇总
  communityDataSummary(data: any) {
    return http("datacompass/community_data_summary", {
      method: "POST",
      data: { ...data },
    });
  },
  // 社群数据列表明细
  communityData(data: any) {
    return http("datacompass/community_data", {
      method: "POST",
      data: { ...data },
    });
  },
  // 社群数据列表导出
  communityDataExport(data: any) {
    return http("datacompass/community_data/export", {
      method: "POST",
      responseType: "blob",
      data: { ...data },
    });
  },
  // 好友数据列表明细
  friendData(data: any) {
    return http("datacompass/friend_data", {
      method: "POST",
      data: { ...data },
    });
  },
  // 好友数据列表导出
  friendDataExport(data: any) {
    return http("datacompass/friend_data/export", {
      method: "POST",
      responseType: "blob",
      data: { ...data },
    });
  },
};
