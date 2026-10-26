<?php

namespace App\Http\Controllers\Service\Cdp;

use App\Http\Controllers\Service\V1\BaseController;
use App\Sdk\CDP\Coach\CustomerInterestedProducts;
use Exception;

class CdpController extends BaseController
{


    /**
     * 获取CDP推荐商品
     *
     * @return array
     * @throws Exception
     */
    public function getCustomerInterestedProducts(): array
    {
        $input  = $this->request->all();
        return (new CustomerInterestedProducts())->get($input);
    }

}