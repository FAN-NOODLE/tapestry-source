<?php


namespace App\Http\Controllers\App\V1\EnterprisePayment;


use App\Http\Controllers\App\V1\BaseController;
use App\Http\Requests\V1\EnterprisePaymentRequest;
use App\Models\EnterprisePayment;
use App\Services\EnterprisePaymentService\EnterprisePaymentService;
use Illuminate\Support\Facades\Redis;


class AddController extends BaseController
{
    public function index()
    {
        $params = $this->request->only([
           'enterprise_id',
           'payment_category',
           'payment_amount',
           'available_amount',
           'collection_amount',
           'server_charge',
           'buy_number',
           'buy_day',
           'payment_time',
           'use_time',
           'payment_type',
           'payment_status',
           'payment_order_sn',
           'voucher_url',
           'remark'
        ]);

        EnterprisePaymentRequest::add($params);
        if ($params['payment_category'] == EnterprisePayment::PAYMENT_CATEGORY_ZHDG) {
            return notice(400104);
        }
        $params['payment_source'] = '手动创建'; // 手动创建
        $result = (new EnterprisePaymentService)->add($params);
        if ($result === true) {
            return success(['data' => []]);
        }
        return notice($result);
    }
}