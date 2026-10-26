import http from "@/utils/axios";

export default {
  getAuthApi() {
    // 获取同步状态
    return http("tag/auth-create", {
      method: "POST"
    });
  },
  // 印象标签组列表
  impressionGroupListApi(data: any) {
    // 获取同步状态
    return http("tag/impression/group/list", {
      method: "POST",
      data: data
    });
  },
  impressionTagListApi(data: any) {
    // 获取同步状态
    return http("tag/impression/tag/list", {
      method: "POST",
      data: data
    });
  },
  labelListApi(data: any) {
    // 获取同步状态
    return http("tag/list", {
      method: "POST",
      data: data
    });
  },
  downListApi(data: any) {
    return http("user/group-user-export", {
      method: "POST",
      responseType: "blob",
      data: data
    });
  },
  // 标签编辑
  tagUpdateApi(data: any) {
    return http("tag/update", {
      method: "POST",
      data: data
    });
  },
  // 单个标签组详情(标签列表)
  tagInfoApi(data: any) {
    return http("tag/tag-info", {
      method: "POST",
      data: data
    });
  },
  // 商品SPU列表
  tagGoodsSpuApi(data: any) {
    return http("tag/goods-spu", {
      method: "POST",
      data: data
    });
  },
  // 优惠券列表
  tagCouponApi(data: any) {
    return http("tag/coupon-list", {
      method: "POST",
      data: data
    });
  },
  // 手动标签-删除标签组
  tagGroupDeleteApi(data: any) {
    return http("tag/tag-group-delete", {
      method: "POST",
      data: data
    });
  },
  // 标签分组排序(自动标签/手动标签)
  tagRankApi(data: any) {
    return http("tag/rank", {
      method: "POST",
      data: data
    });
  },
  // 修改单个手动标签
  tagEditTagApi(data: any) {
    return http("tag/tag-update", {
      method: "POST",
      data: data
    });
  },
  // 新增手动标签
  tagCreateTagApi(data: any) {
    return http("tag/create-tag", {
      method: "POST",
      data: data
    });
  },
  // 新增手动标签组
  tagCreateGroupApi(data: any) {
    return http("tag/create-group", {
      method: "POST",
      data: data
    });
  },
  //客户分组
  groupListApi(data: any) {
    return http("user/group-list", {
      method: "POST",
      data: data
    });
  },
  //客户新建
  groupCreateApi(data: any) {
    return http("user/create-group", {
      method: "POST",
      data: data
    });
  },
  //预计算人数
  groupUsercountApi(data: any) {
    return http("tag/precalculate-usercount", {
      method: "POST",
      data: data
    });
  },
  //标签同步企业微信
  tagSyncQwApi(data: any) {
    return http("tag/sync-qw", {
      method: "POST",
      data: data
    });
  },
  //除子标签
  tagDeleteRowApi(data: any) {
    return http("tag/tag-delete", {
      method: "POST",
      data: data
    });
  },
  //导出标签
  tagExportApi(data: any) {
    return http("tag/export", {
      method: "POST",
      responseType: "blob",
      data: data
    });
  },
  //获取营销列表页面
  marketStrategyListApi(data: any) {
    return http("marketing-strategies/list", {
      method: "POST",
      data: data
    });
  },
  //营销策略详情
  marketStrategyDetailApi(data: any) {
    return http("marketing-strategies/detail", {
      method: "POST",
      data: data
    });
  },
  marketStrategyRecordEditSaveApi(data: any) {
    return http("marketing-strategies/update", {
      method: "POST",
      data: data
    });
  },
  //获取直播列表
  marketStrategyLiveListApi(data: any) {
    return http("marketing-strategies/lives", {
      method: "POST",
      data: data
    });
  },
  //获取策略人群列表
  marketStrategyGroupListApi() {
    return http("marketing-strategies/user-group-list", {
      method: "POST"
    });
  },
  //策略列表记录删除
  marketStrategyRecordDelApi(data: any) {
    return http("marketing-strategies/delete", {
      method: "POST",
      data: data
    });
  },
  //保存营销数据详情
  marketStrategySaveApi(data: any) {
    return http("marketing-strategies/store", {
      method: "POST",
      data: data
    });
  },
  //获取营销数据详情页表格
  marketStrategyDataListApi(data: any) {
    return http("marketing-strategies/data/detail", {
      method: "POST",
      data: data
    });
  },
  //导出营销数据详情页表格
  marketStrategyExportApi(data: any) {
    return http("marketing-strategies/detail/export", {
      method: "POST",
      data: data,
      responseType: "arraybuffer"
    });
  },
  //企业短信数据统计
  smsTemplateTotalApi(data: any) {
    return http("sms-template/total", {
      method: "POST",
      data: data
    });
  },
  //短信模版列表
  smsTemplateListApi(data: any) {
    return http("sms-template/list", {
      method: "POST",
      data: data
    });
  },
  //短信模板删除
  smsTemplateDeleteApi(data: any) {
    return http("sms-template/delete", {
      method: "POST",
      data: data
    });
  },
  //短信模板添加
  smsTemplateCreateApi(data: any) {
    return http("sms-template/create", {
      method: "POST",
      data: data
    });
  },
  //短信模板编辑
  smsTemplateEditApi(data: any) {
    return http("sms-template/edit", {
      method: "POST",
      data: data
    });
  },
  //短信模板提交审核
  smsTemplateArraignmentApi(data: any) {
    return http("sms-template/arraignment", {
      method: "POST",
      data: data
    });
  },
  //短信模版详情
  smsTemplateInfoApi(data: any) {
    return http("sms-template/info", {
      method: "POST",
      data: data
    });
  },
  //场景值列表
  sceneListApi(data: any) {
    return http("wx-channel/wx-channel-list", {
      method: "POST",
      data: data
    });
  },
  //带参渠道列表
  channelListApi(data: any) {
    return http("marketing-strategies/get-chan-list", {
      method: "POST",
      data: data
    });
  },
  //可用优惠券列表
  couponListApi(data: any) {
    return http("coupon/wx-coupon-available-list", {
      method: "POST",
      data: data
    });
  },
  // 自动发券数据列表
  couponDataApi(data: any) {
    return http("marketing-strategies/detail-coupon", {
      method: "POST",
      data: data
    });
  },
  // 下载自动发券数据列表
  exportCouponApi(data: any) {
    return http("marketing-strategies/detail/export-coupon", {
      method: "POST",
      data: data,
      responseType: "arraybuffer"
    });
  },
  //客户列表
  customerListApi(data: any) {
    return http("customer/user-list", {
      method: "POST",
      data: data
    });
  },
  //客户详情基本信息 + 收货地址
  customerDetailListApi(data: any) {
    return http("customer/user-detail", {
      method: "POST",
      data: data
    });
  },
  //客户详情足迹
  customerDetailFootprintListApi(data: any) {
    return http("customer/footprint-list", {
      method: "POST",
      data: data
    });
  },
  //客户详情订单列表
  customerDetailOrderListApi(data: any) {
    return http("customer/order-list", {
      method: "POST",
      data: data
    });
  },
  //商品分页
  customerDetailOrderIdApi(data: any) {
    return http("customer/order-goods", {
      method: "POST",
      data: data
    });
  },
  //客户详情所属群列表
  customerDetailBelongToChatsApi(data: any) {
    return http("customer/belong-to-chats", {
      method: "POST",
      data: data
    });
  },
  //客户详情所属员工列表
  customerDetailBelongToMembersApi(data: any) {
    return http("customer/belong-to-members", {
      method: "POST",
      data: data
    });
  },
  //参数列表
  smsParamsListApi(data: any) {
    return http("sms-template/params-list", {
      method: "POST",
      data: data
    });
  },
  //添加短信签名
  smsCreateSignApi(data: any) {
    return http("sms-template/create-sign", {
      method: "POST",
      data: data
    });
  },
  //短信签名列表
  smsSignListApi(data: any) {
    return http("sms-template/sign-list", {
      method: "POST",
      data: data
    });
  },
  //删除短信签名列表
  smsSignDeleteApi(data: any) {
    return http("sms-template/sign-delete", {
      method: "POST",
      data: data
    });
  },
  //客户列表备注
  customerEditRemark(data: any) {
    return http("customer/edit-remark", {
      method: "POST",
      data: data
    });
  },
  smsSendApi(data: any) {
    return http("sms-template/send", {
      method: "POST",
      data: data
    });
  },
  //获取活动类型
  marketingSelectorTextApi() {
    return http("marketing-strategies/selector-text", {
      method: "POST"
    });
  },
  //获取活动页列表
  marketingActivePageApi(data: any) {
    return http("marketing-strategies/activity-pages", {
      method: "POST",
      data: data
    });
  },

  //获取活动类型下所有活动
  marketingActivitiesApi(data: any) {
    return http("marketing-strategies/activities", {
      method: "POST",
      data: data
    });
  },
  //更新策略状态
  marketingUpdateStatusApi(data: any) {
    return http("marketing-strategies/update-status", {
      method: "POST",
      data: data
    });
  },
  //更新业绩时间
  marketingUpdateSettingApi(data: any) {
    return http("marketing-strategies/update-strategy-data-time", {
      method: "POST",
      data: data
    });
  },
  // 客户标签分组
  customerTagList(data: any) {
    return http("customer/tag-list", {
      method: "POST",
      data: data
    });
  },
  //是否显示公众号 商场
  customerListConfig(data: any) {
    return http("customer/list-config", {
      method: "POST",
      data: data
    });
  },
  // 导出excel
  customerOrderExport(data: any) {
    return http("customer/order-export", {
      method: "POST",
      data: data,
      responseType: "arraybuffer"
    });
  },
  // 客户分组中当前客户人数
  tagGroupUserCount(data: any) {
    return http("marketing-strategies/tag-group-user-count", {
      method: "POST",
      data: data
    });
  },
  //客户分组查询详情
  groupDetailApi(data: any) {
    return http("user/group-detail", {
      method: "POST",
      data: data
    });
  },
  //删除用户分组
  groupDeleteApi(data: any) {
    return http("user/group-delete", {
      method: "POST",
      data: data
    });
  },
  //编辑客户分组保存
  groupEditApi(data: any) {
    return http("user/group-edit", {
      method: "POST",
      data: data
    });
  },
  //会员等级列表
  memberLevelListApi(data: any) {
    return http("customer/member-level-list", {
      method: "POST",
      data: data
    });
  },
  // 获取渠道名称
  geChannelName(data: any) {
    return http("/marketing-strategies/get-chan-list", {
      method: "POST",
      data: data
    });
  },
  // 查询相关渠道的数据
  searchChannelData(data: any) {
    return http("/marketing-strategies/short-url-list", {
      method: "POST",
      data: data
    });
  },
  // 删除当前的这个渠道的相关数据
  deleteChannelData(data: any) {
    return http("/marketing-strategies/short-url-del", {
      method: "POST",
      data: data
    });
  },
  // 获取整个表单的数据
  getFormData(data: any) {
    return http("/marketing-strategies/short-url-data", {
      method: "POST",
      data: data
    });
  },
  // 获取商品详情页的相关参数
  getGoodsParams(data: any) {
    return http("/marketing-strategies/get-spu-detail", {
      method: "POST",
      data: data
    });
  },
  // 保存整个表单的数据
  saveFormData(data: any) {
    return http("/marketing-strategies/short-url-add", {
      method: "POST",
      data: data
    });
  },
  // 获取通用短外链
  getGeneralShortChain(data: any) {
    return http("/marketing-strategies/short-url", {
      method: "POST",
      data: data
    });
  },
  // 导出客户列表
  exportTagDetailApi(data: any) {
    return http("/tag/detail/export", {
      method: "POST",
      responseType: "blob",
      data: data
    });
  },
  getTagCustomersApi(data: any) {
    return http("/tag/detail/list", {
      method: "POST",
      data: data
    });
  }
};
