<?php

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Sdk;
use Illuminate\Http\JsonResponse;

class FiscalMonthController extends BaseController
{


    /**
     * 顾客数据排行榜
     **
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(): JsonResponse
    {
        $params['enterprise_id'] = $this->enterpriseId;

        $res = Sdk::work()->setApi("service/v1/data-compass/customer-statistic/fiscal-month")->post($params);
        
        return success($res);
    }


}