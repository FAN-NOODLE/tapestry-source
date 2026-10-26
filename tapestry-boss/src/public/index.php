<?php

/*
|--------------------------------------------------------------------------
| Create The Application
|--------------------------------------------------------------------------
|
| First we need to get an application instance. This creates an instance
| of the application / container and bootstraps the application so it
| is ready to receive HTTP / Console requests from the environment.
|
*/

$app = require __DIR__.'/../bootstrap/app.php';

/*
|--------------------------------------------------------------------------
| Run The Application
|--------------------------------------------------------------------------
|
| Once we have the application, we can handle the incoming request
| through the kernel, and send the associated response back to
| the client's browser allowing them to enjoy the creative
| and wonderful application we have prepared for them.
|
*/

if (config('app.env') === 'prod') {
    // 生产环境
    $origin = 'https://we.coach.com.cn';
} else {
    // 其他环境
    $origin = $_SERVER['HTTP_ORIGIN'] ?? '*';
}

// 支持跨域访问
header('Access-Control-Allow-Origin: '.$origin);
header('Access-Control-Allow-Headers: Content-Type');
header('Access-Control-Allow-Methods: *');
header('Access-Control-Allow-Credentials: true');

if ($_SERVER['REQUEST_METHOD'] == 'OPTIONS') {
    return true;
}

$app->run();
