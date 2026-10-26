<?php

namespace App\Services;

use App\Models\CrsOrder;
use App\Models\Distribution;
use App\Models\DistributionGoods;
use App\Models\DistributionOrderSettlement;
use App\Models\PartnerEnterpriseMap;
use App\Services\Message\MessageService;
use App\Models\PartnerEnterpriseSettings;
use App\Models\PartnerMember;
use CrsSdk\Boss\Enterprise;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use App\Models\DistributionOrderGoodsSettlement;
use App\Models\ShoppingGuideCashWithdrawal;

class DistributionOrderSettlementService
{

    /**
     * 是否开启二级分销
     * @var
     */
    protected $partner_status;

    /**
     * 推广佣金的分佣方式
     *  1:基于商品分销 2:基于推广佣金总额的分账
     * @var
     */
    protected $commission_type;

    /**
     * 二级分销佣金比例
     * @var
     */
    protected $commission_rate;
    /**
     * 员工状态：1=在职，2=离职
     */
    /**
     * 是否给合伙人发送短信
     * 1 : 开启
     * 2 : 关闭
     * @var int
     */
    protected $sms_switch = 1;
    /**
     * 是否二级分销订单
     */
    const PartnerOrderTure = 1;
    const IN_SERVICE  = 1;
    const OUT_SERVICE = 2;

    const WORKER_PROMOTION_COMMISSION = 1;//导购推广佣金
    const PARTNER_PROMOTION_COMMISSION = 2;//合伙人推广佣金
    const WORKER_SERVICE_COMMISSION = 3;//导购服务佣金

    const APP = '畅鲟智慧导购-品牌合伙人';

    /**
     * 发送分销订单结算消息提醒
     * @param $enterpriseId
     * @param $sendOrderMsg array 结算订单数据
     * @param $settlementStatus integer 结算状态
     * @return bool
     * @throws \Exception
     * @author Calm
     * @time 2020-10-19 15:04:40
     */
    public function sendDistributionOrderMsg($enterpriseId, $sendOrderMsg, $settlementStatus) {
        if (empty($sendOrderMsg) || !isset($settlementStatus)) {
            return true;
        }
        $workerIds = [$sendOrderMsg['shoppingGuideId']];
        switchAppHost(['enterprise_id' => $enterpriseId]);
        $workerList = StoreWorkerService::list([
                'enterprise_id' => $enterpriseId,
                'ids'           => $workerIds,
            ], false)['list'] ?? [];
        if (empty($workerList)) {
            Log::error('导购信息不存在' . '参数:' . json_encode($sendOrderMsg) . '结算状态：' . $settlementStatus);
            return true;
        }
        $workerList = array_column($workerList, null, 'id');
        $content    =
            [
                [
                    'key'   => '结算单号',
                    'value' => $sendOrderMsg['settlement_id'],
                ],
                [
                    'key'   => '结算状态',
                    'value' => !empty($settlementStatus) ? '已结算' : '待结算',
                ],
                [
                    'key'   => '分销商品件数',
                    'value' => $sendOrderMsg['purchase_quantity'],
                ],
                [
                    'key'   => '预计佣金合计',
                    'value' => ($sendOrderMsg['commission'] / 100) . '元',
                ],
                [
                    'key'   => '说明',
                    'value' => !empty($settlementStatus) ? '你的分销订单已经结算，可以查看可提现金额提现咯~' : '耐心等待订单完成就可以结算佣金了~',
                ],
            ];
        $title      = !empty($settlementStatus) ? '分销订单已结算提醒' : '分销订单待结算提醒';
        $memUids    = [$workerList[$sendOrderMsg['shoppingGuideId']]['mem_uid']];
        $page    = 'pagesA/my-brokerage/my-brokerage?currentTab=2&store_id='.$workerList[$sendOrderMsg['shoppingGuideId']]['store_id'];
        (new MessageService())->sendMiniMessage($sendOrderMsg['enterprise_id'], $content, $page, $title, $memUids);
        return true;
    }


