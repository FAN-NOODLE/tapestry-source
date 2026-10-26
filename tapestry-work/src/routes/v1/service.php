<?php

use Laravel\Lumen\Routing\Router;

/**
 * @var Router $router
 */

/**
 * 佣金设置
 */


$router->group(['prefix' => 'commission', 'namespace' => 'Commission'], function () use ($router) {
    // 详情
    $router->get('detail', ['uses' => 'DetailController@index']);
    // 保存
    $router->post('save', ['uses' => 'SaveController@index']);
    // 企业导购佣金设置初始化
    $router->get('init-commission', ['uses' => 'InitCommissionController@index']);
});

/**
 * 客户管理
 */
$router->group(['prefix' => 'customer', 'namespace' => 'Customer'], function () use ($router) {
    // 列表
    $router->get('list', ['uses' => 'ListController@index']);
    // 详情
    $router->get('detail', ['uses' => 'DetailController@index']);
    // 获取客户备注名
    $router->get('remark', ['uses' => 'RemarkController@index']);
    // 离职员工待分配/已分配客户列表
    $router->get('handover-list', ['uses' => 'HandoverListController@index']);
    // 在职员工待分配客户/待确认客户
    $router->post('handover-count', ['uses' => 'HandoverCountController@index']);
    // 交接客户
    $router->post('handover', ['uses' => 'HandoverController@index']);
    // 在职员工待分配客户/待确认客户
    $router->post('handover-count', ['uses' => 'HandoverCountController@index']);
    // 离职员工客户交接状态
    $router->post('handover-status', ['uses' => 'HandoverStatusController@index']);
    // 客户图表
    $router->post('chart', ['uses' => 'ChartController@index']);
    // 统计客户数
    $router->post('statistical', ['uses' => 'StatisticalController@index']);
    // 员工标记外部联系人信息数据获取
    $router->post('member-relation', ['uses' => 'MemberRelationController@index']);
    // 获取客户分配未确认
    $router->post('already-allocated', ['uses' => 'AlreadyAllocatedController@index']);
    // 根据数据权限列表
    $router->post('store', ['uses' => 'StoreCustomerController@index']);
    // 根据客户信息查找所属导购(员工)
    $router->post('relative_member', ['uses' => 'RelativeMemberController@index']);
    // 待添加客户列表
    $router->post('addition_list', ['uses' => 'AdditionListController@index']);
    // 标记客户手机已被复制
    $router->post('mark_user', ['uses' => 'MarkUserController@index']);
    // 获取user_id
    $router->post('uids', ['uses' => 'ListController@userIdList']);

    // 客户列表
    $router->post('member-list', ['uses' => 'MemberListController@memberList']);
    // 员工用户互删数据
    $router->post('user-del', ['uses' => 'StaffUserDelController@staffUserDel']);
    // 客户基本信息
    $router->post('user-detail', ['uses' => 'BaseDetailController@userDetail']);
    // 员工列表接口
    $router->post('staff-list', ['uses' => 'StaffListController@staffList']);
    // 客户统计
    $router->post('member-calc', ['uses' => 'MemberCalcController@memberCalc']);
    // 客户统计导出
    $router->post('member-out-data', ['uses' => 'MemberOutDataController@memberOutData']);
    // 给客户打标签
    $router->post('member-give-tag', ['uses' => 'MemberGiveTagController@memberGiveTag']);
    // 顾客列表（小程序端使用）
//    $router->post('customer-list', ['uses' => 'CustomerListController@index']);
    // 顾客搜索（小程序端使用）
    $router->post('search', ['uses' => 'CustomerListController@search']);
    // 顾客详情
//    $router->get('customer-detail', ['uses' => 'CustomerDetailController@index']);
    // 智慧导购顾客详情
    $router->post('customer-wisdom-detail', ['uses' => 'CustomerWisdomDetailController@index']);

    // (新)顾客列表（小程序端使用）
    $router->post('customer-list', ['uses' => 'CustomerExternalController@index']);
    // (新)顾客详情（小程序端使用）
    $router->get('customer-detail', ['uses' => 'CustomerExternalController@detail']);

    // 智慧导购客户消费记录列表
    $router->post('customer-expense-record', ['uses' => 'CustomerExpenseRecordListController@index']);

    // 导购对应顾客呼起对话详情
    $router->post('conversation-detail', ['uses' => 'ConversationDetailController@index']);
});

/**
 * 任务创建
 */
$router->group(['prefix' => 'task', 'namespace' => 'Task'], function () use ($router) {
    // 添加
    $router->post('add', ['uses' => 'AddController@index']);
    // 列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 删除
    $router->post('delete', ['uses' => 'DeleteController@index']);
    // 详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 获取员工
    $router->post('task-member', ['uses' => 'MemberController@index']);
    // 达标统计
    $router->post('task-statistics', ['uses' => 'StatisticsForAdminController@index']);
    // 客户达标统计
    $router->post('member_count', ['uses' => 'MemberCountController@index']);
    // 任务类型列表
    $router->post('task-type', ['uses' => 'TaskTypeListController@index']);
    // 指标达成统计
    $router->post('target', ['uses' => 'TargetController@index']);
});

/**
 * 部门管理
 */
$router->group(['prefix' => 'department', 'namespace' => 'Department'], function () use ($router) {
    // 列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 新增
    $router->post('add', ['uses' => 'AddController@index']);
    // 详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 编辑
    $router->post('edit', ['uses' => 'EditController@index']);
    // 删除
    $router->post('delete', ['uses' => 'DeleteController@index']);
    // 部门类型
    $router->post('type-list', ['uses' => 'TypeController@index']);
    // 获取门店dptId
    $router->post('store-dptId', ['uses' => 'StoreDptIdController@index']);
    // 部门列表
    $router->post('lists', ['uses' => 'ListController@lists']);
    // 部门列表（micro-scrm表）
    $router->post('list-by-scrm', ['uses' => 'ListByScrmController@index']);
    // 通过部门id获取相关的所有员工
    $router->post('member-list-by-dp-ids', ['uses' => 'MemberListByDpIdsController@index']);
    // 通过员工信息获取部门列表
    $router->post('list-by-member-info', ['uses' => 'ListByMemberInfoController@index']);

});

/**
 * 员工管理
 */
