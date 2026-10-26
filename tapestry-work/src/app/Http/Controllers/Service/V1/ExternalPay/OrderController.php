<?php

namespace App\Http\Controllers\Service\V1\ExternalPay;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\ExternalPayOrderService;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class OrderController extends BaseController
{
    protected $order;

    /**
     * @throws Exception
     */
    public function __construct(ExternalPayOrderService $order)
    {
        parent::__construct();
        $this->order = $order;
    }

    /**
     * 导购订单列表
     *
     * @return JsonResponse
     * @throws ValidationException
     */
    public function orderList(): JsonResponse
    {
        $params  = $this->request->all();
        $rules   = [
            'enterprise_id'  => 'required',
            'mobile'         => 'required',
            'order_state'    => 'sometimes|int',
            'sales_order_no' => 'sometimes',
            'page'           => 'sometimes|int',
            'size'           => 'sometimes|int',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'mobile.required'        => '手机号不能为空',
            'order_state.int'        => '订单状态不正确',
        ];
        $this->validate($this->request, $rules, $message);

        //查询SM信息
        $member = $this->order->getMember($params);
        if (!$member) {
            error("导购不存在");
        }

        //查列表信息
        $params['member_id'] = $member->id;

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
            'enterprise_id'  => 'required',
            'mobile'         => 'required',
            'sales_order_no' => 'sometimes',
            'page'           => 'sometimes|int',
            'size'           => 'sometimes|int',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'mobile.required'        => '手机号不能为空',
        ];
        $this->validate($this->request, $rules, $message);

        //查询SM信息
        $member = $this->order->getMember($params);
        if (!$member) {
            error("导购不存在");
        }

        //查列表信息
        $params['member_id'] = $member->id;

        $res = $this->order->refundOrderList($params);

        return successJson($res);
    }

    /**
     * 导购发货
     *
     * @return JsonResponse
     * @throws ValidationException
     */
    public function deliverOrder(): JsonResponse
    {
        $params  = $this->request->all();
        $rules   = [
            'enterprise_id'  => 'required',
            'mobile'         => 'required',
            'id'             => 'required',
            'logistics_no'   => 'required_if:delivery_type,2',
            'logistics_name' => 'required_if:delivery_type,2',
            'logistics_com'  => 'required_if:delivery_type,2',
            'tickets'        => 'required',
            'delivery_time'  => 'required',
            'sales_order_no' => 'required',
        ];
        $message = [
            'enterprise_id.required'  => '企业ID不能为空',
            'mobile.required'         => '手机号不能为空',
            'id.required'             => '订单ID不能为空',
            'logistics_no.required'   => '发货单号不能为空',
            'logistics_name.required' => '物流公司不能为空',
            'logistics_com.required'  => '物流编码不能为空',
            'tickets.required'        => '请上传小票',
            'delivery_time.required'  => '录单时间不能为空',
            'sales_order_no.required' => '销售单号不能为空',
        ];
        $this->validate($this->request, $rules, $message);

        //查询SM信息
        $member = $this->order->getMember($params);
        if (!$member) {
            error("导购不存在");
        }
        $params['member_id'] = $member->id;

        $res = $this->order->deliverOrder($params);

        if ($res) {
            return successJson();
        } else {
            error("发货失败");
        }
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
            'mobile'        => 'required',
            'id'            => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'mobile.required'        => '手机号不能为空',
            'id.required'            => '订单ID不能为空',
        ];
        $this->validate($this->request, $rules, $message);

        //查询SM信息
        $member = $this->order->getMember($params);
        if (!$member) {
            error("导购不存在");
        }
        $params['member_id'] = $member->id;

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
            'mobile'        => 'required',
            'id'            => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'mobile.required'        => '手机号不能为空',
            'id.required'            => '退款订单ID不能为空',
        ];
        $this->validate($this->request, $rules, $message);

        //查询SM信息
        $member = $this->order->getMember($params);
        if (!$member) {
            error("导购不存在");
        }
        $params['member_id'] = $member->id;

        $res = $this->order->refundOrderDetail($params);

        return successJson($res);
    }
}