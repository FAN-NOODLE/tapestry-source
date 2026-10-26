<?php

use Illuminate\Support\Facades\Route;

Route::get('/demo', ['uses' => 'DemoController']);

Route::get('/test', ['uses' => 'Test\TestController']);
// 公共接口
Route::group(['namespace' => 'Common', 'prefix' => 'common'], function () {
    // 发送短信验证码
    Route::post('send-sms-code', ['uses' => 'SendSmsController']);
    // 店长登录
    Route::post('login', ['uses' => 'LoginController']);
    // 企业微信登陆
    Route::post('qy-login', ['uses' => 'QYLoginController']);
    // 退出登录
    Route::post('exit-login', ['uses' => 'ExitLoginController']);
    // 选择企业
    Route::post('choice-enterprise', ['uses' => 'ChoiceEnterpriseController']);
    // 校验Token
    Route::post('check-login', ['uses' => 'CheckLoginController']);
    //物流列表
    Route::post('logistics-list', ['uses' => 'LogisticsListController@index']);
    //区域门店列表
    Route::post('area-store-list', ['uses' => 'AreaStoreListController@index']);
    // 生成太阳码
    Route::post('create', ['uses' => 'CreateController@index']);
    //企业id 换hash
    Route::post('hash', ['uses' => 'GetEnterpriseHashController@index']);

    //记录引导页浏览记录
    Route::post('boot_log_create', 'BootPageLogController@index');
    //检查是否需要导购引导页
    Route::post('boot_log_check', 'BootPageLogController@check');
    // 用户小程序静默登录
    Route::post("wx-login", ["uses" => "WxMiniLoginController@index"]);

    Route::post("agent-work-login", ["uses" => "AgentQYLoginController@getToken"]);
    // 小程序-首次弹框引导标识
    Route::post('check-first-open', ['uses' => 'FirstOpenController@index']);

});

// 员工管理
Route::group(['namespace' => 'Member', 'prefix' => 'member', 'middleware' => ['auth']], function () {
    // 员工列表
    Route::get('list', ['uses' => 'ListController']);
    // 员工详情
    Route::get('detail', ['uses' => 'DetailController']);
    // 员工离岗
    Route::get('edit-status', ['uses' => 'EditStatusController']);
    // 员工统计
    Route::get('statistics', ['uses' => 'StatisticsController']);
    // 店员修改接收短信状态
    Route::post('update-sms-status', ['uses' => 'UpdateSmsStatusController']);
    // 获取员工业绩
    Route::post('member-statistics', ['uses' => 'MemberStatisticsController']);
    // 店员修改接收短信手机号
    Route::post('update-sms-mobile', ['uses' => 'UpdateSmsMobileController']);
    // 生成海报
    Route::post('generate-worker-poster', ['uses' => 'GenerateWorkerPosterController']);
    // 根据门店或部门获取对应导购列表
    Route::post('guide-list', ['uses' => 'GuideListByStoreController@index']);
    // 获取导购身份和权限类型
    Route::post('guide-position', ['uses' => 'GuidePositionInfoController@index']);

    // 门店员工详情
    Route::get('worker-detail', ['uses' => 'WorkerDetailController']);
    // 获取员工的门店列表
    Route::post('store-lists', ['uses' => 'MemberStoreController@index']);
    // coach 店长(SM)获取管理的员工列表
    Route::post('staff-lists', ['uses' => 'MemberStoreController@staffLists']);
});

// 门店管理
Route::group(['namespace' => 'Store', 'prefix' => 'store', 'middleware' => ['auth']], function () {
    // 门店列表
    Route::post('list', ['uses' => 'ListController@index']);
    // 门店详情
    Route::get('detail', ['uses' => 'DetailController']);
    // 门店统计
    Route::get('statistics', ['uses' => 'StatisticsController']);
    // 修改门店
    Route::post('edit', ['uses' => 'EditController']);
    // 根据经纬度返回省市区
    Route::get('gis-to-address', ['uses' => 'GisToAddressController']);
    // 门店首页业绩
    Route::post('home-statistics', ['uses' => 'HomeStatisticsController']);
    // 订单tab数字
    Route::post('home-order-tab', ['uses' => 'HomeOrderTabController']);
    // 数据罗盘(暂时废弃)
    Route::post('compass-statistics', ['uses' => 'CompassStatisticsController']);
    // 根据导购权限返回门店或部门列表
    Route::post('store-or-department-list', ['uses' => 'StoreOrDepartmentListByGuideController@index']);
    //获取我的订单跳转h5地址
    Route::post('d1m-order-h5', ['uses' => 'D1MOrderH5Controller@index']);

});

// 订单管理
Route::group(['namespace' => 'Order', 'prefix' => 'order', 'middleware' => ['auth']], function () {

    // 订单列表
    Route::post("order-list", ['uses' => 'OrderListController@index']);

    // 订单详情
    Route::post("order-detail", ['uses' => 'OrderDetailController@index']);

    // 未处理订单数量
    Route::post("order-count", ['uses' => 'OrderCountController@index']);

    // 确认发货
    Route::post("order-deliver", ['uses' => 'OrderDeliverController@index']);

    // 确认用户自提
    Route::post("take-confirm", ['uses' => 'TakeConfirmController@index']);

    // 审批售后
    Route::post("aftersale-confirm", ['uses' => 'AfterSaleConfirmController@index']);

    // 确认退款
    Route::post("refund-confirm", ['uses' => 'RefundConfirmController@index']);

    // 售后详情
    Route::post("aftersale", ['uses' => 'AfterSaleController@index']);

    // 一键转人工
    Route::post('refund-to-manual', ['uses' => 'RefundToManualController@index']);

    // 订单校验提货码
    Route::post('check-take-goods-code', ['uses' => 'CheckTakeGoodsCodeController@index']);

    // 推广或服务订单查看最新时间标记接口(暂时废弃)
    Route::post('click-sign', ['uses' => 'ClickSignController@index']);

    // 推广人或服务人列表
    Route::post('member-list', ['uses' => 'MemberListController@index']);

    // 服务人接单
    Route::post('member-sign', ['uses' => 'MemberSignController@index']);

    // 查询多物流订单
    Route::post('order-logistics-info', ['uses' => 'QueryOrderLogisticController@index']);

    // 多物流发货
    Route::post('add-logistics-package', ['uses' => 'SendGoodsBatchController@index']);

    // 订单统计排行
    Route::post('ranking-list', ['uses' => 'OrderStatisticalRankingListController@index']);
    // 订单统计
    Route::post('order-statistical', ['uses' => 'OrderStatisticalController@index']);
    //订单分导购和合伙人统计
    Route::post('order-statistical-classification', ['uses' => 'OrderStatisticalClassificationController@index']);
});