$router->group(['prefix' => 'member', 'namespace' => 'Member'], function () use ($router) {
    // 列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 列表
    $router->post('list-by-dp-ids', ['uses' => 'ListByDpIdsController@index']);
    // 新增
    $router->post('add', ['uses' => 'AddController@index']);
    // 详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 编辑
    $router->post('edit', ['uses' => 'EditController@index']);
    // 删除
    $router->post('delete', ['uses' => 'DeleteController@index']);
    // 员工二维码
    $router->post('qrcode', ['uses' => 'QrCodeController@index']);
    // 员工身份校验
    $router->post('check', ['uses' => 'CheckController@index']);
    // 员工联系我二维码
    $router->post('contact_qr_code', ['uses' => 'ContactWayQrCodeController@index']);
    // 导购分配员工列表
    $router->post('guide_list', ['uses' => 'GuideController@index']);
    // 转交客户给其他导购
    $router->post('customer_transfer', ['uses' => 'CustomerTransferController@index']);
    // 其他导购列表
    $router->post('other_guide_list', ['uses' => 'OtherGuidesController@index']);
    // 导购客户列表
    $router->post('guide_users', ['uses' => 'GuideUsersController@index']);
    // 没有部门员工列表
    $router->post('no-dp-list', ['uses' => 'NoDpListController@index']);
    // 添加二维码到uc以(本地)
    $router->post('add_qr_code', ['uses' => 'AddQrCodeController@index']);
    // 根据userId获取企业人员详情接口
    $router->post('detail_by_userid', ['uses' => 'DetailByUserIdController@index']);
    // 导购数据概览-顾客数据
    $router->post('customer-statistic', ['uses' => 'CustomerStatisticController@index']);
    // 导购顾客数据排行榜
    $router->post('customer-rank', ['uses' => 'CustomerStatisticController@rank']);
    // 数据罗盘筛选获取导购列表
    $router->post('statistic-rank-guides', ['uses' => 'CustomerStatisticController@guideList']);
    // 管理后台数据概览列表 按日期
    $router->post('customer-statistic-date', ['uses' => 'CustomerStatisticController@customerStatisticsByDate']);
    // 管理后台数据概览列表 按导购
    $router->post('customer-statistic-guide', ['uses' => 'CustomerStatisticController@customerStatisticsByGuide']);
    // 管理后台数据概览折线图
    $router->post('customer-statistic-console', ['uses' => 'CustomerStatisticController@customerStatisticsConsole']);

    // 导购和部门列表
    $router->post('get-members-list', ['uses' => 'ListDataController@index']);
    // 员工身份校验
    $router->post('check-source', ['uses' => 'CheckController@CheckSource']);
});

/**
 * 员工角色
 */
$router->group(['prefix' => 'role', 'namespace' => 'Role'], function () use ($router) {
    // 列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 新增
    $router->post('add', ['uses' => 'AddController@index']);
    // 详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 编辑
    $router->post('edit', ['uses' => 'EditController@index']);
    // 删除
    $router->post('delete', ['uses' => 'DeleteController@index']);
});

/**
 * 外部联系人
 */
$router->group(['prefix' => 'external', 'namespace' => 'External'], function () use ($router) {
    // 列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 导入
    $router->post('import', ['uses' => 'ImportController@index']);
    // 详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 更新
    $router->post('update', ['uses' => 'UpdateController@index']);
    // 离职成员的外部联系人再分配
    $router->post('transfer', ['uses' => 'TransferController@index']);
    // 通过外userid获取外部联系人列表
    $router->post('listByUserid', ['uses' => 'ListByUserIdController@index']);
});

/**
 * UC企业相关
 */
$router->group(['prefix' => 'enterprise', 'namespace' => 'Enterprise'], function () use ($router) {
    // 根据corpId查询企业信息
    $router->post('detail-corp-id', 'CorpIdController@index');
    // UC回调地址
    $router->post('uc-callback', 'UcCallbackController@index');
    // 添加企业信息
    $router->post('add', 'AddController@index');
    // 编辑企业信息
    $router->post('update', 'UpdateController@index');
    // 扫码登录 根据auth_code获取信息
    $router->post('auth-code-info', 'AuthCodeInfoController@index');
    // 用token查询登录的企业信息
    $router->post('wx-suite-login-info', 'WxSuiteLoginInfoController@index');
    // 企业微信登陆
    $router->post('qw-login', 'QWLoginController@index');
    // 获取企业配置信息
    $router->post('config-info', 'ConfigInfoController@index');
    // 扫码登录 根据auth_code获取信息
    $router->post('auth-code-info', 'AuthCodeInfoController@index');
    // 用token查询登录的企业信息
    $router->post('auth-admin-info', 'AuthAdminInfoController@index');
});

/**
 * UC企业应用相关
 */
$router->group(['prefix' => 'enterprise-plugin', 'namespace' => 'EnterprisePlugin'], function () use ($router) {
    // 扩展应用安装，包括： 小程序、外部联系人
    $router->post('ext-install', 'ExtInstallController@index');
    // 商城小程序安装
    $router->post('ext-shop-install', 'ExtShopInstallController@index');
    // 扩展应用更新，包括： 小程序、外部联系人
    $router->post('ext-update', 'ExtUpdateController@index');
    // 商城小程序更新
    $router->post('ext-shop-update', 'ExtShopUpdateController@index');
    // 企业微信小程序列表
    $router->post('mini-app-list', 'MiniAppListController@index');
    // 智慧导购安装
    $router->post('dao-gou-install', 'DaoGouInstallController@daoGouInstall');
    // 查询企业安装的应用
    $router->post('list', 'ListController@index');
});

/**
 * 企业扩展
 */
$router->group(['prefix' => 'enterprise-extend', 'namespace' => 'EnterpriseExtend'], function () use ($router) {
    // 更改企业分配客户状态
    $router->post('change_allocation_status', 'ChangeAllocationStatusController@index');
});

/**
 * 员工扩展
 */
$router->group(['prefix' => 'member-extend', 'namespace' => 'MemberExtend'], function () use ($router) {
    // 保存员工头像
    $router->post('save', 'SaveController@index');
    // 批量查询员工头像
    $router->get('list', 'ListController@index');
    // 更新员工分配客户状态/上限
    $router->post('change_allocation_status', 'ChangeAllocationStatusController@index');
});


/**
 * 公共方法
 */
$router->group(['prefix' => 'common', 'namespace' => 'Common'], function () use ($router) {
    // 上传域名校验文件
    $router->post('upload', 'UploadController@upload');
    // 获取企业商城小程序与绑定小程序列表
    $router->post('app_list', 'AppListController@index');
    // 生成海报太阳码
    $router->post('create', 'CreateController@index');
    // 生成太阳码
    $router->post('create-qrcode', 'CreateController@createQrcode');
    // 生成太阳码和短链
    $router->post('get-qrcode-and-link', 'CreateController@getQRCodeAndLink');
    //记录引导页浏览记录
    $router->post('boot_log_create', 'BootPageLogController@index');
    //检查是否需要导购引导页
    $router->post('boot_log_check', 'BootPageLogController@check');
});

/**
 * 欢迎语
 */
$router->group(['prefix' => 'welcome_speech', 'namespace' => 'WelcomeSpeech'], function () use ($router) {
    // 新增欢迎语
    $router->post('add', 'AddController@index');
    // 编辑欢迎语
    $router->post('edit', 'EditController@index');
    // 欢迎语详情
    $router->get('detail', 'DetailController@index');
    // 列表
    $router->get('list', 'ListController@index');
    // 删除
    $router->post('delete', 'DeleteController@index');
});

