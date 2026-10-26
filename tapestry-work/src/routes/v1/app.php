<?php

/**
 * 首页
 */
$router->group(['prefix' => 'home', 'namespace' => 'Home'], function () use ($router) {
    // 首页统计数据
    $router->get('statistics', ['uses' => 'StatisticsController@index']);
    // 统计数据趋势图
    $router->get('statistics_map', ['uses' => 'StatisticsMapController@index']);
    // 首页模块数据统计
    $router->get('count', ['uses' => 'CountController@index']);
    // 首页模块业绩统计更新
    $router->post('order_count_update', ['uses' => 'OrderCountUpdateController@index']);
});

/**
 * 员工相关
 */
$router->group(['prefix' => 'member', 'namespace' => 'Member'], function () use ($router) {
    // 登录接口
    $router->post('login', ['uses' => 'LoginController@index']);
    // 平台小程序(天鲟CRS)登录接口
    $router->post('hub_login', ['uses' => 'HubLoginController@index']);
    // 模拟登录接口，前端无法在模拟器中使用 wx.qy.login 的情况
    $router->post('mock', ['uses' => 'MockController@index']);
});

/**
 * 用户分组相关
 */
$router->group(['prefix' => 'user-group', 'namespace' => 'UserGroup'], function () use ($router) {
    // 新建分组
    $router->post('create', ['uses' => 'CreateController@index']);
    // 分组详情
    $router->get('{groupId}/detail', ['uses' => 'DetailController@index']);
    // 分组列表
    $router->get('list', ['uses' => 'ListController@index']);
    // 分组编辑
    $router->post('{groupId}/edit', ['uses' => 'EditController@index']);
    // 分组删除
    $router->post('{groupId}/delete', ['uses' => 'DeleteController@index']);
    // 分组批量删除
    $router->post('delete', ['uses' => 'DeleteMoreController@index']);
    // 分组添加客户
    $router->post('add-customer', ['uses' => 'AddCustomerController@index']);
    // 分组删除客户
    $router->post('delete-customer', ['uses' => 'DeleteCustomerController@index']);
    // 修改分组昵称
    $router->post('rename', ['uses' => 'RenameController@index']);
});

/**
 * 客户管理
 */
$router->group(['prefix' => 'customer', 'namespace' => 'Customer'], function () use ($router) {
    // 待添加客户列表
    $router->get('addition_list', ['uses' => 'AdditionListController@index']);
    // 已添加客户列表
    $router->get('added_list', ['uses' => 'AddedListController@index']);
    // 待添加客户数量
    $router->get('addition_total', ['uses' => 'AdditionTotalController@index']);
    // 客户详情
    $router->get('detail', ['uses' => 'DetailController@index']);
    // 标记客户手机已被复制
    $router->get('mark_user', ['uses' => 'MarkUserController@index']);

    // 业绩相关
    $router->group(['prefix' => 'achievement'], function () use ($router) {
        // 员工业绩列表
        $router->get('list', ['uses' => 'AchievementListController@index']);
    });
});

/**
 * crs小程序 商城小程序开通状态
 */
$router->group(['prefix' => 'system', 'namespace' => 'System'], function () use ($router) {
    // 添加标签
    $router->post('system-status', ['uses' => 'SystemStatusController@index']);
});

/**
 * 店长端
 */
$router->group(['prefix' => 'manager', 'namespace' => 'Manager'], function () use ($router) {
    // 获取个人信息
    $router->get('me', ['uses' => 'MeController@index']);
    // 获取员工列表
    $router->post('member', ['uses' => 'MemberController@index']);
    // 员工详情接口
    $router->get('member-detail', ['uses' => 'MemberDetailController@index']);
    // 获取二维码
    $router->get('get_qr_code', ['uses' => 'QrCodeController@index']);
});

/**
 * 员工内购券
 */
$router->group(['prefix' => 'staff-coupon', 'namespace' => 'StaffCoupon'], function () use ($router) {
    // 预览分享
    $router->post('preview-share', ['uses' => 'PreviewShareController@index']);
    // 分享
    $router->post('share', ['uses' => 'ShareController@index']);
});

/**
 * work小程序内购订单
 */
$router->group(['prefix' => 'staff_order', 'namespace' => 'StaffOrder'], function () use ($router) {
    // 员工内购订单业绩排行（废弃）
    $router->post('rank', ['uses' => 'AchievementRankController@index']);
    // 当前员工内购订单业绩统计（废弃）
    $router->post('statistics', ['uses' => 'CurrentAchievementController@index']);
    // 订单业绩排行
    $router->post('new_rank', ['uses' => 'AchievementRankListController@index']);
    // 订单列表
    $router->post('list', ['uses' => 'OrderListController@index']);
    // 我的业绩
    $router->post('my_staff', ['uses' => 'MyStaffController@index']);
});

/**
 * work小程序内购员工相关
 */
$router->group(['prefix' => 'staff-right', 'namespace' => 'StaffRight'], function () use ($router) {
    // 判断员工是否有权限
    $router->post('is-right', ['uses' => 'StaffBuyController@index']);
});


# 资金结算
$router->group(['prefix' => 'settlement', 'namespace' => 'Settlement'], function() use($router) {
    // 测试接口
    $router->group(['prefix' => 'test'], function() use($router) {
        // 根据项目编码获取项目信息
        $router->post('get-project-by-code', ['uses' => 'TestController@getProjectByCode']);
        // 获取项目列表
        $router->post('get-project-list', ['uses' => 'TestController@getProjectList']);
        // 用户认证
        $router->post('sign-agreement', ['uses' => 'TestController@signAgreement']);
        // 获取认证结果
        $router->post('query-agreement', ['uses' => 'TestController@queryAgreement']);
        // 发起结算
        $router->post('settle-balance', ['uses' => 'TestController@settleBalance']);
        // 获取结算单信息
        $router->post('get-settle-balance', ['uses' => 'TestController@getSettleBalance']);
        // 同步结算单信息
        $router->post('sync-settlement-info', ['uses' => 'TestController@syncSettlementInfo']);
        // 结算单状态变更回调
        $router->post('settlement-callback', ['uses' => 'TestController@settlementCallback']);
        // 结算单退款回调
        $router->post('settlement-refund-callback', ['uses' => 'TestController@settlementRefundCallback']);
    });
});
$router->group(['prefix' => 'distribution-partner', 'namespace' => 'DistributionPartner'], function() use($router) {
    // 合伙人二维码参数
    $router->post('partner-qr-card-data', 'PartnerQrCardDataController@index');
    // 分销合伙人列表
    $router->post('list', 'ListController@index');
    // 合伙人禁用详情页
    $router->post('disable-detail', 'DisableDetailController@index');
});
