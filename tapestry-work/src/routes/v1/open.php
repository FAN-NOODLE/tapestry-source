<?php

/**
 * 好物集
 */
$router->group(['prefix' => 'guide', 'namespace' => 'Guide'], function () use ($router) {
    // 导购认证
    $router->post('bind-user', ['uses' => 'BindUserController@index']);
});

$router->group(['prefix' => 'share_link', 'namespace' => 'ShareLink'], function () use ($router) {
    // 分销商品链路绑定与解除
    $router->post('bind', ['uses' => 'ShareBindController@index']);
    //查询分销商品链路绑定关系
    $router->post('check', ['uses' => 'ShareCheckController@index']);
});