/**
 * 联系我二维码
 */
$router->group(['prefix' => 'contact', 'namespace' => 'Contact'], function () use ($router) {
    // 二维码列表
    $router->post('list', 'ListController@index');
    // 员工列表
    $router->post('member_tree', 'ContactController@memberTree');
    // 创建二维码
    $router->post('create_qr_code', 'ContactController@createQrCode');
    // 二维码详情
    $router->post('select_qr_code', 'ContactController@selectQrCode');
    // 编辑二维码
    $router->post('update_qr_code', 'ContactController@updateQrCode');
    // 删除二维码
    $router->post('delete_qr_code', 'ContactController@deleteQrCode');
    // 同步组织结构
    $router->post('sync_org', 'SyncController@syncOrg');
    // 上次同步时间
    $router->get('sync_last_time', 'SyncController@syncLastTime');
    // 同步状态
    $router->get('sync_status', 'SyncController@status');
    // 获取员工联系人二维码
    $router->get('get_code', 'SyncController@getContactCode');
    //联系人二维码统计明细
    $router->post('statistics-detail', 'StatisticsDetailController@index');
    // 查询指定二维码
    $router->post('get_qr_code', 'ContactController@getQrCode');
    // 测试接口
    $router->post('test', 'SyncController@test');

    // 更新海报路径
    $router->post("update-poster-path", "ContactController@updatePoster");
});

/**
 * 渠道管理
 */
$router->group(['prefix' => 'channel', 'namespace' => 'Channel'], function () use ($router) {
    // 渠道列表
    $router->post('list', 'ListController@index');
    // 新增渠道
    $router->post('add', 'AddController@index');
    // 编辑渠道
    $router->post('edit', 'EditController@index');
    // 渠道统计列表
    $router->post('statistics-list', 'StatisticsListController@index');
    // 渠道统计明细
    $router->post('statistics-detail', 'StatisticsDetailController@index');
});

/**
 * 用户分组
 */
$router->group(['prefix' => 'customer-group', 'namespace' => 'CustomerGroup'], function () use ($router) {
    // 用户分组列表
    $router->post('list', 'ListController@index');
    // 新增用户分组
    $router->post('add', 'AddController@index');
    // 用户分组详情
    $router->post('detail', 'DetailController@index');
    // 用户分组删除
    $router->post('del', 'DeleteController@index');
    // 编辑用户分组
    $router->post('edit', 'EditController@index');
    // 测试关系更新
    $router->post('fill', 'FillController@index');
});

/**
 * 标签管理
 */
$router->group(['prefix' => 'tag', 'namespace' => 'Tag'], function () use ($router) {
    // 新增标签
    $router->post('add', 'AddController@index');
    // 企业标签组列表
//    $router->post('all-list', 'ListController@allList');
    $router->post('all-list', 'ListController@allGroupList');
    // 企业标签库标签列表
    $router->post('select-list', 'ListController@selectList');
    // 企业标签组列表
    $router->post('list', 'ListController@index');
    // 企业标签列表
    $router->post('page-list', 'ListController@pageList');
    // 标签删除
    $router->post('del', 'DeleteController@index');
    // 编辑标签
    $router->post('edit', 'EditController@index');
    // 关联标签和用户分组
    $router->post('tag-customer', 'TagCustomerGroupRelationController@index');
    // 同步企业标签库
    $router->post('tag-sync', 'SyncTagController@index');
});

/**
 * 企业微信红包
 */
$router->group(['prefix' => 'red_packet', 'namespace' => 'RedPacket'], function () use ($router) {
    // 红包列表
    $router->get('list', 'RedPacketController@list');
    // 红包对账列表
    $router->get('reconciliation_list', 'RedPacketController@reconciliationList');
    // 红包发送人列表
    $router->post('members', 'RedPacketController@members');
    // 红包上下架
    $router->post('sale', 'RedPacketController@sale');
    // 创建红包
    $router->post('create', 'RedPacketController@create');
    // 红包详情
    $router->post('detail', 'RedPacketController@detail');
    // 更新红包
    $router->post('update', 'RedPacketController@update');
    // 删除红包
    $router->post('delete', 'RedPacketController@delete');
});

/**
 * 员工内购订单
 */
$router->group(['prefix' => 'staff_order', 'namespace' => 'StaffOrder'], function () use ($router) {
    // 新增员工内购订单
    $router->post('add', 'AddController@index');
    // 更新员工内购订单
    $router->post('update', 'UpdateController@index');
    // 员工内购订单列表
    $router->post('list', 'ListController@index');
    // 员工内购订单统计
    $router->post('statistics', 'StatisticsForAdminController@index');
    // 员工内购业绩列表
    $router->post('achievement_list', 'AchievementListController@index');
    // 员工内购部门业绩列表
    $router->post('achievement_dp_list', 'AchievementDpListController@index');
    // 我的业绩
    $router->post('my_staff', ['uses' => 'MyStaffController@index']);
});

/**
 * 员工内购
 */
$router->group(['prefix' => 'staff-manager', 'namespace' => 'StaffBuyManager'], function () use ($router) {
    // 获取内购员工信息列表
    $router->post('member', 'MemberController@index');
    // 获取内购员工/部门列表
    $router->post('list', 'ListController@index');
    // 获取内购员工/部门设置信息
    $router->post('setting', 'SettingController@index');
    // 添加内购员工/部门列表及设置
    $router->post('add', 'AddController@index');
    // 修改内购员工/部门列表及设置
    $router->post('save', 'UpdateController@index');
    // 判断员工是否是内购员工
    $router->post('is-right', 'StaffIsBuyController@index');
    // 查询内购员工对应部门(取第一个)
    $router->post('staff_department_relation', 'RelationController@index');
});
/**
 * 员工内购订单
 */
$router->group(['prefix' => 'staff_order', 'namespace' => 'StaffOrder'], function () use ($router) {
    // 新增员工内购订单
    $router->post('add', 'AddController@index');
    // 更新员工内购订单
    $router->post('update', 'UpdateController@index');
    // 员工内购订单列表
    $router->post('list', 'ListController@index');
    // 员工内购订单统计
    $router->post('statistics', 'StatisticsForAdminController@index');
    // 员工内购业绩列表
    $router->post('achievement_list', 'AchievementListController@index');
    // 员工内购部门业绩列表
    $router->post('achievement_dp_list', 'AchievementDpListController@index');
    // 我的业绩
    $router->post('my_staff', ['uses' => 'MyStaffController@index']);
});

/**
 * 员工内购
 */
