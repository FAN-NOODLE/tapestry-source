<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\Service\V1\EnterpriseExtend;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseExtend;
use App\Exceptions\Error;
use Illuminate\Http\Request;

class ChangeAllocationStatusController extends BaseController
{
    /**
     * 更新企业分配客户状态
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $status       = $request->post('status', 0);
        $enterpriseId = $request->post('enterprise_id');
        $couponStatus = $request->post('coupon_status', 0);

        if (empty($couponStatus) && empty($status)) {
            throw new Error(110002, '更新状态不能为空');
        }

        if (empty($enterpriseId)) {
            throw new Error(110003, '企业ID不能为空');
        }

        if ($status && !in_array($status, [EnterpriseExtend::OPEN, EnterpriseExtend::CLOSE])) {
            throw new Error(110004, '分配开关状态不正确');
        }

        if ($couponStatus && !in_array($couponStatus, [EnterpriseExtend::OPEN, EnterpriseExtend::CLOSE])) {
            throw new Error(120004, '红包开关状态不正确');
        }

        $enterpriseExtend = EnterpriseExtend::query()->where('enterprise_id', $enterpriseId)->first();

        $update = [];
        $status       && $update['allocation_status']  = $status;
        $couponStatus && $update['cash_coupon_status'] = $couponStatus;

        if ($enterpriseExtend) {
            EnterpriseExtend::query()->where('enterprise_id', $enterpriseId)->update($update);
        } else {
            EnterpriseExtend::query()->create([
                'enterprise_id'      => $enterpriseId,
                'allocation_status'  => $status ?: 2, // 默认关闭分配
                'cash_coupon_status' => $couponStatus ?: 2, // 默认关闭红包
            ]);
        }

        return success();
    }
}
