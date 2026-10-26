<?php

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Common\RedisKey;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Enterprise;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\Request;

class MerchantSettingSaveController extends BaseController
{
    /**
     * MerchantSettingSaveController constructor.
     * @throws \Exception
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 保存商户信息
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id', 'logo', 'short_name']);

        // 企业ID不能为空
        if (!isset($params['enterprise_id']) || empty($params['enterprise_id'])) {
            return error(201101);
        }

        // 企业简称不能为空
        if (!isset($params['short_name']) || empty($params['short_name'])) {
            return error(700001);
        }

        // 查询企业是否存在
        $enterprise = Enterprise::query()->where('id', $params['enterprise_id'])->first();

        if (!$enterprise) {
            return error(700004);
        }

        // 查询企业logo及企业简称
        $enterpriseThirdPartAccount = EnterpriseThirdPartAccount::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('platform', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
            ->where('type', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->first();

        if ($enterpriseThirdPartAccount) {
            $updateData = [
                'short_name' => $params['short_name'],
            ];

            if (isset($params['logo']) && !empty($params['logo'])) {
                $updateData['logo'] = $params['logo'];
            }

            EnterpriseThirdPartAccount::query()
                ->where('enterprise_id', $params['enterprise_id'])
                ->where('platform', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
                ->where('type', EnterpriseThirdPartAccount::TYPE_SHOP)
                ->update($updateData);
            delRedis([RedisKey::MINI_PROGRAM_AUTH_INFO_BY_ENTERPRISE_ID . $params['enterprise_id']]);
        } else {
            return error(700005);
        }

        return success(['data' => []]);
    }
}