$router->group(['prefix' => 'staff-manager', 'namespace' => 'StaffBuyManager'], function () use ($router) {
    // 获取内购员工信息列表
    $router->post('member', 'MemberController@index');
    // 获取内购员工/部门列表
    $router->post('list', 'ListController@index');
    // 获取内购员工/部门设置信息
    $router->post('setting', 'SettingController@index');
    // 添加内购员工/部门列表及设置
    $router->post('add', 'AddController@index');
    // 修改内购员工/部门列表及设置
    $router->post('save', 'UpdateController@index');
    // 判断员工是否是内购员工
    $router->post('is-right', 'StaffIsBuyController@index');
    // 查询内购员工对应部门(取第一个)
    $router->post('staff_department_relation', 'RelationController@index');
});

/**
 * 导购管理
 */
$router->group(['prefix' => 'guide-manage', 'namespace' => 'GuideManage'], function () use ($router) {

    // 添加导购
    $router->post('add', 'AddController@index');
    // 移除导购
    $router->post('delete', 'DeleteController@index');
    // 导购基础信息列表
    $router->post('guide-base-list', 'GuideBaseLIstController@index');

    $router->post('get-relation-by-memId', 'GetMemberRelationController@index');
    // 获取第三方企业认证状态
    $router->post('third-worker-auth-status', 'ThirdWorkerAuthStatusController@index');
});

/**
 * 企业付款到零钱
 */
$router->group(['prefix' => 'enterprise-payment', 'namespace' => 'EnterprisePayment'], function () use ($router) {

    // 充值
    $router->post('recharge', 'RechargeController@index');
    // 充值明细
    $router->post('recharge-record', 'RechargeRecordController@index');
    // 发送红包
    $router->post('send-red-pack', 'SendRedPackController@index');
    // 发送红包详情
    $router->post('receive-detail', 'ReceiveDetailController@index');
    // 红包对账统计
    $router->post('budget-total', 'BudgetTotalController@index');
    // 新增红包活动
    $router->post('add', 'AddController@index');
    // 编辑红包活动
    $router->post('save', 'SaveController@index');
    // 活动详情
    $router->post('detail', 'DetailController@index');
    // 红包活动列表
    $router->post('list', 'ListController@index');
    // 红包活动对应对账记录
    $router->post('reconciliation-list', 'ReconciliationListController@index');
    // 更新红包活动状态
    $router->post('status', 'StatusController@index');
    // 获取活动对应已发送红包员工
    $router->post('reconciliation-member', 'ReconciliationMemberController@index');
    // 红包对账
    $router->post('budget-by-enterprise', 'BudgetByEnterpriseController@index');
});

/**
 * 素材管理
 */
$router->group(['prefix' => 'source', 'namespace' => 'Source'], function () use ($router) {
    // 更新素材临时Id
    $router->post('update_media_id', ['uses' => 'UpDataMediaIdController@index']);
});

/**
 * 营销推广
 */
$router->group(['prefix' => 'marketing-promotion', 'namespace' => 'MarketingPromotion'], function () use ($router) {
    // 获取海报
    $router->post('poster', ['uses' => 'PosterController@index']);
    // 发送应用消息
    $router->post('send-mini-message', ['uses' => 'SendMiniMessageController@index']);
});


/**
 * 管理员
 */
$router->group(['prefix' => 'admin', 'namespace' => 'Admin'], function () use ($router) {
    // 详情
    $router->post('detail', ['uses' => 'DetailController@index']);

    $router->post('list', ['uses' => 'DetailController@list']);
});

/**
 * 订单
 */
$router->group(['prefix' => 'order', 'namespace' => 'Order'], function () use ($router) {
    // 详情
    $router->post('sync', ['uses' => 'SyncController@index']);
    // 订单统计排行榜
    $router->post('order-ranking', ['uses' => 'OrderStatisticalRankingController@index']);
    // 订单数据统计
    $router->post('order-statistical', ['uses' => 'OrderStatisticalController@index']);
    // 订单列表
    $router->post('status-list', ['uses' => 'StatusListController@index']);
    // 查询导购与合伙人的分类订单数据
    $router->post('order-statistical-classification', ['uses' => 'OrderStatisticalClassificationController@index']);
});
/**
 * 导购业绩设置
 */
$router->group(['prefix' => 'guide-achievement-setting', 'namespace' => 'GuideAchievementSetting'], function () use ($router) {
    // 设置详情
    $router->post('detail', 'DetailController@index');
    // 设置保存
    $router->post('save', 'SaveController@index');
});


/**
 * Auth
 */
$router->group(['prefix' => 'auth', 'namespace' => 'Auth'], function () use ($router) {
    // 获取配置信息
    $router->post('agent_config', ['uses' => 'AgentConfigController@index', 'as' => 'service-agent-config']);
});

/**
 *  用户自定义标签
 */
$router->group(['prefix' => 'customer-tags', 'namespace' => 'CustomerTags'], function () use ($router) {
    // 绑定用户自定义标签
    $router->post('add-relation', ['uses' => 'AddController@addNewRelation']);
    // 添加用户自定义标签
    $router->post('add-tag', ['uses' => 'AddController@addNewTag']);
    // 解绑自定义标签
    $router->post('delete-relation', ['uses' => 'DeleteController@DeleteRelation']);
    // 删除店员自定义标签
    $router->post('delete-tag', ['uses' => 'DeleteController@DeleteTag']);
    // 获取用户绑定标签
    $router->post('get-relation', ['uses' => 'ListController@getRelation']);
    // 获取用户自定义标签列表
    $router->post('get-tags', ['uses' => 'ListController@getTags']);
    // 获取店员所有自定义标签
    $router->post('get-list', ['uses' => 'ListController@getList']);
    // 搜索已有的自定义标签
    $router->post('search-tags', ['uses' => 'ListController@searchTags']);
    // 获取标签分组列表
    $router->post('get-group-list', ['uses' => 'ListController@getGroupList']);
});

// 仪表盘数据
$router->group(['prefix' => 'dashboard', 'namespace' => 'Dashboard'], function () use ($router) {
    // 仪表盘数据
    $router->post('index', ['uses' => 'IndexController@index']);
});

// 分销商品
$router->group(['prefix' => 'distribution-goods', 'namespace' => 'DistributionGoods'], function () use ($router) {
    // 单个商品保存
    $router->post('save', ['uses' => 'SaveController@index']);
    // 批量商品保存
    $router->post('batch-save', ['uses' => 'BatchSaveController@index']);
    // 批量商品设置
    $router->post('batch-setting', ['uses' => 'BatchSettingController@index']);
    // 批量商品删除
    $router->post('batch-delete', ['uses' => 'BatchDeleteController@index']);
    // 批量商品保存
    $router->post('repeat-inspect', ['uses' => 'RepeatInspectController@index']);
    // 分销商品列表
    $router->post('goods-list', ['uses' => 'ListController@index']);
    // 小程序端分销商品列表
    $router->post('app-goods-list', ['uses' => 'AppListController@index']);
    // 小程序分销商品列表-sku维度
    $router->post('app-goods-sku-list', ['uses' => 'AppSkuListController@index']);
});

