<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/8 0008
 * Time: 18:08
 */

namespace App\Http\Controllers\App\V1\Order;

use App\Events\OrderServiceNotice;
use App\Exceptions\Error;
use CrsSdk\Shop\NewOrder;
use Illuminate\Support\Facades\Auth;

class MemberSignController extends OrderBaseController
{
    /**
     * 服务人接单
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        // 验证输入参数
        $rules   = [
            'store_id' => 'required',
            'order_no'     => 'required'
        ];
        $message = [
            'store_id.required' => '店铺ID不能为空',
            'order_no.required'     => '订单号不能为空',
        ];
        $this->validate($this->request, $rules, $message);

        $storeId = $this->request->post('store_id', 0);
        $orderNo    = $this->request->post('order_no', '');

        $enterpriseId = Auth::user()['enterprise_id'];
        $workerDetail = $this->storeWorkerDetail($enterpriseId, $storeId);

        $params = [
            'enterprise_id'    => $enterpriseId,
            'service_store_id' => $storeId,
            'order_no'         => $orderNo,
            'service_staff_id' => $workerDetail['id']
        ];


        try {
            (new NewOrder())->receiveServiceOrder($params);
            event(new OrderServiceNotice($enterpriseId, $workerDetail['id'], $orderNo));
        } catch (Error $error) {
            \error($error->getCode(), $error->getMessage());
        }

        return success([]);
    }
}
