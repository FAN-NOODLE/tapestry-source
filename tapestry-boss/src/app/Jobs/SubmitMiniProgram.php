<?php

namespace App\Jobs;

use App\Services\EnterpriseLogService;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Support\Facades\Log;

class SubmitMiniProgram extends Job
{
    protected $checked;

    protected $thirdPartAccounts;

    /**
     * SubmitMiniProgram constructor.
     *
     * @param $thirdPartAccounts
     * @param $checked
     */
    public function __construct($thirdPartAccounts, $checked)
    {
        $this->checked = $checked;
        $this->thirdPartAccounts = $thirdPartAccounts;
    }

    /**
     * @param MiniProgramService $miniProgramService
     */
    public function handle(MiniProgramService $miniProgramService)
    {
        $checked = $this->checked;
        $thirdPartAccounts = $this->thirdPartAccounts;

        // 记录未发布的企业和错误
        $cannot = [];
        foreach ($thirdPartAccounts as $thirdPartAccount) {
            try {
                $miniProgramService->submitVersion($thirdPartAccount, $checked);
            } catch (\Exception $e) {
                $cannot[] = [
                    'enterprise_id' => $thirdPartAccount->enterprise_id,
                    'account_id'    => $thirdPartAccount->id,
                    'error'         => [
                        'message' => $e->getMessage(),
                        'line'    => $e->getLine(),
                        'file'    => $e->getFile(),
                    ],
                ];
                EnterpriseLogService::error('全网发布代码失败', $cannot);
                Log::info('错误数据:' . var_export($cannot, true));
            }
        }
    }
}
