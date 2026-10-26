<?php

namespace App\Http\Controllers\Service\V1\EnterprisePayment;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterprisePayment;
use App\Models\EnterprisePaymentRandom;
use App\Models\EnterprisePaymentRecord;
use App\Models\MemberUserRelation;
use App\Services\EnterprisePaymentService;
use Exception;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use Illuminate\Support\Facades\DB;

class SendRedPackController extends BaseController
{

    /**
     * 发送红包
     */
    public function index()
    {
        $userId = request()->post('user_id', 0);
        $unionId = request()->post('union_id', '');
        $openId = request()->post('open_id', '');
        $memUid = request()->post('mem_uid', '');
        $activityId = request()->post('activity_id', 0);
        $enterpriseId = request()->post('enterprise_id', 0);

        if (empty($userId) || empty($memUid) || empty($activityId) || empty($enterpriseId)) {
            throw new Error(120001, '参数错误',210);
        }

        $activity = EnterprisePayment::query()
            ->where('id', $activityId)
            ->first();

        if ($activity && $activity->is_sale == EnterprisePayment::IS_SALE_DEFAULT) {
            throw new Error(120002, '活动还未开始，请稍后再开红包',210);
        }

        if (!$activity || $activity->is_sale != EnterprisePayment::IS_SALE_UP) {
            throw new Error(120002, '活动太火爆，红包已经抢完了',210);
        }

        $memUserRelation = MemberUserRelation::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('union_id', $unionId)
            ->where('delete_type', 0)
            ->orderByDesc('relation_at')
            ->first(['relation_at']);
        if (!$memUserRelation) {
            throw new Error(120003, '抱歉，您暂无权限开启该红包',210);
        }

        if ($activity->is_old_user_receive == EnterprisePayment::IS_OLD_USER_RECEIVE_NO && $memUserRelation->relation_at < $activity->created_at->timestamp) {
            throw new Error(120004, '抱歉，该活动仅限新用户参加哦',210);
        }

        $lockKey = EnterprisePayment::SEND_PACK_BY_USER_KEY . $activityId . '_' . intval($userId);

        // 查询锁
        $getCache = Redis::get($lockKey);

        // 如果已经发起直接跳出
        if ($getCache) {
            throw new Error(120006, '每人限领一个红包哦',210);
        }

        Redis::setex($lockKey, 600, time());

        $recordCountKey = EnterprisePayment::SEND_PACK_KEY . $activityId;
        // 查询锁
        $recordCount = Redis::get($recordCountKey);
        if (is_null($recordCount)) {
            $recordCount = (new EnterprisePaymentService())->setCount($activityId, $recordCountKey);
        }

        if ($activity->quantity <= $recordCount) {
            EnterprisePayment::query()
                ->where('id', $activityId)
                ->update(['is_sale' => EnterprisePayment::IS_SALE_OVER]);

            throw new Error(120002, '活动太火爆，红包已经抢完了',210);
        } else {
            $recordCount = Redis::incr($recordCountKey);
            if ($activity->quantity < $recordCount) {
                Redis::set($recordCountKey, $activity->quantity);
                EnterprisePayment::query()
                    ->where('id', $activityId)
                    ->update(['is_sale' => EnterprisePayment::IS_SALE_OVER]);

                throw new Error(120002, '活动太火爆，红包已经抢完了',210);
            }
        }

        $amount = $activity->fixed_amount;
        $randomId = 0;
        if ($activity->type == EnterprisePayment::RANDOM_AMOUNT_TYPE) {
            // 随机金额放进redis列表
            $lkey = EnterprisePayment::RANDOM_LIST_KEY . $activityId;

            $llen = Redis::llen($lkey);
            if ($llen == 0) {
                try {
                    (new EnterprisePaymentService())->setRandomList($activityId, $lkey);
                }catch (Exception $exception) {
                    $this->delCache($lockKey, $recordCountKey, $activity);
                    throw new Error($exception->getCode(), $exception->getMessage(),210);
                }
            }

            $randomDataString = Redis::lpop($lkey);
            $randomData = explode('-', $randomDataString);

            $randomId = $randomData[0] ?? 0;
            $amount = $randomData[1] ?? 0;
        }

        if ($amount == 0) {
            throw new Error(120002, '活动太火爆，红包已经抢完了',210);
        }

        try {
            DB::beginTransaction();

            $params = [
                'enterprise_id' => $enterpriseId,
                'openid' => $openId,
                'bill_no' => EnterprisePaymentRecord::generateBillNo(),
                'total_amount' => $amount,
                'remark' => $activity->name
            ];

            $result = (new \CrsSdk\MicroPay\EnterprisePayment())->sendRedPack($params);

            // 增加红包发送记录
            if (!empty($result['data']) && isset($result['data']['result_code']) && $result['data']['result_code'] == 'SUCCESS') {

                $billNo = $result['data']['partner_trade_no'];
                $detailId = $result['data']['payment_no'];

                EnterprisePaymentRecord::query()->create([
                    'enterprise_id' => $enterpriseId,
                    'enterprise_payment_id' => $activityId,
                    'bill_no' => $billNo,
                    'detail_id' => $detailId,
                    'enterprise_payment_random_id' => $randomId,
                    'amount' => $amount,
                    'mem_uid' => $memUid,
                    'user_id' => $userId,
                    'union_id' => $unionId,
                    'status' => EnterprisePaymentRecord::RECORD_STATUS_SENDED,
                    'send_time' => time(),
                    'receive_refund_time' => time()
                ]);
                if ($randomId) {
                    EnterprisePaymentRandom::query()
                        ->where('id', $randomId)
                        ->update(['is_receive' => EnterprisePaymentRandom::IS_RECEIVE_YES]);
                }

                // 红包发完
                if ($activity->quantity <= $recordCount) {
                    EnterprisePayment::query()
                        ->where('id', $activityId)
                        ->update(['is_sale' => EnterprisePayment::IS_SALE_OVER]);
                }

                DB::commit();
            } else {
                    throw new Error(120007, '开启失败，再次点击试试',210);
            }

        } catch (Exception $exception) {
            DB::rollBack();

            $this->delCache($lockKey, $recordCountKey, $activity, $lkey ?? '', $randomDataString ?? '');

            Log::info("发送红包异常:".$exception->getCode().$exception->getMessage());
            throw new Error($exception->getCode(), $exception->getMessage());
        }

        $return = [];
        $return['name'] = $activity->name ?? '';
        $return['qr_code_url'] = $activity->qr_code_url ?? '';
        $return['send_name'] = $activity->send_name ?? '';
        $return['wishing'] = $activity->wishing ?? '';
        $return['amount'] = $amount;
        $return['receive_refund_time'] = time();

        return success(['data' => $return]);
    }

    /**
     * @param $lockKey
     * @param $recordCountKey
     * @param $recordCount
     * @param $activity
     * @param string $lkey
     * @param string $randomDataString
     * @return bool
     */
    public function delCache($lockKey, $recordCountKey, $activity, $lkey = '', $randomDataString = '')
    {
        // 删除锁
        Redis::del([$lockKey]);
        $recordCount = Redis::decr($recordCountKey);
        if ($activity->quantity > $recordCount) {
            EnterprisePayment::query()
                ->where('id', $activity->id)
                ->update(['is_sale' => EnterprisePayment::IS_SALE_UP]);
        }
        if ($activity->type == EnterprisePayment::RANDOM_AMOUNT_TYPE && $lkey && $randomDataString) {
            Redis::lpush($lkey, $randomDataString);
        }

        return true;
    }
}
