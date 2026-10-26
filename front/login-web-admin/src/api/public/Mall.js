/* eslint-disable no-dupe-keys */
import Vue from "vue";

const http = Vue.axios;

let apis = {
  // 查看佣金设置接口
  commissionDetailApi(data) {
    return http("commission/detail", {
      method: "GET",
      params: data
    });
  },
  // 保存佣金设置接口
  commissionSaveApi(data) {
    return http("commission/save", {
      method: "POST",
      data: data
    });
  },
  // 会员等级列表
  levelListApi(data) {
    return http("user-level/list", {
      method: "GET",
      params: data
    });
  },
  // 会员等级详情
  levelDetailApi(data) {
    return http("user-level/detail", {
      method: "GET",
      params: data
    });
  },
  // 变更等级启用状态
  levelSetStatusApi(data) {
    return http("user-level/set-status", {
      method: "POST",
      data: data
    });
  },
  // 编辑会员等级
  levelSetUpdateApi(data) {
    return http("user-level/update", {
      method: "POST",
      data: data
    });
  },
  // 平台管理员 - 平台管理员列表
  customerListApi(data) {
    return http("bossAdmin/list", {
      method: "Get",
      params: data
    });
  },
  /************************* 商品*************************/
  // 商品分类 - 分类列表
  categoryListApis(data) {
    return http("category/list", {
      method: "POST",
      data: data
    });
  },
  // 商品分类 - 全部分类列表
  categoryRecursionListApi(data) {
    return http("category/recursion-list", {
      method: "POST",
      data: data
    });
  },
  // 商品分类-新增
  addCategoryApi(data) {
    return http("category/add", {
      method: "POST",
      data: data
    });
  },
  // 商品分类-编辑更新
  editCategoryApi(data) {
    return http("category/update", {
      method: "POST",
      data: data
    });
  },
  // 商品分类 - 详情
  categoryDetailApi(data) {
    return http(`category/${data.id}`, {
      method: "GET"
    });
  },
  // 商品分类 - 删除
  deleteCategoryApi(data) {
    return http("category/delete", {
      method: "POST",
      data: data
    });
  },
  /************************* 商品*************************/
  // 商品 - 上下架
  goodsSaleApi(data) {
    return http("goods/sale", {
      method: "POST",
      data: data
    });
  },
  //商品定时上下架
  goodsDelaySaleApi(data) {
    return http("goods/delay-sale", {
      method: "POST",
      data: data
    });
  },
  //商品上下架检查
  goodsSaleCheck(data) {
    return http("goods/sale-check", {
      method: "POST",
      data: data
    });
  },
  //待上下架撤销
  goodsCancelApi(data) {
    return http("goods/cancel-sale", {
      method: "POST",
      data: data
    });
  },
  // 商品 - 列表
  goodsListApi(data) {
    return http("goods", {
      method: "POST",
      data: data
    });
  },
  // 商品 - 新增
  addGoodsApi(data) {
    return http("goods/save", {
      method: "POST",
      data: data
    });
  },
  // 商品 - 编辑更新
  editGoodsApi(data) {
    return http("goods/update", {
      method: "POST",
      data: data
    });
  },
  // 商品 - 详情
  goodsDetailApi(data) {
    return http("goods/detail", {
      method: "POST",
      data: data
    });
  },
  // 商品 - 删除
  deleteGoodsApi(data) {
    return http("goods/delete", {
      method: "POST",
      data: data
    });
  },
  /************************* 订单 *************************/
  //订单导出
  exportOrderApi(data) {
    return http("order/order-export", {
      method: "GET",
      params: data
    });
  },
  // 订单管理 - 订单列表
  orderListApi(data) {
    return http("order/list", {
      method: "POST",
      data: data
    });
  },
  // 订单管理 - 订单删除
  orderDeleteApi(id) {
    return http(`order/${id}`, {
      method: "DELETE"
    });
  },
  // 订单管理 - 订单详情
  orderDetailApi(id, refund_id) {
    if (refund_id) {
      return http(`order/${id}/${refund_id}`, {
        method: "GET"
      });
    } else {
      return http(`order/${id}`, {
        method: "GET"
      });
    }
  },
  // 拼团订单 -- 拼团详情
  orderUsersApi(data) {
    return http("order/user-detail", {
      method: "POST",
      data: data
    });
  },
  // 订单管理 - 订单售后详情
  orderAfterSaleDetailApi(data) {
    return http("new-order/refund-detail", {
      method: "POST",
      data: data
    });
  },
  // 订单管理 - 订单售后审核状态修改
  orderAfterSaleStatusApi(data) {
    return http("new-order/update", {
      method: "POST",
      data: data
    });
  },
  // 订单管理-订单物流信息
  orderLogisticDetail(data) {
    return http("logistic/detail", {
      method: "POST",
      data: data
    });
  },
  /************************* 规格属性 *************************/
  // 商品管理 - 规格分类列表
  categoryListApi(data) {
    return http("attributes/category-list", {
      method: "POST",
      data: data
    });
  },
  // 商品管理 - 属性列表
  attributeListApi(data) {
    return http("attributes/list", {
      method: "POST",
      data: data
    });
  },
  //根据分类查询属性
  attributeListCateApi(data) {
    return http("attributes/list", {
      method: "POST",
      data: data
    });
  },
  // 商品管理 - 删除属性
  attributeDelApi(data) {
    return http("attributes/del", {
      method: "POST",
      data: data
    });
  },
  // 商品管理 - 新增修改属性
  attributeSaveApi(data) {
    return http("attributes/save", {
      method: "POST",
      data: data
    });
  },
  /**************优惠券***************/
  // 礼券列表
  couponListApi(data) {
    return http("coupon/list", {
      method: "GET",
      params: data
    });
  },
  // 礼券删除
  couponDeleteApi(data) {
    return http("coupon/delete", {
      method: "POST",
      data: data
    });
  },
  // 礼券添加
  couponAddApi(data) {
    return http("coupon/add", {
      method: "POST",
      data: data
    });
  },
  // 礼券详情
  couponDetailApi(data) {
    return http("coupon/detail", {
      method: "GET",
      params: data
    });
  },
  // 礼券详情
  couponEditApi(data) {
    return http("coupon/save", {
      method: "POST",
      data: data
    });
  },
  /******************************* 页面管理 ****************************/
  // 首页引导区列表
  indexListApi(data) {
    return http("page/guide/list", {
      method: "GET",
      data: data
    });
  },
  // 首页引导区保存
  saveIndexApi(data) {
    return http("page/guide/save", {
      method: "POST",
      data: data
    });
  },
  // 轮播图列表
  swipperListApi(data) {
    return http("page/banner/list", {
      method: "GET",
      data: data
    });
  },
  // 保存轮播图列表
  saveSwipperApi(data) {
    return http("page/banner/save", {
      method: "POST",
      data: data
    });
  },
  // 爆品列表
  hotListApi(data) {
    return http("page/hot/list", {
      method: "GET",
      data: data
    });
  },
  // 保存爆品
  saveHotApi(data) {
    return http("page/hot/save", {
      method: "POST",
      data: data
    });
  },
  // 猜你喜欢列表
  likeListApi(data) {
    return http("page/guess/list", {
      method: "GET",
      data: data
    });
  },
  // 保存猜你喜欢
  saveLikeApi(data) {
    return http("page/guess/save", {
      method: "POST",
      data: data
    });
  },
  /**
   * 店长推荐、爆品、喜欢、主推列表
   * type 0：爆品推荐，1：店长推荐，2：猜你喜欢，3：主推商品
   */
  recommendListApi(data) {
    return http("page/recommend/list", {
      method: "GET",
      params: data
    });
  },
  // 保存店长推荐
  saveRecommendApi(data) {
    return http("page/recommend/save", {
      method: "POST",
      data: data
    });
  },
  //*******************商品选中后不显示已选中接口 ******************/
  // 新人优惠商品列表
  getNewGoodsApi(data) {
    return http("app-center/novice/goods", {
      method: "GET",
      params: data
    });
  },
  // 拼团商品列表
  getGroupGoodsApi(data) {
    return http("app-center/group/goods", {
      method: "GET",
      params: data
    });
  },
  // 爆品商品列表
  getHotGoodsApi(data) {
    return http("page/hot/goods", {
      method: "POST",
      data: data
    });
  },
  // 限时购商品列表
  getLimitGoodsApi(data) {
    return http("app-center/seckill/goods", {
      method: "POST",
      data: data
    });
  },
  // 超低购活动新增
  cheapAddApi(data) {
    return http("app-center/cheap/add", {
      method: "POST",
      data: data
    });
  },
  // 未参加活动商品接口
  getCouponSpuApi(data) {
    return http("app-center/coupon-spu-list", {
      method: "GET",
      params: data
    });
  },
  // 超低购详情
  getCheapDetailApi(data) {
    return http("app-center/cheap/detail", {
      method: "GET",
      params: data
    });
  },
  // 超低购活动列表
  getCheapListApi(data) {
    return http("app-center/cheap/list", {
      method: "GET",
      params: data
    });
  },
  // 超低购活动开启禁用接口
  enableCheapApi(data) {
    return http("app-center/cheap/enable", {
      method: "POST",
      data: data
    });
  },
  // 超低购活动删除接口
  deleteCheapApi(data) {
    return http("app-center/cheap/delete", {
      method: "POST",
      data: data
    });
  },
  // 超低购列表排序接口
  cheapSortApi(data) {
    return http("app-center/cheap/sort", {
      method: "POST",
      data: data
    });
  },
  // 新增页底 --目前相当于保存
  addPageBottomApi(data) {
    return http("page/footer/add", {
      method: "POST",
      data: data
    });
  },
  // 编辑页底
  editPageBottomApi(data) {
    return http("page/footer/update", {
      method: "POST",
      data: data
    });
  },
  // 获取页底详情
  detailPageBottomApi(data) {
    return http("page/footer/detail", {
      method: "GET",
      params: data
    });
  },
  // 微信优惠券列表
  getWxCouponListApi(data) {
    return http("app-center/wx-coupon/list", {
      method: "GET",
      params: data
    });
  },
  // 微信优惠券详情
  getWxCouponDetailApi(data) {
    return http("app-center/wx-coupon/detail", {
      method: "GET",
      params: data
    });
  },
  // 新增微信优惠券
  addWxCouponApi(data) {
    return http("app-center/wx-coupon/add", {
      method: "POST",
      data: data
    });
  },
  // 编辑微信优惠券
  editWxCouponApi(data) {
    return http("app-center/wx-coupon/edit", {
      method: "POST",
      data: data
    });
  },
  // 改变微信优惠券状态
  wxCouponStatusApi(data) {
    return http("app-center/wx-coupon/changeStatus", {
      method: "POST",
      data: data
    });
  },
  // 删除微信优惠券
  deleteWxCouponApi(data) {
    return http("app-center/wx-coupon/del", {
      method: "POST",
      data: data
    });
  },
  // 占位符选择礼劵列表
  templateCouponList(data) {
    return http("template/home/coupon-list", {
      method: "GET",
      params: data
    });
  }
};

export default apis;
