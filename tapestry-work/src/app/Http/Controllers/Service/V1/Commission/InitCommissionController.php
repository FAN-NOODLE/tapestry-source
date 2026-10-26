<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/10
 * Time: 17:56
 */

namespace App\Http\Controllers\Service\V1\Commission;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\CommissionSetting;
use Illuminate\Http\Request;

class InitCommissionController extends BaseController
{
    /**
     * 企业导购佣金设置初始化
     * （boss平台创建企业成功后调用）
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        // 企业id
        $enterpriseId = $request->get('enterprise_id', 0);

        if ($enterpriseId) {
            $data = [
                [
                    'enterprise_id' => $enterpriseId,
                    'type'          => 1,
                    'value'         => 8,
                    'created_at'    => time()
                ],
                [
                    'enterprise_id' => $enterpriseId,
                    'type'          => 2,
                    'value'         => 2,
                    'created_at'    => time()
                ],
                [
                    'enterprise_id' => $enterpriseId,
                    'type'          => 3,
                    'value'         => 4,
                    'created_at'    => time()
                ]
            ];

            CommissionSetting::insert($data);
        }

        return success(['data' => []]);
    }
}
