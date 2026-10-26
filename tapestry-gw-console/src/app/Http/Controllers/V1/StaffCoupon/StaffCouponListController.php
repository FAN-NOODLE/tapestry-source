<?php


namespace App\Http\Controllers\V1\StaffCoupon;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StaffCoupon;

class StaffCouponListController extends BaseController
{
    /**
     * 列表
     * @return mixed
     * @throws \Exception
     */
    public function staffCouponList()
    {
        $all = $this->request->all();
        $all['enterprise_id'] = $this->enterpriseId;

        $response = (new StaffCoupon())->staffCouponList($all);
        $this->writeOperatorLog('访问内购券列表页');

        return $response;
    }
}