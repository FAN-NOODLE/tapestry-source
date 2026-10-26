<?php

namespace App\Http\Controllers\V1\StaffCoupon;

use CrsSdk\Shop\StaffCoupon;
use App\Http\Controllers\V1\BaseController;

class StaffCouponChangeStatusController extends BaseController
{
    /**
     *  上下架
     * @return mixed
     * @throws \Exception
     */
    public function staffCouponChangeStatus()
    {
        $all                  = $this->request->all();
        $all['enterprise_id'] = $this->enterpriseId;

        $response = (new StaffCoupon())->staffCouponChangeStatus($all);
        $this->writeOperatorLog('访问内购券列表页');

        return $response;
    }
}