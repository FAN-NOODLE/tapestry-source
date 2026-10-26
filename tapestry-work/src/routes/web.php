<?php

$router->get('/', ['uses' => 'IndexController@index']);
$router->get('doc/{name}', ['uses' => 'DocController@index']);
$router->get('test', ['uses' => 'TestController@index']);
// 用户行为
$router->post('micro-ssg/v1/callback/user', ['uses' => 'Service\V1\UserAction\OpenController@index']);


$router->get('/health', ['uses' => 'HealthController@index']);