    /**
     * 设置分销订单结算和订单商品结算
     * @param array 结算数据结构 $data
     * @return mixed
     */
    public function sync($data) {
        if (empty($data['goods'])) {
            Log::info('添加结算订单失败,错误原因：商品数据为空' . '参数：' . json_encode($data));
            return true;
        }
        $data['partner_id'] = 0;
        if ($data['promotion_worker_id'] > 0) {
            Log::info('记录订单promotion_worker_id' . $data['promotion_worker_id'] ?? '');
            $partnerInfo = (new PartnerMemberService())->getPartnerInfo([
                'enterprise_id' => $data['enterprise_id'],
                'id' => $data['promotion_worker_id']
            ]);
            Log::info('记录订单$partnerInfo' . json_encode($partnerInfo));
            //存在合伙人ID，通过关系查询关联导购的ID
            if (!empty($partnerInfo)) {
                $data['partner_id'] = $partnerInfo['id'];
                $data['promotion_worker_id'] = $partnerInfo['worker_id'];
                $this->sms_switch = $partnerInfo['sms_switch'];
            }
        }
        // 推广人和服务人为空
        if (empty($data['promotion_worker_id']) && empty($data['service_shopping_guid_id'])) {
            Log::info('添加结算订单失败,错误原因：推广人或服务人为空' . '参数：' . json_encode($data));
            return true;
        }

        $params = ['enterprise_id' => $data['enterprise_id']];
        switchAppHost($params);
        //添加合伙人 start
        //二级分销设置
        $partnerSettings = $this->getPartnerSettings($params);
        if(!empty($partnerSettings)){
            $this->partner_status = $partnerSettings['partner_status'];
            $this->commission_type = $partnerSettings['commission_type'];
            $this->commission_rate = $partnerSettings['commission_rate'];
        }
        //end

        $withdrawalData               = [];
        $orderSelect['enterprise_id'] = $orderUpdate['enterprise_id'] = $sendOrderMsg['enterprise_id'] = $data['enterprise_id'];
        //1.查看订单是否已有结算记录
        $orderSelect['order_no'] = $orderUpdate['order_no'] = $data['order_no'];
        $orderSettlementField    = [
            'id',
            'order_no',
            'order_status',
            'settlement_status',
            'settlement_standard',
            'promotion_worker_id',
            'service_shopping_guid_id',
            'promotion_commission',
            'service_commission',
            //添加合伙人 start
            'partner_enterprise_id',
            'partner_id',
            'partner_commission',
            'enterprise_id',
            //end
        ];
        $orderSettlementDetail   = $this->orderSettlementDetail($orderSelect,
            $orderSettlementField);


        //2.如果存在结算记录则更新状态及字段再识别部分退款商品扣除结算佣金
        $orderSettlementUpdate = [];

        if (empty($orderSettlementDetail)) {
            /*****新增结算逻辑处理*****/
            $this->addOrderSettlement($data);
            return true;
        }

        /*****以下为更新结算逻辑*****/
        $sendOrderMsg['settlement_id'] = $orderSettlementDetail['id'];
        //全部退款
        if(CrsOrder::STATUS_COMPLETE == $data['order_status']){
           $this->allRefund($data,$orderUpdate,$orderSettlementDetail);
           return  true;
        }

        $withdrawal=false;
        if ($orderSettlementDetail['settlement_standard'] == $data['order_status']) {
            $orderSettlementUpdate['settlement_status'] = DistributionOrderSettlement::ADOPT_SETTLE;
            //当达到已完成标准标记提现申请
            $withdrawal = true;
        }
        //推广导购ID
        if (empty($orderSettlementDetail['promotion_worker_id']) && !empty($data['promotion_worker_id'])) {
            $orderSettlementUpdate['promotion_worker_id'] = $data['promotion_worker_id'];
        }
        //服务导购ID
        if (empty($orderSettlementDetail['service_shopping_guid_id']) && !empty($data['service_shopping_guid_id'])) {
            $orderSettlementUpdate['service_shopping_guid_id'] = $data['service_shopping_guid_id'];
        }
        //添加合伙人 start
        if (empty($orderSettlementDetail['partner_id']) && !empty($data['partner_id'])) {
            $orderSettlementUpdate['partner_id'] = $data['partner_id'];
        }
        //end
        $orderSettlementUpdate['order_status'] = $data['order_status'];
        //部分退款 start
        //goods表添加字段  refund_status 0:默认 1:部分退款 2:全部退款  refund_quantity 退款件数
        if(CrsOrder::PART_REFUND == $data['order_status']){
            $this->partRefund($data,$orderSettlementDetail);
            return  true;
        }
        //end
        DB::beginTransaction();
        try {
            $this->updateOrderSettlement($orderUpdate, $orderSettlementUpdate);

            //添加服务人余额
            if(empty($orderSettlementDetail['service_shopping_guid_id']) && !empty($data['service_shopping_guid_id'])){
                $workerInfo = ShoppingGuideCashWithdrawal::query()
                    ->where('worker_id', $data['service_shopping_guid_id'])
                    ->first();
                $this->saveShoppingGuideCashWithdrawal($data['enterprise_id'], $orderSettlementDetail['service_commission'], $data['service_shopping_guid_id'], $workerInfo, ShoppingGuideCashWithdrawal::CASH_TYPE_PRE_WITHDRAWAL);
            }

            if (!empty($withdrawal)) {
                Log::info('setOrderSettlement已结算插入/更新可提现金额表，订单编号:' . $data['order_no'], $orderSettlementUpdate);
                //已结算则插入/更新导购可提现金额表
                $withdrawalData = [
                    'enterprise_id'            => $data['enterprise_id'],
                    'promotion_worker_id'      => $orderSettlementDetail['promotion_worker_id'],
                    'promotion_commission'     => $orderSettlementDetail['promotion_commission'],
                    'service_shopping_guid_id' => $orderSettlementDetail['service_shopping_guid_id'],
                    'service_commission'       => $orderSettlementDetail['service_commission'],
                    //添加合伙人 start
                    'partner_id'               => $orderSettlementDetail['partner_id'],
                    'partner_commission'       => $orderSettlementDetail['partner_commission'],
                    //end
                ];
                $this->operateShoppingGuideCashWithdrawal($withdrawalData, ShoppingGuideCashWithdrawal::CASH_TYPE_WITHDRAWAL);

                //资金转移
                $withdrawalData = [
                    'enterprise_id'            => $data['enterprise_id'],
                    'promotion_worker_id'      => $orderSettlementDetail['promotion_worker_id'],
                    'promotion_commission'     => - $orderSettlementDetail['promotion_commission'],
                    'service_shopping_guid_id' => $orderSettlementDetail['service_shopping_guid_id'],
                    'service_commission'       => - $orderSettlementDetail['service_commission'],
                    //添加合伙人 start
                    'partner_id'               => $orderSettlementDetail['partner_id'],
                    'partner_commission'       => - $orderSettlementDetail['partner_commission'],
                    //end
                ];
                $this->operateShoppingGuideCashWithdrawal($withdrawalData, ShoppingGuideCashWithdrawal::CASH_TYPE_PRE_WITHDRAWAL);
                //推送已结算佣金提醒
            }
            DB::commit();
            //获取结算订单商品数量
            $sendOrderMsg['purchase_quantity'] = DistributionOrderGoodsSettlement::query()
                ->where('enterprise_id', $data['enterprise_id'])
                ->where('order_no', $data['order_no'])
                ->sum('purchase_quantity');
            if (!empty($withdrawal)) {
                //推送已结算消息提醒 消息推送需要的订单推广导购ID、服务导购ID、推广佣金、服务佣金、结算状态
                $settledWorkerId = $orderSettlementDetail['promotion_worker_id'] ?? 0;
                $settledWorkerId = $orderSettlementUpdate['promotion_worker_id'] ?? $settledWorkerId;
                if (!empty($settledWorkerId)) {
                    $sendOrderMsg['shoppingGuideId'] = $settledWorkerId;
                    $sendOrderMsg['commission']      = $orderSettlementDetail['promotion_commission'];
                    $this->sendDistributionOrderMsg($data['enterprise_id'], $sendOrderMsg, $orderSettlementUpdate['settlement_status']);
                }
                $settledServiceId = $orderSettlementDetail['service_shopping_guid_id'] ?? 0;
                $settledServiceId = $orderSettlementUpdate['service_shopping_guid_id'] ?? $settledServiceId;
                if (!empty($settledServiceId)) {
                    $sendOrderMsg['shoppingGuideId'] = $settledServiceId;
                    $sendOrderMsg['commission']      = $orderSettlementDetail['service_commission'];;
                    $this->sendDistributionOrderMsg($data['enterprise_id'], $sendOrderMsg, $orderSettlementUpdate['settlement_status']);
                }
                //添加合伙人 start 发短信 已结算 可提现变更
                if(!empty($orderSettlementDetail['partner_id']) && $this->partner_status == PartnerEnterpriseSettings::PARTNER_STATUS_ENABLE){
                    $orderSettlementDetail['settlement_status'] = DistributionOrderSettlement::ADOPT_SETTLE;
                    //订单状态变更
                    $this->sendSmsOrderStatus($orderSettlementDetail);
                    //余额
                    $this->sendSmsCashBalance($orderSettlementDetail);
                }
                //end
            } else {
                //推送待结算消息提醒 消息推送需要的订单推广导购ID、服务导购ID、推广佣金、服务佣金、结算状态
                if (!empty($orderSettlementUpdate['promotion_worker_id'])) {
                    $sendOrderMsg['shoppingGuideId'] = $orderSettlementUpdate['promotion_worker_id'];
                    $sendOrderMsg['commission']      = $orderSettlementDetail['promotion_commission'];
                    $this->sendDistributionOrderMsg($data['enterprise_id'], $sendOrderMsg, $orderSettlementDetail['settlement_status']);
                }
                if (!empty($orderSettlementUpdate['service_shopping_guid_id'])) {
                    $sendOrderMsg['shoppingGuideId'] = $orderSettlementUpdate['service_shopping_guid_id'];
                    $sendOrderMsg['commission']      = $orderSettlementDetail['service_commission'];;
                    $this->sendDistributionOrderMsg($data['enterprise_id'], $sendOrderMsg, $orderSettlementDetail['settlement_status']);
                }
            }
        } catch (\Exception $exception) {
            DB::rollBack();
            Log::info('订单结算表更新失败：' . $exception->getFile() . ',行数：' . $exception->getLine() . ',错误原因：' . $exception->getMessage() . ',条件：' . json_encode($orderUpdate) . ',结算参数：' . json_encode($orderSettlementUpdate) . ',提现参数：' . json_encode($withdrawalData).'.getTrace'.json_encode($exception->getTrace()));
        }

       /* //识别退款商品并更新
        if (empty($data['refundGoods'])) {
            return true;
        }*/
//        $skuIds = array_column($data['refundGoods'], 'sku_id');
//        $fields = [
//            'sku_id',
//            'purchase_quantity',
//            'order_goods_price',
//            'promotion_commission',
//            'service_commission',
//            'set_type_ratio_json',
//        ];
//        //检索当前退款单号及商品是否存在结算记录
//        $goodsSettledList = $this->getOrderGoodsSettlementList($data['enterprise_id'], $data['order_no'], $skuIds, $fields);
//        if (empty($goodsSettledList)) {
//            return true;
//        }
//        $refundGoods = array_column($data['refundGoods'], NULL, 'sku_id');
//        //定义及计算需要减去的推广佣金和服务佣金
//        $updateRefundData['promotion_commission'] = $updateRefundData['service_commission'] = 0;
//        foreach ($goodsSettledList as $goodsItem) {
//            $setTypeRatio = json_decode($goodsItem['set_type_ratio_json'], true);
//            //计算已退款商品推广佣金需减去的金额
//            $updateRefundData['promotion_commission'] = $this->calcCommissionRefund(
//                $setTypeRatio['promotion_commission_type'],
//                $setTypeRatio['promotion_commission_amount'],
//                $goodsItem['order_goods_price'],
//                $refundGoods[$goodsItem['sku_id']]['current_apply_quantity']
//            );
//
//            //计算已退款商品服务佣金需减去的金额
//            $updateRefundData['service_commission'] = $this->calcCommissionRefund(
//                $setTypeRatio['service_commission_type'],
//                $setTypeRatio['service_commission_amount'],
//                $goodsItem['order_goods_price'],
//                $refundGoods[$goodsItem['sku_id']]['current_apply_quantity']
//            );
//        }
//        $updateRefundData['promotion_commission'] = $orderSettlementDetail['promotion_commission'] - $updateRefundData['promotion_commission'];
//        $updateRefundData['service_commission']   = $orderSettlementDetail['service_commission'] - $updateRefundData['service_commission'];
//        DistributionOrderSettlement::query()
//            ->where('enterprise_id', $data['enterprise_id'])
//            ->where('order_no', $data['order_no'])
//            ->update($updateRefundData);

        return true;
    }

