<?php

/**
 * 佣金设置
 */
$router->group(['prefix' => 'customer', 'namespace' => 'Customer'], function () use ($router) {
    // 添加外部联系人成功回调
    $router->post('add', ['uses' => 'AddController@index']);
    $router->get('callback-other/{id}', ['uses' => 'OtherController@index']);

});
// 回调新路由，仅用于解耦迭代的测试、验收和灰度期间
$router->group(['prefix' => 'customer-new', 'namespace' => 'Customer'], function () use ($router) {
    // 添加外部联系人成功回调
    $router->post('add', ['uses' => 'AddController@index']);

});
