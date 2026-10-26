<?php


namespace App\Http\Controllers\Service\V1\EnterprisePayment;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseSmsStatistics;
use App\Services\EnterprisePaymentService\EnterpriseSmsStatisticsService;
use Illuminate\Http\Request;

class DepleteSmsController extends BaseController
{
    public function index()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'required | integer | min:0',
            'number' => 'required | integer | min:0'
        ]);

        $input = $this->request->all();
        $result = (new EnterpriseSmsStatisticsService)->depleteSms($input);
        if ($result) {
            return success([
                'sms_number' => $result->sms_number,
                'sms_total_number' => $result->sms_total_number,
            ]);
        }
        return error(400112);
    }
}