// 分销设置
$router->group(['prefix' => 'distribution-config', 'namespace' => 'DistributionConfig'], function () use ($router) {
    // 基础配置
    $router->post('save', ['uses' => 'SaveController@index']);
    // 查看基础设置
    $router->post('detail', ['uses' => 'DetailController@index']);
    //二级分销设置
    $router->post('partner-save', ['uses' => 'PartnerSaveController@index']);
    //二级分销查看
    $router->post('partner-detail', ['uses' => 'PartnerDetailController@index']);
});

// 分销统计
$router->group(['prefix' => 'distribution-statistics', 'namespace' => 'DistributionStatistics'], function () use ($router) {

    // 管理后台数据概览列表 按日期
    $router->post('customer-statistic-date', ['uses' => 'OverviewController@statisticsByDate']);
    // 管理后台数据概览列表 按导购
    $router->post('customer-statistic-guide', ['uses' => 'OverviewController@statisticsByGuide']);
    // 管理后台数据概览折线图
    $router->post('statistic-console', ['uses' => 'OverviewController@statisticsConsole']);

    // 智慧导购统计(概览)
    $router->post('distribution-statistic', ['uses' => 'OverviewController@distributionStatistic']);
    // 智慧导购统计排行榜
    $router->post('statistic-ranking', ['uses' => 'OverviewController@statisticRanking']);
});


// 结算记录
$router->group(['prefix' => 'distribution-settlement', 'namespace' => 'DistributionSettlement'], function () use ($router) {
    // 后台列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 后台详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 导购小程序列表
    $router->post('list-for-app', ['uses' => 'ListForAppController@index']);
});


// 分销订单结算
$router->group(['prefix' => 'distribution-order-settlement', 'namespace' => 'DistributionOrderSettlement'], function () use ($router) {
    // 单个商品保存
    $router->post('sync', ['uses' => 'SyncController@index']);
});

// 导购提现相关路由
$router->group([
    'namespace' => 'Withdrawal',
    'prefix'    => 'withdrawal'
], function () use ($router) {
    // 提现认证 - 添加
    $router->post('certification-add', ['uses' => 'CertificationAddController@index']);
    // 提现认证 - 更新
    $router->post('certification-update', ['uses' => 'CertificationUpdateController@index']);
    // 提现认证 - 审核
    $router->post('certification-audit', ['uses' => 'CertificationAuditController@index']);
    // 提现认证 - 详情
    $router->post('certification-detail', ['uses' => 'CertificationDetailController@index']);
    // 提现认证 - 获取认证文件
    $router->post('certification-file-url', ['uses' => 'CertificationDetailController@certificationFileUrl']);

    // 提现 - 申请
    $router->post('withdrawal-add', ['uses' => 'WithdrawalAddController@index']);
    // 提现 - 申请填写页面数据
    $router->post('withdrawal-home', ['uses' => 'WithdrawalHomeController@index']);
    // 提现 - 记录列表
    $router->post('withdrawal-logs', ['uses' => 'WithdrawalLogsController@index']);
    // 提现 - 银行列表
    $router->post('withdrawal-bank-list', ['uses' => 'WithdrawalBankListController@index']);
});

// 余额管理
$router->group(['prefix' => 'distribution-balance', 'namespace' => 'DistributionBalance'], function () use ($router) {

    // 列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 列表导出
    $router->post('list-export', ['uses' => 'ListExportController@index']);
    // 详情
    $router->post('details-list', ['uses' => 'DetailsListController@index']);
    // 详情导出
    $router->post('details-export', ['uses' => 'DetailsExportController@index']);
    // 导购收支相关信息
    $router->post('cash-withdrawal-info', ['uses' => 'CashWithdrawalInfoController@index']);
});

# 资金结算
$router->group(['prefix' => 'settlement', 'namespace' => 'Settlement'], function () use ($router) {
    // 结算回调
    $router->group(['prefix' => 'callback'], function () use ($router) {
        // 结算单状态变更通知
        $router->post('settlement', 'CallbackController@settleBalance');
        // 结算单退款结果通知
        $router->post('settlement-refund', 'CallbackController@refundBalance');
    });
});


$router->group(['prefix' => 'withdrawal-approval', 'namespace' => 'Withdrawal'], function () use ($router) {
    // 审批列表
    $router->post('approval-list', ['uses' => 'ApprovalListController@index']);
    //审批
    $router->post('approval-check', ['uses' => 'ApprovalCheckController@index']);
});

$router->group(['prefix' => 'withdrawal-setting', 'namespace' => 'WithdrawaSetting'], function () use ($router) {
    //提现设置
    $router->post('add', ['uses' => 'AddController@index']);
    //提现设置查看
    $router->post('detail', ['uses' => 'DetailController@index']);
    //自动提现设置
    $router->post('auto_add', ['uses' => 'AutoAddController@index']);
    //自动提现设置查看
    $router->post('auto_detail', ['uses' => 'AutoDetailController@index']);
    //根据项目编码获取项目信息
    $router->post('check', ['uses' => 'CheckController@index']);
});

/**
 * 分销合伙人管理
 */
$router->group(['prefix' => 'distribution-partner', 'namespace' => 'DistributionPartner'], function () use ($router) {
    // 分销合伙人列表
    $router->post('list', 'ListController@index');
    // 合伙人禁用详情页
    $router->post('disable-detail', 'DisableDetailController@index');
    // 合伙人我的信息
    $router->post('my', 'MyController@index');
    //合伙人短信开关
    $router->post('sms-switch', 'PartnerIsAcceptSmsController@index');
    //合伙人手机号是否脱敏展示
    $router->post('mobile-switch', 'MobileSensitiveController@index');
});

$router->group(['prefix' => 'partner-member', 'namespace' => 'PartnerMember'], function () use ($router) {
    // 合伙人基本信息表详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 扫导购邀请码 check
    $router->post('check', ['uses' => 'CheckController@index']);
    // 合伙人详情
    $router->post('detail-all', ['uses' => 'DetailAllController@index']);
    // 合伙人insert
    $router->post('insert', ['uses' => 'InsertController@index']);
});

$router->group(['prefix' => 'partner-enterprise-map', 'namespace' => 'PartnerEnterpriseMap'], function () use ($router) {
    // 合伙人导购关联表insert
    $router->post('insert', ['uses' => 'InsertController@index']);
});


$router->group(['prefix' => 'partner', 'namespace' => 'Partner'], function () use ($router) {
    // 合伙人列表
    $router->post('list', ['uses' => 'listController@index']);
    // 合伙人详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 合伙人启用禁用
    $router->post('switch', ['uses' => 'SwitchController@index']);
    // 合伙人联系导购
    $router->post('contact', ['uses' => 'partnerContactGuideController@index']);
    // 合伙人当前企业信息查询
    $router->post('check', ['uses' => 'checkPartneStatusByMobileController@index']);
});