// 顾客管理
Route::group(['namespace' => 'Customer', 'prefix' => 'customer', 'middleware' => ['auth']], function () {
    // 顾客列表
//    Route::get('customer-list', ['uses' => 'ListController@index']);

    // 获取外部联系人于员工关系
    Route::post('get-user-member-relation', ['uses' => 'RelationController@index']);
    // 顾客列表(新)
    Route::get('customer-list', ['uses' => 'ListController@list']);
    // 顾客详情(新)
    Route::get('customer-detail', ['uses' => 'DetailController@detail']);
    // 标签列表(新)
    Route::get('tags-group-list', ['uses' => 'TagController@tagsGroupList']);
    // 顾客详情
    Route::get('detail', ['uses' => 'DetailController@index']);
    // 标签列表
    Route::get('all-tag', ['uses' => 'TagController@allTag']);
    // 顾客数据罗盘
    Route::post('data-compass', ['uses' => 'DataCompassController@index']);
    // 智慧导购顾客详情
    Route::post('customer-wisdom-detail', ['uses' => 'WisdomDetailController@index']);
    // 智慧导购顾客消费记录列表
    Route::post('customer-expense-record', ['uses' => 'ExpenseRecordListController@index']);
    // 智慧导购顾客优惠券列表
    Route::post('customer-coupon-list', ['uses' => 'CouponListController@index']);
    // 智慧导购顾客优惠券商品列表
    Route::post('customer-coupon-goods-list', ['uses' => 'CouponGoodsListController@index']);
    // 解除用户字段定义标签绑定
    Route::post('usertag-delete-relation', ['uses' => 'UserTagController@deleteRelation']);
    // 删除店员自定义标签
    Route::post('usertag-delete-tag', ['uses' => 'UserTagController@deleteTag']);
    // 添加用户自定义标签
    Route::post('usertag-add-tag', ['uses' => 'UserTagController@addTag']);
    // 绑定用户自定义标签
    Route::post('usertag-add-relation', ['uses' => 'UserTagController@addRelation']);
    // 获取绑定自定义标签列表
    Route::post('usertag-get-relation', ['uses' => 'UserTagController@getRelation']);
    // 获取员工自定义标签列表(未绑定)
    Route::post('usertag-get-tags', ['uses' => 'UserTagController@getTags']);
    // 获取员工所有自定义标签列表
    Route::post('usertag-get-list', ['uses' => 'UserTagController@getList']);
    // 获取员工所有自定义标签列表
    Route::post('usertag-search-tags', ['uses' => 'UserTagController@searchTags']);

    // 数据罗盘-顾客排行榜
    Route::post('rank', ['uses' => 'RankController@index']);
    // 数据罗盘-新顾客排行榜
    Route::post('rank-data', ['uses' => 'RankDataController@index']);
    // 数据罗盘-导购排行榜-财月财年筛选条件
    Route::post('fiscal-month', ['uses' => 'FiscalMonthController@index']);
    // 数据罗盘-顾客数据概览
    Route::post('statistic', ['uses' => 'StatisticController@index']);
    // 小程序数据罗盘的数据统计，（社群数据&&好友数据）
    Route::post('get-compass-summary', ['uses' => 'StatisticController@getCompassSummary']);
    //小程序数据罗盘的数据统计，（社群数据&&好友数据） Ck查询
    Route::post('get-compass-summary-ck', ['uses' => 'StatisticController@getCompassSummaryCk']);
    // 小程序数据罗盘的数据统计，社群数据统计
    Route::post("chat-group-summary", ['uses' => 'ChatGroupController@getSummary']);
    // 小程序数据罗盘的数据统计，好友数据统计
    Route::post("friends-summary", ['uses' => 'FriendsController@getSummary']);
    // 数据罗盘-顾客数据提示记录
    Route::post('set-view-tab', ['uses' => 'StatisticController@setViewTab']);

    // 丢失客户记录
    Route::post('lost-record', ['uses' => 'LostRecordController@index']);
    // 丢失客户记录额外数据
    Route::post('lost-option', ['uses' => 'LostRecordController@option']);

    // 导购端
    // 小程序-顾客列表
    Route::post('external-user-list', ['uses' => 'ListController@userList']);
    Route::post('external-user-list-ck', ['uses' => 'ListController@userListCk']);
    // 小程序-顾客列表的筛选条件
    Route::post('get-choice-option', ['uses' => 'ListController@getChoiceOption']);
    // 小程序-顾客详情
    Route::post('external-user-detail', ['uses' => 'ListController@userDetail']);
    // 小程序-顾客手机号
    Route::post('external-user-mobile', ['uses' => 'ListController@userMobile']);
    // 小程序-顾客订单
    Route::post('user-crm-order-list', ['uses' => 'ListController@userCrmOrderList']);
    // 小程序-顾客卡劵
    Route::post('user-crm-coupons-list', ['uses' => 'ListController@userCrmCouponsList']);
    // Sw小程序-顾客卡劵
    Route::post('sw-user-crm-coupons-list', ['uses' => 'ListController@swUserCrmCouponsList']);
    // 小程序-获取顾客邀约和播种数量
    Route::post('user-detail-tips-num', ['uses' => 'ListController@userDetailTipsNum']);
    // 小程序-顾客和导购的绑定关系
    Route::post('user-bind-relation', ['uses' => 'ListController@userBindRelation']);
    // 小程序-修改顾客备注姓名
    Route::post('update-remark', ['uses' => 'ListController@userUpdateRemark']);
    // 小程序-修改顾客的描述
    Route::post('update-desc', ['uses' => 'ListController@userUpdateDescription']);
    // 小程序-操作顾客的印象标签
    Route::post('impression-tags-operate', ['uses' => 'ListController@userImpressionTagsOperate']);

    // 小程序-会员聚合-查询条件
    Route::post('customer-option', ['uses' => 'ListController@customerOption']);
    // 小程序-会员聚合-查询条件
    Route::post('customer-list', ['uses' => 'ListController@customerList']);
    // 小程序-会员聚合-新对接cdp
//    Route::post('customer-list', ['uses' => 'ListController@customerListForCDP']);
    // 小程序-会员聚合查询-白名单
    Route::post('white-list', ['uses' => 'ListController@mobileWhiteList']);
    // 小程序-修改顾客地址生日鞋码
    Route::post('save-info-add', ['uses' => 'ListController@saveInfoAdd']);
    // 小程序-顾客列表创建群发全选的人数统计
    Route::post('select-all-num', ['uses' => 'ListController@selectAllNum']);
    // 小程序-顾客列表创建群发返回全选的外部联系人external_user_id
    Route::post('get-external-ids', ['uses' => 'ListController@selectAllList']);

    // 小程序-会员聚合-全选的人数统计
    Route::post('customer-all-num', ['uses' => 'ListController@customerAllNum']);
    // 小程序-会员聚合-筛选要发送的外部联系人id
    Route::post('customer-external-ids', ['uses' => 'ListController@customerAllList']);
    // 小程序-三方接口会员聚合查询
//    Route::post('customer-list-cdp', ['uses' => 'ListController@customerListForCDP']);

    // 小程序-操作顾客的企微标签
    Route::post('work-tags-operate', ['uses' => 'ListController@userWorkTagsOperate']);
    // 播种任务
    Route::group(['prefix' => 'sow-task'], function () {
        Route::post('list', ['uses' => 'SowTaskController@list']);
        Route::post('create', ['uses' => 'SowTaskController@create']);
        Route::post('update', ['uses' => 'SowTaskController@update']);
        Route::post('finished', ['uses' => 'SowTaskController@finished']);
        Route::post('delete', ['uses' => 'SowTaskController@delete']);
        Route::post("template-list", ['uses' => 'SowTaskController@templateList']);
        Route::post("new-create", ['uses' => 'SowTaskController@newCreate']);
    });
    Route::group(['prefix' => 'external-search'], function () {
        // 小程序-顾客、顾客关系检索大宽表新版20241126
        Route::post('list', ['uses' => 'SearchController@list']);
        // 小程序-顾客、顾客关系id检索大宽表新版20241126
        Route::post('get-open-ids', ['uses' => 'SearchController@getOpenIds']);
        // 小程序-顾客、顾客关系检索条件20241126
        Route::post('get-choice-option', ['uses' => 'SearchController@getChoiceOption']);
        // 小程序-顾客详情
        Route::post('detail', ['uses' => 'SearchController@detail']);
        // 小程序-三方顾客详情
        Route::post('detail-by-crm', ['uses' => 'SearchController@detailByCrm']);
        // 小程序-CDP顾客商品推荐
        Route::post('customer-interested-products', ['uses' => 'SearchController@cdpInterestedProductsByCdp']);

    });

    Route::group(['prefix' => 'external-search-sw'], function () {
        // 小程序-顾客、顾客关系检索大宽表新版20241126
        Route::post('list', ['uses' => 'SearchSwController@list']);
        // 小程序-顾客、顾客关系id检索大宽表新版20241126
        Route::post('get-open-ids', ['uses' => 'SearchSwController@getOpenIds']);
        // 小程序-顾客、顾客关系检索条件20241126
        Route::post('get-choice-option', ['uses' => 'SearchSwController@getChoiceOption']);
    });
    Route::group(['prefix' => 'external-crowd'], function () {
        // 外部联系人人群包添加
        Route::post('add', ['uses' => 'CrowdController@addCrowd']);
        // 外部联系人人群包列表
        Route::post('list', ['uses' => 'CrowdController@crowdList']);
        // 外部联系人人群包详情
        Route::post('info', ['uses' => 'CrowdController@crowdInfo']);
        // 外部联系人人群包Openid集合
        Route::post('get-open-ids', ['uses' => 'CrowdController@crowdOpenIds']);
        // 外部联系人人群包编辑
        Route::post('edit', ['uses' => 'CrowdController@crowdEdit']);
        // 外部联系人人群包编辑
        Route::post('del', ['uses' => 'CrowdController@crowdDel']);
        // 外部联系人人群包排序
        Route::post('update-sort', ['uses' => 'CrowdController@crowdUpdateSort']);
        // 外部联系人人群包更新联系状态
        Route::post('update-contact', ['uses' => 'CrowdController@crowdUpdateContact']);
    });
    Route::group(['prefix' => 'external-crowd-sw'], function () {
        // 外部联系人人群包添加
        Route::post('add', ['uses' => 'CrowdSwController@addCrowd']);
        // 外部联系人人群包列表
        Route::post('list', ['uses' => 'CrowdSwController@crowdList']);
        // 外部联系人人群包详情
        Route::post('info', ['uses' => 'CrowdSwController@crowdInfo']);
        // 外部联系人人群包Openid集合
        Route::post('get-open-ids', ['uses' => 'CrowdSwController@crowdOpenIds']);
        // 外部联系人人群包编辑
        Route::post('edit', ['uses' => 'CrowdSwController@crowdEdit']);
        // 外部联系人人群包编辑
        Route::post('del', ['uses' => 'CrowdSwController@crowdDel']);
        // 外部联系人人群包排序
        Route::post('update-sort', ['uses' => 'CrowdSwController@crowdUpdateSort']);
        // 外部联系人人群包更新联系状态
        Route::post('update-contact', ['uses' => 'CrowdSwController@crowdUpdateContact']);
        // 外部联系人人群包更新联系状态
        Route::post('batch-update-contact', ['uses' => 'CrowdSwController@crowdBatchUpdateContact']);
    });
});

