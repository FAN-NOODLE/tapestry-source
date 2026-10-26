<?php

/**
 * @var \Laravel\Lumen\Routing\Router $router
 */

/**
 * 管理员登录接口
 */
$router->group(['prefix' => 'admin', 'namespace' => 'Admin'], function () use ($router) {
    // 登录接口
    $router->post('login', ['uses' => 'LoginController@index']);
    // 退出登录
    $router->post('logout', ['uses' => 'LogoutController@index']);
    // 调试接口
    $router->get('test', ['uses' => 'TestController@index']);
    // 模拟登录
    $router->get('mock-login', ['uses' => 'MockLoginController@index']);
    // 获取图像验证码
    $router->get('image-code', ['uses' => 'LoginController@image']);
    // 检测手机验证码
    $router->post('check-code', ['uses' => 'LoginController@checkCode']);
    // 发送短信验证码
    $router->post('send-code', ['uses' => 'LoginController@sendCode']);
});

/**
 * 企业管理
 */
$router->group(['prefix' => 'enterprise', 'namespace' => 'Enterprise'], function () use ($router) {
    // 企业列表
    $router->get('list', ['uses' => 'ListController@index']);
    // 企业详情
    $router->get('detail', ['uses' => 'DetailController@index']);
    // 增加企业
    $router->post('add', ['uses' => 'AddController@index']);
    // 编辑企业
    $router->post('save', ['uses' => 'SaveController@index']);
});

/**
 * 企业管理员
 */
$router->group(['prefix' => 'enterprise_admin', 'namespace' => 'EnterpriseAdmin'], function () use ($router) {
    // 企业管理员列表
    $router->get('list', ['uses' => 'ListController@index']);
    // 企业管理员详情
    $router->get('detail', ['uses' => 'DetailController@index']);
    // 增加企业管理员
    $router->post('add', ['uses' => 'AddController@index']);
    // 编辑企业管理员
    $router->post('save', ['uses' => 'SaveController@index']);
    // 重置企业管理员密码
    $router->post('password', ['uses' => 'PasswordController@index']);
    // 重置企业管理员密码
    $router->post('reset-password', ['uses' => 'ResetPasswordController@index']);
});

/**
 * 应用管理
 */
$router->group(['prefix' => 'app-manage', 'namespace' => 'AppManage'], function () use ($router) {
    // 应用列表
    $router->get('list', ['uses' => 'ListController@index']);
    // 添加应用
    $router->post('install', ['uses' => 'InstallController@index']);
});

/**
 * 小程序版本管理
 */
$router->group(['prefix' => 'mini-program', 'namespace' => 'MiniProgram'], function () use ($router) {
    // 列表
    $router->post('list', 'ListController@index');
    // 详情
    $router->post('detail', 'DetailController@index');
    // 提审
    $router->post('submit-audit', 'SubmitAuditController@index');
    // 撤审
    $router->post('withdraw-audit', 'WithdrawAuditController@index');
    // 上线
    $router->post('release', 'ReleaseController@index');
    // 体验二维码
    $router->post('experience-qrcode', 'QrCodeController@index');
    // 版本回退
    $router->post('revert-code-release', 'RevertCodeReleaseController@index');
    // 取消分阶段发布
    $router->post('revert-gray-release', 'RevertGrayReleaseController@index');

    // 全网提审
    $router->get('fullProcedure', ['uses' => 'FullProcedureController@index']);
    //小程序手动提审
    $router->get('trig-submit', ['uses' => 'TrigSubmitController@index']);
    // 查询小程序提审状态
    $router->get('get-audit-status', ['uses' => 'GetAuditStatusController@index']);
    // 添加体验者
    $router->get('create-tester', ['uses' => 'CreateTesterController@index']);
    // 下线小程序
    $router->get('offline', ['uses' => 'OfflineController@offline']);
});

/**
 * 开放平台相关
 */
$router->group(['prefix' => 'open-platform', 'namespace' => 'OpenPlatform'], function () use ($router) {
    // 服务商提审额度次数查询
    $router->get('audit-quota', ['uses' => 'AuditQuotaController@index']);
    // 加急审核申请
    $router->get('speed-up-audit', ['uses' => 'SpeedUpAuditController@index']);
});

