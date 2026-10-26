<?php


namespace App\Http\Controllers\Service\V1\Settlement;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\SettlementCallbackService;

class CallbackController extends BaseController
{

    /**
     * @var SettlementCallbackService
     */
    private $service;

    /**
     * 回调数据
     * @var string
     */
    private $callbackData;

    /**
     * CallbackController constructor.
     */
    public function __construct()
    {
        $this->callbackData = request()->input('callback_data', '');
        if (empty($this->callbackData)) {
            throw new Error(110007, "回调数据不能为空");
        }
        $this->service = new SettlementCallbackService();
    }

    /**
     * 结算单状态变更
     * @return array
     */
    public function settleBalance()
    {
        $this->service->settleBalanceCallback($this->callbackData);
        return success([]);
    }

    /**
     * 结算单退款
     */
    public function refundBalance()
    {
        $this->service->refundBalanceCallback($this->callbackData);
        return success();
    }
}