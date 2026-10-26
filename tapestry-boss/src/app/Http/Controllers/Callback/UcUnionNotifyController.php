<?php
/**
 * 小程序开发码回调
 */

namespace App\Http\Controllers\Callback;

use App\Jobs\EnterprisePaymentJob;
use App\Jobs\SendEmailJob;
use App\Library\Uc\AppStore;
use App\Models\Enterprise;
use App\Models\EnterpriseBuyRecord;
use Illuminate\Support\Facades\Log;
use App\Services\EnterpriseLogService;

class UcUnionNotifyController extends BaseController
{
    public function index()
    {
        $params = request()->input();

        EnterpriseLogService::info('企业账号续费回调来了', $params);
        Log::info('企业账号续费回调来了:' . json_encode($params));

        if (empty($params)) {
            EnterpriseLogService::info('企业账号续费回调failed:empty request', $params);
            return 'failed:empty request';
        }

        if (empty($params['corpId']) || empty($params['suiteId']) || empty($params['orderId'])) {
            EnterpriseLogService::info('企业账号续费回调failed:empty corpId or suiteId or orderId', $params);
            return 'success';
        }

        if (empty($params['eventType']) || 'pay_for_app_success' <> $params['eventType'] ) {
            EnterpriseLogService::info('企业账号续费回调failed:non_pay_for_app_success', $params);
            return 'success';
        }

//        $params            = json_decode('{"timeStamp":1595779621,"suiteId":"tjc4b15e6695e09895","corpId":"wwcf40db843ccd31b6","orderId":"T0000085884735C6664174A103DA4","eventType":"pay_for_app_success","plIdentifier":"yuanquan","operatorId":"www"}', 1);
//        $params['corp_id'] = 'wwcf40db843ccd31b6';  // 畅展crs的corp_id

        $enterprise = Enterprise::query()->where('corp_id', $params['corpId'])->first();

        if (empty($enterprise)) {
            EnterpriseLogService::info('企业账号续费回调failed:empty enterprise', $params);
            return 'success';
        }

        $records = EnterpriseBuyRecord::query()->where([
            'suiteid'  => $params['suiteId'],
            'order_id' => $params['orderId']
        ])->get()->toArray();

        if ($records) {
            EnterpriseLogService::info('企业账号续费回调failed:record duplicated', $params);
            return 'success';
        }

        try {
            $order = (new AppStore())->order($enterprise['domain'], [
                    'suiteId' => $params['suiteId'],
                    'orderId' => $params['orderId']
                ])['data'] ?? [];
        } catch (\Exception $exception) {
            EnterpriseLogService::info('企业账号续费回调failed:' . $exception->getMessage(), $params);
            return 'failed:' . $exception->getMessage();
        }

        EnterpriseBuyRecord::query()->create([
            'enterprise_id'         => $enterprise['id'] ?? 0,
            'order_id'              => $params['orderId'],
            'order_status'          => $order['orderStatus'] ?? 0,
            'order_type'            => $order['orderType'] ?? 0,
            'paid_corpid'           => $order['paidCorpid'],
            'operator_id'           => $order['operatorId'] ?? '',
            'suiteid'               => $params['suiteId'],
            'appid'                 => $order['appid'] ?? 0,
            'edition_id'            => $order['editionId'] ?? '',
            'edition_name'          => $order['editionName'] ?? '',
            'price'                 => $order['price'] ?? 0,
            'user_count'            => $order['userCount'] ?? 0,
            'order_period'          => $order['orderPeriod'] ?? 0,
            'order_time'            => $order['orderTime'] ?? 0,
            'paid_time'             => $order['paidTime'] ?? 0,
            'begin_time'            => $order['beginTime'] ?? 0,
            'end_time'              => $order['endTime'] ?? 0,
            'order_from'            => $order['orderFrom'] ?? 0,
            'operator_corpid'       => $order['operatorCorpid'] ?? '',
            'service_share_amount'  => $order['serviceShareAmount'] ?? 0,
            'platform_share_amount' => $order['platformAhareAmount'] ?? 0,
            'dealer_share_amount'   => $order['dealerAhareAmount'] ?? 0,
            'dealer_corp_info'      => $order['dealerCorpInfo'] ?? '',
        ]);
        (new SendEmailJob($enterprise, $order['userCount'] ?? 0))->handle();
        (new EnterprisePaymentJob([
            'enterprise_id'     => $enterprise['id'] ?? 0,
            'buy_day'           => $order['orderPeriod'] ?? 0, // 购买天数
            'buy_number'        => $order['userCount'] ?? 0, // 购买账号数
            'payment_time'      => $order['paidTime'] ?? 0, // 支付时间
            'payment_order_sn'  => $params['orderId'], // 支付流水单号
            'payment_amount'    => $order['price'] ?? 0, // 实付金额
            'use_time'          => $order['beginTime'] ?? 0, // 生效时间
            'fail_time'         => $order['endTime'] ?? 0, // 失效时间
            'order_status'      => $order['orderStatus'] ?? 0,
        ]))->handle();
        EnterpriseLogService::info('企业账号续费回调success', $params);
        return 'success';
    }
}