// 小程序首页dashboard
Route::group(['namespace' => 'Dashboard', 'prefix' => 'dashboard', 'middleware' => ['auth']], function () {
    // 小程序首页dashboard
    Route::get('index', ['uses' => 'IndexController@index']);
});

// 全员分销介绍
Route::group(['namespace' => 'DistributionConfig', 'prefix' => 'distribution-config', 'middleware' => ['auth']], function () {
    // 全员分销介绍
    Route::post('detail', ['uses' => 'DetailController@index']);
    // 二级分销信息
    Route::post('partner-detail', ['uses' => 'DetailPartnerController@index']);
});

// 分销商品
Route::group(['namespace' => 'DistributionGoods', 'prefix' => 'distribution-goods', 'middleware' => ['auth']], function () {
    // 分类列表
    Route::post('category-list', ['uses' => 'CategoryListController@index']);
    // 商品列表
    Route::post('goods-list', ['uses' => 'GoodsListController@index']);
    // 商品二维码
    Route::post('goods-qr-code', ['uses' => 'GoodsQrCodeController@index']);
});

// 分销统计
Route::group(['namespace' => 'DistributionStatistics', 'prefix' => 'distribution-statistics', 'middleware' => ['auth']], function () {
    // 分销统计首页
    Route::post('statistical', ['uses' => 'StatisticalController@index']);
    // 分销排行榜
    Route::post('statistical-ranking-list', ['uses' => 'StatisticalRankingListController@index']);
});

