import http from "@/utils/axios";
// 社群详情
export const groupChatInfoApi = (group_chat_id: string) =>
  http("group-chat/info", {
    method: "POST",
    data: { group_chat_id },
  });
// 社群用户列表
export const groupChatMemberApi = (data: IGroupChatMemberApiParams) =>
  http("group-chat/group-chat-member", {
    method: "POST",
    data,
  });
// 用户添加的员工列表
export const userMemberListApi = (data: IUserMemberApiParams) =>
  http("group-chat/group-users-add-employees", {
    method: "POST",
    data,
  });
// 客户收货地址
export const customerAddrApi = (data: ICustomerBaseParams) =>
  http("customer/address-list", {
    method: "POST",
    data,
  });
// 客户备注信息保存
export const customerRemarkSaveApi = (data: ICustomerRemarkSaveParams) =>
  http("customer/remark-save", {
    method: "POST",
    data,
  });
// 购物订单列表
export const customerOrderListApi = (data: ICustomerOrderListParams) =>
  http("customer/order-list", {
    method: "POST",
    data,
  });
// 购物订单详情
export const customerOrderDetailApi = (order_no: number) =>
  http("customer/order-detail", {
    method: "POST",
    data: { order_no },
  });
// 客户详情
export const customerDetailApi = (external_userid?: string) =>
  http("customer/detail", {
    method: "POST",
    data: { external_userid },
  });
// h5修改客户备注
export const updateRemark = (data: any) =>
  http("customer/update-remark", {
    method: "POST",
    data,
  });
// 用户优惠券列表
export const customerCouponListApi = (data: ICustomerBaseParams) =>
  http("customer/coupon-list", {
    method: "POST",
    data,
  });
// 用户足迹列表
export const customerFootprintListApi = (data: ICustomerBaseParams) =>
  http("customer/footprint-list", {
    method: "POST",
    data,
  });
// 用户足迹列表
export const customerCouponGoodsListApi = (
  data: ICustomerCouponGoodsListParams
) =>
  http("customer/coupon-goods-list", {
    method: "POST",
    data,
  });
// 加入社群的列表
export const customerGroupListApi = (data: ICustomerrGroupListParams) =>
  http("customer/group-list", {
    method: "POST",
    data,
  });
// 企业标签库
export const tagListApi = (relation_id: number) =>
  http("tag/list", {
    method: "POST",
    data: { relation_id },
  });
// 设置客户企业标签
export const setTagApi = (data: ISetTagParams) =>
  http("customer/set-external-tags", {
    method: "POST",
    data,
  });
// 获取订单列表 旧
export const getOrderApiOld = (data: GetOrderParams) =>
  http("customer/dl-order-list", {
    method: "POST",
    data,
  });
// 获取订单列表
export const getOrderApi = (data: any) =>
  http("customer/crm-order-list", {
    method: "POST",
    data,
  });
// 获取订单详情
export const getOrderDetail = (data: GetOrderDetailParams) =>
  http("customer/dl-order-goods", {
    method: "POST",
    data,
  });
// 获取参与活动
export const getJoinActivity = (data: GetJoinActivityParams) =>
  http("group-chat/view-join-activity", {
    method: "POST",
    data,
  });
// 客户积分明细列表
export const integralListApi = (data: GetIntegralParams) =>
  http("customer/dl-points-list", {
    method: "POST",
    data,
  });
// 客户卡券列表 旧
export const couponListApiOld = (data: GetCouponParams) =>
  http("customer/dl-coupons-list", {
    method: "POST",
    data,
  });
// 客户卡券列表
export const couponListApi = (data: any) =>
  http("customer/crm-coupons-list", {
    method: "POST",
    data,
  });
// 客户卡券列表
export const getSwCouponListApi = (data: any) =>
  http("customer/sw-user-crm-coupons-list", {
    method: "POST",
    data,
  });
// 操作印象标签关系
export const tagsImpressionOperate = (data: any) =>
  http("customer/tags/impression/operate", {
    method: "POST",
    data,
  });
// 社群画像列表
export const chatIndexList = (data: IChatIndexParams) =>
  http("chat-group/index_list", {
    method: "POST",
    data,
  });
// 群成员列表/未入群用户/会员列表/非会员列表（分页列表）
export const chatUserList = (data: IChatUserParams) =>
  http("chat-group/members_guide", {
    method: "POST",
    data,
  });
// 群成员参与活动列表
export const chatActsList = (data: IChatActsParams) =>
  http("chat-group/acts_list", {
    method: "POST",
    data,
  });
// 群成员列表
export const chatGroupMembers = (data: IChatGroupMembers) =>
  http("chat-group/members_info", {
    method: "POST",
    data,
  });
// 查看用户活动
export const chatGroupActInfo = (data: IChatGroupActInfo) =>
  http("/chat-group/activity-info", {
    method: "POST",
    data,
  });
// 获取数据更新时间
export const getRefreshTime = (data: any) =>
  http("/chat-group/data-update-time", {
    method: "POST",
    data,
  });
// 编辑客户地址、生日、鞋码
export const saveInfoAdd = (data: any) =>
  http("/customer/save-info-add", {
    method: "POST",
    data,
  });

// 顾客详情
export const getCustomerDetailApi = (data: any) =>
  http("/customer/external-user-detail", {
    method: "POST",
    data,
  });

// 操作顾客的印象标签
export const impressionTagsOperateApi = (data: any) =>
  http("/customer/tags/impression/operate", {
    method: "POST",
    data,
  });

// 顾客和导购的绑定关系
export const userBindRelationApi = (data: any) => {
  return http("/customer/user-bind-relation", {
    method: "POST",
    data,
  });
};

// 操作顾客的印象标签
export const userDetailTipsNumApi = (data: any) =>
  http("/customer/user-detail-tips-num", {
    method: "POST",
    data,
  });

// 获取用户权限
export const getUserAuthApi = (data: any) =>
  http("/member/guide-position", {
    method: "POST",
    data,
  });
