<?php


namespace App\Http\Controllers\App\V1\EnterprisePayment;


use App\Http\Controllers\App\V1\BaseController;
use App\Services\EnterprisePaymentService\EnterprisePaymentService;
use App\Services\EnterprisePaymentService\EnterpriseSmsStatisticsService;
use CrsSdk\Work\EnterprisePayment AS SDKEnterprisePayment;


class StatisticsController extends BaseController
{
    public function index()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'required | integer'
        ],[
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer' => '企业ID类型错误',
        ]);

        $input = $this->request->all();

        // 初始化数据
        $result = [
            'red_pack' => [
                'surplus_money' => 0,
                'total_money' => 0,
            ],
            'sms' => [
                'surplus_number' => 0,
                'total_number' => 0,
                'total_money' => 0
            ],
        ];
        // 红包
        $redPack = (new SDKEnterprisePayment)->budgetTotal($input);
        $result['red_pack']['surplus_money'] = $redPack['data']['surplus_money'];

        // 短信
        $smsStatistics = (new EnterpriseSmsStatisticsService)->getSms($input);
        if ($smsStatistics) {
            $result['sms']['surplus_number'] = $smsStatistics->sms_number;
            $result['sms']['total_number'] = $smsStatistics->sms_total_number;
            $result['sms']['total_money'] = $smsStatistics->sms_payment_total;
        }

        // 支付、账号统计
        $enterprisePaymentData = (new EnterprisePaymentService)->statistics($input);
        $result['other'] = $enterprisePaymentData;
        unset($result['other']['red_pack_total_money']);
        $result['red_pack']['total_money'] = $enterprisePaymentData['red_pack_total_money'];
        return success(['data' => $result]);
    }
}