// 分销收支
Route::group(['namespace' => 'DistributionBalance', 'prefix' => 'distribution-balance', 'middleware' => ['auth']], function () {
    // 分销收支统计
    Route::post('cash-withdrawal-info', ['uses' => 'CashWithdrawalInfoController@index']);
});

// 导购提现相关路由
Route::group([
    'namespace' => 'Withdrawal',
    'prefix' => 'withdrawal',
    'middleware' => ['auth']
], function () {
    // 提现认证 - 添加
    Route::post('certification-add', ['uses' => 'CertificationAddController@index']);
    // 提现认证 - 更新
    Route::post('certification-update', ['uses' => 'CertificationUpdateController@index']);
    // 提现认证 - 审核
    Route::post('certification-audit', ['uses' => 'CertificationAuditController@index']);
    // 提现认证 - 详情
    Route::post('certification-detail', ['uses' => 'CertificationDetailController@index']);
    // 提现认证 - 获取认证文件
    Route::post('certification-file-url', ['uses' => 'CertificationDetailController@certificationFileUrl']);

    // 提现 - 申请
    Route::post('withdrawal-add', ['uses' => 'WithdrawalAddController@index']);
    // 提现 - 申请填写页面数据
    Route::post('withdrawal-home', ['uses' => 'WithdrawalHomeController@index']);
    // 提现 - 记录列表
    Route::post('withdrawal-logs', ['uses' => 'WithdrawalLogsController@index']);
    // 提现 - 银行列表
    Route::post('withdrawal-bank-list', ['uses' => 'WithdrawalBankListController@index']);
});


// 结算记录
Route::group(['namespace' => 'DistributionSettlement', 'prefix' => 'distribution-settlement', 'middleware' => ['auth']], function () {
    // 分销收支统计
    Route::post('list', ['uses' => 'ListController@index']);
});

