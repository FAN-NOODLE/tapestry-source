<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-01-28
 * Time: 18:19
 */

namespace App\Console\Commands;

use App\Models\EnterpriseMiniProgramCodeLibrary;
use App\Services\EnterpriseLogService;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class UpdateAuditStatus extends Command
{
    /**
     * @var string
     */
    protected $signature = 'auto:update_audit_status';

    /**
     * @var string
     */
    protected $description = '定时更新小程序审核状态';

    public function handle(MiniProgramService $miniProgramService)
    {
        // 查询审核中状态
        $codeLibraryList = EnterpriseMiniProgramCodeLibrary::query()->where('audit_status', 'AUDITING')
            ->with('auditLog')->get()->toArray();
        // 如果没有需要审核的
        if (empty($codeLibraryList)) {
            return true;
        }

        EnterpriseLogService::info('定时更新小程序审核状态-审核中的小程序', [
            'list' => $codeLibraryList
        ]);

        // 循环列表
        foreach ($codeLibraryList as $value) {
            if (empty($value['extra_info'])) {
                continue;
            }
            // json解析
            $extraInfo = json_decode($value['extra_info'], true);
            $appId     = $extraInfo['extAppid'] ?? '';
            // 获取审核ID
            $auditId = $value['audit_log']['audit_id'] ?? '';

            if (!empty($auditId) && !empty($appId)) {
                Log::info('开始查询审核状态 app_id:' . $appId . ' audit_id:' . $auditId);
                try {
                    $data = $miniProgramService->getAuditInfoByWXAuditId($appId, $auditId);
                } catch (\Exception $exception) {
                    Log::info('app_id:' . $appId . ' audit_id:' . $auditId . '查询异常原因:' . $exception->getMessage());

                    EnterpriseLogService::error('定时更新小程序审核状态失败', [
                        'app_id'   => $appId,
                        'audit_id' => $auditId,
                        'message'  => $exception->getMessage(),
                        'file'     => $exception->getFile(),
                        'line'     => $exception->getLine()
                    ]);

                    continue;
                }
                Log::info('app_id:' . $appId . ' audit_id:' . $auditId . '已查询审核状态:' . $data['status']);
            }
        }

        return true;
    }
}
