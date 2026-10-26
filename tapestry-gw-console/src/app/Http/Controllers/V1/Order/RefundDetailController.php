<?php

namespace App\Http\Controllers\V1\Order;

use CrsSdk\Shop\Order;
use App\Http\Controllers\V1\BaseController;

class RefundDetailController extends BaseController
{
    /**
     * 售后详情接口
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->input();

        $params = [
            'enterprise_id'   => $this->enterpriseId,
            'refund_apply_id' => $input['refund_apply_id']
        ];
        $response = (new Order())->refundDetail($params);
        return success($response);
    }
}
