<?php

use Freyo\LaravelQueueCMQ\LaravelQueueCMQServiceProvider;
use Illuminate\Filesystem\FilesystemServiceProvider;
use Illuminate\Mail\MailServiceProvider;
use Illuminate\Redis\RedisServiceProvider;
use SimpleSoftwareIO\QrCode\QrCodeServiceProvider;

require_once __DIR__ . '/../vendor/autoload.php';

try {
    (new Dotenv\Dotenv(__DIR__ . '/../'))->load();
} catch (Dotenv\Exception\InvalidPathException $e) {
    //
}
define('LOG_ID', md5(uniqid()) . rand(10000, 99999));
define('START_TIME', microtime(true));
/*
|--------------------------------------------------------------------------
| Create The Application
|--------------------------------------------------------------------------
|
| Here we will load the environment and create the application instance
| that serves as the central piece of this framework. We'll use this
| application as an "IoC" container and router for this framework.
|
*/

$app = new Laravel\Lumen\Application(
    realpath(__DIR__ . '/../')
);

$app->withFacades();

$app->withEloquent();

$app->configure('app');

$app->configure('code');

$app->configure('client');

$app->configure('database');

$app->configure('allocate');

$app->configure('cmq_setting');
// 企业微信小程序相关配置
$app->configure('work_plugin');

$app->configure('frontend');

// 云存储配置
$app->configure('cloud');
$app->configure('third_party');
//open
$app->configure('open');
//sms
$app->configure('easy-sms');
// 腾讯云配置
$app->configure('third_party');

// blob配置
$app->configure('filesystems');

// mail配置
$app->configure('mail');

/*
|--------------------------------------------------------------------------
| Register Container Bindings
|--------------------------------------------------------------------------
|
| Now we will register a few bindings in the service container. We will
| register the exception handler and the console kernel. You may add
| your own bindings here if you like or you can make another file.
|
*/

$app->singleton(
    Illuminate\Contracts\Debug\ExceptionHandler::class,
    App\Exceptions\Handler::class
);

$app->singleton(
    Illuminate\Contracts\Console\Kernel::class,
    App\Console\Kernel::class
);

/*
|--------------------------------------------------------------------------
| Register Middleware
|--------------------------------------------------------------------------
|
| Next, we will register the middleware with the application. These can
| be global middleware that run before and after each request into a
| route or middleware that'll be assigned to some specific routes.
|
*/

$app->middleware([
    App\Http\Middleware\Cors::class
]);

$app->routeMiddleware([
    'auth'      => App\Http\Middleware\Authenticate::class,
    'sign'      => App\Http\Middleware\SignMiddleware::class,
    'access'    => App\Http\Middleware\Access::class,
    'auth_open' => App\Http\Middleware\AuthOpen::class,

]);

/*
|--------------------------------------------------------------------------
| Register Service Providers
|--------------------------------------------------------------------------
|
| Here we will register all of the application's service providers which
| are used to bind services into the container. Service providers are
| totally optional, so you are not required to uncomment this line.
|
*/

$app->register(App\Providers\AppServiceProvider::class);
// $app->register(App\Providers\AuthServiceProvider::class);
// $app->register(App\Providers\EventServiceProvider::class);

$app->register(RedisServiceProvider::class);
$app->register(QrCodeServiceProvider::class);
$app->register(LaravelQueueCMQServiceProvider::class);
$app->register(FilesystemServiceProvider::class);    // 文件系统
$app->register(App\Providers\AzureBlobServiceProvider::class);              // 微软云blob
$app->register(MailServiceProvider::class); //邮件


/*
|--------------------------------------------------------------------------
| Load The Application Routes
|--------------------------------------------------------------------------
|
| Next we will include the routes file so that they can all be added to
| the application. This will provide all of the URLs the application
| can respond to, as well as the controllers that may handle them.
|
*/

$app->router->group([
    'namespace'  => 'App\Http\Controllers',
    'middleware' => 'access'
], function ($router) {
    require __DIR__ . '/../routes/web.php';
});

/**
 * 引入v1版本业务路由
 * 业务路由统一走access中间件
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\App\V1',
    'prefix'     => 'app/v1',
    'middleware' => 'access'
], function ($router) {
    require __DIR__ . '/../routes/v1/app.php';
});

/**
 * 引入v1版本服务路由
 * 服务路由统一走access、sign中间件
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\Service\V1',
    'prefix'     => 'service/v1',
    'middleware' => ['access', 'sign']
], function ($router) {
    require __DIR__ . '/../routes/v1/service.php';
});

/**
 * 引入brand版本服务路由
 * 服务路由统一走access、sign中间件
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\Service\Brand',
    'prefix'     => 'service/brand',
    'middleware' => ['access', 'sign']
], function ($router) {
    require __DIR__ . '/../routes/brand/ys.php';
});

/**
 * 业务端接口回调
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\Callback\V1',
    'prefix'     => 'callback/v1',
    'middleware' => 'access'
], function ($router) {
    require __DIR__ . '/../routes/v1/callback.php';
});

/**
 * html接口
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\Html\V1',
    'prefix'     => 'html/v1',
    'middleware' => ['access']
], function ($router) {
    require __DIR__ . '/../routes/v1/html.php';
});

/**
 * 对外接口
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\Open\V1',
    'prefix'     => 'open/v1',
    'middleware' => ['access', 'auth_open']
], function ($router) {
    require __DIR__ . '/../routes/v1/open.php';
});

/**
 * crm对接接口
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\Service\Crm',
    'prefix'     => 'service/crm',
    'middleware' => ['access', 'sign']
], function ($router) {
    require __DIR__ . '/../routes/crm/crm.php';
});


/**
 * ma对接接口
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\Service\Ma',
    'prefix'     => 'service/ma',
    'middleware' => []
], function ($router) {
    require __DIR__ . '/../routes/ma/ma_plan.php';
});

/**
 * cdp对接接口
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\Service\Cdp',
    'prefix'     => 'service/cdp',
    'middleware' => ['access', 'sign']
], function ($router) {
    require __DIR__ . '/../routes/cdp/cdp.php';
});


/**
 * ec-hub对接接口
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\Service\EcHub',
    'prefix'     => 'service/ec-hub',
    'middleware' => ['access', 'sign']
], function ($router) {
    require __DIR__ . '/../routes/cdp/ec_hub.php';
});



return $app;