$router->group(['prefix' => 'staff-share', 'namespace' => 'StaffShare'], function () use ($router) {
    // 链路绑定列表查询
    $router->post('list', ['uses' => 'ListController@index']);
    //链路绑定详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 获取父级绑定信息
    $router->post('get-share-parent-info', ['uses' => 'GetShareParentInfoController@index']);
    // 变更绑定关系
    $router->post('change_owner', ['uses' => 'ChangeOwnerController@index']);
    $router->post('handle_bind_user', ['uses' => 'HandleBindUserController@index']);
    // 添加绑定链路关系
    $router->post('create_staff_share_data', ['uses' => 'CreateStaffShareDataController@index']);
    // 员工离职时候解除绑定链路关系
    $router->post('unbind-share', ['uses' => 'UnbindStaffShareController@index']);
    // 成为员工解除绑定关系
    $router->post('unbind-share-own', ['uses' => 'UnbindStaffShareOwnController@index']);
    // 获取userIds
    $router->post('get-staff-share-user-ids', ['uses' => 'GetStaffShareUserIdsController@index']);
    // 查询绑定的导购信息
    $router->post('get-bind-worker', ['uses' => 'GetBindStoreWorkerController@index']);
    // 根据条件获取链路关系
    $router->post('get_staff', ['uses' => 'getStaffShareInfoByConditionController@index']);
    // 获取企业下绑定关系
    $router->post('binding_customer_detail', ['uses' => 'BindingCustomerDetailController@index']);

});

$router->group(['prefix' => 'partner-distribution', 'namespace' => 'PartnerSettlement'], function () use ($router) {
    // 结算记录
    $router->post('list', ['uses' => 'listController@index']);
    // 结算记录详情
    $router->post('detail', ['uses' => 'DetailController@index']);
});

// 事件通知
$router->group(['prefix' => 'event-notify', 'namespace' => 'EventNotify'], function () use ($router) {

    $router->group(['prefix' => 'share-link'], function () use ($router) {
        // 绑定
        $router->post('add', ['uses' => 'ShareLinkChangeEventNotifyController@add']);
        // 更新
        $router->post('update', ['uses' => 'ShareLinkChangeEventNotifyController@update']);
        // 解绑
        $router->post('delete', ['uses' => 'ShareLinkChangeEventNotifyController@delete']);
    });

});

$router->group(['prefix' => 'store-worker', 'namespace' => 'StoreWorker'], function () use ($router) {
    // 获取导购详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 获取导购详情列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 获取导购列表
    $router->post('all-list', ['uses' => 'AllListController@index']);
    // 更新导购信息
    $router->post('update', ['uses' => 'UpdateController@index']);
    // 在职离职
    $router->post('job-edit', ['uses' => 'JobEditController@index']);
    // 获取导购分享链路信息
    $router->post('share-link', ['uses' => 'ShareLinkController@index']);
    // 导购ID查询mem_uid
    $router->post('mem-uid-list', ['uses' => 'MemUidListController@index']);
    // 手动生成导购海报
    $router->post('generate-worker-poster', ['uses' => 'GenerateWorkerPosterController@index']);
    // 批量生成导购海报
    $router->post('make-poster-batch', ['uses' => 'MakePosterController@index']);
    // 根据手机号获取门店Ids
    $router->post('storeids', ['uses' => 'GetStoreIdsController@index']);
    // 新增门店导购
    $router->post('add', ['uses' => 'AddController@index']);
    // 编辑门店导购
    $router->post('edit', ['uses' => 'EditController@index']);
    // 根据条件查询对应用户信息
    $router->post('get-worker-member-info', ['uses' => 'GetWorkMemberInfoController@index']);
    // 添加导购(企业导购,社群导购)
    $router->post('add-guide', ['uses' => 'AddGuideController@index']);
    // 移除导购(企业导购,社群导购)
    $router->post('delete-guide', ['uses' => 'DeleteGuideController@index']);
    // 获取导购信息
    $router->post('detail-by-id-or-mem-uid', ['uses' => 'DetailByIdOrMemUidController@index']);
    // 商城客户与员工绑定
    $router->post('bind-user', ['uses' => 'BindUserController@index']);
    // 根据mem_uid获取员工信息
    $router->post('get-worker-mobile', ['uses' => 'GetWorkMobileController@index']);
    // 通过联系我二维码id获取导购信息
    $router->post('store-worker-detail', ['uses' => 'StoreWorkerDetailController@index']);
    // 获取导购信息
    $router->post('store-by-mobile', ['uses' => 'StoreByMobileController@index']);
    // 用户ID查询绑定导购信息
    $router->post('worker-by-user', ['uses' => 'WorkerByUserController@index']);
    // 通过导购名称查导购信息
    $router->post('get-workers-by-name', ['uses' => 'SettlementController@getWorkersByName']);
    // 批量查询导购信息
    $router->post('get-batch-workers-info', ['uses' => 'SettlementController@getBatchWorkersInfo']);
    // 根据门店或部门获取导购列表
    $router->post('guide-list', ['uses' => 'GuideListByStoreController@index']);
    // 定时更新导购信息
    $router->post('update-worker-info', ['uses' => 'UpdateWorkerInfoController@index']);
    // 根据ids获取导购列表
    $router->post('get-worker-by-ids', ['uses' => 'GetWorkerByIdsController@index']);
    // 升级企业管理
    $router->post('upgrade-admin', ['uses' => 'UpgradeAdminController@index']);
    // 同步组织架构中成员为企业导购
    $router->post('batch-add-guide', ['uses' => 'BatchAddGuideController@index']);
    // 根据职位更新用户管理权限范围
    $router->post('update-store-manager-by-position', ['uses' => 'UpdateStoreManagerByPositionController@index']);

    // open调用接口
    $router->group(['prefix' => 'open', 'namespace' => 'Open'], function () use ($router) {
        // 在职离职
        $router->post('job-edit', ['uses' => 'JobEditController@index']);
        // 更新
        $router->post('edit', ['uses' => 'EditController@index']);
        // 新增
        $router->post('add', ['uses' => 'AddController@index']);
        // 详情
        $router->post('detail', ['uses' => 'DetailController@index']);

    });

    // 业绩集成
    $router->group(['prefix' => 'performance-inherit', 'namespace' => 'PerformanceInherit'], function () use ($router) {
        // 继承导购客户列表
        $router->post('customer-list', ['uses' => 'CustomerListController@index']);
        // 全部继承
        $router->post('all-inherit', ['uses' => 'AllInheritController@index']);
        // 部分继承
        $router->post('partial-inherit', ['uses' => 'PartialInheritController@index']);
    });

    // 用户分享-解绑逻辑-员工变导购
    $router->post('staff-share-unbinding-change', ['uses' => 'StaffShareUnbindingController@unbindingChange']);
    // 用户分享-解绑逻辑-员工离职
    $router->post('staff-share-unbinding-leave', ['uses' => 'StaffShareUnbindingController@unbindingLeave']);
    // 选择门店进行绑定用户门店员工关系
    $router->post('bind-store-mem', ['uses' => 'BindStoreMemController@index']);
    // 绑定门店员工关系(带参数访问门店)
    $router->post('bind-store-mem-with-scene', 'BindStoreMemWithSceneController@index');
    // 获取员工对应导购信息
    $router->post('customer-binding-worker-detail', ['uses' => 'CustomerBindingWorkerDetailController@index']);
    // 员工内购详情
    $router->post('staff-purchase-detail', ['uses' => 'StaffPurchaseDetailController@index']);

});