    public function allRefund($data,$orderUpdate,$orderSettlementDetail)
    {
        //查看订单全部退款之前是否部分退款过
        $goodsData = DistributionOrderGoodsSettlement::query()
            ->where('enterprise_id',$data['enterprise_id'])
            ->where('order_no',$data['order_no'])
            ->get();
        $refund_status = 0;
        foreach ($goodsData as $k=>$v){
            if($v['refund_status'] != DistributionOrderGoodsSettlement::REFUND_STATUS_M){
                $refund_status = 1;
            }
        }
        //已部分退款过
        if(!empty($refund_status)){
            $this->partRefund($data,$orderSettlementDetail);
            return  true;
        }
        /*//同一商品买两件 退一件
        $isAllRefund = false;
        $fields = [
            'sku_id',
            'purchase_quantity',
            'order_goods_price',
            'promotion_commission',
            'service_commission',
            'set_type_ratio_json',
            'partner_commission',
            'refund_status',
            'refund_quantity',
        ];
        $refundGoods = $data['refundGoods'];
        if (empty($refundGoods)) {
            return true;
        }
        $skuIds = array_column($refundGoods, 'sku_id');
        //商品
        $goodsSettledList = $this->getOrderGoodsSettlementList($data['enterprise_id'], $data['order_no'], $skuIds, $fields);
        if (empty($goodsSettledList)) {
            return true;
        }
        $goodsSettledList = array_column($goodsSettledList,null,'sku_id');
        //退款的商品数量和购买的商品数量不一致就走部分退款
        foreach ($refundGoods as $key => $val){
            //退款商品数量
            $refund_goods_num = $val['current_apply_quantity'];
            //购买商品数量
            $order_goods_num = $goodsSettledList[$val['sku_id']]['purchase_quantity'] ?? 0 ;
            if($refund_goods_num != $order_goods_num){
                $isAllRefund = true;
            }
        }
        //退款数和购买数量不一致 走部分退款逻辑
        if($isAllRefund){
            $this->partRefund($data,$orderSettlementDetail,$isAllRefund);
            return  true;
        }*/
        //如果推送的订单状态为全部退款则取消当前订单的结算
        if (CrsOrder::STATUS_COMPLETE == $data['order_status']
            && in_array($orderSettlementDetail['settlement_status'],[DistributionOrderSettlement::ADOPT_SETTLE, DistributionOrderSettlement::WAIT_SETTLE])
        ) {
            //添加服务人余额,退款修改余额表
            $cashType = DistributionOrderSettlement::ADOPT_SETTLE == $orderSettlementDetail['settlement_status'] ? ShoppingGuideCashWithdrawal::CASH_TYPE_WITHDRAWAL : ShoppingGuideCashWithdrawal::CASH_TYPE_PRE_WITHDRAWAL;
            $withdrawalData = [
                'enterprise_id'            => $data['enterprise_id'],
                'promotion_worker_id'      => $orderSettlementDetail['promotion_worker_id'],
                'promotion_commission'     => - $orderSettlementDetail['promotion_commission'],
                'service_shopping_guid_id' => $orderSettlementDetail['service_shopping_guid_id'],
                'service_commission'       => - $orderSettlementDetail['service_commission'],
                //添加合伙人 start
                'partner_id'               => $orderSettlementDetail['partner_id'],
                'partner_commission'       => - $orderSettlementDetail['partner_commission'],
                //end
            ];
            $this->operateShoppingGuideCashWithdrawal($withdrawalData, $cashType);

            $orderAllRefundUpdate['settlement_status'] = DistributionOrderSettlement::CANCEL_SETTLE;
            $orderAllRefundUpdate['order_status']      = $data['order_status'];
            $this->updateOrderSettlement($orderUpdate, $orderAllRefundUpdate);
            //结算订单取消 给合伙人发消息
            if(!empty($orderSettlementDetail['partner_id']) && $this->partner_status = PartnerEnterpriseSettings::PARTNER_STATUS_ENABLE){
                $orderSettlementDetail['settlement_status'] = DistributionOrderSettlement::CANCEL_SETTLE;
                $this->sendSmsOrderStatus($orderSettlementDetail);
            }

            return true;
        }
        return true;
    }

    /**
     * 获得当前退款的商品
     * 因每次退款商品信息是当前订单的全部退款商品
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/12/14 17:24
     */
    public function currentRefundGoods($data,$fields)
    {
        //检索当前退款单号及商品是否存在结算记录
        $skuIds = array_column($data['refundGoods'], 'sku_id');
        $goodsSettledList = $this->getOrderGoodsSettlementList($data['enterprise_id'], $data['order_no'], $skuIds, $fields);
        if (empty($goodsSettledList)) {
            return [];
        }

        $refundGoodsSettledList = array_column($goodsSettledList, NULL, 'sku_id');
        $refundGoods = array_column($data['refundGoods'], NULL, 'sku_id');

        //获得当前退款的商品
        $currentRefundGoods = [];
        foreach ($refundGoodsSettledList as $value)
        {
            if(in_array($value['refund_status'],[DistributionOrderGoodsSettlement::REFUND_STATUS_M]) ){
                $currentRefundGoods[] = [
                    'sku_id' => $value['sku_id'],
                    //只退实际退了商品数量
                    //'current_refund_quantity' => $refundGoods[$value['sku_id']]['current_apply_quantity'] - $value['refund_quantity'],
                    //单个商品退一件,此商品全部退
                    'current_refund_quantity' => $value['purchase_quantity'],
                ];
            }
        }
        return  $currentRefundGoods;
    }

