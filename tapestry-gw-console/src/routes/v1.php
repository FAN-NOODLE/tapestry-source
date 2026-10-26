<?php

/**
 * 路由定义格式!!!
 * $router->[method]('xxx', ['uses' => 'xxxController@xxx', 'describe' => '路由描述', 'page' => '对应前端页面');
 */
/**
 * @var Router $router
 */

// {domain}/gw-console/v1
use Laravel\Lumen\Routing\Router;

$router->get('', function () {
    return 'v1版本网关路由';
});

// 公用上传接口
$router->post('upload', ['uses' => 'UploadController@index']);

/**
 * 管理后台配置
 */
$router->group([
    'prefix'    => 'setting',
    'namespace' => 'Setting'
], function () use ($router) {
    // UC在企业微信配置的回调地址
    $router->get('callback', [
        'uses'     => 'CallbackController@index',
        'describe' => '',
        'page'     => ''
    ]);

    $router->post('base-config', [
        'uses'     => 'BaseConfigController@index',
        'describe' => '管理后台配置.基本配置.检查',
        'page'     => '管理后台配置.基本配置.检查'
    ]);

    /**
     * 角色管理
     */
    $router->group([
        'prefix'    => 'role',
        'namespace' => 'Role'
    ], function () use ($router) {
        $router->post('list', [
            'uses'     => 'ListController@index',
            'describe' => '管理后台配置.角色管理.列表',
            'page'     => '设置.账号权限.角色'
        ]);
        $router->post('all-list', [
            'uses'     => 'AllListController@index',
            'describe' => '管理后台配置.角色管理.不分页列表',
            'page'     => '设置.账号权限.角色'
        ]);
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '管理后台配置.角色管理.详情',
            'page'     => '设置.账号权限.角色.编辑'
        ]);
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '管理后台配置.角色管理.新增',
            'page'     => '设置.账号权限.角色.添加'
        ]);
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '管理后台配置.角色管理.编辑',
            'page'     => '设置.账号权限.角色.编辑'
        ]);
        $router->get('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '管理后台配置.角色管理.删除',
            'page'     => '设置.账号权限.角色'
        ]);
    });

    /**
     * 账号管理--新版-auth微服务
     */
    $router->group([
        'prefix'    => 'account',
        'namespace' => 'Account'
    ], function () use ($router) {
        // 列表
        $router->post('list', [
            'uses'     => 'ListController@index',
            'describe' => '管理后台配置.管理员管理.列表',
            'page'     => '设置.账号权限.账号'
        ]);
        // 详情
        $router->post('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '管理后台配置.管理员管理.详情',
            'page'     => '设置.账号权限.账号.编辑'
        ]);
        // 新增
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '管理后台配置.管理员管理.新增',
            'page'     => '设置.账号权限.账号.新增账号'
        ]);
        // 编辑
        $router->post('edit', [
            'uses'     => 'EditController@index',
            'describe' => '管理后台配置.管理员管理.编辑',
            'page'     => '设置.账号权限.账号.编辑'
        ]);
        // 编辑
        $router->post('edit-password', [
            'uses'     => 'EditPsController@index',
            'describe' => '管理后台配置.管理员管理.编辑',
            'page'     => '设置.账号权限.账号.编辑'
        ]);
        // 编辑
        $router->post('reset-password', [
            'uses'     => 'ResetPasswordController@index',
            'describe' => '管理后台配置.管理员管理.重置密码',
            'page'     => '设置.账号权限.账号.重置密码'
        ]);
        // 状态修改
        $router->post('status', [
            'uses'     => 'StatusController@index',
            'describe' => '管理后台配置.管理员管理.状态修改',
            'page'     => '设置.账号权限.账号.编辑'
        ]);
        // 删除
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '管理后台配置.管理员管理.删除',
            'page'     => '设置.账号权限.账号'
        ]);
        // 完善信息-包含修改密码修改手机号账户名称等
        $router->post('perfect', [
            'uses'     => 'PerfectController@index',
            'describe' => '管理后台配置.完善账户信息',
            'page'     => '设置.账号权限.账号'
        ]);

        // 判断企业子账号数量
        $router->post('check-count', [
            'uses'     => 'CountCheckController@index',
            'describe' => '管理后台配置.新建账号',
            'page'     => '设置.账号权限.账号'
        ]);

        // 获取创建账号可用系统选项
        $router->post('account-system', [
            'uses'     => 'AccountSystemController@index',
            'describe' => '可用系统选项',
            'page'     => '设置.账号权限.账号'
        ]);

        // 完善信息部分详情
        $router->post('perfect-detail', [
            'uses'     => 'PerfectDetailController@index',
            'describe' => '完善账户信息详情',
            'page'     => '设置.账号权限.账号'
        ]);

    });


    /**
     * 操作日志管理
     */
    $router->group([
        'prefix'    => 'operation-log',
        'namespace' => 'OperationLog'
    ], function () use ($router) {
        $router->get('module', [
            'uses'     => 'ModuleController@index',
            'describe' => '管理后台配置.操作日志.全部模块列表',
            'page'     => '设置.账号权限.操作日志'
        ]);
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '管理后台配置.操作日志.日志列表',
            'page'     => '设置.账号权限.操作日志'
        ]);
    });

    /**
     * 商户信息
     */
    $router->group([
        'prefix'    => 'merchant_setting',
        'namespace' => 'Merchant'
    ], function () use ($router) {
        $router->post('detail', [
            'uses'     => 'SettingDetailController@index',
            'describe' => '管理后台配置.商户信息.获取商户信息',
            'page'     => '设置.基本设置.商户信息'
        ]);
        $router->post('save', [
            'uses'     => 'SettingSaveController@index',
            'describe' => '管理后台配置.商户信息.保存商户信息',
            'page'     => '设置.基本设置.商户信息'
        ]);
    });
});

/**
 * 小程序授权
 */
$router->group([
    'prefix'    => 'auth',
    'namespace' => 'Auth'
], function () use ($router) {
    $router->post('shop-auth-detail', [
        'uses'     => 'ShopAuthDetailController@index',
        'describe' => '小程序授权.商城小程序授权详情',
        'page'     => '设置.小程序授权.商城小程序授权'
    ]);
    $router->post('work-auth-detail', [
        'uses'     => 'WorkAuthDetailController@index',
        'describe' => '小程序授权.企业小程序授权详情',
        'page'     => '设置.小程序授权.企业微信授权'
    ]);
    $router->post('auth-contact', [
        'uses'     => 'ContactAuthDetailController@index',
        'describe' => '小程序授权.通讯录授权详情',
        'page'     => '设置.小程序授权.客户联系配置'
    ]);
    $router->post('work-auth-complete', [
        'uses'     => 'WorkCompleteController@index',
        'describe' => '小程序授权.企业微信小程序完善信息',
        'page'     => '设置.小程序授权.完善资料'
    ]);
    $router->post('shop-auth-complete', [
        'uses'     => 'ShopCompleteController@index',
        'describe' => '小程序授权.商城小程序完善信息',
        'page'     => '设置.小程序授权.完善资料'
    ]);
    $router->post('auth-complete', [
        'uses'     => 'AuthCompleteController@index',
        'describe' => '小程序授权.统一完善信息',
        'page'     => ''
    ]);
    $router->post('shop-save', [
        'uses'     => 'ShopSaveController@index',
        'describe' => '小程序授权.商城小程序信息编辑',
        'page'     => ''
    ]);
    $router->post('work-save', [
        'uses'     => 'WorkSaveController@index',
        'describe' => '小程序授权.企业小程序信息编辑',
        'page'     => ''
    ]);
    $router->post('mini-detail', [
        'uses'     => 'AuthDetailController@index',
        'describe' => '小程序授权.小程序信息详情',
        'page'     => ''
    ]);
    $router->post('work-auth-complete-new', [
        'uses'     => 'WorkCompleteNewController@index',
        'describe' => '小程序授权.企业微信完善信息-新',
        'page'     => '设置.小程序授权.完善资料'
    ]);
    $router->post('work-save-new', [
        'uses'     => 'WorkSaveNewController@index',
        'describe' => '小程序授权.企业微信信息编辑-新',
        'page'     => ''
    ]);
    $router->post('work-mini-list', [
        'uses'     => 'GuideMiniListController@index',
        'describe' => '企业微信小程序授权.小程序列表',
        'page'     => ''
    ]);
    $router->post('install-app', [
        'uses'     => 'InstallAppController@install',
        'describe' => '企业微信.自建应用安装',
        'page'     => ''
    ]);
    $router->post('dao-gou-config', [
        'uses'     => 'DaoGouConfigController@daoGouConfig',
        'describe' => '企业微信.导购自建应用配置',
        'page'     => ''
    ]);
    $router->post('get-dao-gou-config', [
        'uses'     => 'DaoGouConfigController@getDaoGouConfig',
        'describe' => '企业微信.获取导购自建应用配置',
        'page'     => ''
    ]);
    $router->get('get-work-plugin', [
        'uses'     => 'InstallAppController@getWorkPlugin',
        'describe' => '企业微信.获取自建应用安装信息',
        'page'     => ''
    ]);
    $router->post('work-plugin-save', [
        'uses'     => 'InstallAppController@update',
        'describe' => '企业微信.自建应用更新',
        'page'     => ''
    ]);
    $router->post('update-refresh-token', [
        'uses'     => 'RefreshTokenController@updateMiniProgramToken',
        'describe' => '企业授权信息.更新Token',
        'page'     => ''
    ]);
});

/**
 * 公共
 */
$router->group([
    'prefix'    => 'common',
    'namespace' => 'Common'
], function () use ($router) {
    $router->post('export-excel', [
        'uses'     => 'ImportFailExportExcelController@index',
        'describe' => '导入失败列表导出Excel',
        'page'     => ''
    ]);
    $router->post('send-code', [
        'uses'     => 'SendCodeController@index',
        'describe' => '发送验证码',
        'page'     => ''
    ]);
    $router->post('check-code', [
        'uses'     => 'CheckCodeController@index',
        'describe' => '验证验证码',
        'page'     => ''
    ]);
    $router->post('send-sms-code', [
        'uses'     => 'SmsController@sendCode',
        'describe' => '发送短信验证码',
        'page'     => ''
    ]);
    $router->post('mini-info', [
        'uses'     => 'MiniInfoController@index',
        'describe' => '获取小程序信息',
        'page'     => ''
    ]);
});

/**
 * 管理员相关接口
 */
$router->group([
    'prefix'    => 'admin',
    'namespace' => 'Admin'
], function () use ($router) {
    $router->post('login', [
        'uses'     => 'LoginController@index',
        'describe' => '管理员登陆',
        'page'     => ''
    ]);
    $router->post('send-sms-code', [
        'uses'     => 'SendSmsCodeController@index',
        'describe' => '发送验证码',
        'page'     => ''
    ]);
    $router->post('mobile-login', [
        'uses'     => 'MobileLoginController@index',
        'describe' => '手机验证码登陆',
        'page'     => ''
    ]);
    $router->post('logout', [
        'uses'     => 'LogoutController@index',
        'describe' => '管理员退出',
        'page'     => ''
    ]);
    $router->get('mock-login', [
        'uses'     => 'MockLoginController@index',
        'describe' => '管理员模拟登陆',
        'page'     => ''
    ]);
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '管理员详情',
        'page'     => ''
    ]);
    $router->post('password', [
        'uses'     => 'PasswordController@index',
        'describe' => '管理员修改密码',
        'page'     => ''
    ]);

    $router->post('get-token', [
        'uses'     => 'GetTokenController@index',
        'describe' => '获取登录token',
        'page'     => ''
    ]);

    $router->post('get-enterprise-config', [
        'uses'     => 'GetEnterpriseConfigController@index',
        'describe' => '获取企业配置',
        'page'     => ''
    ]);


    // 获取V1版本管理员数据
    $router->post('adminInfo', ['uses' => 'ScrmGetTokenController@index']);
    $router->post('get-tokenv1', ['uses' => 'AutoAccountController@getInfoV1']);

    $router->post('js-api-config', [
        'uses'     => 'JsApiConfigController@index',
        'describe' => '获取js api config',
        'page'     => ''
    ]);

    $router->post('change-enterprise-list', [
        'uses'     => 'ChangeEnterpriseListController@index',
        'describe' => '切换企业时的列表',
        'page'     => ''
    ]);

    $router->get('work-wechat-login-jump', [
        'uses'     => 'WorkWechatLoginController@jump',
        'describe' => '企业微信扫码登录重定向',
        'page'     => '登录页'
    ]);
    $router->post('work-wechat-login', [
        'uses'     => 'WorkWechatLoginController@login',
        'describe' => '企业微信扫码登录',
        'page'     => '登录页'
    ]);
    $router->post('test', [
        'uses'     => 'TestController@index',
        'describe' => '测试登录',
        'page'     => ''
    ]);
    $router->post('enterprise-list', [
        'uses'     => 'LoginSelectEnterpriseListController@index',
        'describe' => '登录时选择企业列表',
        'page'     => ''
    ]);
    $router->post('check-reset-password', [
        'uses'     => 'ResetPasswordController@checkResetPassword',
        'describe' => '重置密码验证',
        'page'     => ''
    ]);
    $router->post('reset-password', [
        'uses'     => 'ResetPasswordController@index',
        'describe' => '重置密码',
        'page'     => ''
    ]);
});

/**
 * 页面模板
 * 1. 4月8日迁移 店铺风格、店铺导航、提审管理、发布管理、版本管理
 */
