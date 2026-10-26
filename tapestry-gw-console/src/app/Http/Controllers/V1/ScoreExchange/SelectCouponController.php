<?php

namespace App\Http\Controllers\V1\ScoreExchange;

use Illuminate\Http\Request;
use CrsSdk\Shop\ScoreExchange;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class SelectCouponController extends BaseController
{
    /**
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function __invoke(Request $request)
    {
        $params = $request->all();

        $params['enterprise_id'] = $this->enterpriseId;

        $response = (new ScoreExchange)->selectCoupon($params);

        $this->writeOperatorLog('添加优惠券');

        return success($response);
    }
}