    /**
     * 部分退款
     * @param $data
     * @param $orderSettlementDetail
     * @return bool
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/12/14 17:02
     */
    public function partRefund($data,$orderSettlementDetail,$isAllRefund = '')
    {
        //order_status = 83 部分退款 205 全部退款
        //因商城两件商品退一件，商城状态是全部退 故兼容
        /*if(!empty($isAllRefund)){
            $data['order_status'] = CrsOrder::PART_REFUND;
        }*/
        //余额表处理部分退款
        if ((CrsOrder::PART_REFUND == $data['order_status'] || CrsOrder::STATUS_COMPLETE == $data['order_status'])
            && in_array($orderSettlementDetail['settlement_status'],[DistributionOrderSettlement::ADOPT_SETTLE, DistributionOrderSettlement::WAIT_SETTLE]
            )
        ) {
            //识别退款商品并更新
            if (empty($data['refundGoods'])) {
                return true;
            }
            $fields = [
                'sku_id',
                'purchase_quantity',
                'order_goods_price',
                'promotion_commission',
                'service_commission',
                'set_type_ratio_json',
                'partner_commission',
                'refund_status',
                'refund_quantity',
            ];

            $currentRefundGoods = $this->currentRefundGoods($data,$fields);
            if (empty($currentRefundGoods)) {
                return true;
            }
            $skuIds = array_column($currentRefundGoods, 'sku_id');
            //检索当前退款单号及商品是否存在结算记录
            $goodsSettledList = $this->getOrderGoodsSettlementList($data['enterprise_id'], $data['order_no'], $skuIds, $fields);
            if (empty($goodsSettledList)) {
                return true;
            }

            $refundGoods = array_column($currentRefundGoods, NULL, 'sku_id');

            //定义及计算需要减去的推广佣金和服务佣金
            $refund_promotion_commission = $refund_service_commission = $refund_partner_commission = 0;
            foreach ($goodsSettledList as $goodsItem) {
                $setTypeRatio = json_decode($goodsItem['set_type_ratio_json'], true);
                //计算已退款商品推广佣金需减去的金额
                $commission = [];
                //全分销的订单
                if(empty($data['partner_id'])){
                    $commission['promotion_commission'] = $this->calcCommissionRefund(
                        $setTypeRatio['promotion_commission_type'],
                        $setTypeRatio['promotion_commission_amount'],
                        $goodsItem['order_goods_price'],
                        1//$refundGoods[$goodsItem['sku_id']]['current_refund_quantity']
                    );

                    //计算已退款商品服务佣金需减去的金额
                    $commission['service_commission'] = $this->calcCommissionRefund(
                        $setTypeRatio['service_commission_type'],
                        $setTypeRatio['service_commission_amount'],
                        $goodsItem['order_goods_price'],
                        1//$refundGoods[$goodsItem['sku_id']]['current_refund_quantity']
                    );
                }
                //二级分销的订单
                if(!empty($data['partner_id'])){
                    //开启二级分销 && 是二级分销订单
                    if($setTypeRatio['partner_status'] == PartnerEnterpriseSettings::PARTNER_STATUS_ENABLE
                        && $setTypeRatio['partner_order'] == self::PartnerOrderTure
                    ){
                        if (DistributionGoods::PROMOTION_COMMISSION_TYPE_PERCENTAGE == $setTypeRatio['promotion_commission_type']) {
                            //基于商品分销
                            if($setTypeRatio['commission_type'] == PartnerEnterpriseSettings::COMMISSION_TYPE_GOODS){
                                //合伙人推广佣金
                                $commission['partner_commission'] =  $this->priceCalculation($goodsItem['order_goods_price'],$setTypeRatio['promotion_commission_amount']);
                                //导购推广佣金
                                $commission['promotion_commission'] =  $this->priceCalculation($goodsItem['order_goods_price'],$setTypeRatio['commission_rate']);
                            }
                            //基于推广佣金总额的分账
                            if($setTypeRatio['commission_type'] == PartnerEnterpriseSettings::COMMISSION_TYPE_PROMOTION){
                                //总的推广佣金
                                $promotion_commission =  $this->priceCalculation($goodsItem['order_goods_price'],$setTypeRatio['promotion_commission_amount']);
                                //合伙人推广佣金
                                $commission['partner_commission'] = $this->priceCalculation($promotion_commission,$setTypeRatio['commission_rate']);
                                //导购推广佣金 = 总推广佣金-合伙人佣金
                                $commission_rate_guide = 10000 - $setTypeRatio['commission_rate'];
                                $commission['promotion_commission'] = $this->priceCalculation($promotion_commission,$commission_rate_guide);
                                //(int)floor($promotion_commission - $commission['partner_commission']);
                            }
                        }
                    }
                    //未设置或二级分销关闭
                    if(empty($setTypeRatio['partner_status']) || PartnerEnterpriseSettings::PARTNER_STATUS_DISABLE){
                        $commission['promotion_commission'] =  $this->priceCalculation($goodsItem['order_goods_price'],$setTypeRatio['promotion_commission_amount']);
                    }
                    if (DistributionGoods::PROMOTION_COMMISSION_TYPE_PERCENTAGE == $setTypeRatio['service_commission_type']) {//导购服务佣金
                        //导购服务佣金
                        $commission['service_commission'] =  $this->priceCalculation($goodsItem['order_goods_price'],$setTypeRatio['service_commission_amount']);
                    }
                }
                //导购推广佣金
                $refund_promotion_commission  += ($commission['promotion_commission'] ?? 0)* $refundGoods[$goodsItem['sku_id']]['current_refund_quantity'];
                //导购服务佣金
                $refund_service_commission  += ($commission['service_commission'] ?? 0) * $refundGoods[$goodsItem['sku_id']]['current_refund_quantity'];
                //合伙人推广佣金
                $refund_partner_commission += ($commission['partner_commission'] ?? 0) * $refundGoods[$goodsItem['sku_id']]['current_refund_quantity'];
            }
            $updateRefundData['promotion_commission'] = $orderSettlementDetail['promotion_commission'] - $refund_promotion_commission;
            $updateRefundData['service_commission']   = $orderSettlementDetail['service_commission'] - $refund_service_commission;
            $updateRefundData['partner_commission'] =  $orderSettlementDetail['partner_commission'] - $refund_partner_commission;
            //为全员分销商品订单，防止佣金未负
            if($orderSettlementDetail['partner_commission'] == 0){
                $updateRefundData['partner_commission'] = 0;
                $refund_partner_commission = 0;
            }
            $cashType = ($orderSettlementDetail['settlement_status']) ? ShoppingGuideCashWithdrawal::CASH_TYPE_WITHDRAWAL : ShoppingGuideCashWithdrawal::CASH_TYPE_PRE_WITHDRAWAL;
            DB::beginTransaction();
            try {
                //修改退款商品状态
                $this->updateDistributionOrderGoodsSettlement($data, $fields);
                //更新订单表
                DistributionOrderSettlement::query()
                    ->where('enterprise_id', $data['enterprise_id'])
                    ->where('order_no', $data['order_no'])
                    ->update($updateRefundData);
                //更新余额表 预估可提现金额
                $withdrawalData = [
                    'enterprise_id'            => $data['enterprise_id'],
                    'promotion_worker_id'      => $orderSettlementDetail['promotion_worker_id'],
                    'promotion_commission'     => - $refund_promotion_commission,
                    'service_shopping_guid_id' => $orderSettlementDetail['service_shopping_guid_id'],
                    'service_commission'       => - $refund_service_commission,
                    'partner_id'               => $orderSettlementDetail['partner_id'],
                    'partner_commission'       => - $refund_partner_commission,
                ];
                $this->operateShoppingGuideCashWithdrawal($withdrawalData, $cashType);
                //全部退款修改结算状态
                if(CrsOrder::STATUS_COMPLETE == $data['order_status']){
                    $orderAllRefundUpdate['settlement_status'] = DistributionOrderSettlement::CANCEL_SETTLE;
                    $orderAllRefundUpdate['order_status']      = $data['order_status'];
                    $this->updateOrderSettlement(['order_no'=>$data['order_no'],'enterprise_id'=>$data['enterprise_id']], $orderAllRefundUpdate);
                }
                DB::commit();
            } catch (\Exception $exception) {
                DB::rollBack();
                Log::info('订单结算部分退款：' . $exception->getFile() . ',行数：' . $exception->getLine() . ',错误原因：' . $exception->getMessage() . ',条件：' . json_encode(['order_no'=>$data['order_no'],'enterprise_id'=>$data['enterprise_id']]) . ',提现参数：' . json_encode($withdrawalData).'.getTrace'.json_encode($exception->getTrace()));
            }
            $msg = '部分退款';
            if(CrsOrder::STATUS_COMPLETE == $data['order_status']){
                $msg = '已取消';
            }
            //发送短信
            if(!empty($orderSettlementDetail['partner_id']) && $this->partner_status == PartnerEnterpriseSettings::PARTNER_STATUS_ENABLE){
                $orderSettlementDetail['settlement_status'] = $cashType;
                //订单状态变更
                $this->sendSmsOrderStatus($orderSettlementDetail,$msg);
                //余额
                $this->sendSmsCashBalance($orderSettlementDetail);
            }
            return true;
        }
        return true;
    }

