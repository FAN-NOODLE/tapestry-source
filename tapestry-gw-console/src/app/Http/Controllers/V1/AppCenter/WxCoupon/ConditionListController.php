<?php
namespace App\Http\Controllers\V1\AppCenter\WxCoupon;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\WxCoupon;

class ConditionListController extends BaseController
{
    public function index()
    {
        $input['enterprise_id'] = $this->enterpriseId;

        $list = (new WxCoupon())->wxCouponConditionList($input);

        return success($list);
    }
}