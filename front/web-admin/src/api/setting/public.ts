import http from "@/utils/axios";
export default {
  //设置保存接口
  accountList(data: any) {
    return http("config/save", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 设置详情接口
  configInfo(data: any) {
    return http("config/info", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 设置详情接口
  goodsSales(data: any) {
    return http("goods/goods-sales", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 浮动销量保存接口
  floatSave(data: any) {
    return http("goods/goods-sales-save", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 浮动销量详情接口
  floatInfo(data: any) {
    return http("goods/goods-sales-info", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 保存商户信息
  merchantSave(data: any) {
    return http("setting/merchant_setting/save", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 获取商户信息
  merchantDetail(data?: any) {
    return http("setting/merchant_setting/detail", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  // 获取商品销量
  spuSales(data: any) {
    return http("goods/spu-sales", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
};
