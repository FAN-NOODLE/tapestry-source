<?php

/** @var $router */
$router->group([], function () use ($router) {
    $router->post('send-result', [
        'uses' => 'MaPlanController@sendResult',
        'desc' => 'ma任务的结果回执'
    ]);
    $router->post('send-finish', [
        'uses' => 'MaPlanController@sendFinish',
        'desc' => 'ma任务的完成回执'
    ]);
    $router->post('moment-send-result', [
        'uses' => 'MaPlanController@momentSendResult',
        'desc' => '朋友圈任务群发助手导购发送后回执可见范围'
    ]);

    $router->post('moment-comments', [
        'uses' => 'MaPlanController@momentComments',
        'desc' => '朋友圈任务回执评论和点赞'
    ]);
});