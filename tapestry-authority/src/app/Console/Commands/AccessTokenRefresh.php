<?php

namespace App\Console\Commands;

use App\Models\Enterprise;
use App\Models\EnterpriseThirdPartAccount;
use App\WxServices\AccessTokenWxService;
use Illuminate\Console\Command;

class AccessTokenRefresh extends Command
{
    protected $signature = 'refresh_access_token';

    protected $description = '刷新企业access_token';

    // 开放平台
    const TYPE_OPEN = 1;
    // 小程序
    const TYPE_MINI = 2;

    public function handle()
    {
        $enterpriseIDs = EnterpriseThirdPartAccount::query()
            ->where('status', '=', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->where('platform', '=', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
            ->where('type', '=', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->pluck('enterprise_id')
            ->toArray();

        $tokenService = new AccessTokenWxService();

        foreach ($enterpriseIDs as $enterprise) {
            $tokenService->getAccessToken($enterprise, self::TYPE_OPEN, true);
            $tokenService->getAccessToken($enterprise, self::TYPE_MINI, true);
        }
    }
}
