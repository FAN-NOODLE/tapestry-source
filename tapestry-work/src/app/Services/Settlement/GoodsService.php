<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/19 0019
 * Time: 11:10
 */

namespace App\Services\Settlement;


use App\Models\CrsOrderGoods;
use App\Models\DistributionOrderGoodsSettlement;
use App\Models\OrderGoods;
use CrsSdk\MicroGoodsV1\Goods;

class GoodsService
{
    private $settlement;

    private $enterpriseId;

    private $orderGoods;

    private $workId;

    public function __construct(string $enterpriseId)
    {
        $this->enterpriseId = $enterpriseId;

    }

    public function setWorkId(string $workId)
    {
        $this->workId = $workId;
    }

    public function setSettlement(array $settlement)
    {
        $this->settlement = $settlement;
        $this->orderGoods = $this->getOrderGoods();
    }

    public function getGoodsInfo()
    {
        return $this->orderGoods;
    }

    private function getOrderGoods()
    {
        $settlement_goods = DistributionOrderGoodsSettlement::query()
            ->where('order_no', '=', $this->settlement['order_no'])
            ->where('enterprise_id', '=', $this->enterpriseId)
            ->get(['sku_id', 'order_goods_price', 'promotion_commission',
                'service_commission', 'purchase_quantity','partner_commission',
                'refund_status','refund_quantity'
            ])
            ->toArray();
        //根据sku_id查询商品信息
        $goods_param['enterprise_id'] = $this->enterpriseId;
        $goods_param['need_attribute'] = 1;
        $goods_param['search']['sku_id_list'] = array_column($settlement_goods, 'sku_id');
        $result = (new Goods())->list($goods_param);

        if($result['total'] > $result['size']){
            $goods_param['size'] = $result['total'];
            $result = (new Goods())->list($goods_param);
        }
        if (isset($result['list']) && !empty($result['list'])) {
            $goods_list = $result['list'];
            //重组数组 sku_id=>商品名称,商品图片,物料编码
            $goodsArr = [];
            foreach ($goods_list as $goods) {
                $goodsArr[$goods['sku_id']]['spu_name'] = $goods['spu_name'];
                $goodsArr[$goods['sku_id']]['spu_img'] = $goods['spu_img'];
                $goodsArr[$goods['sku_id']]['material_no'] = $goods['material_no'];
                $goodsArr[$goods['sku_id']]['attribute'] = $goods['attribute'];
            }
            foreach ($settlement_goods as &$sgoods) {

                $sgoods['material_no'] = $goodsArr[$sgoods['sku_id']]['material_no'];
                $sgoods['spu_img'] = $goodsArr[$sgoods['sku_id']]['spu_img'];
                $sgoods['spu_name'] = $goodsArr[$sgoods['sku_id']]['spu_name'];
                $attribute = $goodsArr[$sgoods['sku_id']]['attribute'];
                $sgoods['attribute'] = (!empty($attribute) && isset($attribute)) ? array_column($attribute, 'item_value') : [];
                $sgoods['service_commission'] = $this->settlement['service_shopping_guid_id'] ? $sgoods['service_commission'] : 0;
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
                    //$sgoods['spu_name'] = $goodsArr[$sgoods['sku_id']]['spu_name'].'(全部退款)';
                }
            }
            $refundGoodsList = [];
            //部分退款商品展示
            foreach ($settlement_goods as &$refund_sgoods) {
                if($refund_sgoods['refund_status'] == DistributionOrderGoodsSettlement::REFUND_STATUS_PART){
                    $attribute = $goodsArr[$refund_sgoods['sku_id']]['attribute'];
                    $service_commission = $this->settlement['service_shopping_guid_id'] ? $refund_sgoods['service_commission'] : 0;
                    $refundGoodsList[] = array(
                        'material_no' => $goodsArr[$refund_sgoods['sku_id']]['material_no'],
                        'spu_img'     => $goodsArr[$refund_sgoods['sku_id']]['spu_img'],
                        'spu_name'    => $goodsArr[$refund_sgoods['sku_id']]['spu_name'].'(部分退款)',
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
        }

        $this->addSellPrice($settlement_goods);
        return $settlement_goods;
    }

    private function addSellPrice(&$settlement_goods)
    {
        foreach ($settlement_goods as &$goods){
            $goodsInfo = CrsOrderGoods::query()
                ->where('enterprise_id',$this->enterpriseId)
                ->where('order_no', '=', $this->settlement['order_no'])
                ->where('sku_id', '=', $goods['sku_id'])
                ->first();
            if($goodsInfo){
                $goods_pay_price = $goodsInfo->toArray()['goods_pay_price'];
                $original_price = $goodsInfo->toArray()['original_price'];
            }else{
                $goods_pay_price = 0;
                $original_price = 0;
            }
            $goods['order_goods_price'] = $original_price > 0 ? $original_price : $goods['order_goods_price'];
            $goods['goods_pay_price'] = $goods_pay_price;
        }
    }

    public function formatSingle(GoodsAppService $goodsAppService)
    {

        return [
            "id"                    => $this->settlement['id'],
            "settlement_status"     => $this->settlement['settlement_status'],
            "order_no"              => $this->settlement['order_no'],
            "order_status"          => $this->settlement['order_status'],
            "order_pay_time"        => (int)$this->settlement['order_pay_time'],
            "goods_info"            => $goodsAppService->getSkuIdInfo(),
            "settlement_time"       => $this->settlementTime(),
            "settlement_money"      => ($this->settlement['order_status'] == 205) ? 0 : $this->getSettlementMoeny(),
            "partner_id"            => $this->settlement['partner_id'],
            "partner_name"          => $this->settlement['partner_name'],
            "partner_mobile"        => $this->settlement['mobile_show'],
            "partner_commission"    => ($this->settlement['order_status'] == 205) ? 0 : $this->settlement['partner_commission'],
        ];
    }

    private function settlementTime()
    {
        return intval($this->settlement['updated_at']);
    }

    private function getSettlementMoeny()
    {
        $commission = 0;

        if ($this->settlement['promotion_worker_id'] == $this->workId) {
            $commission += intval($this->settlement['promotion_commission']);
        }
        if ($this->settlement['service_shopping_guid_id'] == $this->workId) {
            $commission += intval($this->settlement['service_commission']);
        }
        if ($this->settlement['partner_id'] == $this->workId) {
            $commission += intval($this->settlement['partner_commission']);
        }
        return $commission;
    }
}