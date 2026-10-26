<?php

namespace App\Listeners\WeappAudit;

use App\Events\WeappAudit\SuccessEvent;
use App\Models\EnterpriseMiniProgramAuditLog;
use App\Models\EnterpriseMiniProgramCodeLibrary;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Support\Facades\Log;

class SuccessListener
{
    /**
     * @var MiniProgramService
     */
    protected $miniProgramService;

    /**
     * Create the event listener.
     *
     * @param MiniProgramService $miniProgramService
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
     * @param SuccessEvent $event
     *
     * @return bool
     * @throws \App\Exceptions\Error
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function handle(SuccessEvent $event)
    {
        // 小程序审核成功
        Log::info('小程序审核回调进来了--成功');
        Log::info('接收数据' . var_export($event->payload['app_id']));

        // 获取代码审核id
        $auditInfo = $this->miniProgramService->recentAudit($event->payload['app_id']);

        $log = EnterpriseMiniProgramAuditLog::query()
            ->where('audit_id', $auditInfo['wx_audit_id'])
            ->first();

        // 更新代码审核记录状态/对应代码库审核状态
        $this->miniProgramService->activeRecentAuditByWxAuditId($auditInfo['wx_audit_id'], $log);

        // 小程序审核记录
        $codeLibrary = EnterpriseMiniProgramCodeLibrary::query()
            ->where('id', $log->code_library_id)
            ->first();

        // 上线方式：手动上线
        if (EnterpriseMiniProgramCodeLibrary::ONLINE_TYPE_MANUALLY == $codeLibrary->online_type) {
            return true;
        }

        // 发布小程序代码 (如果不在不自动发布app_id列表中，则进行发布)
        $notAutoReleaseAppList = env('NOT_AUTO_RELEASE_APP_LIST', '');
        $notAutoReleaseAppList = explode(',', $notAutoReleaseAppList);
        $appId = $event->payload['app_id'];
        if (!in_array($appId, $notAutoReleaseAppList)) {
            $this->miniProgramService->release($appId, $auditInfo['wx_audit_id']);
        }

        return true;
    }
}
