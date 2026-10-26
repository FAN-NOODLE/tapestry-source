import http from "@/utils/axios";
// 营销推广-商品
export const goodListApi = (data: IGoodlist) =>
  http("goods/sku-list", {
    method: "POST",
    data,
  });
// 营销推广-活动
export const activityListApi = (data: IMeta) =>
  http("activity/goods-list", {
    method: "POST",
    data,
  });
// 营销推广-微信优惠券
export const wxCouponListApi = (data: IWxCouponList) =>
  http("wx-coupon/list", {
    method: "POST",
    data,
  });
// 营销推广-优惠券部分可用商品
export const wxCouponGoodlistApi = (coupon_id: number) =>
  http("wx-coupon/goods-list", {
    method: "POST",
    data: { coupon_id, page: 1, size: 999 },
  });
// 营销推广-下载海报
export const downloadPoster = (data: IPosterParams) =>
  http("marketing-promotion/poster", {
    method: "POST",
    data,
  });
// 素材列表
export const sourceListApi = (data: ISourceList) =>
  http("source/source-list", {
    method: "POST",
    data,
  });
// 素材列表 （新）
export const selectListApi = (data: any) =>
  http("source/select-list", {
    method: "POST",
    data,
  });
//素材分组列表
export const sourceGroupListApi = () =>
  http("source/source-group-list", {
    method: "POST",
    data: {},
  });
// 根据门店获取素材分组列表
export const getListForStore = () =>
  http("source/group-relation/get-list-for-store", {
    method: "POST",
    data: {},
  });
// 产品类别
export const categoryList = (data: any) =>
  http("source/category/list", {
    method: "POST",
    data,
  });
//生成素材media_id
export const sourceMakeMediaApi = (data: IMakeMediaParams) =>
  http("source/make-media", {
    method: "POST",
    data,
  });
// 获取tab更新数量
export const getNumApi = () =>
  http("/source/get-new-num", {
    method: "POST",
  });
/**素材库上传权限查询 */
export const checkAddApi = () => {
  return http("/source/upload-auth", {
    method: "POST",
  });
};
/**获取cos 签名信息 */
export const getFileCredentialApi = () =>
  http("/source/credential", { method: "POST" });
/** 保存素材 */
export const saveMaterial = (data: ISaveMaterial) => {
  return http("/source/upload-media", { method: "POST", data });
};
/** 查询素材压缩结果 */
export const getResult = (data: {
  compress_id: number;
  source_type: number;
}) => {
  return http("/source/compress", { method: "POST", data });
};
// 判断是否是首次打开页面
export const getCheckFirst = (data) => {
  return http("/source/check-first", { method: "POST", data });
};
// 转发次数上报
export const reportTimes = (data) => {
  return http("/source/report-times", { method: "POST", data });
};
// 获取品牌素材
export const getBrandMaterialListApi = (data) => {
  return http("/source/brand/search-list", { method: "POST", data });
  // return http("/source/brand/collect/search-list", { method: "POST", data });
};
export const getBrandCollectMaterialListApi = (data) => {
  // return http("/source/brand/search-list", { method: "POST", data });
  return http("/source/brand/collect/search-list", { method: "POST", data });
};
// 获取品牌素材详情
export const getBrandMaterialDetailApi = (data) => {
  return http("/source/brand/get-cache-info", { method: "POST", data });
};

// 设置收藏
export const setCollectApi = (data) => {
  return http("/source/collect", { method: "POST", data });
};

// 获取收藏列表
export const getCollectListApi = (data) => {
  return http("/source/collect-select-list", { method: "POST", data });
};