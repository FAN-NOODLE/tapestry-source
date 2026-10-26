<?php
/**
 * 全员分销-基础配置-查看
 * Created by PhpStorm.
 * User: weihua
 * Date: 2020-09-24
 * Time: 18:30
 */

namespace App\Http\Controllers\Service\V1\DistributionConfig;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Distribution;
use App\Services\DistributionService;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    public function index(Request $request)
    {
        $rules = [
            'enterprise_id'   => 'required|integer',
            'distribution_id' => 'integer'
        ];

        $messages = [
            'enterprise_id.*'   => '企业ID不正确',
            'distribution_id.*' => '分销配置ID不合法'
        ];

        $input                       = $this->validate($request, $rules, $messages);
        $distributionId              = $input['distribution_id'] ?? 0;

        // 查询分销配置详情
        $distributionService = new DistributionService();
        $detail = $distributionService->detail($input['enterprise_id'],$distributionId);
        if ($distributionId == 0) {
            $detail = [
                'id'                  => idGenerate(),
                'distribution_title'  => '什么是全员分销',
                'distribution_text'   => "导购人员的推广订单/服务订单内包含指定商品，则在订单完成时按'商品数量*单件商品佣金'算法结算分销佣金",
                'settlement_standard' => Distribution::SETTLEMENT_STANDARD_COMPLETE,
                'poster_img'          => 'https://crs-oss.vchangyi.com/store-manage-wxapp/default.png'
            ];
        }
        return success($detail);
    }
}
