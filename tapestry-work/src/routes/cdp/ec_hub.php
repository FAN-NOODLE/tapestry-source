<?php

/** @var $router */
$router->post('sa-assignment',[
    'uses' => 'EcHubController@saAssignment',
    'desc' => 'coach导购添加顾客'
]);
$router->post('sa-un-assignment',[
    'uses' => 'EcHubController@saUnAssignment',
    'desc' => 'coach导购删除顾客'
]);
$router->post('get-ec-url-scheme',[
    'uses' => 'EcHubController@getEcUrlScheme',
    'desc' => '获取EC小程序h5链接地址'
]);
$router->post('get-ec-qr-code',[
    'uses' => 'EcHubController@getEcQrCode',
    'desc' => '获取EC小程序的太阳码'
]);
$router->post('get-qrcode-and-link',[
    'uses' => 'EcHubController@getQrCodeAndLink',
    'desc' => '获取EC小程序的太阳码和h5链接地址'
]);
