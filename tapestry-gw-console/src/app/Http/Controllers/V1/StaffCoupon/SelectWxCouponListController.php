<?php

namespace App\Http\Controllers\V1\StaffCoupon;

use Exception;
use CrsSdk\Shop\StaffCoupon;
use App\Http\Controllers\V1\BaseController;

class SelectWxCouponListController extends BaseController
{
    /**
     * 删除
     * @return mixed
     * @throws Exception
     */
    public function selectWxCouponList()
    {
        $all                  = $this->request->all();
        $all['enterprise_id'] = $this->enterpriseId;

        $response = (new StaffCoupon())->selectWxCouponList($all);
        $this->writeOperatorLog('查询优惠券列表');

        return $response;
    }
}