<?php

namespace App\Http\Controllers\V1\ScoreExchange;

use Illuminate\Http\Request;
use CrsSdk\Shop\ScoreExchange;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class CreateScoreCouponController extends BaseController
{
    /**
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function __invoke(Request $request)
    {
        $params = $this->request->all();

        $params['enterprise_id'] = $this->enterpriseId;

        $response = (new ScoreExchange())->createScoreCoupon($params);

        $this->writeOperatorLog('创建积分兑换券');

        return success($response);
    }
}