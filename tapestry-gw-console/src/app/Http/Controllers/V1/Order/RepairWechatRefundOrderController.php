<?php

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Manual;

class RepairWechatRefundOrderController extends BaseController
{
    /**
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $params = $this->request->all();

        // 退款报"订单已全额退款"的订单（运营在微信后台操作了全额退款）
        $data = (new Manual())->repairWechatRefundOrder($params);

        return success($data);
    }
}
