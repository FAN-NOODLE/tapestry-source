<?php

/** @var Laravel\Lumen\Routing\Router $router */

$router->get('/', ['uses' => 'IndexController@index']);

// 查看接口文档
$router->get('doc/{name}', ['uses' => 'DocController@index']);
$router->get('heartbeat', ['uses' => 'HeartbeatController@index']);
