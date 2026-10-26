<?php

namespace App\Listeners;

use App\Events\MiniProgram;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Contracts\Queue\ShouldQueue;

class MiniProgramFullProcedure implements ShouldQueue
{
    public $miniProgramService;

    /**
     * Create the event listener.
     *
     * @return void
     */
    public function __construct(MiniProgramService $miniProgramService)
    {
        $this->miniProgramService = $miniProgramService;
    }

    /**
     * Handle the event.
     *
     * @param MiniProgram $event
     *
     * @return bool
     * @throws \App\Services\OpenPlatform\ErrorMiniProgramCodeDraftsException
     */
    public function handle(MiniProgram $event)
    {
        // 防止请求执行超时
        set_time_limit(0);

        // 异常处理
        try {
            $this->miniProgramService->fullProcedure($event->versionType, $event->checked);
        } catch (\Exception $exception) {
            return true;
        };
    }
}
