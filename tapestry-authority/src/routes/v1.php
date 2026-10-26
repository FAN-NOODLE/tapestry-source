<?php

/**
 * @var \Laravel\Lumen\Routing\Router $router
 */

$router->get('', function () {
    return 'v1版本网关路由';
});

/**
 * 公共
 */
$router->group([
    'prefix'    => 'common',
    'namespace' => 'Common'
], function () use ($router) {
    $router->post('send-code', [
        'uses'     => 'SendCodeController@index',
        'describe' => '发送验证码'
    ]);
    $router->post('check-code', [
        'uses'     => 'CheckCodeController@index',
        'describe' => '校验验证码'
    ]);

    $router->post('send-sms-code', [
        'uses'     => 'SmsController@sendCode',
        'describe' => '发送短信验证码'
    ]);
    $router->post('check-sms-code', [
        'uses'     => 'SmsController@checkCode',
        'describe' => '校验短信验证码'
    ]);
    $router->post('validate-code', [
        'uses'     => 'ValidateCodeController@index',
        'describe' => '图片验证码',
        'page'     => '',
    ]);
    $router->post('send-message', [
        'uses'     => 'SendMessageController@index',
        'describe' => '发送短信'
    ]);
});

/**
 * 应用
 */
$router->group([
    'prefix'    => 'application',
    'namespace' => 'Application'
], function () use ($router) {
    $router->post('syncApplication', [
        'uses'     => 'SyncApplicationController@index',
        'describe' => '同步企业应用'
    ]);
});

/**
 * 角色
 */
$router->group([
    'prefix'    => 'role',
    'namespace' => 'Role'
], function () use ($router) {
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '角色列表'
    ]);

    $router->get('all-list', [
        'uses'     => 'AllListController@index',
        'describe' => '角色不分页列表'
    ]);

    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '新增角色'
    ]);

    $router->post('edit', [
        'uses'     => 'EditController@index',
        'describe' => '编辑角色'
    ]);

    $router->get('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '删除角色'
    ]);

    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '角色详情'
    ]);
});

/**
 * 账号
 */
