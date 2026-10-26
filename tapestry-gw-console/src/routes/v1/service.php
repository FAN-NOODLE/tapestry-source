<?php

// 应用及管理权限相关
$router->group(['prefix' => 'application', 'namespace' => 'Application'], function () use ($router) {
    $router->post('syncApplication', ['uses' => 'SyncApplicationController@index']);
    // BOSS关闭应用权限调用
    $router->post('close', ['uses' => 'CloseController@index']);
    // 设置不展示畅移logo
    $router->post('white_list', ['uses' => 'WhiteListController@index']);
});

/**
 * 素材管理
 */
$router->group(['prefix' => 'media', 'namespace' => 'Media'], function () use ($router) {
    // 素材列表
    $router->post('list', ['uses' => 'ListController@index']);
    // 素材详情
    $router->post('detail', ['uses' => 'DetailController@index']);
});