Route::group(['namespace' => 'PartnerMember', 'prefix' => 'partner-member'], function () {
    //上传
    Route::post('upload', ['uses' => 'UploadController@index']);
    //微信登录
    Route::post('wx-login', ['uses' => 'PersonalAppletLoginController@index']);
    //微信授权
    Route::post('wx-auth', ['uses' => 'AuthWxController@index']);
    //注册
    Route::post('register', ['uses' => 'RegisterController@index']);
    //获取手机号
    Route::post('mobile', ['uses' => 'MobileController@index']);
    //合伙人联系导购
    Route::post('contact', ['uses' => 'contactGuideController@index']);

});
// 分销合伙人管理
Route::group(['namespace' => 'DistributionPartner', 'prefix' => 'distribution-partner', 'middleware' => ['auth']], function () {
    // 分销合伙人列表
    Route::post('list', ['uses' => 'ListController@index']);
    // 禁用详情
    Route::post('disable-detail', ['uses' => 'DisableDetailController@index']);
    // 合伙人我的信息
    Route::post('my', ['uses' => 'MyController@index']);
    //合伙人短信开关
    Route::post('sms-switch', ['uses' => 'PartnerIsAcceptSmsController@index']);
    //合伙人手机号是否脱敏展示
    Route::post('mobile-switch', ['uses' => 'MobileSensitiveController@index']);
});
//营销日历
Route::group(['namespace' => 'MarketingPlan', 'prefix' => 'marketing-plan', 'middleware' => ['auth']], function () {
    // 获取营销日历
    Route::post('task-calendar', ['uses' => 'CalendarListController@index']);
    // 上报任务完成情况
    Route::post('task-finish', ['uses' => 'TaskFinishController@index']);
    // 计划列表
    Route::post('list', ['uses' => 'TaskListController@index']);
    // 任务列表
    Route::post('task-list', ['uses' => 'TaskListController@taskList']);
    // 获取任务详情
    Route::post('task-detail', ['uses' => 'TaskListController@taskInfo']);
    // 获取任务标签客户列表
    Route::post('tag-user-list', ['uses' => 'TaskListController@tagUserList']);
    // 设置任务状态
    Route::post('set-status', ['uses' => 'TaskListController@setStatus']);
    Route::post('update-status', ['uses' => 'TaskListController@updateStatus']);
    // 获取
    Route::post('getVideoMedia', ['uses' => 'TaskListController@getVideoMedia']);
    // coach 任务中心的日历列表
    Route::post('month-list', ['uses' => 'PlansController@index']);
    // coach 任务中心获取主任务列表
    Route::post("plan-list", ['uses' => 'PlansController@lists']);
    // cocah 获取主任务信息
    Route::post("plan-detail", ['uses' => 'PlansController@detail']);
    // coach 任务中心的未完成tips
    Route::post("plan-tips", ['uses' => 'PlansController@tipsInfo']);
    // coach 退群的计划任务用户列表
    Route::post("leave-group-chat-lists", ['uses' => 'PlansController@leaveGroupChat']);
    // coach sacode&&unionid|customerid的顾客列表
    Route::post("sa-customer-list", ['uses' => 'PlansController@saCustomerList']);
    // 计划任务点击打电话按钮
    Route::post("confirm-call", ['uses' => 'PlansController@confirmCall']);
    // 计划任务点击发消息按钮
    Route::post("confirm-contact", ['uses' => 'PlansController@confirmContact']);
    // 计划任务标记任务完成
    Route::post("confirm-finish", ['uses' => 'PlansController@confirmFinish']);
    // 计划任务填写联系说明
    Route::post("set-remark", ['uses' => 'PlansController@setRemark']);
    // 点击转发顾客按钮的上报
    Route::post("set-forward", ['uses' => 'PlansController@setForward']);
    // 点击添加好友的按钮上报
    Route::post("confirm-friend", ['uses' => 'PlansController@confirmFriend']);
    // 主任务已完成状态关闭未触达提示
    Route::post("close_touch_tips", ['uses' => 'PlansController@closeTouchTips']);
    // 子任务的已完成状态关闭未触达提示
    Route::post("close_task_touch_tips", ['uses' => 'PlansController@closeTaskTouchTips']);
    // coach 计划任务升级改造的顾客列表
    Route::post("search-customer-list", ['uses' => 'PlansController@searchCustomerList']);
    // coach 计划任务升级改造的顾客列表
    Route::post("search-external-list", ['uses' => 'PlansController@searchExternalList']);
    //coach升级的计划任务顾客列表,全选获取openid
    Route::post("search-external-all-openid", ['uses' => 'PlansController@searchExternalAllOpenId']);
    // coach 计划任务的顾客列表移动
    Route::post("customer-move", ['uses' => 'PlansController@customerMove']);
    // coach 计划任务的转发给顾客上报
    Route::post("send-finish", ['uses' => 'PlansController@sendFinish']);
    // coach 顾客详情的任务看板，获取导购和顾客的任务列表
    Route::post("customer-task-list", ['uses' => 'PlansController@customerTaskList']);
});

// BA业绩看板
Route::group([
    'namespace' => 'BAKanBan',
    'prefix' => 'ba-kanban',
    'middleware' => ['auth'],
], function () {
    // 获取当前用户选择的把业绩看板的数据
    Route::post("", ['uses' => 'IndexController@index']);
    Route::post("ba_source_depte", ['uses' => 'BaSourceDeptController@index']);
});

