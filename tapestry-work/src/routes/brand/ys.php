<?php
/**
 * Created By PhpStorm
 * @Date: 2021/7/1
 * @Time: 14:49
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: ys.php
 * @var \Laravel\Lumen\Routing\Router $router
 */

// 社群-数据罗盘
$router->group(['prefix' => 'community-statistic', 'namespace' => 'CommunityStatistic'], function () use ($router) {
    // 社群数据统计(admin)
    $router->post('count-admin', ['uses' => 'StatisticsForAdminController@index']);
    // 社群数据统计(app)
    $router->post('count-app', ['uses' => 'StatisticsForAppController@index']);
    // 未参加群活动人员列表(app)
    $router->post('list', ['uses' => 'NotJoinActivityMemberListController@index']);
    // 社群数据明细列表(admin)(按社群)
    $router->post('community-list-by-group-for-console', ['uses' => 'CommunityListByGroupForConsoleController@index']);
    // 社群数据明细列表(admin)(按导购)
    $router->post('community-list-by-member-for-console', ['uses' => 'CommunityListByMemberForConsoleController@index']);
});