<?php

namespace App\Http\Controllers\V1\AppCenter\MerchantCoupon;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MerchantCoupon;

class BindableListController extends BaseController
{
    public function index()
    {
        $input                  = $this->request->only([
            'page',
            'size',
            'self',
            'coupon_name',
            'coupon_stock_id',
            'app_name',
        ]);
        $input['enterprise_id'] = $this->enterpriseId;
        $list                   = (new MerchantCoupon())->bindableList($input);
        return success($list);
    }
}