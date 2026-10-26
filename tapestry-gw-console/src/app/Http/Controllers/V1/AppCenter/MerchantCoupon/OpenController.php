<?php

namespace App\Http\Controllers\V1\AppCenter\MerchantCoupon;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MerchantCoupon;

class OpenController extends BaseController
{

    public function index()
    {
        $input = [
            'id' => $this->enterpriseId
        ];
        $list  = (new MerchantCoupon())->open($input);
        return success($list);
    }
}
