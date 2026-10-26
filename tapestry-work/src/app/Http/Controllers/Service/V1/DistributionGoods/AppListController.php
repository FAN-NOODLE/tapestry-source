<?php

namespace App\Http\Controllers\Service\V1\DistributionGoods;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\DistributionGoodsService;
use Illuminate\Http\Request;

class AppListController extends BaseController
{
    /**
     * 小程序端分销商品列表
     *
     * @param Request $request
     * @param DistributionGoodsService $distributionGoodsService
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(DistributionGoodsService $distributionGoodsService)
    {
        $params = request()->only(['enterprise_id', 'new_goods']);

        // 企业id不能为空
        if (empty($params['enterprise_id'])) {
            notice(600000);
        }

        $data           = $distributionGoodsService->appGoodsList($params);

        return success($data);
    }
}