$router->group([
    'prefix'    => 'account',
    'namespace' => 'Account'
], function () use ($router) {
    $router->post('login', [
        'uses'     => 'LoginController@index',
        'describe' => '统一登录'
    ]);
    $router->post('mini-program-login', [
        'uses'     => 'LoginController@miniProgramLogin',
        'describe' => '小程序登录'
    ]);
    $router->post('scrm-login', [
        'uses'     => 'ScrmLoginController@index',
        'describe' => 'scrm管理后台登录'
    ]);

    $router->post('scrm-password-login', [
        'uses'     => 'ScrmPasswordLoginController@index',
        'describe' => 'scrm管理后台账号密码登录'
    ]);
    $router->post('scrm-code-login', [
        'uses'     => 'ScrmCodeLoginController@index',
        'describe' => 'scrm管理后台手机验证码登录'
    ]);

    $router->post('scrm-password-ip-white-login', [
        'uses'     => 'ScrmPasswordIpWhiteLoginController@index',
        'describe' => 'scrm管理后台账号密码ip白名单登录'
    ]);

    $router->post('sso-login', [
        'uses'     => 'SSOLoginController@index',
        'describe' => 'scrm管理后台SSO登录'
    ]);

    $router->post('check-login', [
        'uses'     => 'CheckLoginController@index',
        'describe' => '统一校验账号密码一致性'
    ]);

    $router->post('get-token', [
        'uses'     => 'GetTokenController@index',
        'describe' => '获取正式token'
    ]);

    $router->post('decrypt-token', [
        'uses'     => 'DecryptTokenController@index',
        'describe' => '解密token'
    ]);

    $router->post('check-account', [
        'uses'     => 'CheckAccountController@index',
        'describe' => '校验账号是否存在'
    ]);

    $router->post('account-enterprise-list', [
        'uses'     => 'AccountEnterpriseListController@index',
        'describe' => '账号对应企业列表'
    ]);

    $router->post('account-access-control', [
        'uses'     => 'AccountAccessController@index',
        'describe' => '账号是否有访问权限'
    ]);

    $router->post('data-range', [
        'uses'     => 'DataRangeController@index',
        'describe' => '账号数据范围'
    ]);

    $router->post('data-range-list', [
        'uses'     => 'DataRangeController@getRange',
        'describe' => '账号数据范围-简单列表'
    ]);

    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '创建账号'
    ]);

    $router->post('edit', [
        'uses'     => 'EditController@index',
        'describe' => '编辑账号'
    ]);
    $router->post('edit-password', [
        'uses'     => 'EditPasswordController@index',
        'describe' => '修改密码'
    ]);
    $router->post('reset-password', [
        'uses'     => 'EditPasswordController@resetPassword',
        'describe' => '重置密码'
    ]);
    $router->post('check-reset-password', [
        'uses'     => 'EditPasswordController@checkResetPassword',
        'describe' => '重置密码验证'
    ]);
    $router->post('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '账号详情'
    ]);

    $router->post('del', [
        'uses'     => 'DeleteController@index',
        'describe' => '账号删除'
    ]);

    $router->post('status', [
        'uses'     => 'StatusController@index',
        'describe' => '账号状态修改'
    ]);

    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '账号列表'
    ]);

    $router->post('perfect', [
        'uses'     => 'PerfectController@index',
        'describe' => '完善账号信息'
    ]);
    $router->post('account-info-change', [
        'uses'     => 'InfoChangeController@infoChange',
        'describe' => '完善账号信息-手机号'
    ]);

    $router->post('check', [
        'uses'     => 'CountCheckController@index',
        'describe' => '账号规则验证'
    ]);

    $router->post('account-system', [
        'uses'     => 'AccountSystemController@index',
        'describe' => '账号创建可使用系统'
    ]);

    $router->post('store-account', [
        'uses'     => 'HasAccountController@index',
        'describe' => '判断账号是否已经存在'
    ]);
    $router->post('child-account', [
        'uses'     => 'ChildAccountController@index',
        'describe' => '只获取账号列表'
    ]);

    $router->post('admin-roles', [
        'uses'     => 'AdminRbacController@adminRoles',
        'describe' => '获取管理员的角色'
    ]);

    $router->post('init-store-member', [
        'uses'     => 'InitAccountRoleController@initStoreMember',
        'describe' => '初始化店员'
    ]);
    $router->post('init-store-manager', [
        'uses'     => 'InitAccountRoleController@initStoreManager',
        'describe' => '初始化店长'
    ]);
    $router->post('init-enterprise-guide', [
        'uses'     => 'InitAccountRoleController@initEnterpriseGuide',
        'describe' => '初始化企业导购'
    ]);
    $router->post('init-community-guide', [
        'uses'     => 'InitAccountRoleController@initCommunityGuide',
        'describe' => '初始化企业导购'
    ]);
    $router->post('init-enterprise-admin', [
        'uses'     => 'InitAccountRoleController@initEnterpriseAdmin',
        'describe' => '初始化企业管理'
    ]);

    $router->post('account-by-role', [
        'uses'     => 'AccountRoleController@accountByRole',
        'describe' => '账号角色相关'
    ]);
    $router->post('store-worker-permission', [
        'uses'     => 'AccountRoleController@storeWorkerPermission',
        'describe' => '账号权限相关'
    ]);

    $router->post('base-info-check', [
        'uses'     => 'BaseInfoCheckController@index',
        'describe' => '企业基础配置检查'
    ]);
    $router->post('edit-admin-enterprise-data-range', [
        'uses'     => 'EditAdminEnterpriseDataRangeController@index',
        'describe' => '修改管理员管理部门/门店'
    ]);

    $router->post('forbid-leave-admin', [
        'uses'     => 'ForbidLeaveAdmin@index',
        'describe' => '禁用离职管理员'
    ]);

    $router->post('range-list-by-mobiles', [
        'uses'     => 'RangeListByMobilesController@index',
        'describe' => '根据多个手机号获取管理范围列表'
    ]);

    $router->post('account-permissions', [
        'uses'     => 'AccountPermissionsController@index',
        'describe' => '获取某个导购（管理员）的权限id集合'
    ]);
});

/**
 * 微信
 */
$router->group([
    'prefix'    => 'wx',
    'namespace' => 'Wx'
], function () use ($router) {
    $router->post('access-token', [
        'uses'     => 'GetAccessTokenController@index',
        'describe' => '统一access-token管理'
    ]);

    $router->post('dynamic-message-create', [
        'uses'     => 'DynamicMsgCreateController@index',
        'describe' => '动态消息创建'
    ]);

    $router->post('dynamic-message-update', [
        'uses'     => 'DynamicMsgUpdateController@index',
        'describe' => '动态消息更新'
    ]);

    $router->post('qr-code', [
        'uses'     => 'QrCodeController@index',
        'describe' => '场景二维码'
    ]);
});

$router->group([
    'prefix'    => 'merchant',
    'namespace' => 'Merchant'
], function () use ($router) {

    $router->post('service-merchant-detail', [
        'uses'     => 'ServiceMerchantController@serviceMerchantDetail',
        'describe' => '企业绑定的服务商详情'
    ]);

    $router->post('service-merchant-list', [
        'uses'     => 'ServiceMerchantController@serviceMerchantList',
        'describe' => '企业绑定的服务商列表'
    ]);

    $router->post('service-merchant-default', [
        'uses'     => 'ServiceMerchantController@getDefaultDetail',
        'describe' => '企业绑定的服务商列表'
    ]);
});

/**
 * 角色权限
 */