$router->group(['prefix' => 'template', 'namespace' => 'Template'], function () use ($router) {
    // 首页
    $router->group(['prefix' => 'home', 'namespace' => 'Home'], function () use ($router) {
        // 详情
        $router->get('detail', ['uses' => 'DetailController@index', 'describe' => '首页详情', 'page' => '首页模板']);
        // 优惠券列表（立减金、微信优惠券、站内优惠券）
        $router->get('coupon-list', ['uses' => 'CouponListController@index', 'describe' => '优惠券列表', 'page' => '首页模板']);
        // 新增
        $router->post('add', ['uses' => 'AddController@index', 'describe' => '新增', 'page' => '首页模板']);
        // 编辑
        $router->post('edit', ['uses' => 'UpdateController@index', 'describe' => '编辑', 'page' => '首页模板']);
        // 商品列表
        $router->post('goods-list', ['uses' => 'GoodsController@index', 'describe' => '商品列表', 'page' => '首页模板']);
        // 模板列表
        $router->post('/', ['uses' => 'ListController@index', 'describe' => '模板列表', 'page' => '首页模板']);
        // 删除
        $router->post('delete', ['uses' => 'DeleteController@index', 'describe' => '删除', 'page' => '首页模板删除']);
    });

    // 模板风格
    $router->group(['prefix' => 'shop-style', 'namespace' => 'ShopStyle'], function () use ($router) {
        // 详情
        $router->post('info', ['uses' => 'InfoController@index', 'describe' => '风格详情', 'page' => '风格详情']);
        // 所有企业通用风格配置列表
        $router->post('config', ['uses' => 'ConfigController@index', 'describe' => '企业通用风格配置列表', 'page' => '企业通用风格配置列表']);
        // 设置
        $router->post('set', ['uses' => 'SetController@index', 'describe' => '企业风格设置', 'page' => '企业风格设置']);
        // 添加
        $router->post('add', ['uses' => 'AddController@index', 'describe' => '企业风格添加', 'page' => '企业风格添加']);
        // 删除
        $router->post('delete', ['uses' => 'DeleteController@index', 'describe' => '企业风格删除', 'page' => '企业风格删除']);
    });
    /**
     * 我的模板
     */
    $router->group(['prefix' => 'myTemplate', 'namespace' => 'MyTemplate'], function () use ($router) {
        // 列表
        $router->post('list', 'ListController@index');
        // 新增
        $router->post('add', 'AddController@index');
        // 编辑
        $router->post('update', 'UpdateController@index');
        // 删除
        $router->post('delete', 'DeleteController@index');
        // 详情
        $router->post('detail', 'DetailController@index');
    });

    /**
     * 分类模板
     */
    $router->group(['prefix' => 'category', 'namespace' => 'Category'], function () use ($router) {
        // 列表
        $router->post('', 'ListController@index');
        // 新增
        $router->post('add', 'AddController@index');
        // 编辑
        $router->post('edit', 'EditController@index');
        // 删除
        $router->post('delete', 'DeleteController@index');
        // 详情
        $router->post('detail', 'DetailController@index');
    });

    //版本管理
    $router->group(['prefix' => 'version', 'namespace' => 'Version'], function () use ($router) {
        // 版本更新列表
        $router->get('list', ['uses' => 'ListController@index', 'describe' => '版本更新列表', 'page' => '版本更新列表']);

        // 版本历史记录
        $router->get('history-list', ['uses' => 'HistoryListController@index', 'describe' => '版本历史记录', 'page' => '版本历史记录']);

        // 版本记录详情
        $router->get('history-detail', ['uses' => 'HistoryDetailController@index', 'describe' => '版本记录详情', 'page' => '版本记录详情']);

        // 确认发布版本
        $router->post('release', ['uses' => 'ReleaseController@index', 'describe' => '确认发布版本', 'page' => '确认发布版本']);

        // 查看模板列表数据
        $router->get('template-list', ['uses' => 'TemplateListController@index', 'describe' => '查看模板列表数据', 'page' => '查看模板列表数据']);

        // 选择模板编辑保存
        $router->post('save', ['uses' => 'SaveController@index', 'describe' => '选择模板编辑保存', 'page' => '选择模板编辑保存']);

        // 模拟boss发布版本
        $router->get('save-status', ['uses' => 'SaveStatusController@index', 'describe' => '模拟boss审核版本', 'page' => '模拟boss审核版本']);

        // 自动化打包 提交审核
        $router->post('submit-audit', ['uses' => 'SubmitAuditController@index', 'describe' => '提交审核', 'page' => '提交审核']);

    });

    /**
     * 自动化发布
     */
    $router->group(['prefix' => 'auto-release', 'namespace' => 'AutoRelease'], function () use ($router) {
        // 撤销审核
        $router->post('cancel-audit', ['uses' => 'CancelAuditController@index', 'describe' => '撤销审核', 'page' => '版本管理-']);
        // 定时发布-撤销发布
        $router->post('cancel-release', ['uses' => 'CancelReleaseController@index', 'describe' => '撤销定时发布', 'page' => '版本管理']);
        // 预览二维码
        $router->post('qr-code', ['uses' => 'QrCodeController@index', 'describe' => '体验二维码', 'page' => '版本管理']);
        // 发布
        $router->post('release', ['uses' => 'ReleaseController@index', 'describe' => '发布小程序', 'page' => '版本管理']);
        // 版本回退
        $router->post('revert-release', ['uses' => 'RevertReleaseController@index', 'describe' => '回退小程序版本', 'page' => '版本管理']);
        // 提审
        $router->post('submit-audit', ['uses' => 'SubmitAuditController@index', 'describe' => '提交小程序审核', 'page' => '版本管理']);
        // 版本更新列表
        $router->post('update-list', ['uses' => 'UpdateListController@index', 'describe' => '提审管理详情', 'page' => '版本管理']);
        // 版本管理列表
        $router->post('version-list', ['uses' => 'VersionListController@index', 'describe' => '小程序发布列表', 'page' => '版本管理']);
    });

    /**
     * 小程序开发码
     */
    $router->group(['prefix' => 'develop', 'namespace' => 'Develop'], function () use ($router) {
        // 开发码列表
        $router->get('list', ['uses' => 'ListController@index', 'describe' => '小程序开发码列表', 'page' => '版本管理']);
        // 开发码生成 请求提交
        $router->post('submit', ['uses' => 'SubmitController@index', 'describe' => '提交小程序开发码请求', 'page' => '版本管理']);
    });

    /**
     * 页面版本
     */
    $router->group(['prefix' => 'page-version', 'namespace' => 'PageVersion'], function () use ($router) {
        $router->post('list', ['uses' => 'ListController@index', 'describe' => '版本列表', 'page' => '店铺.店铺模板.版本列表']);
        $router->post('detail', ['uses' => 'DetailController@index', 'describe' => '版本详情', 'page' => '店铺.店铺模板.版本']);
        $router->post('edit', ['uses' => 'EditController@index', 'describe' => '版本编辑', 'page' => '店铺.店铺模板.版本列表.编辑']);
        $router->post('add', ['uses' => 'AddController@index', 'describe' => '版本新增复制', 'page' => '店铺.店铺模板.版本列表']);
        $router->post('delete', ['uses' => 'DeleteController@index', 'describe' => '版本删除', 'page' => '店铺.店铺模板.版本列表']);
    });
});

/**
 * 店铺导航
 */
$router->group(['prefix' => 'tab-bar', 'namespace' => 'TabBar'], function () use ($router) {
    // 编辑
    $router->post('update', ['uses' => 'UpdateController@index', 'describe' => '添加编辑底部按钮', 'page' => '添加编辑底部按钮']);
    // 详情
    $router->get('detail', ['uses' => 'DetailController@index', 'describe' => '底部按钮详情', 'page' => '底部按钮详情']);
});

/**
 * 配置管理
 */
$router->group(['prefix' => 'config', 'namespace' => 'Config'], function () use ($router) {
    //保存配置
    $router->post('save', 'SaveController@index');
    // 详情
    $router->post('info', 'InfoController@index');
});

/**
 * 消息
 */
$router->group(['prefix' => 'message', 'namespace' => 'Message'], function () use ($router) {
    // 列表
    $router->post('subscribe-strategy', 'SubscribeStrategyController@index');
    // 编辑
    $router->post('update-subscribe-strategy', 'SubscribeStrategyController@update');


    // 活动开始提醒列表
    $router->post('activity-message-list', 'ActivityMessageController@lists');
    // 活动开始提醒保存
    $router->post('activity-message-add', 'ActivityMessageController@add');
    // 活动开始提醒删除
    $router->post('activity-message-delete', 'ActivityMessageController@delete');
});


/**
 * 下载中心
 */
$router->group(['prefix' => 'download-center', 'namespace' => 'DownloadCenter'], function () use ($router) {
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '下载中心列表',
        'page'     => '下载中心'
    ]);
    $router->get('download', [
        'uses'     => 'DownloadController@index',
        'describe' => '下载中心下载',
        'page'     => '下载中心'
    ]);
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '下载中心删除',
        'page'     => '下载中心'
    ]);
});

/**
 * 活动
 */
$router->group(['prefix' => 'activity-center', 'namespace' => 'ActivityCenter'], function () use ($router) {

    // 拼团
    $router->group(['prefix' => 'group', 'namespace' => 'Group'], function () use ($router) {
        $router->post('detail', ['uses' => 'DetailController@index']);
        $router->post('list', ['uses' => 'ListController@index']);
        $router->post('add', ['uses' => 'AddController@index']);
        $router->post('edit', ['uses' => 'EditController@index']);
        $router->post('status', ['uses' => 'StatusController@index']);
        $router->post('delete', ['uses' => 'DeleteController@index']);
    });

    // 砍价
    $router->group(['prefix' => 'haggle', 'namespace' => 'Haggle'], function () use ($router) {
        $router->post('detail', ['uses' => 'DetailController@index']);
        $router->post('list', ['uses' => 'ListController@index']);
        $router->post('add', ['uses' => 'AddController@index']);
        $router->post('edit', ['uses' => 'EditController@index']);
        $router->post('status', ['uses' => 'StatusController@index']);
        $router->post('delete', ['uses' => 'DeleteController@index']);
        $router->post('init-list', ['uses' => 'InitListController@index']);
    });

    // 新人购
    $router->group(['prefix' => 'newcomer', 'namespace' => 'Newcomer'], function () use ($router) {
        $router->post('detail', ['uses' => 'DetailController@index']);
        $router->post('list', ['uses' => 'ListController@index']);
        $router->post('add', ['uses' => 'AddController@index']);
        $router->post('edit', ['uses' => 'EditController@index']);
        $router->post('status', ['uses' => 'StatusController@index']);
        $router->post('delete', ['uses' => 'DeleteController@index']);
    });

    // 限时购
    $router->group(['prefix' => 'flash-sale', 'namespace' => 'FlashSale'], function () use ($router) {
        $router->post('detail', ['uses' => 'DetailController@index']);
        $router->post('list', ['uses' => 'ListController@index']);
        $router->post('add', ['uses' => 'AddController@index']);
        $router->post('edit', ['uses' => 'EditController@index']);
        $router->post('status', ['uses' => 'StatusController@index']);
        $router->post('delete', ['uses' => 'DeleteController@index']);
    });

    // 好友助力
    $router->group(['prefix' => 'help', 'namespace' => 'Help'], function () use ($router) {
        $router->post('detail', ['uses' => 'DetailController@index']);
        $router->post('list', ['uses' => 'ListController@index']);
        $router->post('add', ['uses' => 'AddController@index']);
        $router->post('edit', ['uses' => 'EditController@index']);
        $router->post('status', ['uses' => 'StatusController@index']);
        $router->post('delete', ['uses' => 'DeleteController@index']);
        $router->post('init-user-list', ['uses' => 'InitUserListController@index']);
    });

    // 助力抽奖
    $router->group(['prefix' => 'luck', 'namespace' => 'Luck'], function () use ($router) {
        $router->post('detail', ['uses' => 'DetailController@index']);
        $router->post('list', ['uses' => 'ListController@index']);
        $router->post('add', ['uses' => 'AddController@index']);
        $router->post('edit', ['uses' => 'EditController@index']);
        $router->post('status', ['uses' => 'StatusController@index']);
        $router->post('delete', ['uses' => 'DeleteController@index']);
        $router->post('initiator', ['uses' => 'InitiatorController@index']);
        $router->post('helper', ['uses' => 'HelperController@index']);
        $router->post('winner', ['uses' => 'WinnerController@index']);
        $router->post('export', ['uses' => 'ExportController@index']);
        $router->post('init-rank-list', ['uses' => 'InitRankListController@index']);
        $router->post('join-user-list', ['uses' => 'JoinUserListController@index']);

    });

    // 先享卡
    $router->group(['prefix' => 'wxcard', 'namespace' => 'Wxcard'], function () use ($router) {
        $router->post('detail', ['uses' => 'DetailController@index']);
        $router->post('list', ['uses' => 'ListController@index']);
        $router->post('add', ['uses' => 'AddController@index']);
        $router->post('edit', ['uses' => 'EditController@index']);
        $router->post('status', ['uses' => 'StatusController@index']);
        $router->post('delete', ['uses' => 'DeleteController@index']);
        $router->post('card-order-list', ['uses' => 'CardOrderListController@index']);
        $router->post('card-order-list-export', ['uses' => 'CardOrderListExportController@index']);
        $router->post('open-card-list', ['uses' => 'OpenCardListController@index']);
        $router->post('export-open-card-list', ['uses' => 'ExportOpenCardListController@index']);
    });

    // 口令券
    $router->group(['prefix' => 'password-coupon', 'namespace' => 'PasswordCoupon'], function () use ($router) {
        $router->post('generate-posters', ['uses' => 'GeneratePostersController@index']);
        $router->post('add', ['uses' => 'AddController@index', 'describe' => '口令券添加', 'page' => '应用,口令优惠券']);
        $router->post('list', ['uses' => 'ListController@index', 'describe' => '口令券列表', 'page' => '应用,口令优惠券']);
        $router->post('update', ['uses' => 'UpdateController@index', 'describe' => '口令券编辑', 'page' => '应用,口令优惠券']);
        $router->post('detail', ['uses' => 'DetailController@index', 'describe' => '口令券详情', 'page' => '应用,口令优惠券']);
        $router->post('download-poster', ['uses' => 'GeneratePostersController@index', 'describe' => '口令券海报下载', 'page' => '应用,口令优惠券']);
        $router->get('download-password', ['uses' => 'DownLoadPasswordListController@index', 'describe' => '口令券下载所有口令', 'page' => '应用,口令优惠券']);
        $router->post('download-password', ['uses' => 'DownLoadPasswordListController@index', 'describe' => '口令券下载所有口令', 'page' => '应用,口令优惠券']);
        $router->post('download-sun-code', ['uses' => 'DownLoadSunCodeController@index', 'describe' => '口令券下载太阳码', 'page' => '应用,口令优惠券']);
        $router->post('download-poster-more-status', ['uses' => 'DownLoadPosterMoreStatusController@index', 'describe' => '下载一码多券海报状态变更', 'page' => '应用,口令优惠券']);
    });
});

/**
 * 物料
 */
$router->group([
    'prefix'    => 'material',
    'namespace' => 'Material'
], function () use ($router) {
    $router->post('download-add-template', [
        'uses'     => 'AddImportController@downloadAddTemplate',
        'describe' => '下载物料导入模板',
        'page'     => '库存.导入'
    ]);
    $router->post('add-import', [
        'uses'     => 'AddImportController@index',
        'describe' => '物料导入',
        'page'     => '库存.导入'
    ]);
    $router->post('download-edit-template', [
        'uses'     => 'EditStockImportController@downloadAddTemplate',
        'describe' => '下载物料库存编辑模板',
        'page'     => '库存.编辑'
    ]);
    $router->post('edit-stock-import', [
        'uses'     => 'EditStockImportController@index',
        'describe' => '物料库存导入',
        'page'     => '库存.编辑'
    ]);

    //单仓-修改物料库存
    $router->post('modify-stock', [
        'uses'     => 'UpdateSingleStockController@index',
        'describe' => '物料库存修改',
        'page'     => '库存.编辑'
    ]);

    // 获取企业物料列表
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '企业物料库存列表',
        'page'     => '库存'
    ]);

    // 新增物料库存（页面没使用）
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '新增物料库存',
        'page'     => '库存.新增'
    ]);

    // 物料详情
    $router->post('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '物料详情',
        'page'     => '库存.查看'
    ]);

    // 编辑物料
    $router->post('edit', [
        'uses'     => 'EditController@index',
        'describe' => '编辑物料',
        'page'     => '库存.编辑'
    ]);

    // 物料库存修改
    $router->post('update-stock', [
        'uses'     => 'UpdateStockController@index',
        'describe' => '物料库存修改',
        'page'     => '库存.更新库存'
    ]);

    // 总库日志
    $router->post('log-list', [
        'uses'     => 'MaterialLogController@index',
        'describe' => '总库日志',
        'page'     => '库存.详情.总库日志'
    ]);

    // 导出企业物料列表
    $router->post('export-list', [
        'uses'     => 'MaterialExportController@index',
        'describe' => '企业物料库存列表导出',
        'page'     => '库存.列表导出'
    ]);
    // excel分派物料
    $router->post('assign-material', ['uses' => 'AssignMaterialController@index']);
    // 下载excel物料分派模板
    $router->post('download-assign-material', ['uses' => 'AssignMaterialController@downloadAddTemplate']);
    // 变更渠道物料权限状态
    $router->post('update-rule-status', ['uses' => 'UpdateMaterialRuleController@index']);
    // 验证门店的物料权限
    $router->post('check-material-rule', ['uses' => 'CheckMaterialRuleController@index']);
    // 获取物料的分派列表
    $router->post('material-assign-list', ['uses' => 'MaterialAssignListController@index']);
    // 渠道的物料列表
    $router->post('channel-material-list', ['uses' => 'MaterialChannelListController@index']);
    // 商品新增物料选择列表
    $router->post('material-select-list', ['uses' => 'MaterialSelectListController@index']);

});

/**
 * 规格模板
 */
