<?php

/**
 * @var \Laravel\Lumen\Routing\Router $router
 */

$router->get('', function () {
    return 'micro-authority';
});
$router->get('log', ['uses' => 'IndexController@index']);
// 查看运行日志
$router->get('runtime[/{date}]', ['uses' => 'IndexController@runtime']);
// 查看接口文档
$router->get('doc/{name}', ['uses' => 'DocController@index']);
$router->post('upload', ['uses' => 'UploadController@index']);
