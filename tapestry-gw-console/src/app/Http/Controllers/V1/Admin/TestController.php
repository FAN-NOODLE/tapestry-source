<?php

namespace App\Http\Controllers\V1\Admin;

use App\Http\Controllers\V1\BaseController;

class TestController extends BaseController
{
    public function index()
    {
        return success(
            [
                'has_shop_auth' => $this->hasShopAuth,
                'store_ids'     => $this->storeIds,
            ]
        );
    }
}