// 预约单服务
Route::group([
    'namespace' => "Reserve",
    "prefix" => "reserve",
], function () {
    Route::post("login", ["uses" => "WechatMiniLoginController@index"]);

    // 查询特殊活动的任务信息

    Route::post("reserve-activity", ["uses" => "IndexController@getReserveActivity"]);
    // 特殊任务的时间段人数限制
    Route::post("new-time-list", ["uses" => "MemberOpreateController@getNewTimeList"]);

    //用户微信端操作预约单信息
    Route::group(['middleware' => ['reserveAccess']], function () {
        // 用户获取订单信息
        Route::post("order-info", ["uses" => "IndexController@getInfo"]);
        // 用户修改并确认订单
        Route::post("edit-order", ["uses" => "IndexController@editOrder"]);
        // 用户取消订单
        Route::post("cancel-order", ["uses" => "IndexController@cancelOrder"]);
        // 用户订单打分
        Route::post("make-score", ["uses" => "IndexController@makeScore"]);
        // 用户端的时间段邀约记录
        Route::post("user-time-list", ["uses" => "IndexController@timeList"]);
        // 用户端的打分信息
        Route::post("score-info", ["uses" => "IndexController@scoreInfo"]);
        // 用户端的服务列表
        Route::post("user-service-list", ["uses" => "IndexController@serviceList"]);
        // SW端增加签到按钮功能
        Route::post("user-sign-in", ["uses" => "IndexController@signIn"]);
        // coach用户编辑评价
        Route::post("edit-score", ["uses" => "IndexController@editScore"]);
        // coach 用户获取预约单列表
        Route::post("user-order-lists", ["uses" => "IndexController@orderList"]);

    });
    Route::group(['middleware' => ['auth']], function () {
        // 创建预约单
        Route::post("create-order", ["uses" => "MemberOpreateController@createOrder"]);
        // 时间段邀约列表
        Route::post("time-list", ["uses" => "MemberOpreateController@timeList"]);
        // 邀约排行列表
        Route::post("rank-list", ["uses" => "MemberOpreateController@rankList"]);

        // 预约列表
        Route::post("order-list", ["uses" => "MemberOpreateController@orderList"]);
        // coach预约列表
        Route::post("order-list-by-coach", ["uses" => "MemberOpreateController@orderListByCoach"]);

        // 编辑预约
        Route::post("order-edit", ["uses" => "MemberOpreateController@orderEdit"]);
        // 取消预约
        Route::post("order-cancel", ["uses" => "MemberOpreateController@orderCancel"]);
        // 店长分配预约单
        Route::post("order-assign", ["uses" => "MemberOpreateController@orderAssign"]);
        // 分享预约单
        Route::post("order-share", ["uses" => "MemberOpreateController@orderShare"]);
        // 获取用户的基本信息
        Route::post("user-simple-info", ["uses" => "MemberOpreateController@userSimpleInfo"]);
        // 获取导购门店信息
        Route::post("member-info", ["uses" => "MemberOpreateController@memberInfo"]);
        // 导购服务类型列表
        Route::post("service-list", ["uses" => "MemberOpreateController@serviceList"]);
        // 预约详情
        Route::post("order-detail", ["uses" => "MemberOpreateController@orderDetail"]);
        // 完成预约
        Route::post("order-finish", ["uses" => "MemberOpreateController@orderFinish"]);
        // 外部联系人有关系的导购列表手机号
        Route::post("member-list", ["uses" => "MemberOpreateController@memberLists"]);
        // coach 导购确认订单是否是到店购买
        Route::post("change-status", ["uses" => "MemberOpreateController@changeStatus"]);
        // coach 导购确认用户到店
        Route::post("confirm-store", ["uses" => "MemberOpreateController@confirmStore"]);
        // coach 导购邀请用户评论
        Route::post("comment-share", ["uses" => "MemberOpreateController@commentShare"]);
        // 获取门店列表
        Route::post("store-lists", ["uses" => "MemberOpreateController@StoreLists"]);
        // coach 任务中心预约单日历列表
        Route::post("order-month-list", ["uses" => "MemberOpreateController@orderMonthList"]);
        // coach 判断导购身份是否可以给用户发送邀请评价
        Route::post("verify-share", ["uses" => "MemberOpreateController@verifyShare"]);
        // SW&&ks验证同一时间段的邀约限制
        Route::post("check-reserve-time", ["uses" => "MemberOpreateController@checkReserveTime"]);
        // 获取客户的所有邀约单数据
        Route::post("customer-reserve-lists", ["uses" => "MemberOpreateController@customerReserveLists"]);
    });
    //用户微信端操作预约单信息
    Route::group([
        "prefix" => "subscribe-message",
        'middleware' => ['reserveAccess']
    ], function () {
        Route::post("get-config", ["uses" => "SubscribeMessageController@getConfig"]);
        Route::post("subscribe", ["uses" => "SubscribeMessageController@subscribe"]);
    });
    // 用户微信端服务信息
    Route::group(['prefix' => "service"], function () {
        // 服务类型列表
        Route::post("list", ["uses" => "ServiceController@serviceList"]);
    });
});


// 活动
Route::group([
    'namespace' => 'Activity',
    'prefix' => 'activity'
], function () {
    // 刮刮乐活动
    Route::group([
        'prefix' => 'scratch',
        'middleware' => ['reserveAccess']
    ], function () {
        // 活动详情
        Route::post('info', ['uses' => 'ScratchController@Info']);
        // 刮奖
        Route::post('scratch-prize', ['uses' => 'ScratchController@ScratchPrize']);
        // 券列表
        Route::post('coupon-list', ['uses' => 'ScratchController@CouponList']);
        // 检查用户券核销状态
        Route::post('check-coupon', ['uses' => 'ScratchController@CheckCoupon']);
        // 用户分享上报
        Route::post('share-info', ['uses' => 'ScratchController@ShareInfo']);
    });
    Route::group(['prefix' => 'check', 'middleware' => ['auth']], function () {
        // 核销券
        Route::post('write-off-coupon', ['uses' => 'CheckCouponController@index']);
        // 已核销券列表
        Route::post('coupon-list', ['uses' => 'CheckCouponController@CouponList']);
    });

    // Coach的CIIE吊牌预约活动
    Route::group([
        'prefix' => 'reserve',
        'middleware' => ['reserveAccess']
    ], function () {
        // 活动详情
        Route::post('info', ['uses' => 'ReserveCIIEController@Info']);
        // 创建预约
        Route::post('create', ['uses' => 'ReserveCIIEController@CreateReserve']);
        // 取消预约
        Route::post('cancel', ['uses' => 'ReserveCIIEController@CancelReserve']);
        // 编辑预约
        Route::post('edit', ['uses' => 'ReserveCIIEController@EditReserve']);
        // 预约时间段信息列表
        Route::post('time-lists', ['uses' => 'ReserveCIIEController@ReserveLists']);
        // 获取用户手机号
        Route::post('get-phone', ['uses' => 'ReserveCIIEController@GetPhoneNumber']);
    });
    Route::group(['prefix' => 'reserve-check', 'middleware' => ['auth']], function () {
        // 核销
        Route::post('index', ['uses' => 'ReserveCIIECheckController@Index']);
        // 白名单列表
        Route::post('white-list', ['uses' => 'ReserveCIIECheckController@WhiteLists']);
        // 获取预约列表
        Route::post('lists', ['uses' => 'ReserveCIIECheckController@GetLists']);
        // 管理导购处理的时间列表
        Route::post('time-lists', ['uses' => 'ReserveCIIECheckController@GetReserveTimeList']);
        // 管理导购编辑预约时间段
        Route::post('time-config', ['uses' => 'ReserveCIIECheckController@ReserveTimeConfig']);
        // 导购已核销列表
        Route::post('check-lists', ['uses' => 'ReserveCIIECheckController@CheckLists']);
        // 核销券
        Route::post('write-off-coupon', ['uses' => 'ReserveCIIECheckController@index']);
    });
});

