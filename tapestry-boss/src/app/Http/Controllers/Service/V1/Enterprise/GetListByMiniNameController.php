<?php


namespace App\Http\Controllers\Service\V1\Enterprise;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Enterprise;
use App\Models\EnterpriseThirdPartAccount;

class GetListByMiniNameController extends BaseController
{
    public function index()
    {
        $miniName = $this->request->input('name', '');
        $ids = EnterpriseThirdPartAccount::query()
            ->where([
                'authorize_type' => EnterpriseThirdPartAccount::AUTHORIZE_TYPE_THIRD, // 第三方授权
                'status'         => EnterpriseThirdPartAccount::STATUS_ACTIVE, // 有效状态
                'platform'       => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT,// 小程序
                'type'           => EnterpriseThirdPartAccount::TYPE_SHOP // 商城小程序
            ])
            ->where('name', 'like', '%' . $miniName . '%')
            ->pluck('enterprise_id')
            ->toArray();
        $list = [];
        if (!empty($ids)) {
            $list = Enterprise::query()
                ->whereIn('id', $ids)
                ->get()
                ->toArray();
        }

        return success($list);
    }
}