$router->group([
    'prefix'    => 'rbac',
    'namespace' => 'Rbac'
], function () use ($router) {

    $router->group([
        'prefix' => 'role'
    ], function () use ($router) {
        $router->post('init', [
            'uses'     => 'RoleController@init',
            'describe' => '初始化超管'
        ]);
        $router->post('add', [
            'uses'     => 'RoleController@add',
            'describe' => '角色新增'
        ]);
        $router->post('edit', [
            'uses'     => 'RoleController@edit',
            'describe' => '角色编辑'
        ]);
        $router->post('list', [
            'uses'     => 'RoleController@lists',
            'describe' => '角色列表'
        ]);
        $router->post('detail', [
            'uses'     => 'RoleController@detail',
            'describe' => '角色详情'
        ]);
        $router->post('delete', [
            'uses'     => 'RoleController@delete',
            'describe' => '角色删除'
        ]);
    });

    $router->post('all-tree', [
        'uses'     => 'PermissionsController@allTree',
        'describe' => '权限树列表'
    ]);
    $router->post('additional-list', [
        'uses'     => 'PermissionsController@additionalList',
        'describe' => '附加功能列表'
    ]);
    $router->post('add', [
        'uses'     => 'PermissionsController@add',
        'describe' => '添加权限'
    ]);
    $router->post('edit', [
        'uses'     => 'PermissionsController@edit',
        'describe' => '编辑权限'
    ]);
    $router->post('list', [
        'uses'     => 'PermissionsController@lists',
        'describe' => '权限列表'
    ]);
    $router->post('release', [
        'uses'     => 'PermissionsController@release',
        'describe' => '权限发布'
    ]);
    $router->post('system-type-list', [
        'uses'     => 'PermissionsController@systemTypeList',
        'describe' => '系统类型列表'
    ]);
    $router->post('sync', [
        'uses'     => 'PermissionsController@sync',
        'describe' => '同步所有权限'
    ]);
    $router->post('children-id-list', [
        'uses'     => 'PermissionsController@childrenIDList',
        'describe' => '权限下级ID列表'
    ]);
});


$router->group([
    'prefix'    => 'message',
    'namespace' => 'Message'
], function () use ($router) {

    $router->post('account', [
        'uses'     => 'AccountController@send',
        'describe' => '公众号消息推送',
        'page'     => ''
    ]);

    $router->post('mini-subscribe-record', [
        'uses'     => 'MiniSubscribeRecordController@index',
        'describe' => '小程序订阅消息记录',
        'page'     => ''
    ]);

    $router->post('mini-user-subscribe-list', [
        'uses'     => 'MiniSubscribeRecordController@userSubscribeList',
        'describe' => '小程序获取用户订阅记录',
        'page'     => ''
    ]);

    $router->post('mini-subscribe-list', [
        'uses'     => 'MiniSubscribeRecordController@subscribeList',
        'describe' => '小程序获取用户订阅记录',
        'page'     => ''
    ]);

    $router->post('mini-subscribe-status', [
        'uses'     => 'MiniSubscribeRecordStatusController@index',
        'describe' => '小程序获取用户订阅记录状态',
        'page'     => ''
    ]);

    $router->post('activity-message-list', [
        'uses'     => 'ActivityMessageController@lists',
        'describe' => '活动提醒配置列表',
        'page'     => ''
    ]);

    $router->post('activity-message-add', [
        'uses'     => 'ActivityMessageController@add',
        'describe' => '活动提醒配置添加',
        'page'     => ''
    ]);

    $router->post('activity-message-delete', [
        'uses'     => 'ActivityMessageController@delete',
        'describe' => '活动提醒配置删除',
        'page'     => ''
    ]);

    $router->group([
        'middleware' => 'checkMessageStrategy',
    ], function () use ($router) {
        $router->post('mini-sync-templates', [
            'uses'     => 'MiniSyncTemplatesController@index',
            'describe' => '手动更新小程序模板消息',
            'page'     => ''
        ]);
        $router->post('send', [
            'uses'     => 'SendController@index',
            'describe' => '消息推送',
            'page'     => ''
        ]);
        $router->post('mini', [
            'uses'     => 'MiniSubscribeController@send',
            'describe' => '小程序订阅消息',
            'page'     => ''
        ]);
        $router->post('mini-touch-data', [
            'uses'     => 'MiniTouchDataController@index',
            'describe' => '小程序订阅消息对应事件',
            'page'     => ''
        ]);

        $router->post('subscribe-strategy', [
            'uses'     => 'SubscribeStrategyController@index',
            'describe' => '用户召回-配置项列表',
            'page'     => ''
        ]);

        $router->post('update-subscribe-strategy', [
            'uses'     => 'UpdateSubscribeStrategyController@index',
            'describe' => '用户召回-配置项修改',
            'page'     => ''
        ]);

        $router->post('subscribe-notice-status', [
            'uses'     => 'SubscribeNoticeStatusController@index',
            'describe' => '用户召回-首页提示状态',
            'page'     => ''
        ]);

        $router->post('message-config', [
            'uses' => 'MessageConfigController@index'
        ]);

        $router->post('mini-sync-message-num', [
            'uses'     => 'MiniSyncMessageNumberController@index',
            'describe' => '手动同步订阅消息人员数量',
            'page'     => ''
        ]);
    });
});

$router->group([
    'prefix' => 'tester',
], function () use ($router) {
    // 配置默认服务商
    $router->post('set-default-merchant', [
        'uses' => 'TesterController@setDefaultMerchant'
    ]);

    // 配置subscribe_message_setting表记录
    $router->post('set-message-setting', [
        'uses' => 'TesterController@setMessageSetting'
    ]);
});