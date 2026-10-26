<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/8 0008
 * Time: 18:31
 */

namespace App\Http\Controllers\App\V1\Order;

use App\Common\Constant;
use App\Library\Message\FactoryMessage;
use CrsSdk\Shop\NewOrder;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Log;

class OrderDeliverController extends OrderBaseController
{
    /**
     * 确认发货接口
     *
     * @param NewOrder $newOrder
     *
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(NewOrder $newOrder)
    {
        $rules = [
            'order_id'          => 'required',
            'logistic_identity' => 'required',
            'logistic_no'       => 'required',
            'store_id'          => 'integer'
        ];
        $message = [
            'order_id.required'          => '订单id为空',
            'logistic_identity.required' => '物流公司标识为空',
            'logistic_no.required'       => '物流单号为空',
            'store_id.integer'           => '门店ID必须为数字'
        ];
        $input = $this->validate($this->request, $rules, $message);
        // 变更订单状态为确认发货参数
        $changeOrderStatus = [
            'event'                 => Constant::DELIVER_ORDER,
            'id'                    => (int)$input['order_id'],
            'enterprise_id'         => (int)$this->user['enterprise_id'],
            'operate_store_id'      => (int)$input['store_id'],
            'operate_id'            => (int)$this->user['choiceEnterprise']['admin_id'],
            'operate_name'          => $this->user['mobile'],
            'logistic_identity'     => $input['logistic_identity'],
            'logistic_no'           => $input['logistic_no']
        ];
        try {
            $response = $newOrder->changeOrderState($changeOrderStatus);
        } catch (\Exception $e) {
            Log::error("【门店助手】确认订单失败，请求参数:", $changeOrderStatus);
            error($e->getCode(), $e->getMessage());
        }
        return success($response);
    }

}