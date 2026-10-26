<?php

namespace App\Services;

use App\Common\Constant;
use App\Models\Distribution;
use App\Models\DistributionOrderSettlement;
use Illuminate\Support\Facades\Log;

class DistributionStatisticsService
{
    /**
     * 分销概览按时间筛选
     * @param int   $enterpriseId
     * @param array $workerList
     * @param int   $startTime
     * @param int   $endTime
     *
     * @return array
     */
    public function statisticData($enterpriseId = 0, $workerList = [], $startTime = 0, $endTime = 0)
    {
        $historyCommission = $historyCont = $allSettlementNum = $allCommission = $sectionCommission = 0;

        $workerIds = [];
        if (!empty($workerList)) {
            $workerIds = array_column($workerList, 'id');
        }

        $orderList = $this->getOrderList($enterpriseId, $workerIds, $startTime, $endTime);

        Log::info('按时间筛选返回数据:'.json_encode($orderList));

        $orderDayMoneyList = [];
        foreach ($orderList as $value) {
            $date = date('Y-m-d', $value['updated_at']);
            // 获取推广佣金以及服务佣金
            list($serviceCommission, $promotionCommission) = $this->getCommissions($workerIds, $value);

            if (empty($orderDayMoneyList[$date])) {
                $orderDayMoneyList[$date]['settlement_num']   = 1; //初始化数量
                $orderDayMoneyList[$date]['settlement_money'] = $promotionCommission + $serviceCommission;
            } else {
                $orderDayMoneyList[$date]['settlement_money'] += ($promotionCommission + $serviceCommission);
                $orderDayMoneyList[$date]['settlement_num']   += 1; //初始化数量
            }
            $sectionCommission += $promotionCommission + $serviceCommission;
        }

        //分销结算总比数
        $allSettlementNum = $historyCont + count($orderList);
        $allCommission    = $sectionCommission + $historyCommission;
        //初始化统计
        $detailData = $this->dailyDataFormat($startTime, $endTime);
        foreach ($detailData as &$item) {
            $orderDayInfo                     = $orderDayMoneyList[$item['time']] ?? [];
            $item['data']['settlement_money'] = $orderDayInfo['settlement_money'] ?? 0;
            $item['data']['settlement_num']   = $orderDayInfo['settlement_num'] ?? 0;
        }

        return [
            'settlement_num'      => $allSettlementNum,   // 总结算数量
            'settlement_money' => $allCommission,   // 总结算金额
            'detail_data'           => array_column($detailData, 'data','time'),
        ];
    }

    /**
     * 获取对应金额
     * @param array $workerIds
     * @param array $value
     *
     * @return array
     */
    protected function getCommissions($workerIds = [], $value = [])
    {
        $serviceCommission = $promotionCommission = 0;
        // 如果没有选择导购则有导购标识即可
        if (empty($workerIds)) {
            if (!empty($value['service_shopping_guid_id'])) {
                $serviceCommission = $value['service_commission'];
            }
            if (!empty($value['promotion_worker_id'])) {
                $promotionCommission = $value['promotion_commission'];
            }
        }

        // 如果悬在了导购则需要查询对应导购对应的佣金金额
        if (!empty($workerIds)) {
            if (!empty($value['service_shopping_guid_id']) && in_array($value['service_shopping_guid_id'], $workerIds)) {
                $serviceCommission = $value['service_commission'];
            }
            if (!empty($value['promotion_worker_id']) && in_array($value['promotion_worker_id'], $workerIds)) {
                $promotionCommission = $value['promotion_commission'];
            }
        }

        return [$serviceCommission, $promotionCommission];
    }

