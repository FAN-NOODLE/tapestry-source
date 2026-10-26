<?php

require_once __DIR__ . '/../vendor/autoload.php';

try {
    (new Dotenv\Dotenv(__DIR__ . '/../'))->load();
} catch (Dotenv\Exception\InvalidPathException $e) {
    //
}

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
$app->configure('queue');

$app->configure('wechat');

$app->configure('custom');

$app->configure('third_party');

$app->configure('mail');    //加载email配置文件

$app->configure('frontend');

$app->configure('work_plugin');

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

// $app->middleware([
//     App\Http\Middleware\ExampleMiddleware::class
// ]);

$app->routeMiddleware([
    'auth'   => App\Http\Middleware\Authenticate::class,
    'sign'   => App\Http\Middleware\SignMiddleware::class,
    'access' => App\Http\Middleware\Access::class
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
$app->register(App\Providers\EventServiceProvider::class);
$app->register(Overtrue\LaravelWeChat\ServiceProvider::class);
$app->register(Illuminate\Mail\MailServiceProvider::class);
$app->register(Illuminate\Redis\RedisServiceProvider::class);
$app->register(App\Providers\ThirdPartServiceProvider::class);
$app->register(Intervention\Image\ImageServiceProvider::class);

$app->withFacades(true, [
    'Image' => Intervention\Image\Facades\Image::class,
]);

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
    'namespace' => 'App\Http\Controllers',
    'prefix'    => 'boss',
], function ($router) {
    require __DIR__ . '/../routes/web.php';
});

/**
 * 引入v1版本业务路由
 * 业务路由统一走access中间件
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\App\V1',
    'prefix'     => 'boss/app/v1',
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
    'prefix'     => 'boss/service/v1',
    'middleware' => ['access', 'sign']
], function ($router) {
    require __DIR__ . '/../routes/v1/service.php';
});

/**
 * 引入v2版本服务路由
 * 服务路由统一走access、sign中间件
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\Service\V2',
    'prefix'     => 'boss/service/v2',
    'middleware' => ['access', 'sign']
], function ($router) {
    require __DIR__ . '/../routes/v2/service.php';
});

return $app;