//在职继承
Route::group(['prefix' => 'customer-inherit', 'namespace' => 'CustomerInherit', 'middleware' => ['auth']], function () {
    //任务列表
    Route::post('list', ['uses' => 'IndexController@inheritList']);
    //任务详情
    Route::post('info', ['uses' => 'IndexController@inheritInfo']);
    //转移
    Route::post('transfer', ['uses' => 'IndexController@inheritTransfer']);
    //导购列表
    Route::post('member-list', ['uses' => 'IndexController@memberList']);
    //选择外部联系人的列表
    Route::post('customer-list', ['uses' => 'IndexController@customerList']);
    //选择外部联系人的列表(新)
    Route::post('search-list', ['uses' => 'IndexController@searchList']);
    //选择外部联系人的列表(新)
    Route::post('search-list-sw', ['uses' => 'IndexSwController@searchList']);
    //在职分配可接替的导购列表
    Route::post('takeover-member-list', ['uses' => 'IndexController@takeoverMemberList']);
    // 在职分配全选的人员数
    Route::post('select-all-num', ['uses' => 'IndexController@selectAllNum']);
});

//绩效榜单
Route::group(['prefix' => 'performance-list', 'namespace' => 'PerformanceList', 'middleware' => ['auth']], function () {
    //榜单列表活跃榜和顾客榜
    Route::post('active-customers', ['uses' => 'IndexController@activeCustomers']);
    //邀约排行榜
    Route::post('invitation', ['uses' => 'IndexController@invitation']);
    //cache明星榜
    Route::post('star', ['uses' => 'IndexController@star']);
    //业绩榜单(ks和sw)
    Route::post('achievement', ['uses' => 'IndexController@achievement']);
    //sw顾客榜单
    Route::post('sw-customers', ['uses' => 'IndexController@swCustomers']);
    //sw活跃榜单
    Route::post('sw-active', ['uses' => 'IndexController@swActive']);
    //sw明星榜单
    Route::post('sw-star', ['uses' => 'IndexController@swStar']);
    //ks明星榜单
    Route::post('ks-star', ['uses' => 'IndexController@ksStar']);

});

// 素材
Route::group(['prefix' => 'source', 'namespace' => 'Source', 'middleware' => ['auth']], function () {
    // 产品类别列表
    Route::post("category/list", ["uses" => "ListController@categoryList"]);
    // 素材分组归属门店列表
    Route::post("group-relation/list", ["uses" => "ListController@groupList"]);
    // 根据门店获取素材分组列表
    Route::post("group-relation/get-list-for-store", ["uses" => "ListController@getListForStore"]);
    // 小程序端门店员工根据查看分组权限查询各素材列表
    Route::post('select-list', ['uses' => 'ListController@index']);
    // 小程序端门店员工根据查看分组权限查询各素材上新数量
    Route::post('get-new-num', ['uses' => 'ListController@getNewNum']);
    // 小程序端下载，转发，复制上报
    Route::post('report-times', ['uses' => 'ListController@reportTimes']);

    // 小程序端验证员工是否有上传素材的权限
    Route::post("upload-auth", ['uses' => 'UploadController@checkAuth']);
    // 小程序端员工上传素材
    Route::post("upload-media", ["uses" => "UploadController@uploadMedia"]);
    // 小程序端导购获取上传cos的凭证
    Route::post("credential", ["uses" => "UploadController@credential"]);
    // 获取压缩结果
    Route::post("compress", ["uses" => "UploadController@compress"]);
    // 压缩文件
    Route::post("compress-file", ["uses" => "UploadController@compressFile"]);
    // 修改视频创建时间
    Route::post("edit-video-create-time", ["uses" => "UploadController@editVideoCreateTime"]);
    // 验证素材资源是否满足朋友圈
    Route::post("get-media-info", ["uses" => "UploadController@getMediaInfo"]);
    // 判断用户是否首次打开素材库的交互提示
    Route::post("check-first", ["uses" => "ListController@checkFirst"]);
    // 后台处理视频第一帧封面
    Route::post("make-cover", ["uses" => "UploadController@makeCover"]);

    // 品牌商品列表素材搜索
    Route::post("brand/search-list", ["uses" => "BrandSearchController@list"]);
    // 品牌商品列表素材搜索sku
    Route::post("brand/search-list-group-sku", ["uses" => "BrandSearchController@listBySku"]);
    // 根据UpcCode获取商品素材信息
    Route::post("brand/get-cache-info", ["uses" => "BrandSearchController@getCacheInfo"]);
    // 新调整官网素材库查询
    Route::post("brand/collect/search-list", ["uses" => "BrandSearchController@collectSearchList"]);
    // 收藏素材列表
    Route::post("collect-select-list", ['uses' => 'ListController@collectSelectList']);
    // 素材收藏
    Route::post("collect", ['uses' => 'ListController@collect']);
    // AI识图
    Route::post("brand/talkie-image-search", ["uses" => "BrandSearchController@talkieImageSearch"]);

});