    /**
     * 部分退款修改订单商品数据状态
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/12/14 17:56
     */
    public function updateDistributionOrderGoodsSettlement($data,$fields)
    {
        $skuIds = array_column($data['refundGoods'], 'sku_id');
        $goodsSettledList = $this->getOrderGoodsSettlementList($data['enterprise_id'], $data['order_no'], $skuIds, $fields);
        $refundGoods = array_column($data['refundGoods'], NULL, 'sku_id');
        foreach ($goodsSettledList as $key=>$value){
            if($value['purchase_quantity'] == $refundGoods[$value['sku_id']]['current_apply_quantity']){
                $refund_status = DistributionOrderGoodsSettlement::REFUND_STATUS_ALL;
            }else{
                $refund_status = DistributionOrderGoodsSettlement::REFUND_STATUS_PART;
            }
               DistributionOrderGoodsSettlement::query()
                   ->where('enterprise_id', $data['enterprise_id'])
                   ->where('order_no', $data['order_no'])
                   ->where('sku_id',$value['sku_id'])
                   ->update([
                       'refund_status'=>$refund_status,'refund_quantity'=>$refundGoods[$value['sku_id']]['current_apply_quantity']
                   ]);
        }
        return true;
    }
    /**
     * 已退款商品退款金额
     * @param $commissionType
     * @param $commissionAmount
     * @param $orderGoodsPrice
     * @param $currentApplyQuantity
     * @return float|int
     * @author Calm
     * @time 2020-10-19 22:25:44
     */
    private function calcCommissionRefund($commissionType, $commissionAmount, $orderGoodsPrice, $currentApplyQuantity) {
        $commission = 0;
        if (DistributionGoods::PROMOTION_COMMISSION_TYPE_PERCENTAGE == $commissionType) {
            $unit_price    =  $this->priceCalculation($orderGoodsPrice,$commissionAmount);
            $commission    = $unit_price * $currentApplyQuantity;
        } else {
            $commission += $commission * $commission;
        }
        return $commission;
    }

    /**
     * 获取结算订单详情
     * @param $params
     * @param $field array 查询的字段
     * @return mixed
     */
    public function orderSettlementDetail($params, $field = ['*']) {
        $where = [];
        $where = array_merge($where, ['enterprise_id' => $params['enterprise_id']]);
        //订单号
        if (!empty($params['order_no'])) {
            $where = array_merge($where, ['order_no' => $params['order_no']]);
        }

        $dataOld = DistributionOrderSettlement::query()
            ->where($where)
            ->first($field);

        return $dataOld ? $dataOld->toArray() : [];
    }

