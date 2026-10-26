import http from "@/utils/axios";

/**
 * 网仓订单 - 获取历史订单
 * @param data 
 * @returns 
 */
export const getOrderListApi = (data: any) => {
  return http("online_warehouse/list", {
    method: "POST",
    data,
  });
}
/**
 * 网仓订单 - 获取订单详情
 * @param data 
 * @returns 
 */
export const getOrderInfoApi = (data: any) => {
  return http("online_warehouse/info", {
    method: "POST",
    data,
  });
}
/**
 * 网仓订单 - 获取订单详情 无鉴权
 * @param data 
 * @returns 
 */
export const getOrderDetailApi = (data: any) => {
  return http("online_warehouse/customer/order", {
    method: "POST",
    data,
  });
}

/**
 * 网仓订单 - 创建订单
 * @param data 
 * @returns 
 */
export const createOrderApi = (data: any) => {
  return http("online_warehouse/create", {
    method: "POST",
    data,
  });
}

/**
 * 网仓订单 - 编辑订单 无鉴权
 * @param data 
 * @returns 
 */
export const updateOrderApi = (data: any) => {
  return http("online_warehouse/customer/edit", {
    method: "POST",
    data,
  });
}
/**
 * 网仓订单 - 获取库存
 * @param data 
 * @returns 
 */
export const getGoodsStockApi = (data: any) => {
  return http("online_warehouse/inventory/list", {
    method: "POST",
    data,
  });
}
/**
 * 网仓订单 - 获取状态数量
 * @param data 
 * @returns 
 */
export const getOrderNumApi = (data?: any) => {
  return http("online_warehouse/tips_num", {
    method: "POST",
    data,
  });
}
/**
 * 网仓订单 - 店长的sa排行榜
 * @param data 
 * @returns 
 */
export const getSaRankListApi = (data: any) => {
  return http("online_warehouse/sa-rank", {
    method: "POST",
    data,
  });
}
/**
 * 网仓订单 - 店长的门店排行榜
 * @param data 
 * @returns 
 */
export const getStoreRankListApi = (data: any) => {
  return http("online_warehouse/store-rank", {
    method: "POST",
    data,
  });
}

/**
 * 网仓订单 - 导出
 * @param data 
 * @returns 
 */
export const exportEmailApi = (data: any) => {
  return http("online_warehouse/export-email", {
    method: "POST",
    data,
  });
}


/**
 * 网仓订单 - 店长下的门店列表
 * @param data 
 * @returns 
 */
export const getStoreListApi = () => {
  return http("online_warehouse/dep-list", {
    method: "POST",
  });
}

/**
 * 网仓订单 - 店长下的导购列表
 * @param data 
 * @returns 
 */
export const getGuideListApi = (data:any) => {
  return http("online_warehouse/sa-list", {
    method: "POST",
    data
  });
}

/**
 * 网仓订单 - 获取离职导购
 * @param data 
 * @returns 
 */
export const getLeaveGuideListApi = (data:any) => {
  return http("online_warehouse/leave-sa", {
    method: "POST",
    data
  });
}

/**
 * 网仓订单 - 获取财月数据
 * @param data 
 * @returns 
 */
export const getMonthListApi = () => {
  return http("member/fiscal-month", {
    method: "POST",
  });
}
/**
 * 网仓订单 - 取消订单
 * @param data 
 * @returns 
 */
export const cancelOrderApi = (data:any) => {
  return http("online_warehouse/cancel", {
    method: "POST",
    data
  });
}
/**
 * 网仓订单 - 获取图片media
 * @param data 
 * @returns 
 */
export const getImageMediaApi = (data:any) => {
  return http("online_warehouse/poster-media", {
    method: "POST",
    data
  });
}
/**
 * 网仓订单 - 检测销售单号是否重复
 * @param data 
 * @returns 
 */
export const checkOrderRepeatApi = (data:any) => {
  return http("online_warehouse/check-order", {
    method: "POST",
    data
  });
}
/**
 * 网仓订单 - 数据看班
 * @param data 
 * @returns 
 */
export const getDashboardApi = (data:any) => {
  return http("online_warehouse/dashboard", {
    method: "POST",
    data
  });
}
/**
 * 网仓订单 - 商品排行榜
 * @param data 
 * @returns 
 */
export const getGoodsRankApi = (data:any) => {
  return http("online_warehouse/goods-rank", {
    method: "POST",
    data
  });
}