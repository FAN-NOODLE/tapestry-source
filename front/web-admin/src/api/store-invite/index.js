import http from "@/utils/axios";

export default {
  // 预约列表
  orderListApi(data) {
    return http("reserve/order/lists", {
      method: "POST",
      data
    });
  },
  // 企业配置
  getEpConfig(data) {
    return http("enterprise/config", {
      method: "POST",
      data
    });
  },
  // 服务列表 reserve/service/list
  serviceListApi(data) {
    return http("reserve/service/list", {
      method: "POST",
      data
    });
  },
  // 创建服务
  createServiceApi(data) {
    return http("reserve/service/create", {
      method: "POST",
      data
    });
  },
  // 编辑服务
  updateServiceApi(data) {
    return http("reserve/service/update", {
      method: "POST",
      data
    });
  },
  // 删除服务
  deleteServiceApi(data) {
    return http("reserve/service/delete", {
      method: "POST",
      data
    });
  },
  // 统计
  staticDataApi(data) {
    return http("reserve/order/reserve-order-list", {
      method: "POST",
      data
    });
  },
  // 用户评论列表
  userCommonApi(data) {
    return http("reserve/service-score/list", {
      method: "POST",
      data
    });
  },
  // 导出预约单
  serverExportApi(data) {
    return http("reserve/order/export-lists-new", {
      method: "POST",
      data: data,
      timeout: 60 * 60 * 1000
    });
  },
  // 导出评论
  commentExportApi(data) {
    return http("reserve/service-score/export-list", {
      method: "POST",
      data: data,
      responseType: "arraybuffer"
    });
  },
  dataExportApi(data) {
    return http("reserve/order/export-reserve-order-list-new", {
      method: "POST",
      data: data
    });
  },
  dataOrderExportApi(data) {
    return http("reserve/order/export-order-record", {
      method: "POST",
      data: data
    });
  },
  // 预约单详情
  orderDetailApi(data) {
    return http("reserve/order/detail", {
      method: "POST",
      data
    });
  },
  // 获取服务是否可以删除
  checkDelete(data) {
    return http("reserve/service/check-delete", {
      method: "POST",
      data
    });
  },
  // 服务详情
  serviceDetail(data) {
    return http("/reserve/service/service-info", {
      method: "POST",
      data
    });
  }
};
