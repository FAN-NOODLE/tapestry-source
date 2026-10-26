<?php

/**
 * @var \Laravel\Lumen\Routing\Router $router
 */

/**
 * 企业模块
 */
$router->group(['prefix' => 'enterprise', 'namespace' => 'Enterprise'], function () use ($router) {
    // 企业列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 企业列表(不分页)
    $router->post('list-no-page', ['uses' => 'ListNoPageController@index']);
    // 企业详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 编辑企业
    $router->post('save', ['uses' => 'SaveController@index']);
    // 企业微信小程序信息保存接口
    $router->post('work-setting', ['uses' => 'WorkSettingController@index']);
    // 商城小程序信息保存接口
    $router->post('shop-setting', ['uses' => 'ShopSettingController@index']);
    // 查询企业服务信息
    $router->post('enterprise-info', ['uses' => 'EnterpriseInfo@index']);
    // 获取商户信息
    $router->post('merchant_setting_info', ['uses' => 'MerchantSettingDetailController@index']);
    // 保存商户信息
    $router->post('merchant_setting_save', ['uses' => 'MerchantSettingSaveController@index']);
    // 获取所有正式企业信息
    $router->get('list-official', ['uses' => 'ListOfficialController@index']);
    // 添加企业接口
    $router->post('add', ['uses' => 'AddController@index']);
    // 添加企业接口
    $router->post('buy-notice', ['uses' => 'BuyNoticeController@index']);
    // 企业小程序太阳码
    $router->post('miniprogram-code', ['uses' => 'MiniProgramCodeController@index']);
    // 更新小程序Token
    $router->post('update-refresh-token', ['uses' => 'RefreshTokenController@updateMiniProgramToken']);

    // 根据小程序名称获取企业信息列表
    $router->post('list-by-mini-name', ['uses' => 'GetListByMiniNameController@index']);
});

/**
 * 企业管理员模块
 */
$router->group(['prefix' => 'enterprise_admin', 'namespace' => 'EnterpriseAdmin'], function () use ($router) {
    // 企业管理员添加
    $router->post('add', ['uses' => 'AddController@index']);
    // 企业管理员修改
    $router->post('save', ['uses' => 'SaveController@index']);
    // 企业管理员列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 企业管理员
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 企业管理员登陆更新
    $router->post('login_update', ['uses' => 'LoginUpdateController@index']);
    // 企业管理员密码修改
    $router->post('password', ['uses' => 'PasswordController@index']);
    // 企业管理员密码修改
    $router->post('reset-password', ['uses' => 'ResetPasswordController@index']);
    // 企业管理员状态修改
    $router->post('status', ['uses' => 'StatusController@index']);
    // 企业管理员删除
    $router->post('delete', ['uses' => 'DeleteController@index']);
});

/**
 * 企业管理员模块
 */
$router->group(['prefix' => 'enterprise-admin', 'namespace' => 'EnterpriseAdmin'], function () use ($router) {
    // 检查邮箱是否唯一接口
    $router->post('check-email', ['uses' => 'CheckEmailController@index']);
});

/**
 * 授权模块
 */
$router->group(['prefix' => 'auth', 'namespace' => 'Auth'], function () use ($router) {
    // 商城小程序授权详情接口
    $router->get('shop-auth-detail', ['uses' => 'ShopAuthDetailController@index']);
    // 企业微信小程序授权详情接口
    $router->get('work-auth-detail', ['uses' => 'WorkAuthDetailController@index']);
    // 小程序授权详情获取接口
    $router->get('auth-detail', ['uses' => 'AuthAccountDetailController@index']);
    // 企业微信通讯录授权详情获取接口
    $router->get('auth-contact', ['uses' => 'AuthContactDetailController@index']);
    // 获取企业的商城小程序appid
    $router->get('get-shop-app-id', ['uses' => 'GetShopAppIDController@index']);
    // 小程序通过第三方应用ID获取企业ID信息
    $router->post('auth-enterprise', ['uses' => 'AuthEnterpriseController@index']);
    // 小程序信息详情
    $router->post('mini-detail', ['uses' => 'AuthDetailController@index']);
    // 小程序信息详情
    $router->post('get-shop-secret', ['uses' => 'GetShopSecretController@index']);
    // 客户联系配置
    $router->post('uc-callback', ['uses' => 'UcCallbackController@index']);

    // 企业微信自建应用安装
    $router->post('install-app', ['uses' => 'InstallAppController@install']);
    // 企业微信自建应用域名校验文件是否上传
    $router->post('check-file', ['uses' => 'InstallAppController@checkFile']);
    // 获取企业微信自建应用信息
    $router->get('get-work-plugin', ['uses' => 'InstallAppController@getWorkPlugin']);
    // 企业微信自建应用更新
    $router->post('update-app', ['uses' => 'InstallAppController@update']);
    // 获取智慧导购信息
    $router->post('dao-gou-plugin', ['uses' => 'InstallAppController@daoGouPlugin']);
    // 安装自建智慧导购
    $router->post('dao-gou-install', ['uses' => 'InstallAppController@daoGouInstall']);
    // 安装第三方应用（智慧导购等）
    $router->post('install-third-app', ['uses' => 'InstallAppController@thirdInstall']);
    // 导购配置获取
    $router->post('get-dao-gou-config', ['uses' => 'InstallAppController@getDaoGouConfig']);
});

