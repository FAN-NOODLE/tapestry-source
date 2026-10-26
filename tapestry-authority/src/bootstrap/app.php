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

$app->configure('easy-sms');

$app->configure('wechat');

$app->configure('third_party');

$app->configure('message_strategy');


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

$app->singleton('filesystem', function ($app) {
    return $app->loadComponent('filesystems', 'Illuminate\Filesystem\FilesystemServiceProvider', 'filesystem');
});

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
    'access' => App\Http\Middleware\Access::class,
    'checkMessageStrategy' => App\Http\Middleware\CheckMessageStrategy::class
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
$app->register(\Illuminate\Redis\RedisServiceProvider::class);
// $app->register(App\Providers\AuthServiceProvider::class);
$app->register(App\Providers\EventServiceProvider::class);
$app->register(\Maatwebsite\Excel\ExcelServiceProvider::class);
$app->register(\Overtrue\LaravelWeChat\ServiceProvider::class);

$app->register(Appzcoder\LumenRoutesList\RoutesCommandServiceProvider::class);
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
    'prefix'     => 'micro-authority',
    'middleware' => ['access']
], function ($router) {
    require __DIR__ . '/../routes/web.php';
});


/**
 * 引入v1版本网关路由
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\V1',
    'prefix'     => 'micro-authority/v1',
    'middleware' => ['access', 'sign']
], function ($router) {
    require __DIR__ . '/../routes/v1.php';
});

/**
 * 引入brand版本网关路由
 */
$app->router->group([
    'namespace'  => 'App\Http\Controllers\Brand',
    'prefix'     => 'micro-authority/brand',
    'middleware' => ['access', 'sign']
], function ($router) {
    require __DIR__ . '/../routes/brand/ys.php';
});

return $app;
