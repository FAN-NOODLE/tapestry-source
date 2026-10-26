<?php

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Enterprise;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\Request;

class MerchantSettingDetailController extends BaseController
{
    /**
     * MerchantSettingDetailController constructor.
     * @throws \Exception
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 获取商户信息
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id']);

        // 企业ID不能为空
        if (!isset($params['enterprise_id']) || empty($params['enterprise_id'])) {
            return error(201101);
        }

        // 查询企业是否存在
        $enterprise = Enterprise::query()->where('id', $params['enterprise_id'])->first();

        if (!$enterprise) {
            return error(700003);
        }

        $data = [
            'enterprise_name' => $enterprise->name,
            'logo'            => '',
            'short_name'      => '',
        ];

        // 查询企业logo及企业简称
        $enterpriseThirdPartAccount = EnterpriseThirdPartAccount::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('platform', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
            ->where('type', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->first();

        if ($enterpriseThirdPartAccount) {
            $data = [
                'enterprise_name' => $enterprise->name,
                'logo'            => $enterpriseThirdPartAccount->logo,
                'short_name'      => $enterpriseThirdPartAccount->short_name,
            ];
        }

        return success(['data' => $data]);
    }
}
