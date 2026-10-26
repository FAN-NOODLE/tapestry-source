<?php

namespace App\Http\Controllers\V1\ScoreExchange;

use Illuminate\Http\Request;
use CrsSdk\Shop\ScoreExchange;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class SetExchangeScoreController extends BaseController
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

        $response = (new ScoreExchange)->setExchangeScore($params);

        $this->writeOperatorLog('设置积分券兑换积分值');

        return success($response);
    }
}