<?php

namespace App\Http\Controllers\Html\V1\CashCoupon;

use App\Exceptions\Error;
use App\Library\Uc\External;
use App\Models\CashCoupon;
use App\Http\Controllers\Html\V1\BaseController;
use App\Models\CashCouponMember;
use App\Models\CashCouponRecord;
use App\Models\EnterpriseExtend;
use CrsSdk\MicroPay\CashCoupon as CashCouponSdk;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class SendRedPackController extends BaseController
{
    const TOTAL_NUM = 1; // 红包发放人数,必须是1

    /**
     * 发送现金红包
     */
    public function index()
    {
        $userId = request()->post('user_id', '');
        $cashCouponId = request()->post('cash_coupon_id', 0);

        if (empty($userId)) {
            throw new Error(110001, '参数错误');
        }

        if (empty($cashCouponId)) {
            throw new Error(110001, '参数错误');
        }

        $enterpriseId = $this->enterpriseId;

        // 判断红包管理是否开启
        $enterpriseExtend = EnterpriseExtend::query()->where('enterprise_id', $enterpriseId)->first(['cash_coupon_status']);

        if (($enterpriseExtend->cash_coupon_status ?? 0) != EnterpriseExtend::OPEN) {
            throw new Error(110001, '红包管理未开启，请联系管理员。');
        }

        // 查询红包活动
        $cashCoupon = CashCoupon::query()
            ->where('id', $cashCouponId)
            ->where('enterprise_id', $enterpriseId)
            ->where('is_sale', CashCoupon::IS_ON_SALE)
            ->where('is_deleted', CashCoupon::NOT_DELETED)
            ->first();

        if (empty($cashCoupon)) {
            throw new Error(110002, '红包活动不存在');
        }

        // 查询红包活动关联员工
        $cashCouponMember = CashCouponMember::query()
            ->where('cash_coupon_id', $cashCouponId)
            ->where('enterprise_id', $enterpriseId)
            ->where('object', $this->memUid)
            ->where('is_leave', CashCouponMember::MEMBER_IS_NORMAL)
            ->first();

        if (empty($cashCouponMember)) {
            throw new Error(110003, '没有权限发送此红包');
        }

        // 判断发送用户是否是企业微信外部联系人(wo开头)
        $start = substr($userId, 0, 2);

        if ($start == 'wo') {
            throw new Error(900012, '该用户为企业微信外部联系人，不可使用红包功能，如需使用，请让用户使用个人微信扫码后操作。');
        }

        $lockKey = 'send_red_pack_'.$cashCouponId;

        // 查询锁
        $getCache = Redis::get($lockKey);

        // 检查是否已经超出时间 超出时间释放锁（防止死锁）
        if (!empty($getCache) && $getCache < time()) {
            Redis::del([$lockKey]);
        }

        // 开始加锁
        $setCache = Redis::setnx($lockKey, time() + 5);

        // 如果已经发起直接跳出
        if (!$setCache) {
            throw new Error(900013, '当前同时发送人数过多，技术GG扛不住了，请稍后重试。');
        }

        $cashCouponRecordCount = CashCouponRecord::query()
            ->where('cash_coupon_id', $cashCouponId)
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('status', [CashCouponRecord::SENDING, CashCouponRecord::RECEIVED])
            ->count();

        if ($cashCoupon->quantity <= $cashCouponRecordCount) {
            // 删除锁
            Redis::del([$lockKey]);

            throw new Error(900011, '红包已全部发完，如需增加红包数量，请联系管理员。');
        }

        try {
            $openid = (new External())->convertToOpenid($this->domain, ['ecUserid' => $userId]);

            $params = [
                'enterprise_id' => $enterpriseId,
                'openid' => $openid,
                'bill_no' => CashCouponRecord::generateBillNo(),
                'send_name' => $cashCoupon->send_name,
                'total_amount' => $cashCoupon->amount,
                'total_num' => self::TOTAL_NUM,
                'wishing' => $cashCoupon->wishing,
                'act_name' => $cashCoupon->cash_coupon_name,
                'remark' => $cashCoupon->remark
            ];

            $result = (new CashCouponSdk())->sendRedPack($params);

            // 增加红包发送记录
            if (!empty($result['data'])) {
                $billNo = $result['data']['bill_no'];
                $detailId = $result['data']['detail_id'];

                CashCouponRecord::query()->create([
                    'enterprise_id' => $enterpriseId,
                    'cash_coupon_id' => $cashCouponId,
                    'bill_no' => $billNo,
                    'detail_id' => $detailId,
                    'mem_uid' => $this->memUid,
                    'ec_userid' => $userId,
                    'open_id' => $openid,
                    'status' => CashCouponRecord::SENDING,
                    'send_time' => time(),
                ]);
            }
            // 删除锁
            Redis::del([$lockKey]);
        } catch (\Exception $exception) {
            // 删除锁
            Redis::del([$lockKey]);
            Log::info("发送红包异常:".$exception->getCode().$exception->getMessage());
            throw new Error($exception->getCode(), $exception->getMessage());
        }

        return success();
    }
}
