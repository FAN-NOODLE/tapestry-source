<?php

namespace App\Http\Controllers\V1\StaffCoupon;

use CrsSdk\Shop\StaffCoupon;
use App\Http\Controllers\V1\BaseController;

class StaffCouponUpdateController extends BaseController
{
    /**
     * 更新
     * @return mixed
     * @throws \Exception
     */
    public function staffCouponUpdate()
    {
        $all                  = $this->request->all();
        $all['enterprise_id'] = $this->enterpriseId;

        // todo 调用sdk
        $response = (new StaffCoupon())->staffCouponUpdate($all);
        $this->writeOperatorLog('访问内购券列表页');

        return $response;
    }
}