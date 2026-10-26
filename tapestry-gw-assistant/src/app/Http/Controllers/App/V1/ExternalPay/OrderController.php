<?php

namespace App\Http\Controllers\App\V1\ExternalPay;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Work\ExternalPay\ExternalPay;
use Exception;
use Illuminate\Http\JsonResponse;
use Symfony\Component\HttpFoundation\Response as ResponseEntity;

class OrderController extends BaseController
{
    /**
     * @var $inherit
     */
    private $order;

    public function __construct()
    {
        parent::__construct();
        $this->order = new ExternalPay();
    }

    /**
     * 订单列表
     * @return JsonResponse
     * @throws Exception
     */
    public function orderList():JsonResponse
    {
        $params                  = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $result = $this->order->orderList($params);
        return success($result);
    }

    /**
     * 退款列表
     * @return JsonResponse
     * @throws Exception
     */
    public function refundOrderList():JsonResponse
    {
        $params                  = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $result = $this->order->refundOrderList($params);
        return success($result);
    }

    /**
     * 发货
     * @return JsonResponse
     */
    public function deliverOrder():JsonResponse
    {
        $params                  = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;

        $result = $this->order->deliverOrder($params);
        return success($result);
    }

    /**
     * 订单详情
     * @return void
     */
    public function orderDetail():JsonResponse
    {
        $params                  = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;

        $result = $this->order->orderDetail($params);
        return success($result);
    }

    /**
     * 退款详情
     * @return void
     */
    public function refundOrderDetail():JsonResponse
    {
        $params                  = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;

        $result = $this->order->refundOrderDetail($params);
        return success($result);
    }
}