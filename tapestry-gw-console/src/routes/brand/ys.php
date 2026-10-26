<?php

/**
 * @var \Laravel\Lumen\Routing\Router $router
 */

/**
 * 社群统计
 */
$router->group(['prefix' => 'community-statistics', 'namespace' => 'CommunityStatistics'], function () use ($router) {
    // 社群数据
    $router->post('index', 'StatisticsController@index');

    // 社群明细列表(按社群)
    $router->post('list', 'StatisticsController@communityListByGroup');

    // 导购社群明细列表(按导购)
    $router->post('list-by-member', 'StatisticsController@communityListByMember');

    // 导出社群明细列表
    $router->post('export', 'ExportController@index');
});