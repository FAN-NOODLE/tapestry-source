<?php
/**
 * @var $router \Illuminate\Support\Facades\Route
 */
/**
 * @var \Laravel\Lumen\Routing\Router $router
 */

/**
 * 授权模块
 */
$router->group(['prefix' => 'auth', 'namespace' => 'Auth'], function () use ($router) {
    // 企业微信通讯录授权详情获取接口
    $router->post('auth-contact', ['uses' => 'AuthContactDetailController@index']);
    // 客户联系配置
    $router->post('uc-callback', ['uses' => 'UcCallbackController@index']);
    $router->post('set-wxqy-token', ['uses' => 'UcCallbackController@setToken']);
    // 企业微信小程序信息保存接口
    $router->post('work-setting', ['uses' => 'WorkSettingController@index']);

    // 企业微信自建应用安装
    $router->post('install-app', ['uses' => 'InstallAppController@install']);
    // 企业微信自建应用更新
    $router->post('update-app', ['uses' => 'InstallAppController@update']);
    // 获取企业微信自建应用信息
    $router->post('get-work-plugin', ['uses' => 'InstallAppController@getWorkPlugin']);
    // 企业微信会话存档自建应用安装
    $router->post('install-session-archiving-app', ['uses' => 'InstallSessionArchivingAppController@install']);
    // 企业微信会话存档自建应用更新
    $router->post('update-session-archiving-app', ['uses' => 'InstallSessionArchivingAppController@update']);
});