$router->group([
    'prefix'    => 'attribute',
    'namespace' => 'Attribute'
], function () use ($router) {
    $router->post('template-add', [
        'uses'     => 'TemplateAddController@index',
        'describe' => '规格模板新增',
        'page'     => '规格模板.新增'
    ]);
    $router->post('template-delete', [
        'uses'     => 'TemplateDeleteController@index',
        'describe' => '规格模板删除',
        'page'     => '规格模板.删除'
    ]);
    $router->post('template-detail', [
        'uses'     => 'TemplateDetailController@index',
        'describe' => '规格模板详情',
        'page'     => '规格模板.详情'
    ]);
    $router->post('template-list', [
        'uses'     => 'TemplateListController@index',
        'describe' => '规格模板列表',
        'page'     => '规格模板.列表'
    ]);
    $router->get('template-select-list', [
        'uses'     => 'TemplateGroupByTagListController@index',
        'describe' => '商品新增',
        'page'     => '商品管理.新增'
    ]);
    $router->post('template-update', [
        'uses'     => 'TemplateUpdateController@index',
        'describe' => '规格模板编辑',
        'page'     => '规格模板.编辑'
    ]);
    $router->post('template-move', [
        'uses'     => 'TemplateMoveController@index',
        'describe' => '规格模板移动',
        'page'     => '规格模板.移动'
    ]);
    $router->post('template-tag-add', [
        'uses'     => 'TemplateTagAddController@index',
        'describe' => '规格模板标签新增',
        'page'     => '规格模板.标签新增'
    ]);
    $router->post('template-tag-list', [
        'uses'     => 'TemplateTagListController@index',
        'describe' => '规格模板标签列表',
        'page'     => '规格模板.标签列表'
    ]);
    $router->post('template-goods-list', [
        'uses'     => 'TemplateGoodsListController@index',
        'describe' => '规格模板关联商品列表',
        'page'     => '规格模板.关联商品列表'
    ]);
    $router->post('item-use-status', [
        'uses'     => 'ItemUseStatusController@index',
        'describe' => '规格是否被使用',
        'page'     => '规格模板.规格是否被使用'
    ]);
});

/**
 * 商品
 */
$router->group([
    'prefix'    => 'goods',
    'namespace' => 'Goods'
], function () use ($router) {
    $router->post('save', [
        'uses'     => 'SaveController@index',
        'describe' => '商品新增编辑发布',
        'page'     => '商品操作'
    ]);
    $router->post('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '商品详情',
        'page'     => '商品操作'
    ]);
    $router->post('update-spu', [
        'uses'     => 'UpdateSpuController@index',
        'describe' => '修改商品spu',
        'page'     => '商品操作'
    ]);
    $router->post('delete-sku', [
        'uses'     => 'DelSkuController@index',
        'describe' => '修改商品spu',
        'page'     => '商品操作'
    ]);
    $router->post('delete-spu', [
        'uses'     => 'DelSpuController@index',
        'describe' => '修改商品spu',
        'page'     => '商品操作'
    ]);
    $router->post('copy-goods', [
        'uses'     => 'CopyGoodsController@index',
        'describe' => '复制商品至门店',
        'page'     => '商品操作'
    ]);
    $router->post('check-import-data', [
        'uses'     => 'CheckImportDataController@index',
        'describe' => '核对批量修改商品的数据',
        'page'     => '商品操作'
    ]);
    $router->post('export-template', [
        'uses'     => 'ExportTemplateController@index',
        'describe' => '下载批量修改商品的模板',
        'page'     => '商品操作'
    ]);
    $router->post('export-errors', [
        'uses'     => 'ExportErrorsController@index',
        'describe' => '下载批量修改商品的错误列表',
        'page'     => '商品操作'
    ]);
    $router->post('confirm-batch-update', [
        'uses'     => 'ConfirmBatchUpdateController@index',
        'describe' => '确认批量修改商品',
        'page'     => '商品操作'
    ]);
    $router->post('imported-data-list', [
        'uses'     => 'ImportedDataListController@index',
        'describe' => '确认批量修改商品',
        'page'     => '商品操作'
    ]);
    $router->post('sale', [
        'uses'     => 'SaleController@index',
        'describe' => '上下架',
        'page'     => '商品操作'
    ]);

    $router->post('sale-all', [
        'uses'     => 'SaleAllController@index',
        'describe' => '全部上下架',
        'page'     => '商品操作'
    ]);
//    $router->post('list', [
//        'uses'     => 'ListController@index',
//        'describe' => '商品列表',
//        'page'     => '商品操作'
//    ]);
    $router->post('goods-list', [
        'uses'     => 'ConsoleGoodsListController@index',
        'describe' => '商品列表',
        'page'     => '商品操作'
    ]);
    $router->post('spu-option-list', [
        'uses'     => 'SpuOptionListController@index',
        'describe' => '商品SPU选择框列表',
        'page'     => '商品操作'
    ]);
    $router->post('spu-list', [
        'uses'     => 'SpuListController@index',
        'describe' => '商品SPU列表',
        'page'     => '商品操作'
    ]);
    $router->post('spu-list-new', [
        'uses'     => 'SpuListController@lists',
        'describe' => '商品SPU列表新',
        'page'     => '商品操作'
    ]);
    $router->post('sku-list', [
        'uses'     => 'CateSkuListController@index',
        'describe' => '带分类的商品列表',
        'page'     => '商品操作'
    ]);
    $router->post('spu-draft-list', [
        'uses'     => 'SpuDraftListController@index',
        'describe' => '商品spu草稿列表',
        'page'     => '商品操作'
    ]);
    $router->post('export-goods', [
        'uses'     => 'ExportGoodsController@index',
        'describe' => '商品导出',
        'page'     => '商品操作'
    ]);
    $router->post('shop-goods-list', [
        'uses'     => 'ShopGoodsListController@index',
        'describe' => '获取商城的商品列表',
        'page'     => '商品操作'
    ]);
    $router->post('edit-price', [
        'uses'     => 'EditPriceController@index',
        'describe' => '修改商品价格',
        'page'     => '商品操作'
    ]);
    $router->post('sku-pop-window', [
        'uses'     => 'SkuPopWindowController@skuPopWindow',
        'describe' => 'SKU弹窗',
        'page'     => '店铺'
    ]);

    // 待上架商品列表
    $router->post('not-sale-list', [
        'uses'     => 'NotSaleListController@index',
        'describe' => '待上架商品信息',
        'page'     => '店铺'
    ]);

    // 批量改价失败列表
    $router->post('multi-change-price-fail-download', [
        'uses'     => 'MultiChangePriceController@failList',
        'describe' => '商品批量改价失败列表下载',
        'page'     => '商品操作'
    ]);

    // 批量改价上传接口
    $router->post('multi-change-price-upload', [
        'uses'     => 'MultiChangePriceController@index',
        'describe' => '商品批量改价数据上传',
        'page'     => '商品操作'
    ]);

    // 是否可以进行批量改价操作
    $router->post('multi-change-price-check', [
        'uses'     => 'MultiChangePriceController@operCheck',
        'describe' => '商品批量改价数据上传',
        'page'     => '商品操作'
    ]);

    // 商品批量改价模板下载
    $router->get('multi-change-price-template-download', [
        'uses'     => 'MultiChangePriceController@getTemplate',
        'describe' => '商品批量改价数据上传',
        'page'     => '商品操作'
    ]);

    // 批量改价进度
    $router->get('multi-change-price-upload', [
        'uses'     => 'MultiChangePriceController@process',
        'describe' => '商品批量改价进度接口',
        'page'     => '商品操作'
    ]);

    // 批量改价重置redis
    $router->get('multi-change-price-reset', [
        'uses'     => 'MultiChangePriceController@reset',
        'describe' => '商品批量改价进度接口',
        'page'     => '商品操作'
    ]);

    // 批量改价redis查看
    $router->get('multi-change-price-showredis', [
        'uses'     => 'MultiChangePriceController@showredis',
        'describe' => '商品批量改价进度接口',
        'page'     => '商品操作'
    ]);
    $router->post('spu-sales', [
        'uses'     => 'SpuSalesController@index',
        'describe' => 'Spu销量',
        'page'     => '商品详情页设置'
    ]);
    // 批量更改商品信息
    $router->post('goods-batch-change', [
        'uses'     => 'GoodsBatchChangeController@index',
        'describe' => '批量更改商品信息',
        'page'     => '商品批量操作'
    ]);
    // 批量更改商品信息查询
    $router->post('goods-batch-change-select', [
        'uses'     => 'GoodsBatchChangeSelectController@index',
        'describe' => '批量更改商品信息查询',
        'page'     => '商品批量操作'
    ]);
    // 批量更改商品信息进度
    $router->post('goods-batch-change-progress', [
        'uses'     => 'GoodsBatchChangeProgressController@index',
        'describe' => '批量更改商品信息进度',
        'page'     => '商品批量操作'
    ]);
    // 批量更改商品信息结果
    $router->post('goods-batch-change-result', [
        'uses'     => 'GoodsBatchChangeResultController@index',
        'describe' => '批量更改商品信息结果',
        'page'     => '商品批量操作'
    ]);
    // 批量更改商品信息模板下载
    $router->post('goods-batch-change-template-download', [
        'uses'     => 'GoodsBatchChangeTemplateController@index',
        'describe' => '批量更改商品信息模板下载',
        'page'     => '商品批量操作'
    ]);
    // 批量更改商品原有信息下载
    $router->post('goods-batch-change-info-download', [
        'uses'     => 'GoodsBatchChangeInfoController@index',
        'describe' => '批量更改商品信息下载',
        'page'     => '商品批量操作'
    ]);
    // 获取要复制的已失效的商品详情
    $router->post('get-delete-goods-detail', [
        'uses'     => 'GetDeleteGoodsDetailController@index',
        'describe' => '商品详情',
        'page'     => '复制商品'
    ]);
    // 商品预警
    $router->group([
        'prefix'    => 'stock-warning',
        'namespace' => 'StockWarning'
    ], function () use ($router) {
        // 获取要复制的已失效的商品详情
        $router->post('console-close-warning', [
            'uses'     => 'ConsoleCloseWarningController@index',
            'describe' => '商品预警',
            'page'     => '关闭后台首页商品预警信息'
        ]);
        // 获取要复制的已失效的商品详情
        $router->post('console-index', [
            'uses'     => 'ConsoleIndexController@index',
            'describe' => '商品预警',
            'page'     => '后台首页商品预警信息'
        ]);
        // 获取要复制的已失效的商品详情
        $router->post('export', [
            'uses'     => 'ExportController@index',
            'describe' => '商品预警',
            'page'     => '导出商品预警'
        ]);
        // 获取要复制的已失效的商品详情
        $router->post('list', [
            'uses'     => 'ListController@index',
            'describe' => '商品预警',
            'page'     => '商品预警列表'
        ]);
        // 获取要复制的已失效的商品详情
        $router->post('update-stock-warning-value', [
            'uses'     => 'UpdateStockWarningValueController@index',
            'describe' => '商品预警',
            'page'     => '修改商品预警值'
        ]);
    });
});


/**
 * 门店员工管理
 */
$router->group([
    'prefix'    => 'store/worker',
    'namespace' => 'Store\Worker'
], function () use ($router) {
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '员工新增',
        'page'     => '员工管理'
    ]);
    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '员工列表',
        'page'     => '员工管理'
    ]);
    $router->post('edit', [
        'uses'     => 'EditController@index',
        'describe' => '修改员工',
        'page'     => '员工管理'
    ]);
    $router->post('edit-status', [
        'uses'     => 'JobEditController@index',
        'describe' => '在职状态',
        'page'     => '员工管理'
    ]);
    // 隐藏接口（同步组织架构）
    $router->post('sync-org', [
        'uses'     => 'SyncOrgController@index',
        'describe' => '同步组织架构',
        'page'     => '员工管理'
    ]);
    $router->post('add-worker-poster', [
        'uses'     => 'AddWorkerPosterController@index',
        'describe' => '申请添加员工海报',
        'page'     => '员工管理'
    ]);
    $router->post('member-select', [
        'uses'     => 'MemberSelectController@index',
        'describe' => '员工精确搜索',
        'page'     => '员工管理'
    ]);
});

/**
 * 门店区域管理
 */
$router->group([
    'prefix'    => 'store-area',
    'namespace' => 'StoreArea',
], function () use ($router) {
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '区域列表',
        'page'     => '区域管理'
    ]);
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '区域删除',
        'page'     => '区域管理',
    ]);
    $router->post('save', [
        'uses'     => 'SaveController@index',
        'describe' => '区域保存',
        'page'     => '区域管理',
    ]);
    $router->post('batch-save', [
        'uses'     => 'BatchSaveController@index',
        'describe' => '区域批量保存',
        'page'     => '区域管理',
    ]);
    $router->post('sort', [
        'uses'     => 'SortController@index',
        'describe' => '区域排序',
        'page'     => '区域管理',
    ]);
    $router->post('move', [
        'uses'     => 'MoveController@index',
        'describe' => '区域迁移',
        'page'     => '区域管理',
    ]);
    $router->post('get-area', [
        'uses'     => 'GetAreaController@index',
        'describe' => '通过PID获取区域',
        'page'     => '区域管理',
    ]);
    $router->post('get-store', [
        'uses'     => 'GetStoreController@index',
        'describe' => '获取区域下门店',
        'page'     => '区域管理',
    ]);
    $router->post('store-migration', [
        'uses'     => 'StoreMigrationController@index',
        'describe' => '门店迁移',
        'page'     => '区域管理',
    ]);
    $router->post('areaStore-list', [
        'uses'     => 'AreaStoreListController@index',
        'describe' => '区域门店列表',
        'page'     => '区域管理',
    ]);
});

/**
 * 门店分类
 */
$router->group([
    'prefix'    => 'store/category',
    'namespace' => 'Store\Category'
], function () use ($router) {
    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '门店分类列表',
        'page'     => '门店分类列表'
    ]);
});

/**
 * 门店装修
 */
$router->group([
    'prefix'    => 'store/decorate',
    'namespace' => 'Store\Decorate'
], function () use ($router) {
    $router->post('index', [
        'uses'     => 'IndexController@index',
        'describe' => '门店装修',
        'page'     => '门店首页装修'
    ]);
});

/**
 * 分类
 */
$router->group([
    'prefix'    => 'category',
    'namespace' => 'Category'
], function () use ($router) {
    // 这接口文件都丢了
    $router->post('list', [
        'uses'     => 'CategoryListController@categoryList',
        'describe' => '分类列表',
        'page'     => '分类管理'
    ]);
    // 模板用分类树
    $router->post('category-tree', [
        'uses'     => 'CategoryTreeController@categoryTree',
        'describe' => '分类列表',
        'page'     => '分类管理'
    ]);
});

/**
 * 渠道
 */
$router->group(['prefix' => 'material-channel', 'namespace' => 'MaterialChannel'], function () use ($router) {
    // 列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 修改渠道库存
    $router->post('update-channel-stock', ['uses' => 'UpdateChannelStockController@index']);
    // 日志列表
    $router->post('log-list', ['uses' => 'LogListController@index']);
    // 渠道类型列表
    $router->post('type-list', ['uses' => 'TypeListController@index']);
    // 查物料渠道库存列表
    $router->post('source-stock-list', ['uses' => 'SourceStockListController@index']);
    // 详情
    $router->post('detail', ['uses' => 'DetailController@index']);
});

/**
 * 邮费
 */
$router->group([
    'prefix'    => 'postage',
    'namespace' => 'Postage'
], function () use ($router) {
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '邮费模板新增',
        'page'     => '邮费'
    ]);
    $router->post('edit', [
        'uses'     => 'EditController@index',
        'describe' => '邮费模板编辑',
        'page'     => '邮费'
    ]);
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '邮费模板列表',
        'page'     => '邮费'
    ]);
    $router->post('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '邮费模板详情',
        'page'     => '邮费'
    ]);
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '邮费模板删除',
        'page'     => '邮费'
    ]);
    $router->post('relation-activity-list', [
        'uses'     => 'RelationActivityListController@index',
        'describe' => '邮费模板关联活动',
        'page'     => '邮费'
    ]);
    $router->post('relation-sku-list', [
        'uses'     => 'RelationSkuListController@index',
        'describe' => '邮费模板关联商品',
        'page'     => '邮费'
    ]);
    $router->post('relation-store-list', [
        'uses'     => 'RelationStoreListController@index',
        'describe' => '邮费模板关联门店',
        'page'     => '邮费'
    ]);
    $router->post('create-relation', [
        'uses'     => 'PostageTypeRelationController@index',
        'describe' => '创建模板关联',
        'page'     => '邮费'
    ]);
});

/**
 * 热词搜索
 */
