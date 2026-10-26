<?php

namespace App\Http\Controllers\V1\Order;

use CrsSdk\Shop\Order;
use App\Http\Controllers\V1\BaseController;

class RefundGoodsEditController extends BaseController
{
    /**
     * 售后退款商品修改
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();

        $input['enterprise_id'] = $this->enterpriseId;
        $data = (new Order())->refundGoodsEdit($input);

        $this->writeOperatorLog('修改售后商品');
        // 已改
        return success();
    }
}
