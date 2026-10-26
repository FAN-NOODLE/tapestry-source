<?php

namespace App\Http\Controllers\App\V1\ExternalPay;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Work\ExternalPay\ExternalPay;
use Illuminate\Http\JsonResponse;

class LogisticsController extends BaseController
{
    /**
     * 物流公司
     * @return JsonResponse
     * @throws \Exception
     */
     public function logisticsCompany():JsonResponse
     {
         $params = request()->all();
         $result = (new ExternalPay())->logisticsCompany($params);

         return success($result);
     }


    /**
     * @return void
     * 物流信息
     */
     public function logisticsTrack():JsonResponse
     {
         $params                  = request()->all();
         $params['enterprise_id'] = $this->enterpriseId;

         $result = (new ExternalPay())->logisticsTrack($params);

         return success($result);
     }
}