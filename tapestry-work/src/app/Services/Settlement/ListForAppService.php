<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/19 0019
 * Time: 11:03
 */

namespace App\Services\Settlement;


use App\Models\DistributionGoods;
use App\Models\DistributionOrderGoodsSettlement;
use App\Models\PartnerEnterpriseSettings;
use App\Services\BaseService;
use App\Models\DistributionOrderSettlement;
use App\Services\DistributionOrderSettlementService;
use Illuminate\Database\Eloquent\Builder;

class ListForAppService extends BaseService
{
    private $page;

    private $size;

    private $workId;

    private $enterpriseId;

    private $goodsService;

    public function __construct($param)
    {
        $this->page = $param['page'];
        $this->size = $param['size'];
        $this->workId = $param['work_id'];
        $this->enterpriseId = $param['enterprise_id'];
        $this->goodsService = new GoodsService($param['enterprise_id']);
        $this->goodsService->setWorkId($this->workId);
    }

    public function list()
    {

        $query = DistributionOrderSettlement::query()
            ->leftJoin('partner_member','partner_member.id','=','distribution_order_settlement.partner_id')
            ->join('partner_enterprise_map',function($join){
                $join->on('partner_enterprise_map.partner_id','=','partner_member.id')
                    ->where('partner_enterprise_map.enterprise_id','=',$this->enterpriseId);
            }, null,null,'left')
            ->where('distribution_order_settlement.enterprise_id', $this->enterpriseId)
            ->where(function (Builder $query) {
                $query->where('distribution_order_settlement.promotion_worker_id', $this->workId)
                    ->orWhere('distribution_order_settlement.service_shopping_guid_id', $this->workId)
                    ->orWhere('distribution_order_settlement.partner_id', $this->workId);
                    })

            ->orderBy('distribution_order_settlement.created_at','desc');
        $pageData =  $query->paginate($this->size, $this->getField(), '', $this->page);
        //兼容部分退款后，佣金为零
        foreach ($pageData as $key=>$value){
            if($value['settlement_status'] == DistributionOrderSettlement::CANCEL_SETTLE){
                $value['promotion_commission'] = $value['promotion_commission_snapshot'];
                $value['service_commission'] = $value['service_commission_snapshot'];
                $value['partner_commission'] = $value['partner_commission_snapshot'];
                $pageData[$key]=$value;
            }
        }
        return $this->formData($pageData);
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
            'distribution_order_settlement.service_shopping_guid_id',
            'distribution_order_settlement.updated_at',
            'distribution_order_settlement.partner_id',
            'distribution_order_settlement.partner_commission',
            'partner_member.partner_name',
            'partner_member.partner_mobile',
            'partner_enterprise_map.mobile_show'
        ];
    }

    public function formData($pageData)
    {
        return [
            'list'       => $this->compat($pageData->items()),
            'total'      => $pageData->total() ?? 0,
            'page'       => (int)$this->page,
            'size'       => (int)$this->size
        ];
    }

    public function compat(array $settlementList)
    {
        $formatSettlement = [];
        $goodsAppService = new GoodsAppService($this->enterpriseId,$this->workId,$settlementList);
        foreach ($settlementList as $settlement) {
            $settlement = $settlement->toArray();
            $goodsAppService->setSettlement($settlement);
            $this->goodsService->setSettlement($settlement);

            $formatSettlement[] = $this->goodsService->formatSingle($goodsAppService);
        }
        return $formatSettlement;
    }

}
