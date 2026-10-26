<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/22 22:45
 */

namespace App\Http\Controllers\Service\V1\DistributionGoods;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\DistributionGoodsService;

class AppSkuListController extends BaseController
{
    public function index(DistributionGoodsService $distributionGoodsService)
    {
        $params = request()->only(['enterprise_id', 'new_goods']);

        // 企业id不能为空
        if (empty($params['enterprise_id'])) {
            notice(600000);
        }

        $data           = $distributionGoodsService->appGoodsSkuList($params);

        return success($data);
    }
}