    /**
     * 新增分销订单结算和订单商品结算
     * @param array 结算数据结构 $data
     * @return mixed
     */
    public function addOrderSettlement($data) {
        //只有从已支付状态开始才会处理新增结算逻辑
        if (CrsOrder::STATUS_PRE_DELIVER != $data['order_status']) {
            return true;
        }

        $dataOne['enterprise_id'] = $sendOrderMsg['enterprise_id'] = $data['enterprise_id'];
        $time                     = time();

        // 1.获取结算标准
        $distribution = (new DistributionService())->getDistribution($dataOne);
        if (empty($distribution)) {
            $distribution['settlement_standard'] = Distribution::SETTLEMENT_STANDARD_COMPLETE;
        }
        $dataTwo['enterprise_id'] = $data['enterprise_id'];
        $dataTwo['sku_ids']       = array_column($data['goods'], 'sku_id');
        //2.获取分销商品配置信息
        $distributionGoods = (new DistributionGoodsService())->getDistributionGoodsList($dataTwo);
        if (empty($distributionGoods)) {
            return true;
        }
        //3.计算商品推广佣金和服务佣金
        $newGoods = array_column($data['goods'], NULL, 'sku_id');
        foreach ($distributionGoods as $num => &$vo) {//过滤不在分销时间范围类的
            //是否在设置的分销时间区间内
            if ($time < $vo['distribution_start'] || $time > $vo['distribution_end']) {
                unset($distributionGoods[$num]);
                continue;
            }
            $vo['purchase_quantity'] = $newGoods[$vo['sku_id']]['purchase_quantity'];
            $vo['order_goods_price'] = $newGoods[$vo['sku_id']]['order_goods_price'];
        }
        if (empty($distributionGoods)) {
            return true;
        }
        $distributionGoods        = array_values($distributionGoods);
        $goodsSettlement          = [];
        $orderPromotionCommission = $orderServiceCommission = $sendOrderMsg['purchase_quantity'] =  $orderPartnerCommission = 0;
        //4.组装订单商品结算表数据
        foreach ($distributionGoods as $key => $vaule) {

            $setTypeRatioJson = [
                'promotion_commission_type'   => $vaule['promotion_commission_type'],
                'promotion_commission_amount' => $vaule['promotion_commission_amount'],
                'service_commission_type'     => $vaule['service_commission_type'],
                'service_commission_amount'   => $vaule['service_commission_amount'],
                //是否为二级分销订单
                'partner_order'              => !empty($data['partner_id']) ? 1 : 0,
                'partner_status'              =>  $this->partner_status,
                'commission_type'             => $this->commission_type ?? 0,
                'commission_rate'             => $this->commission_rate ?? 0,
            ];

            $goodsSettlement[$key]['set_type_ratio_json'] = json_encode($setTypeRatioJson);
            $goodsSettlement[$key]['enterprise_id']       = $data['enterprise_id'];
            $goodsSettlement[$key]['order_no']            = $data['order_no'];
            $goodsSettlement[$key]['sku_id']              = $vaule['sku_id'];
            $goodsSettlement[$key]['purchase_quantity']   = $vaule['purchase_quantity'];
            $goodsSettlement[$key]['order_goods_price']   = $vaule['order_goods_price'];
            $goodsSettlement[$key]['created_at']          = $time;
            $goodsSettlement[$key]['updated_at']          = $time;
            //计算商品推广佣金&服务佣金
            //导购分享商品
            if(empty($data['partner_id'])){
                $goodsSettlement[$key]['promotion_commission'] = $this->calcCommission($vaule['promotion_commission_type'], $vaule['promotion_commission_amount'], $vaule['order_goods_price']);
                $goodsSettlement[$key]['service_commission']   = $this->calcCommission($vaule['service_commission_type'], $vaule['service_commission_amount'], $vaule['order_goods_price']);
            }
            //添加合伙人 start  合伙人分享商品
            if(!empty($data['partner_id'])) {
                $commissionInfo = $this->getCommissionInfo($vaule);
                //导购推广佣金
                $goodsSettlement[$key]['promotion_commission'] = $commissionInfo['promotion_commission'] ?? 0;
                //导购服务佣金
                $goodsSettlement[$key]['service_commission'] = $commissionInfo['service_commission'] ?? 0;
                //合伙人推广佣金
                $goodsSettlement[$key]['partner_commission'] = $commissionInfo['partner_commission'] ?? 0;
                //添加合伙人 start 订单合伙人推广佣金
                $orderPartnerCommission +=  (int)$goodsSettlement[$key]['partner_commission'] * $vaule['purchase_quantity'];
            }
            //查询导购离职情况
            $workerInfoByWorkerID = StoreWorkerService::detail(
                [
                    'enterprise_id' => $data['enterprise_id'],
                    'worker_id' => $data['promotion_worker_id'],
                ]
            );
            if(empty($workerInfoByWorkerID) ||  $workerInfoByWorkerID['worker_status'] == self::OUT_SERVICE){ //导购离职
                $goodsSettlement[$key]['promotion_commission']  = 0;
            }
            // end
            //计算订单导购推广佣金
            $orderPromotionCommission += (int)$goodsSettlement[$key]['promotion_commission'] * $vaule['purchase_quantity'];
            //计算订单导购服务佣金
            $orderServiceCommission += (int)$goodsSettlement[$key]['service_commission'] * $vaule['purchase_quantity'];
            //消息推送需要的结算订单商品总数
            $sendOrderMsg['purchase_quantity'] += $vaule['purchase_quantity'];
        }

        //4.组装订单结算表数据
        $orderPromotionCommission = $data['promotion_worker_id'] ? $orderPromotionCommission : 0;
        $orderSettlement = [
            'enterprise_id'            => $data['enterprise_id'],
            'order_no'                 => $data['order_no'],
            'order_status'             => $data['order_status'],
            'settlement_standard'      => $distribution['settlement_standard'],
            'settlement_status'        => $distribution['settlement_standard'] == $data['order_status'] ? DistributionOrderSettlement::ADOPT_SETTLE : DistributionOrderSettlement::WAIT_SETTLE,
            'order_pay_time'           => $data['order_pay_time'],
            'promotion_worker_id'      => $data['promotion_worker_id'] ?? 0,
            'promotion_commission'     => $orderPromotionCommission,
            'service_shopping_guid_id' => $data['service_shopping_guid_id'] ?? 0,
            'service_commission'       => $orderServiceCommission,
            //添加合伙人 start
            'partner_id'               => $data['partner_id'] ?? 0,
            'partner_commission'       => $orderPartnerCommission,
            //end
            'created_at'               => $time,
            'updated_at'               => $time,
            'promotion_commission_snapshot' => $orderPromotionCommission,
            'partner_commission_snapshot'   => $orderPartnerCommission,
            'service_commission_snapshot'   => $orderServiceCommission
        ];
        //5.插入订单商品结算表和订单结算表
        DB::beginTransaction();
        try {
            DistributionOrderGoodsSettlement::query()->insert($goodsSettlement);
            $orderSettlement['id'] = idGenerate();
            $settlementTable = DistributionOrderSettlement::query()->create($orderSettlement);
            if (DistributionOrderSettlement::ADOPT_SETTLE == $orderSettlement['settlement_status']) {
                Log::error('addOrderSettlement已结算插入/更新可提现金额表', $orderSettlement);
                //已结算则插入/更新导购可提现金额表
                $withdrawalData = [
                    'enterprise_id'            => $data['enterprise_id'],
                    'promotion_worker_id'      => $orderSettlement['promotion_worker_id'],
                    'promotion_commission'     => $orderSettlement['promotion_commission'],
                    'service_shopping_guid_id' => $orderSettlement['service_shopping_guid_id'],
                    'service_commission'       => $orderSettlement['service_commission'],
                    'partner_id'               => $orderSettlement['partner_id'],
                    'partner_commission'       => $orderSettlement['partner_commission'],
                ];
                $this->operateShoppingGuideCashWithdrawal($withdrawalData, ShoppingGuideCashWithdrawal::CASH_TYPE_WITHDRAWAL);
            }elseif(DistributionOrderSettlement::WAIT_SETTLE == $orderSettlement['settlement_status']){
                Log::error('addOrderSettlement已结算插入/更新待提现金额表', $orderSettlement);
                //已结算则插入/更新导购可提现金额表
                $preWithdrawalData = [
                    'enterprise_id'            => $data['enterprise_id'],
                    'promotion_worker_id'      => $orderSettlement['promotion_worker_id'],
                    'promotion_commission'     => $orderSettlement['promotion_commission'],
                    'service_shopping_guid_id' => $orderSettlement['service_shopping_guid_id'],
                    'service_commission'       => $orderSettlement['service_commission'],
                    'partner_id'               => $orderSettlement['partner_id'],
                    'partner_commission'       => $orderSettlement['partner_commission'],
                ];
                $this->operateShoppingGuideCashWithdrawal($preWithdrawalData, ShoppingGuideCashWithdrawal::CASH_TYPE_PRE_WITHDRAWAL);
            }
            $sendOrderMsg['settlement_id'] = $settlementTable->id;
            DB::commit();
            //消息推送需要的订单推广导购ID、服务导购ID、推广佣金、服务佣金、结算状态
            if (!empty($orderSettlement['promotion_worker_id'])) {
                $sendOrderMsg['shoppingGuideId'] = $orderSettlement['promotion_worker_id'];
                $sendOrderMsg['commission']      = $orderPromotionCommission;
                $this->sendDistributionOrderMsg($data['enterprise_id'], $sendOrderMsg, $orderSettlement['settlement_status']);
            }
            if (!empty($orderSettlement['service_shopping_guid_id'])) {
                $sendOrderMsg['shoppingGuideId'] = $orderSettlement['service_shopping_guid_id'];
                $sendOrderMsg['commission']      = $orderServiceCommission;
                $this->sendDistributionOrderMsg($data['enterprise_id'], $sendOrderMsg, $orderSettlement['settlement_status']);
            }
            //添加合伙人start 发短信
            if(!empty($orderSettlement['partner_id']) && $this->partner_status == PartnerEnterpriseSettings::PARTNER_STATUS_ENABLE){
                $this->sendSmsOrderStatus($orderSettlement);
            }
            //end
        } catch (\Exception $exception) {
            DB::rollBack();
            Log::info('结算添加失败文件：' . $exception->getFile() . ',行数：' . $exception->getLine() . ',错误原因：' . $exception->getMessage() . ',参数：' . json_encode($goodsSettlement));
        }

        return true;
    }

