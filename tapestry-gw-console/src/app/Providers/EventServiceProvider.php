<?php

namespace App\Providers;

use App\Events\AddTaskEvent;
use App\Listeners\AddTaskListener;
use Laravel\Lumen\Providers\EventServiceProvider as ServiceProvider;

class EventServiceProvider extends ServiceProvider
{
    /**
     * The event listener mappings for the application.
     *
     * @var array
     */
    protected $listen = [
        'App\Events\SomeEvent' => [
            'App\Listeners\EventListener',
        ],
        // 任务添加
        'App\Events\AddTaskEvent' => [
            'App\Listeners\AddTaskListener',
        ],
    ];
}
