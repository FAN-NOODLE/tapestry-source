<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/20 0020
 * Time: 10:10
 */

namespace App\Services\Settlement;


use App\Models\DistributionOrderGoodsSettlement;
use CrsSdk\MicroGoodsV1\Goods;

class GoodsAppService
{
    private $enterpriseId;

    private $settlementList=[];

    private $settlementGoods=[];

    private $orderNos=[];

    private $orderNo;

    private $orderSkuIdInfo;

    private $skuInfoMap = [];

    private $settlement;

    private $workId;

    public function __construct(string $enterpriseId, string $workId,array $settlementList=[])
    {
        $this->enterpriseId = $enterpriseId;
        $this->settlementList = $settlementList;
        $this->workId = $workId;
        foreach ($settlementList as $settlement){
            $settlement = $settlement->toArray();
            $this->orderNos[] = $settlement['order_no'];
        }

        $this->setSkuInfoMap();
    }

    public function setSettlement(array $settlement)
    {
        $this->settlement = $settlement;
        $this->orderNo = $this->settlement['order_no'];
        $this->orderSkuIdInfo = $this->setOrderSkuIds();
    }

    public function getSkuIdInfo()
    {
        return $this->orderSkuIdInfo;
    }

    public function setSkuInfoMap()
    {
        $this->settlementGoods = [];

        $settlement_goods = DistributionOrderGoodsSettlement::query()
            ->whereIn('order_no', $this->orderNos)
            ->where('enterprise_id', '=', $this->enterpriseId)
            ->get(['sku_id', 'partner_commission', 'order_goods_price', 'promotion_commission', 'service_commission',
                'purchase_quantity', 'order_no','refund_status','refund_quantity'])
            ->toArray();

        $this->settlementGoods = $settlement_goods;
        $goods_param['enterprise_id'] = $this->enterpriseId;
        $goods_param['need_attribute'] = 1;
        $uniqueSkuIds = array_unique(array_column($settlement_goods, 'sku_id'));
        $goods_param['search']['sku_id_list'] = $uniqueSkuIds;
        $goods_param['size'] = count($uniqueSkuIds)==0 ? 1 : count($uniqueSkuIds);
        $result = (new Goods())->list($goods_param);
        if(!empty($result['total']) && $result['total'] > $goods_param['size']){
            $goods_param['size'] = $result['total'];
            $result = (new Goods())->list($goods_param);
        }
        if($result['list']){
            foreach ($result['list'] as $skuInfo){
                $this->skuInfoMap[$skuInfo['sku_id']] = $skuInfo;
            }
        }

    }

    public function setOrderSkuIds()
    {
        $settlement_goods = [];
        foreach ($this->settlementGoods as $aSettlementGoods){

            if($aSettlementGoods['order_no'] == $this->orderNo){
                $settlement_goods[] = $aSettlementGoods;
            }
        }

        foreach ($settlement_goods as &$sgoods) {
            $skuId = $sgoods['sku_id'];
            $sgoods['material_no'] = $this->skuInfoMap[$skuId]['material_no'] ?? "";
            $sgoods['spu_img'] = $this->skuInfoMap[$skuId]['sku_img'] ?? "";
            $sgoods['spu_name'] = $this->skuInfoMap[$skuId]['sku_name'] ?? "";
            $attribute = $this->skuInfoMap[$skuId]['attribute']??[];
            $sgoods['attribute'] = (!empty($attribute) && isset($attribute)) ? array_column($attribute, 'item_value') : [];
            //福分退款
            if($sgoods['refund_status'] == DistributionOrderGoodsSettlement::REFUND_STATUS_PART){
                $purchase_quantity = $sgoods['purchase_quantity'] - $sgoods['refund_quantity'];
                $sgoods['promotion_commission'] =  (int)floor(($sgoods['promotion_commission']/$sgoods['purchase_quantity'])*$purchase_quantity);
                $sgoods['service_commission'] =    (int)floor(($sgoods['service_commission']/$sgoods['purchase_quantity'])*$purchase_quantity);
                $sgoods['partner_commission'] =    (int)floor(($sgoods['partner_commission']/$sgoods['purchase_quantity'])*$purchase_quantity);
                $sgoods['purchase_quantity']  =   $purchase_quantity;
            }
            //全部退款
            if($sgoods['refund_status'] == DistributionOrderGoodsSettlement::REFUND_STATUS_ALL){
                $sgoods['spu_name'] =$sgoods['spu_name'].'(全部退款)';
            }
            //unset($sgoods['sku_id']);
        }
        $refundGoodsList = [];
        //部分退款商品展示
        foreach ($settlement_goods as &$refund_sgoods) {
            if($refund_sgoods['refund_status'] == DistributionOrderGoodsSettlement::REFUND_STATUS_PART){
                $attribute = $this->skuInfoMap[$skuId]['attribute']??[];
                $spu_name = $this->skuInfoMap[$skuId]['sku_name'] ?? "";
                $service_commission = $this->settlement['service_shopping_guid_id'] ? $refund_sgoods['service_commission'] : 0;
                $refundGoodsList[] = array(
                    'material_no' => $this->skuInfoMap[$skuId]['material_no'] ?? [],
                    'spu_img'     => $this->skuInfoMap[$skuId]['sku_img'] ?? "",
                    'spu_name'    => $spu_name.'(部分退款)',
                    'attribute'   => (!empty($attribute) && isset($attribute)) ? array_column($attribute, 'item_value') : [],
                    'service_commission' => (int)floor($service_commission/$refund_sgoods['purchase_quantity']*$refund_sgoods['refund_quantity']),
                    'promotion_commission' => (int)floor($refund_sgoods['promotion_commission']/$refund_sgoods['purchase_quantity']* $refund_sgoods['refund_quantity']),
                    'partner_commission'  => (int)floor($refund_sgoods['partner_commission']/$refund_sgoods['purchase_quantity']* $refund_sgoods['refund_quantity']),
                    'order_goods_price'  => $refund_sgoods['order_goods_price'],
                    'purchase_quantity' => $refund_sgoods['refund_quantity'],
                    'sku_id'=> $refund_sgoods['sku_id'],
                    'refund_status'=>$refund_sgoods['refund_status'],
                    'refund_quantity'=>$refund_sgoods['refund_quantity']
                );
            }
        }
        unset($sgoods);
        $settlement_goods = array_merge($settlement_goods,$refundGoodsList);
        $settlement_goods = $this->sortByCommission($settlement_goods);
        $commission = array_column($settlement_goods, 'commission');
        if (!empty($commission)) {
            // 根据佣金金额列降序排列
            array_multisort($commission, SORT_DESC, $settlement_goods);
        }
        return $settlement_goods;
    }

    private function sortByCommission($settlement_goods)
    {

        foreach ($settlement_goods as &$sku){
            $commission = 0;
            if($this->workId == $this->settlement['promotion_worker_id']){
                $commission += intval($sku['promotion_commission']) * $sku['purchase_quantity'];
            }

            if($this->workId == $this->settlement['service_shopping_guid_id']){
                $commission += intval($sku['service_commission']) * $sku['purchase_quantity'];
            }

            if($this->workId == $this->settlement['partner_id']){
                $commission += intval($sku['partner_commission']) * $sku['purchase_quantity'];
            }
            $sku['commission'] = $commission;
        }
       // array_multisort(array_column($settlement_goods,'commission'), SORT_DESC, $settlement_goods);
        return $settlement_goods;
    }



}