    /**
     * 佣金计算
     * @param $commissionType
     * @param $commissionAmount
     * @param $orderGoodsPrice
     * @return int
     * @author Zhujiaqi
     * @time 2020-10-19 22:21:52
     */
    private function calcCommission($commissionType, $commissionAmount, $orderGoodsPrice) {

        //计算商品推广佣金
        if (DistributionGoods::PROMOTION_COMMISSION_TYPE_PERCENTAGE == $commissionType) {
            $commissionAmount = $this->priceCalculation($orderGoodsPrice,$commissionAmount);
        }

        return $commissionAmount;
    }

    /**
     * 价格精度计算
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/12/25 15:01
     */
    public function priceCalculation($order_goods_price,$rate)
    {
        if($order_goods_price < 0){
            return 0;
        }
       return (int)floor(bcmul($order_goods_price,$rate,2)/10000);
    }
    /**
     * 添加合伙人 佣金计算
     * @param $value
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/11/23 16:48
     */
    public function getCommissionInfo($value)
    {
        $commission = [];
        //开启二级分销
        if($this->partner_status == PartnerEnterpriseSettings::PARTNER_STATUS_ENABLE){
            if (DistributionGoods::PROMOTION_COMMISSION_TYPE_PERCENTAGE == $value['promotion_commission_type']) {
                //基于商品分销
                if($this->commission_type == PartnerEnterpriseSettings::COMMISSION_TYPE_GOODS){
                    //合伙人推广佣金
                    $commission['partner_commission'] = $this->priceCalculation($value['order_goods_price'],$value['promotion_commission_amount']);
                    //导购推广佣金
                    $commission['promotion_commission'] = $this->priceCalculation($value['order_goods_price'],$this->commission_rate);
                }
                //基于推广佣金总额的分账
                if($this->commission_type == PartnerEnterpriseSettings::COMMISSION_TYPE_PROMOTION){
                    //总的推广佣金
                    $promotion_commission = $this->priceCalculation($value['order_goods_price'],$value['promotion_commission_amount']);
                    //合伙人推广佣金
                    $commission['partner_commission'] = $this->priceCalculation($promotion_commission,$this->commission_rate);
                    //导购推广佣金 = 总推广佣金-合伙人佣金
                    $commission_rate_guide = 10000 - $this->commission_rate;
                    $commission['promotion_commission'] =  $this->priceCalculation($promotion_commission,$commission_rate_guide);
                }
            }
        }
        //未设置或二级分销关闭
        if(empty($this->partner_status) || PartnerEnterpriseSettings::PARTNER_STATUS_DISABLE){
            $commission['promotion_commission'] =  $this->priceCalculation($value['order_goods_price'],$value['promotion_commission_amount']);
        }
        if (DistributionGoods::PROMOTION_COMMISSION_TYPE_PERCENTAGE == $value['service_commission_type']) {//导购服务佣金
            //导购服务佣金
            $commission['service_commission'] =  $this->priceCalculation($value['order_goods_price'],$value['service_commission_amount']);
        }
        return $commission;
    }

    /**
     * 更新结算订单
     * @param $params
     * @param $update array
     * @return mixed
     */
    public function updateOrderSettlement($params, $update) {
        return DistributionOrderSettlement::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('order_no', $params['order_no'])
            ->update($update);
    }

    /**
     * @param array $withdrawalData 提现申请数据结构
     * @param int $cashType
     * @return bool
     * @author Calm
     * @time 2020-10-19 21:43:43
     */
    public function operateShoppingGuideCashWithdrawal($withdrawalData, $cashType = ShoppingGuideCashWithdrawal::CASH_TYPE_PRE_WITHDRAWAL) {
        //都不存在直接返回
        if (empty($withdrawalData['promotion_worker_id']) && empty($withdrawalData['service_shopping_guid_id'])) {
            return true;
        }
        $equation = false;

        //添加合伙人 start
        //合伙人推广
        if (!empty($withdrawalData['partner_id'])) {
            Log::info('合伙人推广:' . $withdrawalData['partner_id']);
            $partnerInfo = ShoppingGuideCashWithdrawal::query()
                ->where('worker_id', $withdrawalData['partner_id'])
                ->where('enterprise_id', $withdrawalData['enterprise_id'])
                ->where('user_type',ShoppingGuideCashWithdrawal::USER_TYPE_PARTNER)
                ->first();
            $this->saveShoppingGuideCashWithdrawalPartner($withdrawalData['enterprise_id'], $withdrawalData['partner_commission'], $withdrawalData['partner_id'], $partnerInfo, $cashType);
        }
        //end

        //推广导购和服务导购相等
        if ($withdrawalData['promotion_worker_id'] == $withdrawalData['service_shopping_guid_id']) {
            $equation = true;
        }
        //推广导购和服务导购相等逻辑处理
        if (!empty($equation)) {
            $workerInfo = ShoppingGuideCashWithdrawal::query()
                ->where('worker_id', $withdrawalData['promotion_worker_id'])
                ->where('enterprise_id', $withdrawalData['enterprise_id'])
                ->where('user_type',ShoppingGuideCashWithdrawal::USER_TYPE_WORKER)
                ->first();
            //总佣金 = 推广佣金+服务佣金
            $commissionAmount = $withdrawalData['promotion_commission'] + $withdrawalData['service_commission'];
            $this->saveShoppingGuideCashWithdrawal($withdrawalData['enterprise_id'], $commissionAmount, $withdrawalData['promotion_worker_id'], $workerInfo, $cashType);
            return true;
        }

        //推广导购
        if (!empty($withdrawalData['promotion_worker_id'])) {
            Log::info('推广导购:' . $withdrawalData['promotion_worker_id']);
            $workerInfo = ShoppingGuideCashWithdrawal::query()
                ->where('worker_id', $withdrawalData['promotion_worker_id'])
                ->where('enterprise_id', $withdrawalData['enterprise_id'])
                ->where('user_type',ShoppingGuideCashWithdrawal::USER_TYPE_WORKER)
                ->first();

            $this->saveShoppingGuideCashWithdrawal($withdrawalData['enterprise_id'], $withdrawalData['promotion_commission'], $withdrawalData['promotion_worker_id'], $workerInfo, $cashType);

        }

        //服务导购
        if (!empty($withdrawalData['service_shopping_guid_id'])) {
            Log::info('服务导购:' . $withdrawalData['service_shopping_guid_id']);
            $workerServiceInfo = ShoppingGuideCashWithdrawal::query()
                ->where('worker_id', $withdrawalData['service_shopping_guid_id'])
                ->where('enterprise_id', $withdrawalData['enterprise_id'])
                ->where('user_type',ShoppingGuideCashWithdrawal::USER_TYPE_WORKER)
                ->first();
            $this->saveShoppingGuideCashWithdrawal($withdrawalData['enterprise_id'], $withdrawalData['service_commission'], $withdrawalData['service_shopping_guid_id'], $workerServiceInfo, $cashType);
        }

        return true;
    }

