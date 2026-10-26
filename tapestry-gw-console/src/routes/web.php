<?php

$router->get('/', function () {
    return 'gw-console-v1';
});
$router->get('log', ['uses' => 'IndexController@index']);
// 查看运行日志
$router->get('runtime[/{date}]', ['uses' => 'IndexController@runtime']);
// 查看接口文档
$router->get('doc/{name}', ['uses' => 'DocController@index']);
$router->post('upload', ['uses' => 'UploadController@index']);
$router->post('test', ['uses' => 'TestController@index']);
$router->get('test', ['uses' => 'TestController@index']);
$router->get('heartbeat', ['uses' => 'HeartbeatController@index']);