$router->group(['prefix' => 'hot-words', 'namespace' => 'HotWords'], function () use ($router) {
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '热词添加',
        'page'     => '店铺.热词搜搜.新增'
    ]);
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '热词列表',
        'page'     => '店铺.热词搜搜'
    ]);
    $router->post('edit', [
        'uses'     => 'EditController@index',
        'describe' => '热词编辑',
        'page'     => '店铺.热词搜搜.编辑'
    ]);
    $router->post('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '热词详情',
        'page'     => '店铺.热词搜搜'
    ]);
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '热词详情',
        'page'     => '店铺.热词搜搜'
    ]);
});

/**
 * 模板管理
 */
$router->group(['prefix' => 'new-template', 'namespace' => 'NewTemplate'], function () use ($router) {
    // 列表
    $router->get('list', ['uses' => 'ListController@index', 'describe' => '模板列表', 'page' => '模板列表']);
    // 新增
    $router->post('add', ['uses' => 'AddController@index', 'describe' => '添加模板', 'page' => '添加模板']);
    // 编辑
    $router->post('update', ['uses' => 'UpdateController@index', 'describe' => '修改模板', 'page' => '修改模板']);
    // 详情
    $router->get('detail', ['uses' => 'DetailController@index', 'describe' => '模板详情', 'page' => '模板详情']);
    // 删除
    $router->post('delete', ['uses' => 'DeleteController@index', 'describe' => '删除模板', 'page' => '删除模板']);
    //获取骨架模板
    $router->post('tpl-skeleton', ['uses' => 'TplSkeletonController@index', 'describe' => '获取门店首页模板骨架', 'page' => '获取门店首页模板骨架']);
    //获取线上使用模板版本
    $router->post('online-template-version', ['uses' => 'OnlineTplVersionController@index', 'describe' => '获取线上使用模板版本', 'page' => '获取线上使用模板版本']);
    // 我的详情
    $router->post('my-detail', ['uses' => 'MyTemplateController@myDetail', 'describe' => '我的详情', 'page' => '我的详情']);
    // 我的保存
    $router->post('my-save', ['uses' => 'MyTemplateController@mySave', 'describe' => '我的保存', 'page' => '我的保存']);
    // 分类详情
    $router->post('category-detail', ['uses' => 'CategoryTemplateController@categoryDetail', 'describe' => '我的保存', 'page' => '我的保存']);
    // 分类编辑保存
    $router->post('category-save', ['uses' => 'CategoryTemplateController@categorySave', 'describe' => '我的保存', 'page' => '我的保存']);
    // 授权页详情
    $router->post('authorization-detail', ['uses' => 'AuthorizationTemplateController@authorizationDetail', 'describe' => '授权页详情', 'page' => '授权页详情']);
    // 授权页保存
    $router->post('authorization-save', ['uses' => 'AuthorizationTemplateController@authorizationSave', 'describe' => '授权页保存', 'page' => '授权页保存']);
});

/**
 * 分类模板管理
 */
$router->group(['prefix' => 'category-template', 'namespace' => 'CategoryTemplate'], function () use ($router) {

    //分类模板发布
    $router->post('release', ['uses' => 'CategoryReleaseController@index', 'describe' => '发布分类模板', 'page' => '发布分类模板']);

    //根据分类模板id获取树形分类列表
    $router->post('recursion-list', ['uses' => 'CategoryRecursionListController@index', 'describe' => '树形分类列表', 'page' => '树形分类列表']);

    //分类上移下移
    $router->post('sort-category', ['uses' => 'SortCategoryController@index', 'describe' => '分类移动', 'page' => '分类移动']);

    //分类添加商品
    $router->post('add-goods', ['uses' => 'AddGoodsController@index', 'describe' => '分类添加商品', 'page' => '分类添加商品']);

    //导入分类和分类下关联的商品
    $router->post('import-category', [
        'uses'     => 'ImportCategoryController@index',
        'describe' => '导入分类',
        'page'     => '导入分类'
    ]);
    //下载导入分类文件
    $router->get('download-import-file', [
        'uses'     => 'ImportCategoryController@downloadTemplate',
        'describe' => '下载导入分类文件',
        'page'     => '下载导入分类文件'
    ]);

    //获取某门店分类树形结构
    $router->post('store-category-list', [
        'uses'     => 'StoreCategoryListController@index',
        'describe' => '下载导入分类文件',
        'page'     => '下载导入分类文件'
    ]);

    // 添加模版分类
    $router->post('add-category-template', [
        'uses'     => 'AddCategoryTemplateController@index',
        'describe' => '添加模版分类',
        'page'     => '分类模版管理'
    ]);
    // 模版分类列表
    $router->post('category-template-list', [
        'uses'     => 'ListCategoryTemplateController@index',
        'describe' => '模版分类列表',
        'page'     => '分类模版管理'
    ]);
    // 模版分类详情
    $router->post('category-template-detail', [
        'uses'     => 'DetailCategoryTemplateController@index',
        'describe' => '模版分类详情',
        'page'     => '分类模版管理'
    ]);
    // 删除模版分类
    $router->post('category-template-delete', [
        'uses'     => 'DeleteCategoryTemplateController@index',
        'describe' => '删除模版分类',
        'page'     => '分类模版管理'
    ]);
    // 添加分类
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '添加分类',
        'page'     => '分类模版管理'
    ]);
    // 修改分类
    $router->post('edit', [
        'uses'     => 'EditController@index',
        'describe' => '修改分类',
        'page'     => '分类模版管理'
    ]);
    // 分类详情
    $router->post('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '分类详情',
        'page'     => '分类模版管理'
    ]);
    // 删除分类
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '删除分类',
        'page'     => '分类模版管理'
    ]);
    // 分类关联商品列表
    $router->post('category-goods-list', [
        'uses'     => 'CategoryGoodsListController@index',
        'describe' => '分类关联商品列表',
        'page'     => '分类模版管理'
    ]);
    // 删除分类关联商品
    $router->post('delete-category-goods', [
        'uses'     => 'DeleteCategoryGoodsController@index',
        'describe' => '删除分类关联商品',
        'page'     => '分类模版管理'
    ]);
    // 设置分类关联商品权重
    $router->post('update-category-goods-sort', [
        'uses'     => 'UpdateCategoryGoodsSortController@index',
        'describe' => '设置分类关联商品权重',
        'page'     => '分类模版管理'
    ]);
    // 获取商城已发布的分类模板id
    $router->post('get-release-template-id', ['uses' => 'GetReleaseCategoryTemplateIdController@index', 'describe' => '获取商城已发布的分类模板id', 'page' => '评论选择分类']);

});

/**
 * 页面管理
 */
$router->group([
    'prefix'    => 'page',
    'namespace' => 'Page'
], function () use ($router) {
    //    /**
    //     * 轮播图
    //     */
    //    $router->group([
    //        'prefix'    => 'banner',
    //        'namespace' => 'Banner'
    //    ], function () use ($router) {
    //        $router->get('list', [
    //            'uses'     => 'ListController@index',
    //            'describe' => '页面管理.轮播图.列表',
    //            'page'     => '店铺.页面.轮播图'
    //        ]);
    //        $router->post('save', [
    //            'uses'     => 'SaveController@index',
    //            'describe' => '页面管理.轮播图.保存',
    //            'page'     => '店铺.页面.轮播图保存'
    //        ]);
    //    });
    //
    //    /**
    //     * 引导区
    //     */
    //    $router->group([
    //        'prefix'    => 'guide',
    //        'namespace' => 'Guide'
    //    ], function () use ($router) {
    //        $router->get('list', [
    //            'uses'     => 'ListController@index',
    //            'describe' => '页面管理.引导区.列表',
    //            'page'     => '店铺.页面.首页引导区列表'
    //        ]);
    //        $router->post('save', [
    //            'uses'     => 'SaveController@index',
    //            'describe' => '页面管理.引导区.保存',
    //            'page'     => '店铺.页面.首页引导区保存'
    //        ]);
    //    });
    //
    //    /**
    //     * 爆品推荐
    //     */
    //    $router->group([
    //        'prefix'    => 'hot',
    //        'namespace' => 'Hot'
    //    ], function () use ($router) {
    //        $router->get('list', [
    //            'uses'     => 'ListController@index',
    //            'describe' => '页面管理.爆品推荐.列表',
    //            'page'     => '店铺.页面.爆品推荐列表'
    //        ]);
    //        $router->post('save', [
    //            'uses'     => 'SaveController@index',
    //            'describe' => '页面管理.爆品推荐.保存',
    //            'page'     => '店铺.页面.爆品推荐保存'
    //        ]);
    //        $router->post('goods', [
    //            'uses'     => 'GoodsController@index',
    //            'describe' => '页面管理.爆品推荐.爆品推荐商品列表',
    //            'page'     => '店铺.页面.爆品推荐选择商品'
    //        ]);
    //    });

    /**
     * 猜你喜欢/店长推荐
     */
    $router->group([
        'prefix'    => 'recommend',
        'namespace' => 'Recommend'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '页面管理.店长推荐.列表',
            'page'     => '店铺.页面.店长推荐列表'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '页面管理.店长推荐.保存',
            'page'     => '店铺.页面.店长推荐保存'
        ]);
        $router->post('goods', [
            'uses'     => 'GoodsController@index',
            'describe' => '页面管理.商品推荐.商品列表',
            'page'     => '店铺.页面.商品推荐选择商品'
        ]);
    });

    //    /**
    //     * 猜你喜欢
    //     */
    //    $router->group([
    //        'prefix'    => 'guess',
    //        'namespace' => 'Guess'
    //    ], function () use ($router) {
    //        $router->get('list', [
    //            'uses'     => 'ListController@index',
    //            'describe' => '页面管理.猜你喜欢.列表',
    //            'page'     => '店铺.页面.猜你喜欢'
    //        ]);
    //        $router->post('save', [
    //            'uses'     => 'SaveController@index',
    //            'describe' => '页面管理.猜你喜欢.保存',
    //            'page'     => '店铺.页面.猜你喜欢保存'
    //        ]);
    //    });

    /**
     * 邮费管理
     */
    $router->group([
        'prefix'    => 'postage',
        'namespace' => 'Postage'
    ], function () use ($router) {
        $router->post('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '页面管理.邮费管理.列表',
            'page'     => '店铺.页面.全场包邮详情'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '页面管理.邮费管理.保存',
            'page'     => '店铺.页面.全场包邮保存'
        ]);
    });
    //
    //    /**
    //     * 立减金管理
    //     */
    //    $router->group([
    //        'prefix'    => 'reduce-money',
    //        'namespace' => 'ReduceMoney'
    //    ], function () use ($router) {
    //        $router->get('list', [
    //            'uses'     => 'ListController@index',
    //            'describe' => '页面管理.立减金管理.列表',
    //            'page'     => '店铺.页面.立减金管理列表'
    //        ]);
    //        $router->post('add', [
    //            'uses'     => 'AddController@index',
    //            'describe' => '页面管理.立减金管理.添加',
    //            'page'     => '店铺.页面.立减金管理添加'
    //        ]);
    //        $router->post('save', [
    //            'uses'     => 'SaveController@index',
    //            'describe' => '页面管理.立减金管理.保存',
    //            'page'     => '店铺.页面.立减金管理保存'
    //        ]);
    //        $router->get('delete', [
    //            'uses'     => 'DeleteController@index',
    //            'describe' => '页面管理.立减金管理.删除',
    //            'page'     => '店铺.页面.立减金管理删除'
    //        ]);
    //        $router->get('detail', [
    //            'uses'     => 'DetailController@index',
    //            'describe' => '页面管理.立减金管理.详情',
    //            'page'     => '店铺.页面.立减金管理详情'
    //        ]);
    //    });

    /**
     * 页底
     */
    $router->group([
        'prefix'    => 'footer',
        'namespace' => 'Footer'
    ], function () use ($router) {
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '页面管理.页底管理.添加',
            'page'     => '店铺.页面.页底管理添加'
        ]);
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '页面管理.页底管理.详情',
            'page'     => '店铺.页面.页底管理详情'
        ]);
    });

    /**
     * 小程序收藏
     */
    $router->group([
        'prefix'    => 'pop-content',
        'namespace' => 'PopContent'
    ], function () use ($router) {
        $router->post('save', [
            'uses'     => 'PopContentController@contentSave',
            'describe' => '页面管理.页底管理.添加',
            'page'     => '店铺.页面.小程序收藏添加'
        ]);
        $router->post('detail', [
            'uses'     => 'PopContentController@contentDetail',
            'describe' => '页面管理.页底管理.详情',
            'page'     => '店铺.页面.小程序收藏详情'
        ]);

    });
});

/**
 * 微信省市区地址
 */
$router->group([
    'prefix'    => 'wx-area',
    'namespace' => 'WxArea'
], function () use ($router) {
    // 获取省市区地址信息
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '获取省市区地址信息',
        'page'     => '邮费.区域选择'
    ]);
});

/**
 * 门店绩效
 */
$router->group([
    'prefix'    => 'store-performance',
    'namespace' => 'StorePerformance',
], function () use ($router) {
    $router->post('fetch', [
        'uses'     => 'PerformanceFetchController@index',
        'describe' => '绩效列表',
        'page'     => '绩效管理'
    ]);

    $router->post('export', [
        'uses'     => 'PerformanceExportController@index',
        'describe' => '绩效导出',
        'page'     => '绩效管理'
    ]);

    $router->post('detail', [
        'uses'     => 'PerformanceDetailController@index',
        'describe' => '绩效详情',
        'page'     => '绩效管理'
    ]);
});

// ------- 以下开始 企业微信 标签栏 -------
/**
 * 组织管理
 */
$router->group([
    'prefix'    => 'department',
    'namespace' => 'Department'
], function () use ($router) {
    $router->get('list', [
        'uses'     => 'ListController@thirdList',
        'describe' => '组织管理.列表',
        'page'     => '客户.组织'
    ]);
    $router->get('page-list', [
        'uses'     => 'PageListController@index',
        'describe' => '组织管理.分页列表',
        'page'     => '客户.组织'
    ]);
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '组织管理.新增',
        'page'     => '客户.组织.新增部门'
    ]);
    $router->post('update', [
        'uses'     => 'UpdateController@index',
        'describe' => '组织管理.编辑',
        'page'     => '客户.组织.编辑部门'
    ]);
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '组织管理.删除',
        'page'     => '客户.组织'
    ]);
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '组织管理.详情',
        'page'     => '客户.组织'
    ]);
    $router->post('sort', [
        'uses'     => 'SortController@index',
        'describe' => '客户.组织',
        'page'     => '组织管理.排序'
    ]);
    $router->post('tree-list', [
        'uses'     => 'TreeListController@index',
        'describe' => '组织管理.树形列表',
        'page'     => '客户.组织'
    ]);
    $router->post('member-lists-by-dp-ids', [
        'uses'     => 'MemberListByScrmController@index',
        'describe' => '组织管理.部门员工列表',
        'page'     => '组织管理'
    ]);
});

/**
 * 员工管理
 */
