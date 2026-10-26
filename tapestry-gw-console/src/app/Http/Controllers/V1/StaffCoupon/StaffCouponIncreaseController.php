<?php

namespace App\Http\Controllers\V1\StaffCoupon;

use CrsSdk\Shop\StaffCoupon;
use App\Http\Controllers\V1\BaseController;

class StaffCouponIncreaseController extends BaseController
{
    /**
     * 新增
     * @return mixed
     * @throws \Exception
     */
    public function staffCouponIncrease()
    {
        $all                  = $this->request->all();
        $all['enterprise_id'] = $this->enterpriseId;

        $response = (new StaffCoupon())->staffCouponIncrease($all);
        $this->writeOperatorLog('访问内购券列表页');

        return $response;
    }
}