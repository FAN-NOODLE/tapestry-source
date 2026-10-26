<?php

namespace App\Http\Controllers\Service\V1\ExternalPay;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\ExternalPayUserOrderService;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class UserOrderController extends BaseController
{
    protected $order;

    /**
     * @throws \Exception
     */
    public function __construct(ExternalPayUserOrderService $order)
    {
        parent::__construct();
        $this->order = $order;
    }

    /**
     * 导购订单列表
     *
     * @return JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function orderList(): JsonResponse
    {
        $params  = $this->request->all();
        $rules   = [
            'enterprise_id' => 'required',
            'union_id'      => 'required',
            'order_state'   => 'sometimes|int',
            'page'          => 'sometimes|int',
            'size'          => 'sometimes|int',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'union_id.required'      => 'union_id不能为空',
            'order_state.int'        => '订单状态不正确',
        ];
        $this->validate($this->request, $rules, $message);

        //查询用户信息
        $member = $this->order->getUser($params);
        if (!$member) {
            error("用户不存在");
        }

        //查列表信息
        $params['external_userid'] = $member->openid;

        $res = $this->order->orderList($params);

        return successJson($res);
    }

    /**
     * 退款列表
     *
     * @return void
     */
    public function refundOrderList(): JsonResponse
    {
        $params  = $this->request->all();
        $rules   = [
            'enterprise_id' => 'required',
            'union_id'      => 'required',
            'page'          => 'sometimes|int',
            'size'          => 'sometimes|int',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'union_id.required'      => 'union_id不能为空',
        ];
        $this->validate($this->request, $rules, $message);

        //查询用户信息
        $member = $this->order->getUser($params);
        if (!$member) {
            error("用户不存在");
        }

        //查列表信息
        $params['external_userid'] = $member->openid;

        $res = $this->order->refundOrderList($params);

        return successJson($res);
    }

    /**
     * 订单详情
     *
     * @return JsonResponse
     * @throws ValidationException
     */
    public function orderDetail()
    {
        $params  = $this->request->all();
        $rules   = [
            'enterprise_id' => 'required',
            'union_id'      => 'required',
            'id'            => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'union_id.required'      => 'union_id不能为空',
            'id.required'            => '订单ID不能为空',
        ];
        $this->validate($this->request, $rules, $message);

        //查询用户信息
        $member = $this->order->getUser($params);
        if (!$member) {
            error("用户不存在");
        }

        $params['external_userid'] = $member->openid;

        $res = $this->order->orderDetail($params);

        return successJson($res);
    }

    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function refundOrderDetail()
    {
        $params  = $this->request->all();
        $rules   = [
            'enterprise_id' => 'required',
            'union_id'      => 'required',
            'id'            => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'union_id.required'      => 'union_id不能为空',
            'id.required'            => '退款订单ID不能为空',
        ];
        $this->validate($this->request, $rules, $message);

        //查询用户信息
        //查询用户信息
        $member = $this->order->getUser($params);
        if (!$member) {
            error("用户不存在");
        }

        $params['external_userid'] = $member->openid;

        $res = $this->order->refundOrderDetail($params);

        return successJson($res);
    }
}