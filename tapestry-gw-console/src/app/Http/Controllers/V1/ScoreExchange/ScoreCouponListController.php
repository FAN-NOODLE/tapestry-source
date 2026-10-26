<?php

namespace App\Http\Controllers\V1\ScoreExchange;

use Illuminate\Http\Request;
use CrsSdk\Shop\ScoreExchange;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class ScoreCouponListController extends BaseController
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

        $response = (new ScoreExchange)->scoreCouponList($params);

        $this->writeOperatorLog('获取积分兑券列表');

        return success($response);
    }
}