$router->group([
    'prefix'    => 'user',
    'namespace' => 'User'
], function () use ($router) {
    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '员工管理.列表',
        'page'     => '组织.在职员工列表.在职员工列表'
    ]);
    $router->post('guide-v-user-list', [
        'uses'     => 'GuideVUserListController@index',
        'describe' => '员工管理.列表',
        'page'     => '组织.员工列表'
    ]);
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '员工管理.新增',
        'page'     => '组织.在职员工列表.添加员工'
    ]);
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '员工管理.详情',
        'page'     => '组织.在职员工列表.编辑员工'
    ]);
    $router->post('update', [
        'uses'     => 'UpdateController@index',
        'describe' => '员工管理.编辑',
        'page'     => '组织.在职员工列表.编辑员工'
    ]);
    $router->post('leave', [
        'uses'     => 'LeaveController@index',
        'describe' => '员工管理.离职',
        'page'     => '组织.在职员工列表.在职员工列表'
    ]);
    $router->get('handover-list', [
        'uses'     => 'HandoverListController@index',
        'describe' => '员工管理.交接客户列表',
        'page'     => '组织.离职员工列表.已交接客户列表'
    ]);
    $router->post('handover', [
        'uses'     => 'HandoverController@index',
        'describe' => '员工管理.交接客户',
        'page'     => '组织.离职员工列表.已交接客户列表'
    ]);
    $router->post('department-user-all-list', [
        'uses'     => 'DepartmentUserAllListController@index',
        'describe' => '员工管理.部门分层列表',
        'page'     => '组织'
    ]);
    $router->post('change_enterprise_allocation_status', [
        'uses'     => 'ChangeEnterpriseAllocationStatusController@index',
        'describe' => '员工管理.更改企业分配客户状态',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('change_allocation_status', [
        'uses'     => 'ChangeAllocationStatusController@index',
        'describe' => '员工管理.更改员工分配客户状态',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('guide_list', [
        'uses'     => 'GuideController@index',
        'describe' => '员工管理.导购分配员工列表',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('other_guides', [
        'uses'     => 'OtherGuidesController@index',
        'describe' => '员工管理.转交导购列表',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('guide_users', [
        'uses'     => 'GuideUsersController@index',
        'describe' => '员工管理.导购客户列表',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('customer_transfer', [
        'uses'     => 'CustomerTransferController@index',
        'describe' => '员工管理.转交客户',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('guides', [
        'uses'     => 'GuideController@list',
        'describe' => '导购列表(演示用)',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('get_dept_list', [
        'uses'     => 'GuideController@getDeptList',
        'describe' => '门店对应部门列表(演示用)',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('tree-list', [
        'uses'     => 'TreeListController@index',
        'describe' => '员工管理.组织树列表',
        'page'     => ''
    ]);
});

/**
 * 导购管理
 */
$router->group([
    'prefix'    => 'guide-manage',
    'namespace' => 'GuideManage'
], function () use ($router) {

    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '导购管理.新增',
        'page'     => '导购管理.添加导购'
    ]);
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '导购管理.删除',
        'page'     => '导购管理.删除导购'
    ]);
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
    // 联系我二维码列表
    $router->get('list', 'ListController@index');
    // 联系我二维码列表(下拉框选择)
    $router->get('selection_list', 'SelectionListController@index');
    // 员工树形列表
    $router->post('member_tree', 'ContactController@memberTree');
    // 创建二维码
    $router->post('create_qr_code', 'ContactController@createQrCode');
    // 编辑二维码
    $router->post('update_qr_code', 'ContactController@updateQrCode');
    // 查看二维码
    $router->post('select_qr_code', 'ContactController@selectQrCode');
    // 删除二维码
    $router->post('delete_qr_code', 'ContactController@deleteQrCode');
    // 同步组织结构
    $router->post('sync_org', 'SyncController@syncOrg');
    // 上次同步时间
    $router->get('sync_last_time', 'SyncController@syncLastTime');
    // 同步状态
    $router->get('sync_status', 'SyncController@status');
    // 同步状态
    $router->post('statistics-detail', 'StatisticsDetailController@index');
});

/**
 * 企业微信
 */
$router->group(['prefix' => 'qywx', 'namespace' => 'Qywx'], function () use ($router) {
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
        // 渠道列表
        $router->post('list', 'ListController@index');
        // 新增用户分组
        $router->post('add', 'AddController@index');
        // 用户分组详情
        $router->post('detail', 'DetailController@index');
        // 编辑用户分组
        $router->post('edit', 'EditController@index');
        // 删除用户分组
        $router->post('del', 'DeleteController@index');
    });

    /**
     * 标签管理
     */
    $router->group(['prefix' => 'tag', 'namespace' => 'Tag'], function () use ($router) {
        // 新增标签标签组
        $router->post('add', 'AddController@index');
        // 编辑标签标签组
        $router->post('edit', 'EditController@index');
        // 删除标签标签组
        $router->post('del', 'DeleteController@index');
        // 标签组列表
        $router->post('list', 'ListController@index');
        // 标签列表-根据标签组id获取
        $router->post('tag-list', 'ListController@tagList');
        // 企业标签库列表
        $router->post('all-list', 'ListController@allList');
        // 标签列表-分页-所有
        $router->post('select-list', 'ListController@selectList');
        // 关联标签和用户分组
        $router->post('tag-customer', 'TagCustomerGroupController@index');
        // 关联标签和用户分组
        $router->post('tag-sync', 'SyncTagController@index');
    });
});

/**
 * 客户管理
 */
$router->group([
    'prefix'    => 'member',
    'namespace' => 'Member'
], function () use ($router) {
    $router->post('export-customer-list', [
        'uses'     => 'ExportMemberListController@index',
        'describe' => '客户列表导出',
        'page'     => '企业微信'
    ]);

    $router->post('list', [
        'uses'     => 'MemberListController@memberList',
        'describe' => '客户列表',
        'page'     => '企业微信'
    ]);
    $router->post('user-del', [
        'uses'     => 'StaffWithUserDelController@staffWithUserDel',
        'describe' => '删除客户',
        'page'     => '企业微信'
    ]);
    $router->post('user-detail', [
        'uses'     => 'MemberDetailController@memberDetail',
        'describe' => '客户详情',
        'page'     => '企业微信'
    ]);
    $router->post('staff-list', [
        'uses'     => 'StaffListController@staffList',
        'describe' => '员工列表',
        'page'     => '企业微信'
    ]);
    $router->post('member-calc', [
        'uses'     => 'MemberCalcController@memberCalc',
        'describe' => '客户统计',
        'page'     => '企业微信'
    ]);
    $router->post('member-out-data', [
        'uses'     => 'OutDataController@outData',
        'describe' => '客户统计导出',
        'page'     => '企业微信'
    ]);
    $router->post('member-give-tag', [
        'uses'     => 'MemberGiveTagController@memberGiveTag',
        'describe' => '单独批量打标签',
        'page'     => '企业微信'
    ]);
});

/**
 * 客户管理
 */
$router->group([
    'prefix'    => 'customer',
    'namespace' => 'Customer'
], function () use ($router) {
    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '客户管理.列表(暂且废弃)',
        'page'     => ''
    ]);
    $router->post('member-list', [
        'uses'     => 'MemberListController@index',
        'describe' => '客户管理.会员客户列表',
        'page'     => '客户.会员客户列表'
    ]);
    $router->post('order-member', [
        'uses'     => 'OrderMemberListController@index',
        'describe' => '客户管理.下单客户列表',
        'page'     => '客户.下单客户列表'
    ]);
    $router->post('all-list', [
        'uses'     => 'AllCustomerListController@index',
        'describe' => '客户管理.全部客户列表',
        'page'     => '客户.全部客户列表'
    ]);
    $router->post('friend-list', [
        'uses'     => 'MemberCustomerListController@index',
        'describe' => '客户管理.导购好友客户列表',
        'page'     => '客户.导购好友客户列表'
    ]);
    $router->get('all-list-export', [
        'uses'     => 'AllCustomerListExportController@index',
        'describe' => '客户管理.导出全部客户列表',
        'page'     => '客户.全部客户列表'
    ]);
    $router->get('friend-list-export', [
        'uses'     => 'MemberCustomerListExportController@index',
        'describe' => '客户管理.导出导购好友客户列表',
        'page'     => '客户.导购好友客户列表'
    ]);
    $router->get('order-member-export', [
        'uses'     => 'OrderMemberExportController@index',
        'describe' => '客户管理.导出下单客户列表',
        'page'     => '客户.下单客户列表'
    ]);
    $router->get('member-list-export', [
        'uses'     => 'MemberListExportController@index',
        'describe' => '客户管理.导出会员客户列表',
        'page'     => '客户.会员客户列表'
    ]);
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '客户管理.基础信息',
        'page'     => '客户.客户详情'
    ]);
    $router->get('statistics', [
        'uses'     => 'StatisticsController@index',
        'describe' => '客户管理.统计数据',
        'page'     => '客户.客户详情'
    ]);
    $router->get('order-list', [
        'uses'     => 'OrderListController@index',
        'describe' => '客户管理.订单列表',
        'page'     => '客户.订单列表'
    ]);
    $router->get('order-detail', [
        'uses'     => 'OrderDetailController@index',
        'describe' => '客户管理.订单详情',
        'page'     => '客户.订单详情'
    ]);
    $router->get('addition_list', [
        'uses'     => 'AdditionListController@index',
        'describe' => '客户管理.待添加客户列表',
        'page'     => '客户.待添加客户列表'
    ]);
    $router->get('mark_user', [
        'uses'     => 'MarkUserController@index',
        'describe' => '客户管理.复制手机号标记',
        'page'     => '客户.复制手机号标记'
    ]);

    $router->post('statistic-guide-customer', [
        'uses'     => 'StatisticsController@guideCustomer',
        'describe' => '数据罗盘.顾客数据概览',
        'page'     => '数据罗盘.顾客数据概览'
    ]);
    $router->post('statistic-guide-customer-date', [
        'uses'     => 'StatisticsController@guideCustomerDate',
        'describe' => '数据罗盘.顾客数据按时间列表',
        'page'     => '数据罗盘.顾客数据按时间列表'
    ]);
    $router->post('statistic-guide-customer-list', [
        'uses'     => 'StatisticsController@guideCustomerList',
        'describe' => '数据罗盘.顾客数据按导购列表',
        'page'     => '数据罗盘.顾客数据按导购列表'
    ]);
    $router->post('statistic-guide-customer-export', [
        'uses'     => 'StatisticsController@guideCustomerListExport',
        'describe' => '数据罗盘.顾客数据导出',
        'page'     => '数据罗盘.顾客数据导出'
    ]);
    $router->post('rank', [
        'uses'     => 'RankController@index',
        'describe' => '首页.导购排行',
        'page'     => '首页.导购排行'
    ]);
    $router->post('rank-export', [
        'uses'     => 'RankController@export',
        'describe' => '首页.导购排行',
        'page'     => '首页.导购排行'
    ]);
    $router->post('rank-data', [
        'uses'     => 'RankDataController@index',
        'describe' => '首页.导购排行',
        'page'     => '首页.导购排行'
    ]);
    $router->post('rank-data-export', [
        'uses'     => 'RankDataController@export',
        'describe' => '首页.导购排行-导出',
        'page'     => '首页.导购排行-导出'
    ]);
});

/**
 * 企业微信-客户维护 Customer maintenance
 */
$router->group([
    'prefix'    => 'customer-maintenance',
    'namespace' => 'CustomerMaintenance'
], function () use ($router) {
    /**
     * 企业群发
     */
    $router->group(['prefix' => 'msg', 'namespace' => 'Msg'], function () use ($router) {
        //--添加消息群发--
        $router->post('add', ['uses' => 'AddController@index']);
        //--群发列表----
        $router->post('list', ['uses' => 'ListController@index']);
        //--推送对象全部下拉获取
        $router->get('getObject', ['uses' => 'ListController@getObjectSelect']);
        //--推送对象下拉获取  type:1所属员工|2渠道|3标签
        $router->get('getObject/{type}', ['uses' => 'ListController@getObjectSelect']);
        //--查看h5图文页面 - 素材和消息通用
        $router->get('article/{id}', ['uses' => 'ArticleController@index']);
        //--获取群发结果详情列表
        $router->post('getInfo/{tsid}', ['uses' => 'ListController@getInfo']);
        //--获取群发消息详情用于编辑
        $router->get('getEditInfo/{tsid}', ['uses' => 'ListController@getEditInfo']);
        //--获取预计触达数量
        $router->post('getEst', ['uses' => 'ListController@getEst']);
        //--更新企业群发--
        $router->post('update/{tsid}', ['uses' => 'UpdateController@index']);
        //--更新群发状态和实际触发人数--
        $router->post('updateStatus/{tsid}', ['uses' => 'UpdateStatusController@index']);
        //--删除推送消息-仅草稿状态下
        $router->get('delete/{tsid}', ['uses' => 'DeleteController@index']);
    });
    /**
     * 素材管理
     */
    $router->group(['prefix' => 'media', 'namespace' => 'Media'], function () use ($router) {
        //--用于生成media_id的临时素材
        $router->post('mediaIMG', ['uses' => 'UploadController@upMediaIMG']);
        //--用于富文本编辑的永久素材
        $router->post('upIMG', ['uses' => 'UploadController@upIMG']);
        //--获取商城小程序卡片的默认appid
        $router->get('getShopAppid', ['uses' => 'ListController@getShopAppid']);
        //--添加素材
        $router->post('add', ['uses' => 'AddController@index']);
        //--素材列表 - 所有素材
        $router->post('list', ['uses' => 'ListController@index']);
        //--素材列表 - 指定素材 type:0文本素材 1图片素材|2图文链接|3小程序卡片
        $router->post('list/{type}', ['uses' => 'ListController@index']);
        //--更新素材--
        $router->post('update/{scid}', ['uses' => 'UpdateController@index']);
        //--删除素材
        $router->get('delete/{scid}', ['uses' => 'DeleteController@index']);
    });
    /**
     * 群管理
     */
    $router->group(['prefix' => 'cluster', 'namespace' => 'Cluster'], function () use ($router) {

        $router->post('list', ['uses' => 'ListController@index']);

    });
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
    // 切换应用开关
    $router->post('switch', 'RedPacketController@switch');
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
 * 员工内购
 */
$router->group(['prefix' => 'staff-buy', 'namespace' => 'StaffBuyManager'], function () use ($router) {

    // 内购员工列表获取
    $router->post('list', ['uses' => 'ListController@index', 'describe' => '员工亲友内购列表获取', 'page' => '组织']);

    // 内购员工设置获取
    $router->post('setting', ['uses' => 'SettingController@index', 'describe' => '员工亲友内购设置获取', 'page' => '组织']);

    // 内购员工添加
    $router->post('add', ['uses' => 'AddController@index', 'describe' => '员工亲友内购添加', 'page' => '组织']);

    // 内购员工修改
    $router->post('save', ['uses' => 'UpdateController@index', 'describe' => '员工亲友内购修改', 'page' => '组织']);

    // 运营设置数据添加
    $router->post('fill', ['uses' => 'FillController@index', 'describe' => '员工亲友内购运营设置数据添加', 'page' => '组织']);

    // 运营设置数据详情
    $router->post('ext-detail', ['uses' => 'StaffExtDetailController@index', 'describe' => '员工亲友内购运营设置数据详情', 'page' => '组织']);
});

/**
 * 内购订单
 */
$router->group([
    'prefix'    => 'staff-order',
    'namespace' => 'StaffOrder'
], function () use ($router) {
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '内购订单.列表',
        'page'     => '内购订单.列表'
    ]);
    $router->post('statistics', [
        'uses'     => 'StatisticsController@index',
        'describe' => '内购订单.统计数据',
        'page'     => '内购订单.统计数据'
    ]);
    $router->post('export-order-list', [
        'uses'     => 'ExportOrderListController@index',
        'describe' => '创建内购订单.列表导出任务',
        'page'     => '创建内购订单.列表导出任务'
    ]);

    $router->post('achievement-list', [
        'uses'     => 'AchievementListController@index',
        'describe' => '员工业绩.列表',
        'page'     => '员工业绩.列表'
    ]);

    $router->post('export-achievement-list', [
        'uses'     => 'ExportAchievementListController@index',
        'describe' => '创建内购业绩列表.导出任务',
        'page'     => '创建内购业绩列表.导出任务'
    ]);

    $router->post('department-achievement-list', [
        'uses'     => 'AchievementDpListController@index',
        'describe' => '部门业绩.列表',
        'page'     => '部门业绩.列表'
    ]);

    $router->post('export-department-achievement-list', [
        'uses'     => 'ExportAchievementDpListController@index',
        'describe' => '创建内购部门业绩列表.导出任务',
        'page'     => '创建内购部门业绩列表.导出任务'
    ]);
});

/**
 * 内购订单
 */
$router->group([
    'prefix'    => 'material',
    'namespace' => 'Material'
], function () use ($router) {
    // 新增物料接口（暂时不用）
    $router->post('add', 'AddController@index');

    // 获取企业物料列表
    $router->post('list', 'ListController@index');

    // 编辑物料
    $router->post('edit', 'EditController@index');

    // 物料详情
    $router->post('detail', 'DetailController@index');

    // 物料库存修改
    $router->post('update-stock', 'UpdateStockController@index');

    // 配送方式修改
    $router->post('update-sale-type', 'ChangeSaleTypeController@index');
});

/**
 * 标签管理
 */
$router->group(['prefix' => 'label-goods', 'namespace' => 'Label\Goods'], function () use ($router) {

    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '添加标签',
        'page'     => '标签管理'
    ]);
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '删除标签',
        'page'     => '标签管理'
    ]);
    $router->post('update', [
        'uses'     => 'UpdateController@index',
        'describe' => '修改标签',
        'page'     => '标签管理'
    ]);
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '标签列表',
        'page'     => '标签管理'
    ]);
});

