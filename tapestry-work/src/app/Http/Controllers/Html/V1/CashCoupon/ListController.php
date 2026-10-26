<?php

namespace App\Http\Controllers\Html\V1\CashCoupon;

use App\Exceptions\Error;
use App\Models\CashCoupon;
use App\Http\Controllers\Html\V1\BaseController;
use App\Models\CashCouponMember;
use App\Models\CashCouponRecord;
use App\Models\EnterpriseExtend;
use CrsSdk\MicroPay\CashCoupon as CashCouponSdk;
use Illuminate\Support\Facades\DB;

class ListController extends BaseController
{
    /**
     * 现金红包活动列表
     */
    public function index()
    {
        $size = request()->input('size', 20);
        $page = request()->input('page', 1);

        // 企业ID
        $enterpriseId = $this->enterpriseId;
        // 员工ID
        $memUid = $this->memUid;

        // 判断红包管理是否开启
        $enterpriseExtend = EnterpriseExtend::query()->where('enterprise_id', $enterpriseId)->first(['cash_coupon_status']);

        if (($enterpriseExtend->cash_coupon_status ?? 0) != EnterpriseExtend::OPEN) {
            return success(['data' => []]);
        }

        // 员工关联的红包活动ID集合
        $cashCouponId = CashCouponMember::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('object', $memUid)
            ->where('is_leave', CashCouponMember::MEMBER_IS_NORMAL)
            ->pluck('cash_coupon_id')
            ->toArray();

        if (empty($cashCouponId)) {
            return success(['data' => []]);
        }

        // 查询红包活动发送记录
        $cashCouponRecord = CashCouponRecord::query()
            ->select(DB::raw('count(cash_coupon_id) as count, cash_coupon_id'))
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('cash_coupon_id', $cashCouponId)
            ->whereIn('status', [CashCouponRecord::SENDING, CashCouponRecord::RECEIVED])
            ->groupBy('cash_coupon_id')
            ->get()
            ->toArray();

        $cashCouponRecord = array_column($cashCouponRecord, null, 'cash_coupon_id');

        // 查询红包活动列表
        $cashCoupon = CashCoupon::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('is_sale', CashCoupon::IS_ON_SALE)
            ->where('is_deleted', CashCoupon::NOT_DELETED)
            ->whereIn('id', $cashCouponId)
            ->orderBy('id', 'DESC')
            ->get()
            ->toArray();

        // 排序
        if ($cashCoupon) {
            $valid = [];
            $invalid = [];
            foreach ($cashCoupon as $value) {
                $useQuantity = ($cashCouponRecord[$value['id']] ?? [])['count'] ?? 0;
                if ($value['quantity'] > $useQuantity) {
                    $value['remain_quantity'] = $value['quantity'] - $useQuantity;
                    $valid[] = $value;
                } else {
                    $value['remain_quantity'] = 0;
                    $invalid[] = $value;
                }
            }
            $cashCoupon = array_merge($valid, $invalid);
        }

        $start = ($page - 1) * $size;
        $totals = count($cashCoupon);
        $countPage = ceil($totals / $size);
        $pageData = array_slice($cashCoupon, $start, $size);

        return success([
            'data' => $pageData,
            'meta' => [
                'total' => (int)$totals,
                'last_page' => (int)$countPage,
                'size'  => (int)$size,
                'page'  => (int)$page
            ]
        ]);
    }
}
