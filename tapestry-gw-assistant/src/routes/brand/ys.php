<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/24
 * @Time: 19:42
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: ys.php
 */
use Illuminate\Support\Facades\Route;

// 社群统计
Route::group(['namespace' => 'CommunityStatistics', 'prefix' => 'community-statistics'], function () {
    // 社群统计数量
    Route::post('count', ['uses' => 'StatisticController@index']);
    // 社群未参加活动人数统计
    Route::post('list', ['uses' => 'StatisticController@getNotJoinActivityMemberList']);
});