//营销日历
$router->group(['prefix' => 'marketing-plan', 'namespace' => 'MarketingPlan'], function () use ($router) {
    // 创建更新
    $router->post('create', ['uses' => 'CreateController@index']);
    //更新
    $router->post('update', ['uses' => 'UpdateController@index']);
    // 详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    //删除
    $router->post('delete', ['uses' => 'DeleteController@index']);
    //列表
    $router->post('list', ['uses' => 'ListController@index']);
    //获取导购端任务列表
    $router->post('task-list', 'TaskListController@index');
    //获取导购端日历以及任务数
    $router->post('task-calendar', 'TaskCalendarController@index');
    //任务详情
    $router->post('task-detail', 'TaskDetailController@index');
    //任务完成记录
    $router->post('task-finish', 'TaskFinishController@index');
    //未完成名单
    $router->post('not-finished', ['uses' => 'NotFinishedController@index']);
});


// 数据罗盘
$router->group(['prefix' => 'data-compass', 'namespace' => 'DataCompass'], function () use ($router) {
    $router->group(['prefix' => 'customer-statistic', 'namespace' => 'CustomerStatistic'], function () use ($router) {
        // 按日期维度查询顾客数据
        $router->post('list-by-date', ['uses' => 'ListByDateController@index']);
        // 按导购维度查询顾客数据
        $router->post('list-by-worker', ['uses' => 'ListByWorkerController@index']);
        // 顾客数据总览-后台
        $router->post('overview-for-admin', ['uses' => 'OverviewForAdminController@index']);
        // 顾客数据总览-app
        $router->post('overview-for-app', ['uses' => 'OverviewForAppController@index']);
        // 顾客数据排行榜-app
        $router->post('rank-for-app', ['uses' => 'RankForAppController@index']);
        // 管理后台数据排行榜-console
        $router->post('rank-for-console', ['uses' => 'RankForConsoleController@index']);
        //财年财月筛选条件
        $router->post('fiscal-month', ['uses' => 'FiscalMonthController@index']);
    });
});

// 对接DataLake
$router->group(['prefix' => 'data-lake', 'namespace' => 'DataLake'], function () use ($router) {
    $router->post('query-member-grade', ['uses' => 'QueryMemberGradeController@index']);
    $router->post('query-member-grade-calculator', ['uses' => 'QueryMemberGradeCalculatorController@index']);
    $router->post('query-member-info', ['uses' => 'QueryMemberInfoController@index']);
    $router->post('query-member-point', ['uses' => 'QueryMemberPointController@index']);
    $router->post('query-member-points-expire', ['uses' => 'QueryMemberPointsExpireController@index']);
    $router->post('query-member-points-detail-list', ['uses' => 'QueryMemberPointsDetailListController@index']);
    $router->post('query-member-order-list', ['uses' => 'QueryMemberOrderListController@index']);
    $router->post('query-member-order-detail', ['uses' => 'QueryMemberOrderDetailController@index']);
    $router->post('query-tier-list', ['uses' => 'QueryTierListController@index']);
    $router->post('query-store-list', ['uses' => 'QueryStoreListController@index']);
    $router->post('query-ba-list', ['uses' => 'QueryBaListController@index']);
    $router->post('query-ba-list', ['uses' => 'QueryBaListController@index']);
    $router->post('query-integrated-member-info', ['uses' => 'QueryIntegratedMemberInfoController@index']);
    $router->post('member-register', ['uses' => 'MemberRegisterController@index']);

    // LD中企业微信BA/顾客关系
    $router->group(['prefix' => 'ba-bind'], function () use ($router) {
        $router->post('', ['uses' => 'BaBindController@index']);
        $router->post('query', ['uses' => 'BaBindQueryController@index']);
    });
});
// 设置企业对接DL参数
$router->group(['prefix' => 'setting', 'namespace' => 'Setting'], function () use ($router) {
    $router->post('add', ['uses' => 'AddController@index']);
    $router->post('list', ['uses' => 'ListController@index']);
    $router->post('edit', ['uses' => 'EditController@index']);
    $router->post('delete', ['uses' => 'DeleteController@index']);
});

// 任务中心
$router->group(['prefix' => 'mission-center', 'namespace' => 'MissionCenter'], function () use ($router) {
    // 创建更新
    $router->post('create', ['uses' => 'CreateController@index']);
    //更新
    $router->post('update', ['uses' => 'UpdateController@index']);
    // 详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    //删除
    $router->post('delete', ['uses' => 'DeleteController@index']);
    //列表
    $router->post('list', ['uses' => 'ListController@index']);
    //获取导购端任务列表
    $router->post('task-list', 'TaskListController@index');
    //获取导购端日历以及任务数
    $router->post('task-calendar', 'TaskCalendarController@index');
    //任务详情
    $router->post('task-detail', 'TaskDetailController@index');
    //任务完成记录
    $router->post('task-finish', 'TaskFinishController@index');
    //未完成名单
    $router->post('not-finished', ['uses' => 'NotFinishedController@index']);

    // 配置小程序
    $router->post('config-save', ['uses' => 'ConfigSaveController@index']);
    // 获取配置小程序
    $router->post('config-detail', ['uses' => 'ConfigDetailController@index']);
    // 选择客户标签列表
    $router->post('tag-list', ['uses' => 'TagListController@index']);
    // 选择和展示可见范围
    $router->post('range', ['uses' => 'RangeController@index']);
});

// 推送智慧导购小程序消息
$router->group(['prefix' => 'mini-program', 'namespace' => 'MiniProgram'], function () use ($router) {
    // 推送计划任务消息
    $router->post('plan-task', ['uses' => 'SendMessagePlanTask@index']);
    // 发送智汇导购应用消息
    $router->post('send-message', ['uses' => 'SendMessageController@index']);
});