/**
 * 小程序自动发版
 */
$router->group(['prefix' => 'auto-release', 'namespace' => 'AutoRelease'], function () use ($router) {
    // 列表
    $router->post('list', 'ListController@index');
    // 发布版本 - 稳定
    $router->post('release-stable', 'ReleaseStableController@index');
    // 发布版本 - 灰度
    $router->post('release-gray-scale', 'ReleaseGrayScaleController@index');
    // 小程序列表
    $router->post('enterprise-list', 'EnterpriseListController@index');
    // 撤销灰度
    $router->post('cancel', 'CancelController@index');
    // 删除
    $router->post('delete', 'DeleteController@index');
    // 发布项目tag
    $router->post('release-project', 'ReleaseProjectController@index');
});

/**
 * 企业管理--门店商城-新版
 */
$router->group(['prefix' => 'enterprise-new', 'namespace' => 'EnterpriseNew'], function () use ($router) {
    // 企业列表
    $router->get('list', ['uses' => 'ListController@index']);
    // 企业详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 增加企业
    $router->post('add', ['uses' => 'AddController@index']);
    // 编辑企业
    $router->post('edit', ['uses' => 'EditController@index']);
    // 修改基本信息
    $router->post('change-base-info', ['uses' => 'ChangeBaseInfoController@index']);
    // 获取企业附加功能
    $router->post('additional-tree', ['uses' => 'AdditionalTreeController@additionalTree']);
});

/**
 * 企业付款到零钱-充值相关
 */
$router->group(['prefix' => 'enterprise-payment', 'namespace' => 'EnterprisePayment'], function () use ($router) {
    // 充值明细列表
    $router->post('recharge-record', 'RechargeRecordController@index');
    // 充值
    $router->post('recharge', 'RechargeController@index');

    // 企业付费
    // 付费列表
    $router->post('list', 'ListController@index');
    // 添加付费信息
    $router->post('add', 'AddController@index');
    // 付费统计
    $router->post('statistics', 'StatisticsController@index');
    // 付费详情
    $router->post('detail', 'DetailController@index');
    // 收款
    $router->post('receive-payment', 'EditController@receivePayment');
    // 续费列表
    $router->post('renew-list', 'ListController@renewList');
    // 续费
    $router->post('renew', 'EditController@renew');
    // 导出
    $router->post('export', 'ExportController@index');
});

/**
 * 行业相关
 */
$router->group(['prefix' => 'industry', 'namespace' => 'Industry'], function () use ($router) {
    // 行业列表
    $router->post('list', 'ListController@index');
});

/**
 * 角色权限
 */
$router->group([
    'prefix'    => 'rbac',
    'namespace' => 'Rbac'
], function () use ($router) {
    $router->group([
        'prefix' => 'permissions'
    ], function () use ($router) {
        $router->post('all-tree', [
            'uses'     => 'PermissionsController@allTree',
            'describe' => '权限树'
        ]);
        $router->post('add', [
            'uses'     => 'PermissionsController@add',
            'describe' => '新增权限'
        ]);
        $router->post('edit', [
            'uses'     => 'PermissionsController@edit',
            'describe' => '编辑权限'
        ]);
        $router->post('release', [
            'uses'     => 'PermissionsController@release',
            'describe' => '发布权限'
        ]);
        $router->post('sync-and-release', [
            'uses'     => 'PermissionsController@syncAndRelease',
            'describe' => '同步并发布权限'
        ]);
        $router->post('system-type-list', [
            'uses'     => 'PermissionsController@systemTypeList',
            'describe' => '系统类型列表'
        ]);
    });
});

/**
 * 公告相关
 */
$router->group(['prefix' => 'announcement', 'namespace' => 'Announcement'], function () use ($router) {
    // 公告列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 公告详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 公告新增
    $router->post('add', ['uses' => 'AddController@index']);
    // 公告编辑
    $router->post('edit', ['uses' => 'EditController@index']);
    // 公告删除
    $router->post('delete', ['uses' => 'DeleteController@index']);
});

/**
 * 支付服务商相关
 */
$router->group(['prefix' => 'merchant', 'namespace' => 'Merchant'], function () use ($router) {
    // 支付服务商列表
    $router->post('list', ['uses' => 'ListController@index']);

});

