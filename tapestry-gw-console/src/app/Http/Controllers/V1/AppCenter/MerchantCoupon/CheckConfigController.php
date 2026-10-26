<?php

namespace App\Http\Controllers\V1\AppCenter\MerchantCoupon;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MerchantCoupon;

class CheckConfigController extends BaseController
{

    public function index()
    {
        $input = [
            'id' => $this->enterpriseId
        ];
        $list  = (new MerchantCoupon())->check($input);
        return success($list);
    }
}