    /**
     * 按照导购筛选
     * @param int   $enterpriseId
     * @param array $workerList
     * @param int   $startTime
     * @param int   $endTime
     *
     * @return array
     */
    public function statisticDataByGuide($enterpriseId = 0, $workerList = [], $startTime = 0, $endTime = 0)
    {

        $workerIds = array_column($workerList, 'id');

        $orderList = $this->getOrderList($enterpriseId, $workerIds, $startTime, $endTime);

        $guideOrderLIst = $guideOrderPromotionCommission = $guideOrderServiceCommission = [];

        foreach ($orderList as $value) {
            $guideOrderLIst[$value['promotion_worker_id']][]      = $value['order_no'];
            $guideOrderLIst[$value['service_shopping_guid_id']][] = $value['order_no'];

            // 推广订单
            if (!empty($guideOrderPromotionCommission[$value['promotion_worker_id']])) {
                $guideOrderPromotionCommission[$value['promotion_worker_id']] += $value['promotion_commission'];
            } else {
                $guideOrderPromotionCommission[$value['promotion_worker_id']] = $value['promotion_commission'];
            }

            // 服务订单
            if (!empty($guideOrderServiceCommission[$value['service_shopping_guid_id']])) {
                $guideOrderServiceCommission[$value['service_shopping_guid_id']] += $value['service_commission'];
            } else {
                $guideOrderServiceCommission[$value['service_shopping_guid_id']] = $value['service_commission'];
            }
        }


        // 重新组装数据
        foreach ($workerList as &$value) {

            $settlementNum = $customerMoney = 0;
            if (!empty($guideOrderLIst[$value['id']])) {
                $settlementNum = count(array_unique($guideOrderLIst[$value['id']]));
            }

            $customerMoney = ($guideOrderPromotionCommission[$value['id']] ?? 0) + ($guideOrderServiceCommission[$value['id']] ?? 0);

            $data['settlement_money'] = $customerMoney;
            $data['settlement_num']   = $settlementNum;
            $data['id']               = $value['id'];
            $data['worker_no']        = $value['worker_no'];
            $data['name']             = $value['name'];
            $data['store_id']         = $value['store_id'];
            $data['avatar']           = $value['avatar'];
            $data['mem_uid']          = $value['mem_uid'];
            $value                    = $data;
        }

        return $workerList;
    }

    /**
     * 查询结算列表
     *
     * @param $enterpriseId
     * @param $workerIds
     * @param $startTime
     * @param $endTime
     *
     * @return array
     */
    public function getOrderList($enterpriseId, $workerIds, $startTime, $endTime)
    {

       return  DistributionOrderSettlement::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('settlement_status', DistributionOrderSettlement::ADOPT_SETTLE)
            ->when(!empty($workerIds), function ($query) use ($workerIds) {
                return $query->where(function ($query) use ($workerIds){
                    return $query->whereIn('promotion_worker_id', $workerIds)
                        ->orWhereIn('service_shopping_guid_id', $workerIds);
                });
            })
            ->when(!empty($startTime), function ($query) use ($startTime) {
                return $query->where('updated_at', '>=', $startTime);
            })
            ->when(!empty($endTime), function ($query) use ($endTime) {
                return $query->where('updated_at', '<=', $endTime);
            })
            ->get(['order_no', 'promotion_worker_id', 'service_shopping_guid_id', 'order_pay_time', 'promotion_commission', 'service_commission','updated_at'])
            ->toArray();

    }


    /**
     * dailyDataFormat
     *
     * @desc 初始化时间范围内数据
     *
     * @param $startTime
     * @param $endTime
     *
     * @return array
     * @date 2020-08-18
     */
    private function dailyDataFormat($startTime, $endTime)
    {
        $timeIndex = $startTime;
        $result    = [];
        $initData  = [
            'settlement_num'   => 0,  //顾客总数
            'settlement_money' => 0,  //浏览商城顾客数
        ];
        while ($timeIndex < $endTime) {
            $initData['time'] = date('Y-m-d', $timeIndex);
            $result[]         = [
                'time'     => date('Y-m-d', $timeIndex),
                'data'     => $initData,
                'day_time' => date('Ymd', $timeIndex)
            ];
            $timeIndex        += Constant::ONE_DAY_SEC;
        }
        return $result;
    }
}
