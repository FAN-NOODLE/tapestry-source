<?php

namespace App\Micros;

use App\Models\EnterpriseThirdPartAccount;
use CrsSdk\Boss\MiniProgram;
use Illuminate\Support\Facades\Log;

class MiniProgramMicro
{
    public function getMiniProgramName($enterpriseId)
    {
        $enterpriseDetail = [];
        try {
            $enterpriseDetail = (new MiniProgram())->detail(['enterprise_id' => $enterpriseId])['data'];

        } catch (\Exception $e) {
            Log::info('获取小程序详情');
        }

        return $enterpriseDetail['name'] ?? '';
    }

    public function listMiniPrograms($enterpriseIds)
    {

        return (new EnterpriseThirdPartAccount())->listMiniPrograms($enterpriseIds);
    }

    public function getMiniProgram($enterpriseId)
    {

        return (new EnterpriseThirdPartAccount())->getMiniProgram($enterpriseId);
    }
}