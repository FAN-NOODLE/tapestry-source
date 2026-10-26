<?php
/**
 * Created by PhpStorm.
 * User: meijie
 * Date: 2019-06-04
 * Time: 15:42
 */

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;
use Illuminate\Support\Facades\Log;

/**
 * 售前订单,用户申请退款,订单人工介入,同意退款调用此接口直接进行退款
 * Class BeforeSalesDirectRefundController
 * @package App\Http\Controllers\App\V1\Order
 */
class AfterSalesDirectRefundController extends BaseController
{
    public function index()
    {
        // 获取订单号,不存在为null
        $order_no = $this->request->input('order_no', null);

        // 未传递单号抛异常
        if (is_null($order_no)) {
            error(320255);
        }

        try {
            // 请求shop进行订单恢复正常流程
            (new Order())->refundOrder(['order_no' => $order_no, 'status' => true, 'enterprise_id' => $this->enterpriseId]);

            // 记录管理员操作
            $this->writeOperatorLog("订单号:{$order_no} 售后订单同意退款,发起微信退款");

            // 响应数据
            return success(['data' => [], 'msg' => '进行微信退款请求成功']);

        } catch (\Exception $e) {
            Log::error('订单手动退款:请求异常', [
                'line'    => $e->getLine(),
                'file'    => $e->getFile(),
                'code'    => $e->getCode(),
                'message' => $e->getMessage(),
                'params'  => ['order_no' => $order_no, 'enterprise_id' => $this->enterpriseId],
            ]);

            return success(['data' => [], 'msg' => '进行微信退款请求失败']);
        }
    }
}
