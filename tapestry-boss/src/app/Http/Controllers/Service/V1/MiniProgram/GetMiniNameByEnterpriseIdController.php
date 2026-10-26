<?php


namespace App\Http\Controllers\Service\V1\MiniProgram;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseThirdPartAccount;

class GetMiniNameByEnterpriseIdController extends BaseController
{
    public function index()
    {
        $enterpriseIds = $this->request->get('enterprise_ids', []);
        if (empty($enterpriseIds)) {
            return success([]);
        }
        if (!is_array($enterpriseIds)) {
            return success([]);
        }
        $list = EnterpriseThirdPartAccount::query()
            ->where([
                'authorize_type' => EnterpriseThirdPartAccount::AUTHORIZE_TYPE_THIRD, // 第三方授权
                'status'         => EnterpriseThirdPartAccount::STATUS_ACTIVE, // 有效状态
                'platform'       => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT,// 小程序
                'type'           => EnterpriseThirdPartAccount::TYPE_SHOP // 商城小程序
            ])
            ->whereIn('enterprise_id', $enterpriseIds)
            ->get(['enterprise_id', 'name'])
            ->toArray();
        if (empty($list)) {
            return success([]);
        }
        return success($list);
    }
}