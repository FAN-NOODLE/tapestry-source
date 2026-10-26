<?php

namespace App\Http\Controllers\V1\Home;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\Enterprise;

class BuyNoticeController extends BaseController
{
    /**
     * 购买提醒
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $res = (new Enterprise())->buyNotice([
            'enterprise_id' => $this->enterpriseId
        ]);
        return success($res);

    }
}