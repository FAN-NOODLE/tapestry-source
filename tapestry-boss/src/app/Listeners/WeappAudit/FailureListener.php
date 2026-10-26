<?php
/*
 * This file is part of the car/chedianai_bc.
 *
 * (c) chedianai <i@chedianai.com>
 *
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Listeners\WeappAudit;

use App\Events\WeappAudit\FailureEvent;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Support\Facades\Log;

class FailureListener
{
    /**
     * @var MiniProgramService
     */
    protected $miniProgramService;

    /**
     * FailureListener constructor.
     *
     * @param MiniProgramService $miniProgramService
     */
    public function __construct(MiniProgramService $miniProgramService)
    {
        $this->miniProgramService = $miniProgramService;
    }

    /**
     * Handle the event.
     *
     * @param FailureEvent $event
     *
     * @throws \App\Exceptions\Error
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function handle(FailureEvent $event)
    {
        // 小程序审核失败
        Log::info('小程序审核回调进来了--失败');

        Log::info('接收数据' . var_export($event->payload['app_id']));

        // 获取代码审核id
        $auditInfo = $this->miniProgramService->recentAudit($event->payload['app_id']);

        // 更新代码审核记录状态/对应代码库审核状态
        $this->miniProgramService->negativeAuditByWxAuditId($auditInfo['wx_audit_id'], $auditInfo['reason']);
    }
}
