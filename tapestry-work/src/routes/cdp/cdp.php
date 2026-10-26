<?php

/**
 * CDP
 */
/** @var $router */
$router->group([
], function () use ($router) {
    $router->post('get-customer-interested-products', [
        'uses' => 'CdpController@getCustomerInterestedProducts',
        'desc' => '获取CDP推荐商品'
    ]);
});