// 三方门店及员工数据
$router->group(['prefix' => 'store', 'namespace' => 'Store'], function () use ($router) {
    // 门店列表
    $router->post('department-list', ['uses' => 'DepartmentListController@index']);
    // 员工列表
    $router->post('employee-list', ['uses' => 'EmployeeListController@index']);
    // 修改门店信息
    $router->post('save-department-info', ['uses' => 'SaveDepartmentInfoController@index']);
    // 开店闭店
    $router->post('close-or-open-department', ['uses' => 'CloseOrOpenDepartmentController@index']);
    // 删除门店
    $router->post('delete-department', ['uses' => 'DeleteDepartmentController@index']);
    // 导购分配/导购转岗
    $router->post('move-employee-department', ['uses' => 'MoveEmployeeDepartmentController@index']);
    // 同步数据到企微
    $router->post('sync-data-to-wechat', ['uses' => 'SyncDataToWechatController@index']);
    // 文件中错误信息记录下载
    $router->post('sftp-error-record-list', ['uses' => 'SftpErrorRecordListController@index']);
    // 文件拉取记录列表
    $router->post('sftp-log-list', ['uses' => 'SftpLogListController@index']);
    // 导购部门信息
    $router->post('department-detail', ['uses' => 'StoreDepartmentController@detail']);
    // 门店信息导入更新
    $router->post('batch-update-store-info', ['uses' => 'BatchUpdateStoreInfoController@index']);
    // coach业务架构门店数据更新
    $router->post('sync-data-by-departments', ['uses' => 'SyncDataByDepartmentsController@index']);
    // 获取D1M三方订单h5跳转地址
    $router->post('d1m-order-h5-url', ['uses' => 'D1MOrderH5Controller@index']);

});

// 订阅消息
$router->group(['prefix' => 'subscribe-message', 'namespace' => 'SubscribeMessage'], function () use ($router) {
    // 获取配置
    $router->post('get-config', ['uses' => 'SubscribeMessageController@getConfig']);
    // 订阅
    $router->post('subscribe', ['uses' => 'SubscribeMessageController@subscribe']);
    // 发送消息
    $router->post('send', ['uses' => 'SendMessageController@index']);
});

// 活动
$router->group([
    'namespace' => 'Activity',
    'prefix'    => 'activity'
], function () use ($router) {
    $router->group(['prefix' => 'scratch', 'namespace' => 'Scratch'], function () use ($router) {
        // 活动详情
        $router->post('info', ['uses' => 'IndexController@Info']);
        // 刮奖
        $router->post('scratch-prize', ['uses' => 'IndexController@ScratchPrize']);
        // 券列表
        $router->post('coupon-list', ['uses' => 'IndexController@CouponList']);
        // 核销券
        $router->post('write-off-coupon', ['uses' => 'IndexController@WriteOffCoupon']);
        // 用户进入活动上报
        $router->post('user-view', ['uses' => 'IndexController@userView']);
        // 检查用户的核销状态
        $router->post('check-coupon', ['uses' => 'IndexController@CheckCoupon']);
        // 用户的分享记录上报
        $router->post('share-info', ['uses' => 'IndexController@ShareInfo']);
    });
    $router->group(['prefix' => 'reserve', 'namespace' => 'Reserve'], function () use ($router) {
        // 活动详情
        $router->post('info', ['uses' => 'CIIEController@Info']);
        // 创建预约
        $router->post('create', ['uses' => 'CIIEController@Create']);
        // 取消预约
        $router->post('cancel', ['uses' => 'CIIEController@Cancel']);
        // 修改预约
        $router->post('edit', ['uses' => 'CIIEController@Edit']);
        // 用户进入活动上报
        $router->post('user-view', ['uses' => 'CIIEController@userView']);
        // 预约时间段
        $router->post('time-lists', ['uses' => 'CIIEController@Timelists']);
        // 预约时间段
        $router->post('time-config', ['uses' => 'CIIEController@TimeConfig']);
        // 全部的预约列表
        $router->post('lists', ['uses' => 'CIIEController@Lists']);
        // 导购的核销列表
        $router->post('check-lists', ['uses' => 'CIIEController@CheckLists']);
        // 核销券
        $router->post('write-off-coupon', ['uses' => 'CIIEController@WriteOffCoupon']);
        // 导购白名单
        $router->post('white-list', ['uses' => 'CIIEController@WhiteList']);
        // 获取手机号
        $router->post('get-phone-number', ['uses' => 'CIIEController@GetPhoneNumber']);
    });
});

//在职继承
$router->group(['prefix' => 'customer-inherit', 'namespace' => 'CustomerInherit'], function () use ($router) {
    //任务列表
    $router->post('list', ['uses' => 'IndexController@inheritList']);
    //任务详情
    $router->post('info', ['uses' => 'IndexController@inheritInfo']);
    //分配写入
    $router->post('transfer', ['uses' => 'IndexController@inheritTransfer']);
    //导购列表
    $router->post('member-list', ['uses' => 'IndexController@memberList']);
    //查询外部联系人转接信息
    $router->post('customer-inherit-record', ['uses' => 'IndexController@customerInheritRecord']);
    //查询外部联系人对应多个导购的信息
    $router->post('member-record', ['uses' => 'IndexController@memberRecord']);
    // 查询分配顾客成功的记录
    $router->post('record-external_ids', ['uses' => 'IndexController@recordForExternalIds']);
});

//企微对外收款
$router->group(['prefix' => 'external-pay', 'namespace' => 'ExternalPay'], function () use ($router) {
    //物流公司列表
    $router->post('logistics-company', ['uses' => 'LogisticsController@logisticsCompany']);
    //导购订单列表
    $router->post('order-list', ['uses' => 'OrderController@orderList']);
    //导购退款订单列表
    $router->post('refund-order-list', ['uses' => 'OrderController@refundOrderList']);
    //导购订单发货
    $router->post('deliver-order', ['uses' => 'OrderController@deliverOrder']);
    //订单详情
    $router->post('order-detail', ['uses' => 'OrderController@orderDetail']);
    //退款订单详情
    $router->post('refund-order-detail', ['uses' => 'OrderController@refundOrderDetail']);
    //快递信息
    $router->post('order-logistics_track', ['uses' => 'LogisticsController@logisticsTrack']);
});


$router->group([
    'namespace' => 'Work',
    'prefix'    => 'work'
], function () use ($router) {
    // 活动详情
    $router->post('index', ['uses' => 'IndexController@index']);

});

//用户端企微对外收款
$router->group(['prefix' => 'user-external-pay', 'namespace' => 'ExternalPay'], function () use ($router) {
    //导购订单列表
    $router->post('order-list', ['uses' => 'UserOrderController@orderList']);
    //导购退款订单列表
    $router->post('refund-order-list', ['uses' => 'UserOrderController@refundOrderList']);
    //订单详情
    $router->post('order-detail', ['uses' => 'UserOrderController@orderDetail']);
    //退款订单详情
    $router->post('refund-order-detail', ['uses' => 'UserOrderController@refundOrderDetail']);
});

// 神策透传转发
$router->group(['prefix' => 'shenCe', 'namespace' => 'ShenCe'], function () use ($router) {
    $router->post('itemSet', ['uses' => 'ShenCeController@itemSet']);
    $router->post('profileSet', ['uses' => 'ShenCeController@profileSet']);
    $router->post('track', ['uses' => 'ShenCeController@track']);
});