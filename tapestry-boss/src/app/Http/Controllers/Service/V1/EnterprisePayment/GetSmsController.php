<?php


namespace App\Http\Controllers\Service\V1\EnterprisePayment;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\EnterprisePaymentService\EnterpriseSmsStatisticsService;

class GetSmsController extends BaseController
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

        $smsStatistics = (new EnterpriseSmsStatisticsService)->getSms($input);
        if ($smsStatistics) {
            return success([
                'sms_number' => $smsStatistics->sms_number,
                'sms_total_number' => $smsStatistics->sms_total_number,
            ]);
        }
        return error(400110);
    }
}