    public function saveShoppingGuideCashWithdrawalPartner($enterpriseId, $commissionAmount, $partnerId, $partnerInfo, $cashType = ShoppingGuideCashWithdrawal::CASH_TYPE_PRE_WITHDRAWAL)
    {
        $updateField = ShoppingGuideCashWithdrawal::CASH_TYPE[$cashType];
        if (!empty($partnerInfo)) {
            $partnerInfo->$updateField += $commissionAmount;
            ShoppingGuideCashWithdrawal::query()
                ->where('id', $partnerInfo->id)
                ->where('enterprise_id', $enterpriseId)
                ->where('user_type',ShoppingGuideCashWithdrawal::USER_TYPE_PARTNER)
                ->update([$updateField => $partnerInfo->$updateField]);
            //如果可提现金额为负禁用账号
            $tmp = ShoppingGuideCashWithdrawal::CASH_TYPE[ShoppingGuideCashWithdrawal::CASH_TYPE_WITHDRAWAL];
            if($updateField == $tmp && $partnerInfo->$updateField < 0){
                (new PartnerMemberService())->negativeDisable($partnerId, $enterpriseId);
            }
        } else {

            Log::info('请求合伙人参数:'.json_encode([
                    'enterprise_id' => $enterpriseId,
                    'id'            => $partnerId,
                ]));
            $insert = [
                'worker_id'             => $partnerId,
                'enterprise_id'         => $enterpriseId,
                'user_type'            => ShoppingGuideCashWithdrawal::USER_TYPE_PARTNER,
                $updateField => $commissionAmount,
            ];
            ShoppingGuideCashWithdrawal::query()->create($insert);
        }
    }
    /**
     * 保存导购提现金额
     * @param $enterpriseId
     * @param $commissionAmount
     * @param $workerId
     * @param $workerInfo
     * @param int $cashType
     * @throws \Exception
     * @author Zhujiaqi
     * @time 2020-10-19 21:41:12
     */
    private function saveShoppingGuideCashWithdrawal($enterpriseId, $commissionAmount, $workerId, $workerInfo, $cashType = ShoppingGuideCashWithdrawal::CASH_TYPE_PRE_WITHDRAWAL) {

        $updateField = ShoppingGuideCashWithdrawal::CASH_TYPE[$cashType];

        if (!empty($workerInfo)) {
            $workerInfo->$updateField += $commissionAmount;
            ShoppingGuideCashWithdrawal::query()
                ->where('id', $workerInfo->id)
                ->where('enterprise_id', $enterpriseId)
                ->where('user_type',ShoppingGuideCashWithdrawal::USER_TYPE_WORKER)
                ->update([$updateField => $workerInfo->$updateField]);
        } else {

            Log::info('请求导购参数:'.json_encode([
                    'enterprise_id' => $enterpriseId,
                    'id'            => $workerId,
                ]));
            $storeWorker = StoreWorkerService::detail([
                'enterprise_id' => $enterpriseId,
                'worker_id'     => $workerId,
            ]);
            Log::info('导购返回参数:' . json_encode($storeWorker));
            $insert = [
                'worker_id'             => $workerId,
                'enterprise_id'         => $enterpriseId,
                'mem_uid'               => $storeWorker['mem_uid'] ?? 0,
                'user_type'            => ShoppingGuideCashWithdrawal::USER_TYPE_WORKER,
                $updateField => $commissionAmount,
            ];
            ShoppingGuideCashWithdrawal::query()->create($insert);
        }
    }


    /**
     * 获取订单结算商品列表信息
     * @param $enterpriseId string 企业ID
     * @param $orderNo string 订单编号
     * @param $skuIds array 商品sku_id集合
     * @param $fields array 查询的字段集合
     * @return mixed
     */
    public function getOrderGoodsSettlementList($enterpriseId = '', $orderNo = '', $skuIds = [], $fields = ['*']) {
        if (empty($orderNo) || empty($enterpriseId)) {
            return [];
        }
        $where = [];
        $where = array_merge($where, ['enterprise_id' => $enterpriseId]);
        $where = array_merge($where, ['order_no' => $orderNo]);

        return DistributionOrderGoodsSettlement::query()
            ->where($where)
            ->when(!empty($skuIds), function ($query) use ($skuIds) {
                return $query->whereIn('sku_id', $skuIds);
            })
            ->get($fields)
            ->toArray();
    }

    /**
     * 二级分销信息
     * @param $params
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/11/23 15:59
     */
    public function getPartnerSettings($params)
    {
        $detail = PartnerEnterpriseSettings::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->first();
        $detail = !empty($detail) ? $detail->toArray() : [];

        return $detail;
    }

    /**
     *
     * @param $orderInfo
     * @throws \Exception
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/11/24 15:25
     */
    public function sendSmsOrderStatus($orderInfo,$msg = '')
    {
        //合伙人是否开启接受短信
        if($this->sms_switch == PartnerEnterpriseMap::sms_switch_false){
            return true;
        }
        $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $orderInfo['enterprise_id']]);
        $partnerInfo = PartnerMember::query()
            ->where('id',$orderInfo['partner_id'])
            ->select('partner_mobile')
            ->first();
        $mobile = !empty($partnerInfo) ? $partnerInfo->partner_mobile : '';
        Log::info('分销结算发送合伙人短信:'.json_encode([
                'order_no' => $orderInfo['order_no'],
                'partner_id'            => $orderInfo['partner_id'],
            ]));
        $order_state = empty($msg) ? DistributionOrderSettlement::$status[$orderInfo['settlement_status']]: $msg;
        $data = [
            'brand' => $enterpriseInfo['data']['name'] ?? '',
            'amount' => $orderInfo['id'],
            'order_state' => $order_state,
            'app' => self::APP
        ];
        (new SendSmsService())->sendCode($mobile,$data,'distribution_order_status');
        return true;
    }

    /**
     * 合伙人余额变更发送短信
     * @param $orderInfo
     * @throws \Exception
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/11/24 15:25
     */
    public function sendSmsCashBalance($orderInfo)
    {
        //合伙人是否开启接受短信
        if($this->sms_switch == PartnerEnterpriseMap::sms_switch_false){
            return true;
        }
        $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $orderInfo['enterprise_id']]);
        $partnerInfo = PartnerMember::query()
            ->where('id',$orderInfo['partner_id'])
            ->select('partner_mobile')
            ->first();
        $mobile = !empty($partnerInfo) ? $partnerInfo->partner_mobile: '';
        Log::info('分销结算发送合伙人短信:'.json_encode([
                'order_no' => $orderInfo['order_no'],
                'partner_id'            => $orderInfo['partner_id'],
            ]));
        //查询余额情况
        $cashInfo = ShoppingGuideCashWithdrawal::query()
            ->where('enterprise_id',$orderInfo['enterprise_id'])
            ->where('worker_id',$orderInfo['partner_id'])
            ->where('user_type',ShoppingGuideCashWithdrawal::USER_TYPE_PARTNER)
            ->first();

        $data = [
            'brand' => $enterpriseInfo['data']['name'] ?? '',
            'amount' => isset($cashInfo->available_cash_amount) ? ($cashInfo->available_cash_amount/100) : 0,
            'es_amount' => isset($cashInfo->expect_cash_amount) ? ($cashInfo->expect_cash_amount/100)  : 0,
            'app' => self::APP
        ];
        (new SendSmsService())->sendCode($mobile,$data,'balance_status');
        return true;
    }

}