/**
 * 应用管理
 */
$router->group(['prefix' => 'app-manage', 'namespace' => 'AppManage'], function () use ($router) {
    // 应用列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 获取当前企业应用是否已经安装
    $router->post('status', ['uses' => 'StatusController@index']);
});

/**
 * 企业配置模块
 */
$router->group(['prefix' => 'enterprise-config', 'namespace' => 'EnterpriseConfig'], function () use ($router) {
    // 企业配置保存
    $router->post('save', ['uses' => 'SaveController@index']);
    // 获取企业配置
    $router->post('info', ['uses' => 'InfoController@index']);
    // 根据字段状态查询企业ID集合
    $router->post('select-enterprise-id', ['uses' => 'SelectEnterpriseIdController@index']);
    // 设置打印小票二维码
    $router->post('set-print-code', ['uses' => 'SetPrintCodeController@setPrintCode']);
});
/**
 * 小程序版本管理
 */
$router->group(['prefix' => 'mini-program', 'namespace' => 'MiniProgram'], function () use ($router) {
    // 添加体验者
    $router->post('create-tester', ['uses' => 'CreateTesterController@index']);
    // 详情
    $router->get('detail', 'DetailController@index');
    // 查询小程序提审状态
    $router->get('get-audit-status', ['uses' => 'getAuditStatusController@index']);
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
    //小程序手动提审
    $router->get('trig-submit', ['uses' => 'TrigSubmitController@index']);
    // 下线小程序
    $router->get('offline', ['uses' => 'OfflineController@offline']);
    // 自动打包 - 打包
    $router->post('auto-pack', ['uses' => 'AutoPackController@index', 'as' => 'autoPack']);
    // 直播间列表
    $router->post('live-info', ['uses' => 'LiveInfoController@index']);
    // 小程序基础配置信息列表
    $router->post('base-list', ['uses' => 'BaseListController@index']);
    // 小程序 - 数据收集sdk - 配置查询接口
    $router->get('sdk-config', ['uses' => 'SdkConfigController@index']);

    // 小程序 - 获取小程序线上使用版本
    $router->get('get-online-version', ['uses' => 'GetOnlineVersion@index']);

    // 小程序 - 根据企业id获取
    $router->post('get-name-by-enterpriseId', ['uses' => 'GetMiniNameByEnterpriseIdController@index']);
});

/**
 * 企业微信应用
 */
$router->group([
    'prefix'    => 'enterprise-work-plugin',
    'namespace' => 'EnterpriseWorkPlugin'
], function () use ($router) {
    // 应用列表
    $router->post('list', ['uses' => 'ListController@index']);
});

/**
 * 行业
 */
$router->group([
    'prefix'    => 'industry',
    'namespace' => 'Industry'
], function () use ($router) {
    // 行业列表
    $router->post('list', 'ListController@index');
});

/**
 * 公告
 */
$router->group([
    'prefix'    => 'announcement',
    'namespace' => 'Announcement'
], function () use ($router) {
    // 列表
    $router->post('list', 'ListController@index');
    // 详情
    $router->post('detail', 'DetailController@index');
});

/**
 * 灰度企业微信应用
 */
$router->group([
    'prefix'    => 'release-enterprise-work',
    'namespace' => 'ReleaseEnterpriseWork'
], function () use ($router) {
    // 企业信息
    $router->post('info', ['uses' => 'InfoController@index']);
});

/**
 * 企业付费模块
 */
$router->group(['prefix' => 'enterprise-payment', 'namespace' => 'EnterprisePayment'], function () use ($router) {

    // 获取企业短信数
    $router->post('get-sms', ['uses' => 'GetSmsController@index']);
    //
    $router->post('deplete-sms', ['uses' => 'DepleteSmsController@index']);

});

/**
 * 代注册
 */
$router->group(['prefix' => 'fast-register', 'namespace' => 'FastRegisterWeApp'], function () use ($router) {

    // 申请
    $router->post('apply', ['uses' => 'ApplyController@index']);
    // 详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 更新信息
    $router->post('update-info', ['uses' => 'UpdateInfoController@index']);

});
