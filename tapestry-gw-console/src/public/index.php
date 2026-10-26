<?php
header_remove('X-Powered-By');
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

define('LUMEN_START', microtime(true));

$app = require __DIR__ . '/../bootstrap/app.php';

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
$httpOrigin = $_SERVER['HTTP_ORIGIN'] ?? '';
$env = config('app.env');
for ($i = 0; $i < 10; $i++) {
    $localOrigin[] = 'http://localhost:2000' . $i;//本地调试域名
}

if ($env === 'prod') {
    // 生产环境
    $allowOrigin = [
        'https://we.coach.com.cn',
        "https://we-h5.coach.com.cn",
        "https://we-h5.stuartweitzman.cn",
        "https://we.stuartweitzman.cn",
        "https://we-h5.katespade.cn",
        "https://we.katespade.cn",
    ];
} elseif ($env === 'uat') {
    $allowOrigin = [
        'https://tapestry-web-uat.coach.com.cn',
        "https://tapestry-h5-uat.coach.com.cn",
        "https://tapestry-h5-uat.stuartweitzman.cn",
        "https://tapestry-web-uat.stuartweitzman.cn",
        "https://tapestry-h5-uat.katespade.cn",
        "https://tapestry-web-uat.katespade.cn",
    ];
} else {
    // 其他环境
    $allowOrigin = [$_SERVER['HTTP_ORIGIN'] ?? '*'];
}

$origin = '';
if (in_array($httpOrigin, array_merge($allowOrigin, $localOrigin ?? []))) {
    $origin = $httpOrigin;
}
// 支持跨域访问
header('Access-Control-Allow-Origin: ' . $origin);
header('Access-Control-Allow-Headers: Content-Type,Enterprise-Hash,Auth-Token');
header('Access-Control-Allow-Methods: *');
header('Access-Control-Allow-Credentials: true');

if ($_SERVER['REQUEST_METHOD'] == 'OPTIONS') {
    return true;
}

$app->run();
