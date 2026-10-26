<?php

/**
 * 员工相关
 */
$router->group(['prefix' => 'auth', 'namespace' => 'Auth'], function () use ($router) {
    // 登录接口
    $router->post('login', ['uses' => 'LoginController@index']);
    // 获取配置信息
    $router->get('agent_config', ['uses' => 'AgentConfigController@index']);
    // 获取配置信息
    $router->get('config', ['uses' => 'ConfigController@index']);
    // 模拟登录接口，前端无法在模拟器中使用 wx.qy.login 的情况
    $router->get('mock', ['uses' => 'MockController@index']);
});

/**
 * 客户相关
 */
$router->group(['prefix' => 'customer', 'namespace' => 'Customer'], function () use ($router) {
    // 用户详情
    $router->get('detail',  ['uses' => 'CustomerController@detail']);
    // 客户订单详情
    $router->get('order-detail', ['uses' => 'OrderDetailController@index']);
    // 批量查询物流
    $router->post('logistics-query', ['uses' => 'LogisticsQueryController@index']);
    // 标签列表
    $router->post('tag-list',  ['uses' => 'TagController@list']);
    // 打标签
    $router->post('member-give-tag',  ['uses' => 'MemberGiveTagController@memberGiveTag']);
});

/**
 * 现金红包
 */
$router->group(['prefix' => 'cash_coupon', 'namespace' => 'CashCoupon'], function () use ($router) {
    // 现金红包活动列表
    $router->get('list', ['uses' => 'ListController@index']);
    // 发送现金红包
    $router->post('send_red_pack', ['uses' => 'SendRedPackController@index']);
});

/**
 * 素材管理
 */
$router->group(['prefix' => 'source', 'namespace' => 'Source'], function () use ($router) {
    // 素材列表
    $router->get('list', ['uses' => 'ListController@index']);
    // 素材详情
    $router->get('detail', ['uses' => 'DetailController@index']);
});

/**
 * 营销推广
 */
$router->group(['prefix' => 'marketing-promotion', 'namespace' => 'MarketingPromotion'], function () use ($router) {
    // 商品列表
    $router->post('goods-list', ['uses' => 'GoodsListController@index']);
    // 活动列表
    $router->post('activity-list', ['uses' => 'ActivityListController@index']);
    // 优惠券列表
    $router->post('coupon-list', ['uses' => 'CouponListController@index']);
    // 门店列表
    $router->post('store-list', ['uses' => 'StoreListController@index']);
    // 获取海报
    $router->post('poster', ['uses' => 'PosterController@index']);
});
