<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/3
 * Time: 18:33
 */

namespace App\Services;


use App\Exceptions\Error;
use App\Models\EnterprisePayment;
use App\Models\EnterprisePaymentBudget;
use App\Models\EnterprisePaymentRandom;
use App\Models\EnterprisePaymentRecord;
use Illuminate\Support\Facades\Redis;

class EnterprisePaymentService
{
    /**
     * 获取当前活动统计数量
     *
     * @param array $detail
     *
     * @return array
     */
    public function getStatisticalNum($detail = [])
    {
        $data = [
            'send_num'      => 0,   // 发送数量
            'receive_num'   => 0,   // 领取数量
            'surplus_num'   => 0,   // 剩余数量
            'send_money'    => 0,   // 发送金额
            'receive_money' => 0,   // 领取金额
            'surplus_money' => 0    // 剩余金额
        ];

        if (empty($detail)) {
            return $data;
        }

        // 发送数量
        $sendList = EnterprisePaymentRecord::query()
            ->where('enterprise_payment_id', $detail['id'])
            ->get(['amount'])
            ->toArray();

        // 发送数量
        $sendNum = count($sendList);

        // 如果是随机金额
        $sendMoney = array_sum(array_column($sendList, 'amount'));

        $data = [
            'send_num'      => $sendNum,
            'receive_num'   => $sendNum,
            'surplus_num'   => $detail['quantity'] - $sendNum,
            'send_money'    => $sendMoney,
            'receive_money' => $sendMoney,
            'surplus_money' => $detail['enterprise_payment_amount'] - $sendMoney
        ];

        return $data;
    }

    /** 终止活动
     *
     * @param $detail
     */
    public function overActivity($detail)
    {
        // 发送数量
        $sendAmount = EnterprisePaymentRecord::query()
            ->where('enterprise_payment_id', $detail['id'])
            ->sum('amount');

        // 查询企业充值表信息
        $info = EnterprisePaymentBudget::query()
            ->where('enterprise_id', $detail['enterprise_id'])
            ->orderByDesc('id')
            ->first();

        // 查询主表信息
        $paymentInfo = EnterprisePayment::query()
            ->where('id', $detail['id'])
            ->first();

        // 获取剩余金额
        $amount = $paymentInfo['enterprise_payment_amount'] - $sendAmount;

        if ($amount > 0) {
            $BudgetData = [
                'enterprise_id' => $detail['enterprise_id'],
                'number'        => $detail['number'],
                'amount'        => $amount,
                'balance'       => $info['balance'] + $amount,
                'income_type'   => EnterprisePaymentBudget::INCOME_TYPE_IN,
                'type'          => EnterprisePaymentBudget::ACTIVITY_REFUND_TYPE
            ];

            EnterprisePaymentBudget::query()
                ->create($BudgetData);
        }
    }

    /**
     * 设置红包活动缓存
     * @param $activityId
     * @param $recordCountKey
     * @return int
     */
    public function setCount($activityId, $recordCountKey)
    {
        $redis = Redis::connection();
        $setCountKey = 'enterprise_payment_set_count_' . $activityId;
        $setRes = $redis->setnx($setCountKey, time());
        if (!$setRes) {
            throw new Error(120007, '开启失败，再次点击试试',210);
        }

        $recordCount = EnterprisePaymentRecord::query()
            ->where('enterprise_payment_id', $activityId)
            ->where('status', EnterprisePaymentRecord::RECORD_STATUS_SENDED)
            ->count();
        Redis::set($recordCountKey, $recordCount);

        Redis::del([$setCountKey]);

        return $recordCount ?? 0;
    }

    /**
     * 随机金额缓存列表
     * @param $activityId
     * @param $lkey
     * @return mixed
     */
    public function setRandomList($activityId, $lkey)
    {
        $redis = Redis::connection();

        $setCountKey = 'enterprise_payment_set_random_' . $activityId;
        $setRes = $redis->setnx($setCountKey, time());

        if (!$setRes) {
            throw new Error(120007, '开启失败，再次点击试试',210);
        }

        $randomInfo = EnterprisePaymentRandom::query()
            ->where('enterprise_payment_id', $activityId)
            ->where('is_receive', EnterprisePaymentRandom::IS_RECEIVE_NO)
            ->get(['id', 'amount'])
            ->toArray();
        $list = [];
        $llen = Redis::llen($lkey);
        if ($llen == 0 && $randomInfo) {
            foreach ($randomInfo as $value) {
                $list[] = $value['id'] . '-' . $value['amount'];
            }

            if ($list) {
                $llen = Redis::lpush($lkey, $list);
            }
        }

        Redis::del([$setCountKey]);

        return $llen;
    }

}
