<?php

/**
 * @var \Laravel\Lumen\Routing\Router $router
 */

$router->get('/', ['uses' => 'IndexController@index']);

// 接口文档
$router->get('doc/{name}', ['uses' => 'DocController@index']);

$router->post('upload', ['uses' => 'UploadController@index']);

$router->get('edit-split-flow', ['uses' => 'EditSplitFlowController@index']);
$router->post('edit-split-flow', ['uses' => 'EditSplitFlowController@index']);

$router->get('test', ['uses' => 'TestController@index']);

$router->get('health', function (){
    return "true";
});

/**
 * 小程序授权提审相关
 */
$router->group(['prefix' => 'callback', 'namespace' => 'Callback'], function () use ($router) {
    // 小程序代码全网发布
    $router->get('mini_program/code', 'MiniProgramController@index');
    // 小程序代码全网发布post
    //$router->post('mini_program/code', 'MiniProgramController@index');
    // 公众号、小程序审核回调信息
    $router->get('account/{appId}/serve', 'AccountController@index');
    // 公众号、小程序审核回调信息post
    $router->post('account/{appId}/serve', 'AccountController@index');
    // 服务号，小程序授权回调
    $router->get('authorize-access-token', 'MpMiniProgramController@index');
    // UC授权安装成功回调地址
    $router->get('contact-install', 'ContactController@index');
    // 企业微信应用授权回调
    $router->get('work-plugin', 'WorkPluginController@index');
    // uc统一费回调
    $router->post('uc-union-notify', 'UcUnionNotifyController@index');
});

/**
 * 开放平台相关
 */
$router->group(['prefix' => 'open-platform'], function () use ($router) {
    // 接收微信事件消息地址（授权、取消授权、令牌推送等等）
    $router->get('serve', ['uses' => 'ServerController@index']);
    // 接收微信事件消息地址
    $router->post('serve', ['uses' => 'ServerController@index']);
});

/**
 * 调试相关
 */
$router->group(['prefix' => 'callback', 'namespace' => 'Callback'], function () use ($router) {
    // 根据提审id查询小程序提审状态（临时用）
    $router->get('get-uditid-log', 'TesterController@getAuditIdLog');
    // 查询小程序提审状态（临时用）
    $router->get('get-udit-log', 'TesterController@getAuditLog');
    // 添加体验者（临时用）
    $router->get('create-tester', 'TesterController@createTester');
    // 提交小程序审核（临时用）
    $router->get('submit-audit', 'TesterController@submitAudit');
    // 提交指定小程序审核（临时用）
    $router->get('submit-specified', 'TesterController@submitSpecified');
    // 完整提交小程序审核流程（临时用）
    $router->get('submit', 'TesterController@submitVersion');
    // 上线小程序（临时用）
    $router->get('release', 'TesterController@release');
    // 下线小程序（临时用）
    $router->get('offline', 'TesterController@offline');
    // 取消提审（临时用）
    $router->get('cancel-audit', 'TesterController@cancelAudit');
    // 生成体验二维码
    $router->get('qrcode', 'TesterController@qrcode');
    // 更新可信域名
    $router->get('set-domain', 'TesterController@setDomains');
    // 手动打包提审
    $router->post('hand-pack-audit', 'TesterController@handPackAudit');
    // 微信事件验证回调
    $router->get('wx-event', 'TesterController@wxEvent');
    $router->post('wx-event', 'TesterController@wxEvent');
    $router->get('bind-to-open', 'TesterController@bindToOpenPlatform');

    // 新增开放平台配置
    $router->post('deploy-open-config', 'TesterController@deployOpenConfig');

    // 更新小程序域名配置
    $router->post('update-domain-json', 'TesterController@updateDomainJson');
});

/**
 * 小程序自动化打包相关
 */
$router->group(['prefix' => 'callback', 'namespace' => 'Callback'], function () use ($router) {
    // 自动打包回调
    $router->post('auto-pack-audit', ['uses' => 'AutoPackAuditController@index', 'as' => 'autoPackAudit']);
    // 开发码回调
    $router->post('develop-code/{hash}', ['uses' => 'DevelopCodeController@index', 'as' => 'developCode']);
});

$router->get('monitoring_template_num', 'MonitoringTemplate@index');