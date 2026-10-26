import http from "@/utils/axios";

export default {
  // 全员分销保存接口
  updateDistributionListApi(data: any) {
    return http("distribution-config/save", {
      method: "POST",
      data
    });
  },
  // 店铺列表接口
  getFieldList(data: any) {
    return http("field/list", {
      method: "POST",
      data
    });
  },
  // 分销查看上部信息接口（分销描述/海报图片等）
  getDistributionDetail(data: any) {
    return http("distribution-config/detail", {
      method: "POST",
      data
    });
  },
  // 分销商品列表接口
  getDistributionGoodsList(data: any) {
    return http("distribution-goods/goods-list", {
      method: "POST",
      data
    });
  },
  // 单个修改分销商品
  updateDistributionGoods(data: any) {
    return http("distribution-goods/save", {
      method: "POST",
      data
    });
  },
  // 批量修改分销商品的推广佣金比例与服务佣金比例
  batchSettingDistributionGoods(data: any) {
    return http("distribution-goods/batch-setting", {
      method: "POST",
      data
    });
  },
  // 批量添加分销商品
  batchSaveDistributionGoods(data: any) {
    return http("distribution-goods/batch-save", {
      method: "POST",
      data
    });
  },
  // 批量删除分销商品
  batchDeleteDistributionGoods(data: any) {
    return http("distribution-goods/batch-delete", {
      method: "POST",
      data
    });
  },
  // 最新获取sku商品列表
  getSkuGoodList(data: any) {
    return http("goods/sku-pop-window", {
      method: "POST",
      data
    });
  },
  // 分销商品添加重复校验接口
  repeatInspect(data: any) {
    return http("distribution-goods/repeat-inspect", {
      method: "POST",
      data
    });
  },
  // 移除/更新商品列表
  modifyDistributionList(data: any) {
    return http("worker-distribution/distribution-goods-modify", {
      method: "POST",
      data
    });
  },
  // 下载导入模板
  downloadTemplate(data: any) {
    return http("worker-distribution/download-distribution-goods-template", {
      method: "GET",
      params: data,
      responseType: "blob"
    });
  },
  // 下载失败列表
  downloadFail(data: any) {
    return http("worker-distribution/download-distribution-goods-error", {
      method: "GET",
      params: data,
      responseType: "blob"
    });
  }
  // // 新增分销商品信息
  // addGoods(data: any) {
  //   return http("worker-distribution/distribution-goods-modify", {
  //     method: "POST",
  //     data
  //   });
  // }
};
