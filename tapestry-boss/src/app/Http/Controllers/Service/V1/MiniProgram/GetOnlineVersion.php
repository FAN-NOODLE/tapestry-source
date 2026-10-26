<?php


namespace App\Http\Controllers\Service\V1\MiniProgram;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseMiniProgramCodeLibrary;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\Request;

class GetOnlineVersion extends BaseController
{

    public function index(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id');

        // 获取小程序最新版本
        $authInfo = EnterpriseThirdPartAccount::query()
            ->where([
                'enterprise_id' => $enterpriseId,
                'platform'      => 2,
                'type'          => 1
            ])
            ->select(['id', 'release_branch', 'release_tag'])
            ->first()
            ->toArray();

        // 获取小程序线上使用版本
        $onlineVersion = EnterpriseMiniProgramCodeLibrary::query()
            ->where([
                'online_status' => 'ONLINE',
                'using'         => 1,
                'account_id'    => $authInfo['id'],
            ])
            ->orderByDesc('id')
            ->select(['release_branch', 'release_tag'])
            ->first()
            ->toArray();

        $data = [
            'use_tag'    => $onlineVersion['release_tag'],
            'use_branch' => $onlineVersion['release_branch']
        ];

        return success($data);
    }
}
