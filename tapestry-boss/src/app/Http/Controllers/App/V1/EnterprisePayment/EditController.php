<?php


namespace App\Http\Controllers\App\V1\EnterprisePayment;


use App\Http\Controllers\App\V1\BaseController;
use App\Http\Requests\V1\EnterprisePaymentRequest;
use App\Services\EnterprisePaymentService\EnterprisePaymentService;

class EditController extends BaseController
{
    /**
     * 收款
     * @param EnterprisePaymentService $enterprisePaymentService
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     * @date 2020/11/6
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function receivePayment(EnterprisePaymentService $enterprisePaymentService)
    {
        $params = $this->request->only([
            'id',
            'collection_amount',
            'server_charge'
        ]);
        EnterprisePaymentRequest::receivePayment($params);
        $result = $enterprisePaymentService->receivePayment($params);

        if ($result !== true) {
            return notice($result);
        }
        return success(['data' => []]);

    }


    /**
     * 续费
     * @param EnterprisePaymentService $enterprisePaymentService
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     * @date 2020/11/6
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function renew(EnterprisePaymentService $enterprisePaymentService)
    {
        $params = $this->request->only([
            'id',
            'payment_amount',
            'collection_amount',
            'server_charge',
            'buy_day',
            'payment_time',
            'payment_type',
            'payment_status',
            'payment_order_sn',
            'voucher_url',
            'remark',
        ]);

        EnterprisePaymentRequest::renew($params);
        $result = $enterprisePaymentService->renew($params);

        if ($result !== true) {
            return notice($result);
        }
        return success(['data' => []]);
    }
}