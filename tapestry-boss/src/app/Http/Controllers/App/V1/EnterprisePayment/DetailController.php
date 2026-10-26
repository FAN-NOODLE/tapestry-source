<?php


namespace App\Http\Controllers\App\V1\EnterprisePayment;


use App\Http\Controllers\App\V1\BaseController;
use App\Http\Requests\V1\EnterprisePaymentRequest;
use App\Services\EnterpriseNew\EnterpriseNewService;
use App\Services\EnterprisePaymentService\EnterprisePaymentService;

class DetailController extends BaseController
{
    public function index(EnterprisePaymentService $enterprisePaymentService)
    {
        $params = $this->request->only([
            'id',
        ]);
        EnterprisePaymentRequest::detail($params);

        $result = $enterprisePaymentService->detail($params);
        if (!$result) {
            notice(400404);
        }
        $enterprise = (new EnterpriseNewService())->getEnterpriseDetail($result->enterprise_id, true);
        $data = $result->toArray();
        $data['enterprise_name'] = $enterprise['name'];
        return success(['data' => $data]);

    }

}