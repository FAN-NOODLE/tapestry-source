// 社群用户列表接口入参类型
interface IGroupChatMemberApiParams extends IMeta {
  group_chat_id: string; // 社群Id
  type: number; // 排序方式（1：默认成员列表，2：活跃排行榜，3：成交排行榜）
}
// 用户添加的员工列表接口入参类型
interface IUserMemberApiParams extends IMeta {
  external_userid: string; // 外部联系人Id
}
// 客户备注信息保存接口入参类型
interface ICustomerRemarkSaveParams {
  user_id: number; // 用户ID
  age?: string | number; // 年龄
  birthday?: number | string; // 生日
  gender?: number | string; // 性别：0-未知 1-男性 2-女性
  mobile?: number | string; // 电话
  email?: string; // 邮箱
  qq?: string; // QQ号
  describe?: string; // 描述
}
// 购物订单列表接口入参类型
interface ICustomerOrderListParams extends IMeta {
  user_id: number; // 用户ID
  order_no?: string | number; // 订单号
  status?: number; // 订单状态：1 待付款，2 已完成，3 已取消, 4 订单确认中
}
// 用户优惠券列表接口、添加的员工列表、用户足迹列表基础入参类型
interface ICustomerBaseParams extends IMeta {
  user_id: number; // 用户ID
}
// 加入社群的列表列表入参类型
interface ICustomerrGroupListParams extends IMeta {
  user_id: number; // 用户user_id
}
// 优惠券可使用商品列表接口入参类型
interface ICustomerCouponGoodsListParams extends IMeta {
  coupon_id: number; // 优惠券ID
}
// 打标签接口入参类型
interface ISetTagParams extends IMeta {
  relation_id: number;
  add_tag?: Array<string>; // 新增的标签IDS（不包含现有的）
  remove_tag?: Array<string>; // 删除的标签IDS
}
// 获取订单列表
interface GetOrderParams extends IMeta {
  union_id: string
}
// 获取订单详情
interface GetOrderDetailParams extends IMeta {
  union_id: string,
  trade_no: string
}
// 获取订单详情
interface GetOrderDetail extends IMeta {
  union_id: string,
  tradeNo: string,
  goods: Array<string>,
  list: Array<string>,
  tradePrice: string | number
}

// 获取订单详情
interface GetJoinActivityParams extends IMeta {
  page: number;
  size: number;
  union_id: string | number;
}

// 获取积分详情
interface GetIntegralParams extends IMeta {
  union_id:  string | number;
  page_no: string;
  limit: string;
}
// 获取积分详情
interface GetCouponParams extends IMeta {
  consumerID:  string;
  transType: string;
}
// 社群画像列表
interface IChatIndexParams extends IMeta {
  chat_id: string;
  sts: number;	
  ets: number;
}
// 群成员列表/未入群用户/会员列表/非会员列表（分页列表）
interface IChatUserParams extends IMeta {
  chat_id: string;	
  list_type: number;	
  friend: number;
  page: number;
  size: number;
  ets: number;
}
// 群成员参与活动列表
interface IChatActsParams extends IMeta {
  chat_id: string;
}
// 群成员列表
interface IChatGroupMembers extends IMeta {
  chat_id: string;
  page: number;
  size: number;
  friend: number;
  list_type: number;
}
// 查看用户活动
interface IChatGroupActInfo extends IMeta {
  chat_id: string;
  user_id: string;
}