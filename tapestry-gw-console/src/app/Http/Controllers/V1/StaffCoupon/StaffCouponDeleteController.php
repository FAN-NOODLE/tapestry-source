<?php

namespace App\Http\Controllers\V1\StaffCoupon;

use CrsSdk\Shop\StaffCoupon;
use App\Http\Controllers\V1\BaseController;

class StaffCouponDeleteController extends BaseController
{
    /**
     * 删除
     * @return mixed
     * @throws \Exception
     */
    public function staffCouponDelete()
    {
        $all                  = $this->request->all();
        $all['enterprise_id'] = $this->enterpriseId;

        $response = (new StaffCoupon())->staffCouponDelete($all);
        $this->writeOperatorLog('访问内购券列表页');

        return $response;
    }
}