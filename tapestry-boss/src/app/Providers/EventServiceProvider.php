<?php

namespace App\Providers;

use Laravel\Lumen\Providers\EventServiceProvider as ServiceProvider;

class EventServiceProvider extends ServiceProvider
{
    /**
     * The event listener mappings for the application.
     *
     * @var array
     */
    protected $listen = [
        'App\Events\SomeEvent'                                        => [
            'App\Listeners\EventListener',
        ],

        // 微信开放平台授权
        'Overtrue\LaravelWeChat\Events\OpenPlatform\Authorized'       => [
            'App\Listeners\OpenPlatform\Authorized',
        ],
        // 微信开放平台取消授权
        'Overtrue\LaravelWeChat\Events\OpenPlatform\Unauthorized'     => [
            'App\Listeners\OpenPlatform\Unauthorized',
        ],
        // 微信开放平台更新授权
        'Overtrue\LaravelWeChat\Events\OpenPlatform\UpdateAuthorized' => [
            'App\Listeners\OpenPlatform\UpdateAuthorized',
        ],
        // 小程序授权完成
        'App\Events\MiniProgramAuthed'                                => [
            //'App\Listeners\InitPortalData',
            //'App\Listeners\InitAiMaintenanceData',
            'App\Listeners\CommitMiniProgram',  // 提审小程序
            //'App\Listeners\MiniShopSet', // 初始化小程序数据
        ],
        // 小程序审核成功事件
        'App\Events\WeappAudit\SuccessEvent'                          => [
            'App\Listeners\WeappAudit\SuccessListener',
        ],
        // 小程序审核失败事件
        'App\Events\WeappAudit\FailureEvent'                          => [
            'App\Listeners\WeappAudit\FailureListener',
        ],
        // 企业服务号信息初始化
        'App\Events\EnterpriseThirdPartAccountWasInit'                => [
            // 初始化微信菜单
            'App\Listeners\InitWechatMenu',
        ],
        // 企业授权成功
        'App\Events\EnterpriseAuthWorkSucceed'                        => [
            // 注册企业、企业号、应用
            'App\Listeners\RegisterEnterpriseAndWorkAgent',
        ],
        // 小程序
        'App\Events\MiniProgram'                                      => [
            // 全网发布
            'App\Listeners\MiniProgramFullProcedure',
        ],
    ];
}