// 邀请函
Route::group(['prefix' => 'invite', 'namespace' => 'Invite'], function () {
    // 小程序邀请函
    Route::group(['middleware' => ['auth']], function () {
        // 导购端-邀请函模板列表信息
        Route::post('config-list', ['uses' => 'MemberInviteController@configList']);
        // 导购端-邀请函模板配置详情
        Route::post('config-detail', ['uses' => 'MemberInviteController@configDetail']);
        // 导购端-发送邀约单记录
        Route::post('create-record', ['uses' => 'MemberInviteController@createRecord']);
        // 导购端-导购查看邀请函信息
        Route::post('record-info', ['uses' => 'MemberInviteController@getRecordInfo']);

    });

    // 用户微信端操作邀请函信息
    Route::group(['middleware' => ['reserveAccess']], function () {
        // 用户查看邀请函信息
        Route::post("get-record-info", ["uses" => "IndexController@getRecordInfo"]);
        // 用户创建邀请函订单
        Route::post("create-order", ["uses" => "IndexController@createOrder"]);

        // 用户我的页面
        Route::post('my-center', ['uses' => 'IndexController@myCenter']);
        // 我的页面-添加专属导购
        Route::post('my-sa', ['uses' => 'IndexController@mySA']);

        // 用户预览邀请函
        Route::post('info', ['uses' => 'IndexController@inviteInfo']);
        // 用户确认邀请函的时候SW&&KS需要验证到店时间的限制
        Route::post('check-reserve-create', ['uses' => 'IndexController@checkReserveCreate']);
    });


});

//企微对外收款
Route::group(['prefix' => 'external-pay', 'namespace' => 'ExternalPay', 'middleware' => ['auth']], function () {
    //物流列表
    Route::post('logistics-company', ['uses' => 'LogisticsController@logisticsCompany']);
    //导购订单列表
    Route::post('order-list', ['uses' => 'OrderController@orderList']);
    //导购退款订单列表
    Route::post('refund-order-list', ['uses' => 'OrderController@refundOrderList']);
    //订单发货
    Route::post('deliver-order', ['uses' => 'OrderController@deliverOrder']);
    //订单详情
    Route::post('order-detail', ['uses' => 'OrderController@orderDetail']);
    //退款详情
    Route::post('refund-order-detail', ['uses' => 'OrderController@refundOrderDetail']);
    //快递信息
    Route::post('order-logistics_track', ['uses' => 'LogisticsController@logisticsTrack']);
});

//用户端企微对外收款订单
Route::group(['prefix' => 'user-external-pay', 'namespace' => 'ExternalPay', 'middleware' => ['reserveAccess']], function () {
    //单列表
    Route::post('order-list', ['uses' => 'UserOrderController@orderList']);
    //退款订单列表
    Route::post('refund-order-list', ['uses' => 'UserOrderController@refundOrderList']);
    //查看物流
    Route::post('order-logistics_track', ['uses' => 'UserOrderController@logisticsTrack']);
    //订单详情
    Route::post('order-detail', ['uses' => 'UserOrderController@orderDetail']);
    //退款详情
    Route::post('refund-order-detail', ['uses' => 'UserOrderController@refundOrderDetail']);
});


//用户端企微对外收款订单
Route::group(['prefix' => 'chat-group', 'namespace' => 'ChatGroup', 'middleware' => ['auth']], function () {
    //社群成员明细
    Route::post('members_info', ['uses' => 'ChatGroupController@membersInfo']);
    //社群详情
    Route::post('index_list', ['uses' => 'ChatGroupController@indexList']);
    //数据更新时间
    Route::post('data-update-time', ['uses' => 'ChatGroupController@dataUpdateTime']);
});

//成长激励(We Journey)
Route::group(['prefix' => 'excitation', 'namespace' => 'Excitation', 'middleware' => ['auth']], function () {
    //首页
    Route::post('index', ['uses' => 'ExcitationController@index']);
    //导购等级信息
    Route::post('grade-points', ['uses' => 'ExcitationController@gradePoints']);
    //积分明细
    Route::post('points-list', ['uses' => 'ExcitationController@pointsList']);
    //积分规则说明
    Route::post('integral-rule-description', ['uses' => 'ExcitationController@integralRuleDescription']);
    //获取积分规则信息
    Route::post('get-points-rules', ['uses' => 'ExcitationController@getPointsRules']);
    //获取等级变动记录
    Route::post('level-record', ['uses' => 'ExcitationController@levelRecord']);
    //金币明细
    Route::post('coins-list', ['uses' => 'ExcitationController@coinsList']);
    //金币配置相关
    Route::post('coins-settings-info', ['uses' => 'ExcitationController@coinsSettingsInfo']);

    // 奖章
    Route::group(['prefix' => 'medal'], function () {
        // 奖章墙
        Route::post('medal-wall', ['uses' => 'MedalController@MedalWall']);
        // 奖章墙详情
        Route::post('medal-wall-detail', ['uses' => 'MedalController@MedalWallDetail']);
        // 领取奖章数量和弹框
        Route::post('record-already-got', ['uses' => 'MedalController@RecordAlreadyGot']);
        // 领取奖章
        Route::post('record-member-received', ['uses' => 'MedalController@RecordMemberReceived']);
    });
});

// 库存查询
Route::group(['prefix' => 'inventory', 'namespace' => 'Inventory', 'middleware' => ['auth']], function () {
    // 库存查询列表
    Route::post('list', ['uses' => 'IndexController@index']);
    // 库存查询的搜索记录
    Route::post('search-keywords', ['uses' => 'IndexController@searchKeywords']);
    // 删除库存查询的关键字
    Route::post('del-keywords', ['uses' => 'IndexController@delKeywords']);
});

// SW库存查询
Route::group(['prefix' => 'inventory-sw', 'namespace' => 'Inventory', 'middleware' => ['auth']], function () {
    // 库存查询列表
    Route::post('list', ['uses' => 'IndexSwController@index']);
    // 库存查询的搜索记录
    Route::post('search-keywords', ['uses' => 'IndexSwController@searchKeywords']);
    // 删除库存查询的关键字
    Route::post('del-keywords', ['uses' => 'IndexSwController@delKeywords']);
});