<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/16 0016
 * Time: 15:50
 */

namespace App\Services\Settlement;


use App\Models\DistributionOrderSettlement;
use App\Services\BaseService;

class DetailService extends BaseService
{
    private $enterpriseId;

    private $settleId;

    private $guideService;

    private  $goodsService;

    public function __construct($enterpriseId,$settleId,$version)
    {
        $this->enterpriseId = $enterpriseId;
        $this->settleId = $settleId;
        $this->guideService = new GuideService($enterpriseId,$version);
        $this->goodsService = new GoodsService($enterpriseId);
    }



    public function detail()
    {
        $detail = DistributionOrderSettlement::query()
            ->where('enterprise_id', $this->enterpriseId)
            ->where('id',$this->settleId)
            ->first($this->getField())
            ->toArray();
        $this->goodsService->setSettlement($detail);
        return [
            'detail'           => $this->compact($detail) ?? [],
            'settlement_goods' => $this->goodsService->getGoodsInfo()
        ];

    }

    private function compact($detail)
    {
        $this->guideService->setGuideIds(
            (array)$detail['promotion_worker_id'],
            (array)$detail['service_shopping_guid_id']
        );
        $this->guideService->setUser($detail);
        return  $this->guideService->formatSingle();
    }

    private function getField()
    {
        return [
            'distribution_order_settlement.id',
            'distribution_order_settlement.order_no',
            'distribution_order_settlement.order_status',
            'distribution_order_settlement.settlement_status',
            'distribution_order_settlement.order_pay_time',
            'distribution_order_settlement.promotion_worker_id',
            'distribution_order_settlement.promotion_commission',
            'distribution_order_settlement.service_shopping_guid_id',
            'distribution_order_settlement.service_commission',
            'distribution_order_settlement.service_shopping_guid_id'
        ];
    }

}