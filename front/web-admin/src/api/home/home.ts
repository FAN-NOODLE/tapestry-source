import http from "@/utils/axios";

export default {
  //导购业绩排行榜
  getShopers(data: any) {
    return http("index/shoppers/list", {
      method: "POST",
      data: data,
    });
  },
  //基本配置检查
  getNormalCheck(data: any) {
    return http("index/base/config", {
      method: "POST",
      data: data,
    });
  },
  //商城转化
  getShopConversion(data: any) {
    return http("index/shop/conversion", {
      method: "POST",
      data: data,
    });
  },
  // 私域数据
  getQwFriend(data: any) {
    return http("index/qw/friend", {
      method: "POST",
      data: data,
    });
  },
  //商城变现
  getShopFee(data: any) {
    return http("index/shop/fee", {
      method: "POST",
      data: data,
    });
  },
  //门店列表
  getShopList(data: any) {
    return http("index/store/list", {
      method: "POST",
      data: data,
    });
  },
  // 首页企业续费通知
  buyNotice(data?: any) {
    return http("index/buy/notice", {
      method: "POST",
      data,
    });
  },
  // 获取导购业绩排行榜
  customerRankData(data?: any) {
    return http("customer/rank-data", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 获取版本信息
  currentVersionInfo(data?: any) {
    return http("version/current-version-info", {
      method: "POST",
      data,
    });
  },
  // 导出导购排行榜
  shopperListExport(data: any) {
    return http("customer/rank-data-export", {
      method: "POST",
      responseType: "blob",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
};
