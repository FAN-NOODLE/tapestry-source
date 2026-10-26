import http from "@/utils/axios";

export default {
  /********************* 门店管理 *********************/
  getStoreList(data: any) {
    return http("store/list", {
      method: "POST",
      data: data
    });
  },
  storeDefault(data: any) {
    return http("store/default", {
      method: "POST",
      data
    });
  },
  // 启用禁用
  storeStatusUse(data: any) {
    return http("store/status", {
      method: "POST",
      data: data
    });
  },
  // 新增
  StoreAdd(data: any) {
    return http("store/add", {
      method: "POST",
      data: data
    });
  },
  // 下载批量导入门店模版
  uploadTemplate(data: any) {
    return http("store/export", {
      method: "POST",
      data: data,
      responseType: "blob"
    });
  },
  // 门店模板确认更新
  confirmImport(data: any) {
    return http("store/import", {
      method: "POST",
      data
    });
  },
  // 下载失败门店列表
  downloadFail(data: any) {
    return http("store/export-error", {
      method: "POST",
      data,
      responseType: "blob"
    });
  },
  // 编辑
  StoreEdit(data: any) {
    return http("store/edit", {
      method: "POST",
      data: data
    });
  },
  // 详情
  StoreDetail(data: any) {
    return http("store/detail", {
      method: "POST",
      data: data
    });
  },
  getSotrePhone(data: any) {
    return http("store/phone", {
      method: "POST",
      data: data
    });
  },
  // 邮费模板列表
  getFeelist(data: any) {
    return http("postage/list", {
      method: "POST",
      data: data
    });
  },
  /********************* 业绩管理 *********************/
  // 列表
  performanceList(data: any) {
    return http("store-performance/fetch", {
      method: "POST",
      data: data
    });
  },
  // 列表-导出
  performanceExport(data: any) {
    return http("store-performance/export", {
      method: "POST",
      data: data
    });
  },
  // 业绩-明细
  performanceDetail(data: any) {
    return http("store-performance/detail", {
      method: "POST",
      data: data
    });
  },
  /********************* 店员管理 *********************/
  // 员工-列表
  storeClerkList(data: any) {
    return http("store/worker/list", {
      method: "GET",
      params: data
    });
  },
  // 员工-编辑
  storeClerkEdit(data: any) {
    return http("store/worker/edit", {
      method: "POST",
      data: data
    });
  },
  // 员工-新增
  storeClerkAdd(data: any) {
    return http("store/worker/add", {
      method: "POST",
      data: data
    });
  },
  // 员工-复岗离岗
  storeClerkStatus(data: any) {
    return http("store/worker/edit-status", {
      method: "POST",
      data: data
    });
  },
  // 店员新增-海报
  storeClerkNewAdd(data: any) {
    return http("store/worker/add-worker-poster", {
      method: "POST",
      data: data
    });
  },
  /********************* 商品管理 *********************/
  // 商品-商品列表
  getSkuList(data: any) {
    return http("goods/sku-list", {
      method: "POST",
      data: data
    });
  },
  // 商品-批量更新数据模版下载
  goodsTemplateExport(data: any) {
    return http("goods/export-template", {
      method: "POST",
      data: data,
      responseType: "blob"
    });
  },
  // 商品-所属列表
  getTypeList(data: any) {
    return http("category-template/store-category-list", {
      method: "POST",
      data: data
    });
  },
  // 商品-失败模版下载
  downloadGoodsFail(data: any) {
    return http("goods/export-errors", {
      method: "POST",
      data: data,
      responseType: "blob"
    });
  },
  // 商品-确认更新
  confirmGoodsImport(data: any) {
    return http("goods/confirm-batch-update", {
      method: "POST",
      data: data
    });
  },
  // 商品-选择商品列表
  selectGoodsList(data: any) {
    return http("goods/shop-goods-list", {
      method: "POST",
      data: data
    });
  },
  // 商品-修改价格
  changePrice(data: any) {
    return http("goods/edit-price", {
      method: "POST",
      data: data
    });
  },
  // 商品-复制商品
  copyGoods(data: any) {
    return http("goods/copy-goods", {
      method: "POST",
      data: data
    });
  },
  // 商品-上下架
  onSale(data: any) {
    return http("goods/sale", {
      method: "POST",
      data: data
    });
  },
  // 商品-查询可用库存
  getStock(data: any) {
    return http("warehouse-stock/get-stock", {
      method: "POST",
      data: data
    });
  },
  getGoodsStock(data: any) {
    return http("material-channel/detail", {
      method: "POST",
      data: data
    });
  },
  // 商品-更新库存
  updateStock(data: any) {
    return http("material-channel/update-channel-stock", {
      method: "POST",
      data: data
    });
  },
  /********************* 门店装修 *********************/
  // 门店商品（SKU）列表
  getStoreSkuListApi(data: any) {
    return http("store/sku-list", {
      method: "POST",
      data
    });
  },
  // 门店分类列表
  getStoreCategoryListApi(data: any) {
    return http("store/category/list", {
      method: "GET",
      params: data
    });
  },
  // 门店首页商品（SKU）详情列表
  getIndexSkuListApi(data: any) {
    return http("store/index-sku-list", {
      method: "POST",
      data
    });
  },
  // 获取线上使用模板版本
  getOnlineTemplateApi(data: any) {
    return http("new-template/online-template-version", {
      method: "POST",
      data
    });
  },
  // 保存门店首页商品
  setIndexDecorateGoodsApi(data: any) {
    return http("store/decorate/index", {
      method: "POST",
      data
    });
  },
  // 搜索推广人/服务人
  memberSelect(data: any) {
    return http("store/worker/member-select", {
      method: "POST",
      data: data
    });
  },
  // 门店列表
  authStoreList(data: any) {
    return http("store/auth-store-list", {
      method: "POST",
      data: data
    });
  }
};
