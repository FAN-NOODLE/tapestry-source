import http from "@/utils/axios";

export default {
  // 网仓订单 - 列表
  getOnlineOrderListApi(data: any) {
    return http("online_warehouse/list", {
      method: "POST",
      data
    });
  },
  // 网仓订单 - 详情
  getOnlineOrderDetailApi(data: any) {
    return http("online_warehouse/info", {
      method: "POST",
      data
    });
  },
  // 网仓订单 - 提交
  submitOnlineOrderApi(data: any) {
    return http("online_warehouse/submit", {
      method: "POST",
      data
    });
  },
  // 网仓订单 - 编辑
  editOnlineOrderApi(data: any) {
    return http("online_warehouse/edit", {
      method: "POST",
      data
    });
  },
  // 网仓订单 - 待上传的门店列表
  getStoreListApi(data: any) {
    return http("online_warehouse/store_list", {
      method: "POST",
      data
    });
  },
  // 网仓订单 - 待上传的门店订单列表
  getStoreOrderListApi(data: any) {
    return http("online_warehouse/store_order_list", {
      method: "POST",
      data
    });
  },
  // 网仓订单 - 导出数据
  exportOrderListApi(data: any) {
    return http("online_warehouse/export", {
      method: "POST",
      data
    });
  },
  // 网仓订单 - 导入数据
  importOrderListApi(data: any) {
    return http("online_warehouse/import_order", {
      method: "POST",
      data
    });
  },
  // 网仓订单 - 上传转货单
  uploadTransferOrderApi(data: any) {
    return http("online_warehouse/batch_transfer", {
      method: "POST",
      data
    });
  },
  // 网仓订单 - 拒单
  refundOrderApi(data: any) {
    return http("online_warehouse/refund", {
      method: "POST",
      data
    });
  },
  // 网仓订单 - 获取门店列表
  getDepListApi(data?: any) {
    return http("online_warehouse/dep_list", {
      method: "POST",
      data
    });
  },
  // 网仓订单 - 获取带上传转货单模版
  getExportListApi(data?: any) {
    return http("online_warehouse/transfer_export", {
      method: "POST",
      data
    });
  },
  // 网仓订单 - 批量导入
  importDataApi(data?: any) {
    return http("online_warehouse/transfer_import", {
      method: "POST",
      data
    });
  },

}