/**
 * 应用中心
 */

$router->group([
    'prefix'    => 'app-center',
    'namespace' => 'AppCenter'
], function () use ($router) {
    /**
     * 微信好物圈
     */
    $router->group(['prefix' => 'wx-mall', 'namespace' => 'WxMall'], function () use ($router) {
        // 物品类目导入
        $router->post('category-import', 'CategoryImportController@index');
        // 物品类目列表
        $router->post('category-list', 'CategoryListController@index');
        // 商品列表
        $router->post('goods-list', 'GoodsListController@index');
        // 导出全部商品
        $router->get('goods-export', 'GoodsExportController@index');
        // 商品批量设置好物圈分类
        $router->post('set-category', 'SetCategoryController@index');
    });

    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '应用中心',
        'page'     => '应用.全部应用'
    ]);

    /**
     * 新人优惠 - 梁欢
     */
    $router->group([
        'prefix'    => 'novice',
        'namespace' => 'Novice'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.新人优惠.列表',
            'page'     => '应用.我的应用.新人优惠'
        ]);
        $router->get('all', [
            'uses'     => 'AllController@index',
            'describe' => '应用中心.新人优惠.所有商品SPUID',
            'page'     => '应用.我的应用.新人优惠'
        ]);
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.新人优惠.添加',
            'page'     => '应用.我的应用.新人优惠.新增优惠商品'
        ]);
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.新人优惠.详情',
            'page'     => '应用.我的应用.新人优惠.编辑优惠商品'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '应用中心.新人优惠.编辑',
            'page'     => '应用.我的应用.新人优惠.编辑优惠商品'
        ]);
        $router->post('status', [
            'uses'     => 'StatusController@index',
            'describe' => '应用中心.新人优惠.状态修改',
            'page'     => '应用.我的应用.新人优惠'
        ]);
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '应用中心.新人优惠.删除',
            'page'     => '应用.我的应用.新人优惠'
        ]);
        $router->get('goods', [
            'uses'     => 'GoodsController@index',
            'describe' => '应用中心.新人优惠.商品列表',
            'page'     => '应用.我的应用.新人优惠'
        ]);
    });

    /**
     * 拼团 - 党志强
     */
    $router->group([
        'prefix'    => 'group',
        'namespace' => 'Group'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.拼团.列表',
            'page'     => '应用.我的应用.拼团'
        ]);
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.拼团.新增',
            'page'     => '应用.我的应用.拼团.新增拼团商品'
        ]);
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '应用中心.拼团.编辑',
            'page'     => '应用.我的应用.拼团.查看拼团商品'
        ]);
        $router->get('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '应用中心.拼团.删除',
            'page'     => '应用.我的应用.拼团'
        ]);
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.拼团.详情',
            'page'     => '应用.我的应用.拼团.查看拼团商品'
        ]);
        $router->get('group-goods', [
            'uses'     => 'GroupGoodsController@index',
            'describe' => '应用中心.拼团.团购商品',
            'page'     => '应用.我的应用.拼团.查看拼团商品'
        ]);
        $router->get('group-setting', [
            'uses'     => 'GroupSettingController@index',
            'describe' => '应用中心.拼团.商品参加活动修改',
            'page'     => '应用.我的应用.拼团.查看拼团商品'
        ]);
        $router->get('goods', [
            'uses'     => 'GoodsController@index',
            'describe' => '应用中心.拼团.团购商品选择列表',
            'page'     => '应用.我的应用.拼团.查看拼团商品'
        ]);
    });

    /**
     * 收藏有礼
     */
    $router->group([
        'prefix'    => 'collection',
        'namespace' => 'Collection'
    ], function () use ($router) {
        $router->post('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.收藏有礼.用户提交列表',
            'page'     => '应用.我的应用.收藏有礼'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '应用中心.收藏有礼.设置',
            'page'     => '应用.我的应用.收藏有礼.设置'
        ]);
        $router->post('audit', [
            'uses'     => 'AuditController@index',
            'describe' => '应用中心.收藏有礼.审核',
            'page'     => '应用.我的应用.收藏有礼.审核'
        ]);
        $router->post('collection_list', [
            'uses'     => 'CollectionListController@index',
            'describe' => '应用中心.收藏有礼.列表',
            'page'     => '应用.我的应用.收藏有礼.列表'
        ]);
    });

    /**
     * 秒杀
     */
    $router->group([
        'prefix'    => 'seckill',
        'namespace' => 'Seckill'
    ], function () use ($router) {
        $router->post('', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.秒杀.列表',
            'page'     => '应用.我的应用.限时购'
        ]);
        // 新增
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '应用中心.秒杀.新增',
            'page'     => '应用.我的应用.限时购.新增活动场次'
        ]);
        // 编辑
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '应用中心.秒杀.编辑',
            'page'     => '应用.我的应用.限时购.编辑活动场次'
        ]);
        // 详情
        $router->post('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.秒杀.详情',
            'page'     => '应用.我的应用.限时购.编辑活动场次'
        ]);
        // 商品列表详情
        $router->post('goods-list-detail', [
            'uses'     => 'GoodsListDetailController@index',
            'describe' => '应用中心.秒杀.商品列表详情',
            'page'     => '应用.我的应用.限时购.编辑活动场次'
        ]);
        // 商品详情
        $router->post('goods-detail', [
            'uses'     => 'GoodsDetailController@index',
            'describe' => '应用中心.秒杀.商品详情',
            'page'     => '应用.我的应用.限时购.编辑活动场次'
        ]);
        // 删除
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '应用中心.秒杀.删除',
            'page'     => '应用.我的应用.限时购.限时购'
        ]);
        // 启用/禁用
        $router->post('enable', [
            'uses'     => 'EnableController@index',
            'describe' => '应用中心.秒杀.启用/禁用',
            'page'     => '应用.我的应用.限时购.限时购'
        ]);
        // 秒杀商品选择列表
        $router->post('goods', [
            'uses'     => 'GoodsController@index',
            'describe' => '应用中心.秒杀.秒杀商品选择列表',
            'page'     => '应用.我的应用.限时购.编辑活动场次'
        ]);
    });

    /**
     * 超低购
     */
    $router->group([
        'prefix'    => 'cheap',
        'namespace' => 'Cheap',
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.超低购.列表',
            'page'     => '应用.我的应用.超低购'
        ]);
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.超低购.新增',
            'page'     => '应用.我的应用.超低购.新增'
        ]);
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.超低购.详情',
            'page'     => '应用.我的应用.超低购.查看'
        ]);
        $router->post('enable', [
            'uses'     => 'EnableController@index',
            'describe' => '应用中心.超低购.状态修改',
            'page'     => '应用.我的应用.超低购.列表'
        ]);
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '应用中心.超低购.删除',
            'page'     => '应用.我的应用.超低购.列表'
        ]);

        $router->post('sort', [
            'uses'     => 'SortController@index',
            'describe' => '应用中心.超低购.排序',
            'page'     => '应用.我的应用.超低购.排序'
        ]);
    });

    /**
     * 营销活动
     */
    $router->group([
        'prefix'    => 'activity',
        'namespace' => 'Activity'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.营销活动.列表',
            'page'     => '应用.我的应用.营销活动'
        ]);
        // 新增
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.营销活动.新增',
            'page'     => '应用.我的应用.营销活动.新增活动'
        ]);
        // 编辑
        $router->post('edit', [
            'uses'     => 'EditController@index',
            'describe' => '应用中心.营销活动.编辑',
            'page'     => '应用.我的应用.营销活动.编辑活动'
        ]);
        // 详情
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.营销活动.详情',
            'page'     => '应用.我的应用.营销活动.详情'
        ]);
        // 删除
        $router->post('del', [
            'uses'     => 'DelController@index',
            'describe' => '应用中心.营销活动.删除',
            'page'     => '应用.我的应用.营销活动.活动删除'
        ]);
        // 启用/禁用
        $router->post('changeStatus', [
            'uses'     => 'ChangeStatusController@index',
            'describe' => '应用中心.营销活动.启用/禁用',
            'page'     => '应用.我的应用.营销活动.启用/禁用'
        ]);

    });

    /**
     * 微信优惠卷
     */
    $router->group([
        'prefix'    => 'wx-coupon',
        'namespace' => 'WxCoupon'
    ], function () use ($router) {
        $router->post('condition-list', [
            'uses'     => 'ConditionListController@index',
            'describe' => '应用中心.微信优惠卷.优惠条件列表',
            'page'     => '应用.微信优惠卷.微信优惠卷'
        ]);
        $router->post('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.微信优惠卷.列表',
            'page'     => '应用.微信优惠卷.微信优惠卷'
        ]);
        // 新增
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.微信优惠卷.新增',
            'page'     => '应用.微信优惠卷.微信优惠卷.新增优惠卷'
        ]);
        // 编辑
        $router->post('edit', [
            'uses'     => 'EditController@index',
            'describe' => '应用中心.微信优惠卷.编辑',
            'page'     => '应用.我的应用.微信优惠卷.编辑优惠卷'
        ]);
        // 详情
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.微信优惠卷.详情',
            'page'     => '应用.我的应用.微信优惠卷.详情'
        ]);

        // 启用/禁用
        $router->post('changeStatus', [
            'uses'     => 'ChangeStatusController@index',
            'describe' => '应用中心.微信优惠卷.启用/禁用',
            'page'     => '应用.我的应用.微信优惠卷.启用/禁用'
        ]);
        // 启用/禁用
        $router->post('del', [
            'uses'     => 'DelController@index',
            'describe' => '应用中心.微信优惠卷.删除',
            'page'     => '应用.我的应用.微信优惠卷.删除'
        ]);

    });

    /**
     * 关联代金券和商家券
     */
    $router->group([
        'prefix'    => 'coupon-relation',
        'namespace' => 'CouponRelation'
    ], function () use ($router) {
        // 获取可关联的代金券列表
        $router->post('coupon-list', [
            'uses'     => 'CouponListController@index',
            'describe' => '应用中心.商家券.代金券/商家券关联.获取可关联的代金券列表',
            'page'     => '应用.商家券.代金券/商家券关联.新增代金券'
        ]);
        // 获取可关联的商家券列表
        $router->post('merchant-coupon-list', [
            'uses'     => 'MerchantCouponListController@index',
            'describe' => '应用中心.商家券.代金券/商家券关联.关联.获取可关联的商家券列表',
            'page'     => '应用.商家券.代金券/商家券关联.关联.添加商家券'
        ]);
        // 关联列表
        $router->post('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.商家券.代金券/商家券关联.列表',
            'page'     => '应用.商家券.代金券/商家券关联'
        ]);
        // 关联代金券和商家券
        $router->post('relation', [
            'uses'     => 'RelationController@index',
            'describe' => '应用中心.商家券.代金券/商家券关联.关联',
            'page'     => '应用.商家券.代金券/商家券关联.关联'
        ]);
        // 启用或禁用关联优惠券
        $router->post('change-status', [
            'uses'     => 'ChangeStatusController@index',
            'describe' => '应用中心.商家券.代金券/商家券关联.启用/禁用',
            'page'     => '应用.商家券.代金券/商家券关联.启用/禁用'
        ]);
    });

    /**
     * 除去活动已选的商品列表
     */
    $router->get('coupon-spu-list', [
        'uses'     => 'CouponSpuListController@index',
        'describe' => '未参加活动商品列表',
        'page'     => '应用.我的应用.未参加活动商品'
    ]);

    // 抽奖活动
    $router->group(['prefix' => 'lottery', 'namespace' => 'Lottery'], function () use ($router) {
        $router->post('add', ['uses' => 'AddController@index', 'describe' => '抽奖添加', 'page' => '应用,分享抽奖']);
        $router->get('list', ['uses' => 'ListController@index', 'describe' => '抽奖列表', 'page' => '应用,分享抽奖']);
        $router->post('update', ['uses' => 'UpdateController@index', 'describe' => '抽奖编辑', 'page' => '应用,分享抽奖']);
        $router->get('detail', ['uses' => 'DetailController@index', 'describe' => '抽奖详情', 'page' => '应用,分享抽奖']);
        $router->post('delete', ['uses' => 'DeleteController@index', 'describe' => '抽奖删除', 'page' => '应用,分享抽奖']);
        $router->post('status', ['uses' => 'StatusController@index', 'describe' => '抽奖状态', 'page' => '应用,分享抽奖']);
        $router->get('data', ['uses' => 'DataController@index', 'describe' => '抽奖活动数据', 'page' => '应用,分享抽奖']);
    });


    /**
     * 助力
     */
    $router->group([
        'prefix'    => 'help',
        'namespace' => 'Help'
    ], function () use ($router) {

        // 编辑
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '应用中心.助力.编辑',
            'page'     => '应用.助力.编辑'
        ]);

        // 删除
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '应用中心.助力.删除',
            'page'     => '应用.助力.删除'
        ]);

        // 上下架
        $router->post('enable', [
            'uses'     => 'EnableController@index',
            'describe' => '应用中心.助力.上下架',
            'page'     => '应用.助力.上下架'
        ]);
        // 详情
        $router->post('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.助力.详情',
            'page'     => '应用.助力.详情'
        ]);
        // 添加
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.助力.添加',
            'page'     => '应用.助力.添加'
        ]);
        // 发起助力用户
        $router->post('user-list', [
            'uses'     => 'UserListController@index',
            'describe' => '应用中心.助力.发起助力用户列表',
            'page'     => '应用.助力.发起助力用户列表'
        ]);
        // 发起助力用户
        $router->post('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.助力.发起助力用户列表',
            'page'     => '应用.助力.发起助力用户列表'
        ]);
    });

    /**
     * 砍价
     */
    $router->group([
        'prefix'    => 'bargain',
        'namespace' => 'Bargain'
    ], function () use ($router) {

        // 列表
        $router->post('', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.砍价.砍价列表',
            'page'     => '应用.砍价.砍价列表'
        ]);

        // 新增
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.砍价.添加',
            'page'     => '应用.砍价.新增'
        ]);

        // 编辑
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '应用中心.砍价.编辑',
            'page'     => '应用.砍价.编辑'
        ]);

        // 详情
        $router->post('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.砍价.详情',
            'page'     => '应用.砍价.详情'
        ]);

        // 删除
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '应用中心.砍价.删除',
            'page'     => '应用.砍价.删除'
        ]);

        // 上下架
        $router->post('enable', [
            'uses'     => 'EnableController@index',
            'describe' => '应用中心.砍价.上下架',
            'page'     => '应用.砍价.上下架'
        ]);

        // 发起助力用户
        $router->post('user-list', [
            'uses'     => 'UserListController@index',
            'describe' => '应用中心.砍价.发起砍价用户列表',
            'page'     => '应用.砍价.发起砍价用户列表'
        ]);

    });

    /**
     * 双倍券
     */
    $router->group([
        'prefix'    => 'double-coupon',
        'namespace' => 'DoubleCoupon'
    ], function () use ($router) {

        // 列表
        $router->post('putaway', [
            'uses'     => 'PutawayController@index',
            'describe' => '上下架',
            'page'     => '应用.双倍券.上下架'
        ]);

        // 新增
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '添加活动',
            'page'     => '应用.双倍券.添加活动'
        ]);
        // 编辑
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '更新活动',
            'page'     => '应用.双倍券.更新活动'
        ]);
        // 列表
        $router->post('list', [
            'uses'     => 'ListController@index',
            'describe' => '活动列表',
            'page'     => '应用.双倍券.活动列表'
        ]);
        // 详情
        $router->post('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '活动详情',
            'page'     => '应用.双倍券.活动详情'
        ]);
        // 删除
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '活动删除',
            'page'     => '应用.双倍券.活动删除'
        ]);
    });

    /**
     * 社群导流
     */
    $router->group([
        'prefix'    => 'community-diversion',
        'namespace' => 'CommunityDiversion'
    ], function () use ($router) {
        $router->get('get-config', [
            'uses'     => 'GetConfigController@index',
            'describe' => '获取社群导流配置',
            'page'     => '应用.社群导流'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '保存社群导流',
            'page'     => '应用.社群导流.编辑'
        ]);
    });

    /**
     * 商家券
     */
    $router->group([
        'prefix'    => 'merchant-coupon',
        'namespace' => 'MerchantCoupon'
    ], function () use ($router) {
        $router->post('check', [
            'uses'     => 'CheckConfigController@index',
            'describe' => '检查商家券是否开通',
            'page'     => '应用'
        ]);
        $router->post('open', [
            'uses'     => 'OpenController@index',
            'describe' => '开通商家券',
            'page'     => '应用'
        ]);
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '新增商家券',
            'page'     => '应用'
        ]);
        $router->post('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '商家券详情',
            'page'     => '应用'
        ]);
        $router->post('list', [
            'uses'     => 'ListController@index',
            'describe' => '商家券列表',
            'page'     => '应用'
        ]);
        $router->post('bindable-list', [
            'uses'     => 'BindableListController@index',
            'describe' => '可绑定的商家券列表',
            'page'     => '应用'
        ]);
    });

});

