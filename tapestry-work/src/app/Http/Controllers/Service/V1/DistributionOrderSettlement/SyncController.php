<?php

namespace App\Http\Controllers\Service\V1\DistributionOrderSettlement;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\DistributionOrderSettlementService;

class SyncController extends BaseController
{
    /**
     * 设置分销订单结算和订单商品结算
     *
     * @param array 结算数据结构 $data
     *
     * @return mixed
     */
    public function index()
    {
        $orderSettlement                    = $this->request->input('order_settlement');
        $distributionOrderSettlementService = new DistributionOrderSettlementService();
        $distributionOrderSettlementService->sync($orderSettlement);
        return success();
    }
}
