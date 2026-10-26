<?php
/**
 * Created By PhpStorm
 * @Date: 2021/7/1
 * @Time: 11:48
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: ys.php
 * @var \Laravel\Lumen\Routing\Router $router
 */

/**
 * 企业品牌权限
 */
$router->group([
    'prefix'    => 'brandAuth',
    'namespace' => 'BrandAuth'
], function () use ($router) {
    $router->post('auth', [
        'uses'     => 'AuthController@index',
        'describe' => '品牌权限'
    ]);
});