/**
 * 供货关系
 */
$router->group(['prefix' => 'supplier', 'namespace' => 'Supplier'], function () use ($router) {
    // 供货关系列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 供货门店的物料列表
    $router->post('material-list', ['uses' => 'MaterialListController@index']);
    // 导出供货门店的物料列表
    $router->post('download-material-list', ['uses' => 'MaterialListController@downloadMaterial']);
    // 导入并保存供货关系
    $router->post('import-save', ['uses' => 'ImportSaveController@index']);
    // 下载导入模板
    $router->post('download-supplier', ['uses' => 'DownloadSupplierController@index']);
    // 供货对象列表（门店列表）
    $router->post('store-list', ['uses' => 'StoreListController@index']);
});

/**
 * 仓库
 */
$router->group(['prefix' => 'warehouse', 'namespace' => 'Warehouse'], function () use ($router) {
    // 列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 新增编辑
    $router->post('save', ['uses' => 'SaveController@index']);
    // 开启关闭
    $router->post('enable', ['uses' => 'EnableController@index']);
    // 仓库物料列表
    $router->get('material-list', ['uses' => 'MaterialListController@index']);
    // 仓库物料日志
    $router->get('material-log', ['uses' => 'MaterialLogController@index']);
    // 仓库导入物料模板
    $router->post('import-material', ['uses' => 'ImportMaterialController@index']);
    // 下载物料导入模板
    $router->post('download-material', ['uses' => 'ImportMaterialController@downloadAddTemplate']);
    //下载添加物料信息模板
    $router->post('download-insert-material-temp', ['uses' => 'ImportMaterialController@downloadInsertTemplate']);
    //下载修改物料信息模板
    $router->post('download-update-material-temp', ['uses' => 'ImportMaterialController@downloadUpdateTemplate']);
    //新增物料
    $router->post('add-material', ['uses' => 'AddMaterialController@index']);
});

/**
 * 仓库库存
 */
$router->group(['prefix' => 'warehouse-stock', 'namespace' => 'WarehouseStock'], function () use ($router) {
    // 获取仓库库存
    $router->post('get-stock', ['uses' => 'GetStockController@index']);

    //下载修改物料库存模板
    $router->post('download-update-stock-temp', ['uses' => 'EditStockImportController@downloadUpdateTemplate']);
});

/**
 * 内容管理
 */
$router->group(['prefix' => 'content-manage', 'namespace' => 'ContentManage'], function () use ($router) {
    // 保存内容
    $router->post('save', 'SaveController@index');
    // 详情
    $router->post('info', 'InfoController@index');
    // 活动列表
    $router->get('activity-list', 'ActivityListController@index');
});

/**
 * 活动关联的商品接口
 */
$router->group(['prefix' => 'activity-goods', 'namespace' => 'ActivityGoods'], function () use ($router) {
    // 根据spuIds获取skuList
    $router->post('get-sku-by-spu', 'GetSkuListBySpuIdsController@index');
});

/**
 * 优惠活动-满减满折
 * @author 胡龙阳
 */
$router->group(['prefix' => 'discount-activity', 'namespace' => 'DiscountActivity'], function () use ($router) {
    // 列表
    $router->post('list', 'ListController@index');
    // 新增
    $router->post('add', 'CreateController@index');
    // 编辑
    $router->post('update', 'UpdateController@index');
    // 详情
    $router->post('detail', 'DetailController@index');
    // 更新状态
    $router->post('change-status', 'UpdateStatusController@index');
    // 删除
    $router->post('delete', 'DeleteController@index');
    // 修改库存
    $router->post('update-up-stock', 'UpdateController@updateFullGiftUpStock');
    // 修改礼品状态
    $router->post('update-gift-status', 'UpdateController@updateGiftStatus');
});

// 活动公共接口
$router->group(['prefix' => 'activity-common', 'namespace' => 'ActivityCommon'], function () use ($router) {
    $router->post('activity', [
        'uses'     => 'ActivityController@index',
        'describe' => '活动通用(列表/新增/删除/详情/编辑/上架/下架)',
        'page'     => '应用'
    ]);
    $router->post('download-sun-code', [
        'uses'     => 'DownloadSunCodeController@index',
        'describe' => '下载活动太阳码',
        'page'     => '应用'
    ]);
    $router->post('activity-list', [
        'uses'     => 'ActivityListController@index',
        'describe' => '获取活动列表接口',
        'page'     => '应用'
    ]);
    $router->post('activity-spu-list', [
        'uses'     => 'ActivitySpuListController@index',
        'describe' => '获取活动商品列表接口',
        'page'     => '应用'
    ]);
    $router->post('activity-type-list', [
        'uses'     => 'ActivityTypeListController@index',
        'describe' => '获取活动列表接口',
        'page'     => '应用'
    ]);
    $router->post('export', [
        'uses'     => 'ExportController@index',
        'describe' => '活动列表导出',
        'page'     => '应用'
    ]);
    $router->post('increase-stock', [
        'uses'     => 'IncreaseStockController@index',
        'describe' => '活动库存增加',
        'page'     => '应用'
    ]);
    $router->post('decrease-stock', [
        'uses'     => 'DecreaseStockController@index',
        'describe' => '活动库存减少',
        'page'     => '应用'
    ]);
    $router->post('generate-sun-code', [
        'uses'     => 'GenerateSunCodeController@index',
        'describe' => '生成活动太阳码',
        'page'     => '应用'
    ]);
    $router->post('activity-sku', [
        'uses'     => 'ActivitySkuController@index',
        'describe' => '获取活动占用sku列表接口',
        'page'     => '应用'
    ]);
    $router->post('refresh-activity-sku', [
        'uses'     => 'RefreshActivityGoodsListController@index',
        'describe' => '刷新活动活动列表接口',
        'page'     => '应用'
    ]);

    // 砍价活动接口
    $router->group(['prefix' => 'haggle', 'namespace' => 'Haggle'], function () use ($router) {
        $router->post('users', [
            'uses'     => 'UserListController@index',
            'describe' => '砍价用户列表接口',
            'page'     => '应用'
        ]);
    });

    // 先享卡活动接口
    $router->group(['prefix' => 'discount-card', 'namespace' => 'DiscountCard'], function () use ($router) {
        // 开卡信息列表
        $router->post('open-card', [
            'uses'     => 'OpenCardListController@index',
            'describe' => '开卡信息',
            'page'     => '应用'
        ]);
        // 卡订单列表
        $router->post('card-order', [
            'uses'     => 'CardOrderListController@index',
            'describe' => '卡订单列表',
            'page'     => '应用'
        ]);
        // 开卡信息列表导出
        $router->post('open-card-export', [
            'uses'     => 'OpenCardListExportController@index',
            'describe' => '开卡信息',
            'page'     => '应用'
        ]);
        // 卡订单列表导出
        $router->post('card-order-export', [
            'uses'     => 'CardOrderListExportController@index',
            'describe' => '卡订单列表',
            'page'     => '应用'
        ]);
        // 卡状态修改
        $router->post('card-status', [
            'uses'     => 'CardStatusController@index',
            'describe' => '卡订单列表',
            'page'     => '应用'
        ]);
    });
});


/**
 * 企业付款到零钱
 */
$router->group([
    'prefix'    => 'enterprise-payment',
    'namespace' => 'EnterprisePayment'
], function () use ($router) {

    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '企业红包活动.新增',
        'page'     => '企业红包活动.添加'
    ]);
    $router->post('save', [
        'uses'     => 'SaveController@index',
        'describe' => '企业红包活动.编辑',
        'page'     => '企业红包活动.编辑'
    ]);
    $router->post('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '企业红包活动.详情',
        'page'     => '企业红包活动.详情'
    ]);
    $router->post('status', [
        'uses'     => 'StatusController@index',
        'describe' => '企业红包.更新活动状态',
        'page'     => '企业红包.更新活动状态'
    ]);
    $router->post('estimated-amount', [
        'uses'     => 'EstimatedAmountController@index',
        'describe' => '企业红包.企业红包动态分布预计金额',
        'page'     => '企业红包.企业红包动态分布预计金额'
    ]);
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '企业红包.列表',
        'page'     => '企业红包.列表'
    ]);
    $router->post('reconciliation-list', [
        'uses'     => 'ReconciliationListController@index',
        'describe' => '企业红包.对账列表',
        'page'     => '企业红包.对账列表'
    ]);
    $router->post('reconciliation-member', [
        'uses'     => 'ReconciliationMemberController@index',
        'describe' => '企业红包.对账员工列表',
        'page'     => '企业红包.对账员工列表'
    ]);
    $router->post('budget-list', [
        'uses'     => 'BudgetListController@index',
        'describe' => '企业红包.红包平台充值对账列表',
        'page'     => '企业红包.红包平台充值对账列表'
    ]);
    $router->post('create-export-budget-list', [
        'uses'     => 'CreateExportBudgetList@index',
        'describe' => '企业红包.创建导出平台对账列表',
        'page'     => '企业红包.创建导出平台对账列表'
    ]);
    $router->post('budget-total', [
        'uses'     => 'BudgetTotalController@index',
        'describe' => '企业红包.平台对账统计',
        'page'     => '企业红包.平台对账统计'
    ]);
});

/**
 * 行业
 */
$router->group([
    'prefix'    => 'enterprise',
    'namespace' => 'Enterprise'
], function () use ($router) {
    $router->post('check-member', [
        'uses'     => 'CheckMemberController@index',
        'describe' => '检查企业是否可见员工',
        'page'     => '企业微信'
    ]);
});

/**
 * 行业
 */
$router->group([
    'prefix'    => 'industry',
    'namespace' => 'Industry'
], function () use ($router) {
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '行业列表',
        'page'     => ''
    ]);
});

/**
 * 权限角色相关接口
 */
$router->group([
    'prefix'    => 'rbac',
    'namespace' => 'Rbac'
], function () use ($router) {
    $router->group([
        'prefix'    => 'role',
        'namespace' => 'Role'
    ], function () use ($router) {
        $router->post('list', [
            'uses'     => 'RoleController@lists',
            'describe' => '角色列表',
            'page'     => ''
        ]);
        $router->post('add', [
            'uses'     => 'RoleController@add',
            'describe' => '添加角色',
            'page'     => ''
        ]);
        $router->post('edit', [
            'uses'     => 'RoleController@edit',
            'describe' => '修改角色',
            'page'     => ''
        ]);
        $router->post('detail', [
            'uses'     => 'RoleController@detail',
            'describe' => '角色详情',
            'page'     => ''
        ]);
        $router->post('delete', [
            'uses'     => 'RoleController@delete',
            'describe' => '删除角色',
            'page'     => ''
        ]);
    });
});

/**
 * 客服相关
 */
$router->group([
    'prefix'    => 'customer-service',
    'namespace' => 'CustomerService'
], function () use ($router) {
    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '客服列表',
        'page'     => ''
    ]);
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '客服新增',
        'page'     => ''
    ]);
    $router->post('edit', [
        'uses'     => 'EditController@index',
        'describe' => '客服编辑',
        'page'     => ''
    ]);
    $router->post('del', [
        'uses'     => 'DelController@index',
        'describe' => '客服删除',
        'page'     => ''
    ]);
    $router->post('status', [
        'uses'     => 'StatusController@index',
        'describe' => '客服上下架状态',
        'page'     => ''
    ]);
    $router->post('edit-group', [
        'uses'     => 'EditGroupController@index',
        'describe' => '批量编辑客服分组',
        'page'     => ''
    ]);
});

/**
 * 客服分组相关
 */
$router->group([
    'prefix'    => 'customer-service-group',
    'namespace' => 'CustomerServiceGroup'
], function () use ($router) {
    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '客服分组列表',
        'page'     => ''
    ]);
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '客服分组新增',
        'page'     => ''
    ]);
    $router->post('edit', [
        'uses'     => 'EditController@index',
        'describe' => '客服分组编辑',
        'page'     => ''
    ]);
    $router->post('del', [
        'uses'     => 'DelController@index',
        'describe' => '客服分组删除',
        'page'     => ''
    ]);
    $router->post('import', [
        'uses'     => 'ImportGroupController@index',
        'describe' => '批量导入客服分组',
        'page'     => ''
    ]);
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
 * 智慧导购
 */
$router->group([
    'prefix'    => 'wisdom-guide',
    'namespace' => 'WisdomGuide'
], function () use ($router) {

    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '智慧导购.新增',
    ]);
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '智慧导购.删除',
    ]);
//    $router->post('list', [
//        'uses'     => 'ListController@index',
//        'describe' => '智慧导购.列表',
//    ]);
    // 修改新导购列表
    $router->post('list', [
        'uses'     => 'ListController@list',
        'describe' => '智慧导购.列表',
    ]);
    $router->post('export-list', [
        'uses'     => 'ExportListController@index',
        'describe' => '智慧导购.导出导购列表',
    ]);
    $router->post('batch-poster', [
        'uses'     => 'BatchPosterController@index',
        'describe' => '智慧导购.批量生成导购海报',
    ]);
    $router->post('generate-worker-poster', [
        'uses'     => 'GenerateWorkerPosterController@index',
        'describe' => '智慧导购.手动生成导购海报',
    ]);

    $router->post('admin-dept-list', [
        'uses'     => 'AdminDeptListController@index',
        'describe' => '智慧导购.企业管理员管理部门列表',
    ]);

    $router->post('admin-store-list', [
        'uses'     => 'AdminStoreListController@index',
        'describe' => '智慧导购.企业管理员管理门店列表',
    ]);

    $router->post('save-admin-scope', [
        'uses'     => 'SaveAdminScopeController@index',
        'describe' => '智慧导购.修改企业管理员管理部门/门店',
    ]);
    $router->post('user-lost-config', [
        'uses'     => 'LostConfigController@index',
        'describe' => '智慧导购.丢失客户消息推送配置',
    ]);
    $router->post('user-lost-info', [
        'uses'     => 'LostConfigController@info',
        'describe' => '智慧导购.丢失客户消息人员详情',
    ]);
    $router->post('upgrade-admin', [
        'uses'     => 'UpgradeAdminController@index',
        'describe' => '导购升级为企业管理',
    ]);
    $router->post('batch-add-guide', [
        'uses'     => 'BatchAddGuideController@index',
        'describe' => '添加导购',
    ]);
    $router->post('update-store-manager-by-position', [
        'uses'     => 'UpdateStoreManagerByPositionController@index',
        'describe' => '根据用户职位更新企业管理范围',
    ]);
});

/**
 * 导购业绩指标
 */
$router->group([
    'prefix'    => 'ba-performance',
    'namespace' => 'BaPerformance'
], function () use ($router) {
    $router->post('list', 'ListController@index');
    $router->post('export', 'ListController@export');
    $router->post('import', 'ImportController@index');
    $router->get('download', 'ImportController@download');
    $router->get('group-download', 'ImportController@groupDownload');
    $router->get('plan-task-download', 'ImportController@planTaskDownload');
});


/**
 * 积分
 */
$router->group(['prefix' => 'score', 'namespace' => 'Score'], function () use ($router) {

    // 积分规则保存
    $router->post('save', [
        'uses'     => 'SaveController@index',
        'describe' => '增加积分规则',
        'page'     => '客户管理.积分管理'
    ]);

    // 积分规则编辑
    $router->post('update', [
        'uses'     => 'UpdateController@index',
        'describe' => '编辑积分规则',
        'page'     => '客户管理.积分管理'
    ]);

    // 积分规则详情
    $router->post('rule-detail', [
        'uses'     => 'RuleDetailController@index',
        'describe' => '积分规则详情',
        'page'     => '客户管理.积分管理'
    ]);

    // 积分规则删除
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '删除积分规则',
        'page'     => '客户管理.积分管理'
    ]);

    // 积分规则禁用开启
    $router->post('status', [
        'uses'     => 'StatusController@index',
        'describe' => '积分状态',
        'page'     => '客户管理.积分管理'
    ]);

    // 积分规则列表
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '积分规则列表',
        'page'     => '客户管理.积分管理'
    ]);

    // 有效期保存
    $router->post('save-deadline', [
        'uses'     => 'SaveDeadlineController@index',
        'describe' => '积分有效期保存',
        'page'     => '客户管理.积分管理'
    ]);

    // 有效期列表
    $router->post('deadline-list', [
        'uses'     => 'DeadlineListController@index',
        'describe' => '积分有效期列表',
        'page'     => '客户管理.积分管理'
    ]);

    // 积分统计
    $router->post('count', [
        'uses'     => 'CountController@index',
        'describe' => '积分统计',
        'page'     => '客户管理.积分管理'
    ]);

    $router->get('user-score-list', [
        'uses'     => 'UserScoreListController@userScoreList',
        'describe' => '客户积分明细列表',
        'page'     => '客户管理.积分管理'
    ]);

    $router->get('user-score-detail', [
        'uses'     => 'UserScoreDetailController@userScoreDetail',
        'describe' => '客户积分明细详情',
        'page'     => '客户管理.积分管理'
    ]);

    $router->post('user-score-change', [
        'uses'     => 'UserScoreChangeController@userScoreChange',
        'describe' => '管理员录入积分',
        'page'     => '客户管理.积分管理'
    ]);
});

/**
 * 积分乐园
 */
$router->group(['prefix' => 'score-exchange', 'namespace' => 'ScoreExchange'], function () use ($router) {
    // 列表
    $router->post('list', ['uses' => 'ScoreCouponListController', 'describe' => '积分兑券列表', 'page' => '积分兑券']);
    // 设置积分
    $router->post('set-score', ['uses' => 'SetExchangeScoreController', 'describe' => '设置积分券所需积分', 'page' => '积分兑券']);
    // 改变状态
    $router->post('change-status', ['uses' => 'ChangeScoreCouponStatusController', 'describe' => '改变积分券状态', 'page' => '积分兑券']);
    // 选择优惠券
    $router->post('select-coupon', ['uses' => 'SelectCouponController', 'describe' => '选择优惠券', 'page' => '积分兑券']);
    // 创建积分券
    $router->post('create-score-coupon', ['uses' => 'CreateScoreCouponController', 'describe' => '创建积分券', 'page' => '积分兑券']);
});

/**
 * 公告相关
 */
$router->group([
    'prefix'    => 'announcement',
    'namespace' => 'Announcement'
], function () use ($router) {
    // 列表
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '公告列表',
        'page'     => '首页'
    ]);
    // 详情
    $router->post('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '公告详情',
        'page'     => '首页'
    ]);
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
    //分销商品列表
    $router->post('goods-list', ['uses' => 'ListController@index']);
});

//分销商品配置
$router->group(['prefix' => 'distribution-config', 'namespace' => 'DistributionConfig'], function () use ($router) {
    // 保存分销设置
    $router->post('save', ['uses' => 'SaveController@index']);
    // 查看分销设置
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 二级分销设置
    $router->post('second-save', ['uses' => 'PartnerSaveController@index']);
    // 二级分销查看
    $router->post('second-detail', ['uses' => 'PartnerDetailController@index']);
});

// 全员分销统计相关
$router->group(['prefix' => 'distribution-statistics', 'namespace' => 'DistributionStatistics'], function () use ($router) {

    // 管理后台数据概览列表 按日期
    $router->post('statistic-date', ['uses' => 'StatisticDateController@index']);
    // 管理后台数据概览列表 按导购
    $router->post('statistic-guide', ['uses' => 'StatisticGuideController@index']);
    // 管理后台数据概览折线图
    $router->post('statistic-console', ['uses' => 'StatisticsConsoleController@index']);
    // 导出分销数据
    $router->post('statistic-export', ['uses' => 'StatisticalExportController@index']);
});

$router->group([
    'prefix'    => 'worker-distribution',
    'namespace' => 'WorkerDistribution'
], function () use ($router) {
    $router->group([
        'prefix'    => 'settlement',
        'namespace' => 'Settlement'
    ], function () use ($router) {

        $router->post('list', [
            'uses'     => 'ListController@index',
            'describe' => '结算记录.列表',
            'page'     => ''
        ]);
        $router->post('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '结算记录.明细',
            'page'     => ''
        ]);
        // 导出
        $router->post('export', [
            'uses'     => 'ExportController@index',
            'describe' => '导出结算记录',
            'page'     => ''
        ]);
    });

});


// 余额相关
$router->group([
    'prefix'    => 'distribution-balance',
    'namespace' => 'DistributionBalance'
], function () use ($router) {
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '余额相关.列表',
        'page'     => ''
    ]);
    // 导出
    $router->post('list-export', [
        'uses'     => 'ListExportController@index',
        'describe' => '导出余额列表',
        'page'     => ''
    ]);
    $router->post('details', [
        'uses'     => 'DetailsController@index',
        'describe' => '余额相关.明细',
        'page'     => ''
    ]);

    $router->post('details-export', [
        'uses'     => 'DetailsExportController@index',
        'describe' => '余额相关.明细导出',
        'page'     => ''
    ]);

});

$router->group([
    'prefix'    => 'help-center',
    'namespace' => 'HelpCenter'
], function () use ($router) {
    // 列表
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '帮助中心 ',
        'page'     => '侧边栏',
    ]);
    // 详情
    $router->post('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '帮助中心',
        'page'     => '侧边栏',
    ]);
});

/**
 * 商品评论
 */
$router->group(['prefix' => 'goods-comment', 'namespace' => 'GoodsComment'], function () use ($router) {
    // 列表
    $router->post('list', 'ListController@index');
    // 新增
    $router->post('add', 'AddController@index');
    // 编辑
    $router->post('edit', 'EditController@index');
    // 删除
    $router->post('delete', 'DeleteController@index');
    // 上下架
    $router->post('status', 'StatusController@index');
    // 详情
    $router->post('detail', 'DetailController@index');
    // 导入
    $router->post('import', 'ImportController@index');
    // 下载模板
    $router->post('download', 'ImportController@downloadAddTemplate');
});

/**
 * 视频号
 */
$router->group([
    'prefix'    => 'video',
    'namespace' => 'Video'
], function () use ($router) {
    $router->post('get-video-sku-base', [
        'uses'     => 'GetVideoSkuBaseController@index',
        'describe' => '获取视频号sku信息',
        'page'     => '视频号'
    ]);

    $router->post('get-category-from-pid', [
        'uses'     => 'GetCategoryFromIdController@index',
        'describe' => '通过父id获取分类',
        'page'     => '视频号'
    ]);

    $router->post('get-category-tree', [
        'uses'     => 'GetCategoryTreeController@index',
        'describe' => '获取视频号分类树',
        'page'     => '视频号'
    ]);

    $router->post('create-video-spu', [
        'uses'     => 'CreateVideoSpuController@index',
        'describe' => '创建视频号spu',
        'page'     => '视频号'
    ]);

    $router->post('get-video-spu-list', [
        'uses'     => 'GetVideoSpuListController@index',
        'describe' => '获取视频号spu列表',
        'page'     => '视频号'
    ]);

    $router->post('get-video-spu-detail', [
        'uses'     => 'GetVideoSpuDetailController@index',
        'describe' => '获取视频号spu详情',
        'page'     => '视频号'
    ]);

    $router->post('handler-video-spu', [
        'uses'     => 'HandlerVideoSpuController@index',
        'describe' => '修改视频号spu',
        'page'     => '视频号'
    ]);

    $router->post('get-video-brands', [
        'uses'     => 'GetVideoBrandsController@index',
        'describe' => '获取视频号品牌列表',
        'page'     => '视频号'
    ]);

    $router->post('get-video-brand-detail', [
        'uses'     => 'GetVideoBrandDetailController@index',
        'describe' => '获取视频号品牌详情',
        'page'     => '视频号'
    ]);

    $router->post('get-video-categories', [
        'uses'     => 'GetVideoCategoriesController@index',
        'describe' => '获取视频号类目列表',
        'page'     => '视频号'
    ]);

    $router->post('get-video-category-detail', [
        'uses'     => 'GetVideoCategoryDetailController@index',
        'describe' => '获取视频号类目详情',
        'page'     => '视频号'
    ]);

    $router->post('get-video-spu-skus', [
        'uses'     => 'GetVideoSpuSkusController@index',
        'describe' => '获取视频号商品sku列表',
        'page'     => '视频号'
    ]);

    $router->post('create-video-brand', [
        'uses'     => 'CreateVideoBrandController@index',
        'describe' => '添加视频号品牌',
        'page'     => '视频号'
    ]);

    $router->post('create-video-category', [
        'uses'     => 'CreateVideoCategoryController@index',
        'describe' => '添加视频号类目',
        'page'     => '视频号'
    ]);

    $router->post('refresh-video-categories', [
        'uses'     => 'RefreshVideoCategoriesController@index',
        'describe' => '刷新视频号商品',
        'page'     => '视频号'
    ]);
    $router->post('refresh-video-goods', [
        'uses'     => 'RefreshVideoGoodsController@index',
        'describe' => '刷新视频号商品',
        'page'     => '视频号'
    ]);
});

/**
 * 社群导入模板下载
 */
$router->group(['prefix' => 'group-chat', 'namespace' => 'GroupChat'], function () use ($router) {
    $router->get('download', 'GroupChatController@download');
});


// 三方门店及员工数据
$router->group(['prefix' => 'store', 'namespace' => 'Store'], function () use ($router) {
    // 门店列表
    $router->post('department-list', [
        'uses'     => 'DepartmentListController@index',
        'describe' => '门店列表',
        'page'     => '业务架构'
    ]);
    // 员工列表
    $router->post('employee-list', [
        'uses'     => 'EmployeeListController@index',
        'describe' => '员工列表',
        'page'     => '业务架构'
    ]);
    // 修改门店信息
    $router->post('save-department-info', [
        'uses'     => 'SaveDepartmentInfoController@index',
        'describe' => '修改门店信息',
        'page'     => '业务架构'
    ]);
    // 删除门店
    $router->post('delete-department', [
        'uses'     => 'DeleteDepartmentController@index',
        'describe' => '删除门店',
        'page'     => '业务架构'
    ]);
    // 关闭/打开门店
    $router->post('close-or-open-department', [
        'uses'     => 'CloseOrOpenDepartmentController@index',
        'describe' => '开店/闭店',
        'page'     => '业务架构'
    ]);
    // 导购分配/导购转岗
    $router->post('move-employee-department', [
        'uses'     => 'MoveEmployeeDepartmentController@index',
        'describe' => '分配导购',
        'page'     => '业务架构'
    ]);
    // 同步数据到企微
    $router->post('sync-data-to-wechat', [
        'uses'     => 'SyncDataToWechatController@index',
        'describe' => '同步数据到企微',
        'page'     => '业务架构'
    ]);
    // 文件中错误信息记录下载
    $router->post('export-sftp-error-record-list', [
        'uses'     => 'ExportSftpErrorRecordListController@index',
        'describe' => '文件中错误信息记录下载',
        'page'     => '业务架构'
    ]);
    // 文件拉取记录列表
    $router->post('sftp-log-list', [
        'uses'     => 'SftpLogListController@index',
        'describe' => '文件拉取记录列表',
        'page'     => '业务架构'
    ]);

    // 门店列表批量更新
    $router->post('import-store-info', [
        'uses'     => 'ImportStoreInfoController@index',
        'describe' => '门店列表批量更新',
        'page'     => '业务架构'
    ]);

    // 下载门店列表批量更新文件模板
    $router->post('download-store-list', [
        'uses'     => 'DownloadStoreListController@index',
        'describe' => '下载门店列表批量更新文件模板',
        'page'     => '业务架构'
    ]);

    // 手动更新部门到门店列表（coach)
    $router->post('sync-data-by-departments', [
        'uses'     => 'SyncDataByDepartmentsController@index',
        'describe' => '手动更新部门到门店列表（coach)',
        'page'     => '业务架构'
    ]);

});

// 邀请函
$router->group(['prefix' => 'invite', 'namespace' => 'Invite'], function () use ($router) {
    $router->post('create-qrcode', 'InviteController@index');
});

// 成长激励
$router->group(['prefix' => 'excitation', 'namespace' => 'Excitation'], function () use ($router) {
    $router->addRoute(['POST', 'GET'], 'download-crowd-pack', 'MedalController@DownloadCrowdPack');
    $router->post('read-crowd-pack-data', 'MedalController@ReadCrowdPackData');
    $router->addRoute(['POST', 'GET'], 'download-reissue', 'MedalController@DownloadReissue');
    $router->post('read-reissue-data', 'MedalController@ReadReissueData');
    // 下载奖章领取记录
    $router->addRoute(['POST', 'GET'], 'download-medal-record', 'MedalRecordController@download');

    // 下载手动添加导购积分模板
    $router->post('download-manual-points-template', [
        'uses'     => 'DownloadManualPointsTemplateController@index',
        'describe' => '下载手动添加导购积分模板',
        'page'     => 'SA激励'
    ]);

    // 导入手动添加导购积分模板
    $router->post('import-manual-points', [
        'uses'     => 'ImportManualPointsController@index',
        'describe' => '导入手动添加导购积分模板',
        'page'     => 'SA激励'
    ]);

    // 历史积分记录导出
    $router->post('download-points-record', [
        'uses'     => 'DownloadPointsRecordController@index',
        'describe' => '下载手动添加导购积分模板',
        'page'     => 'SA激励'
    ]);

    // 下载手动添加导购金币模板
    $router->get('download-manual-coins-template', [
        'uses'     => 'DownloadManualCoinsTemplateController@index',
        'describe' => '下载手动添加导购金币模板',
        'page'     => 'SA激励'
    ]);

    // 导入手动添加导购积分模板
    $router->post('import-manual-coins', [
        'uses'     => 'ImportManualCoinsController@index',
        'describe' => '导入手动添加导购金币模板',
        'page'     => 'SA激励'
    ]);

});


// 传讯
$router->group(['prefix' => 'notification', 'namespace' => 'Notification'], function () use ($router) {

    // 下载手动添加导购金币模板
    $router->get('download-crowd-template', [
        'uses'     => 'DownloadCrowdTemplateController@index',
        'describe' => '下载传讯白名单导入',
        'page'     => '传讯'
    ]);

    // 导入手动添加导购积分模板
    $router->post('import-crowd', [
        'uses'     => 'ImportCrowdController@index',
        'describe' => '传讯白名单导入',
        'page'     => '传讯'
    ]);

    // 下载失败导入记录
    $router->post('export-fail-import-crowd', [
        'uses'     => 'ExportFailImportCrowdController@index',
        'describe' => '传讯白名单导入',
        'page'     => '传讯'
    ]);

});



