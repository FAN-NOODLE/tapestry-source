<?php

namespace App\Http\Controllers\V1\Score;

use CrsSdk\Shop\Score;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class UserScoreDetailController extends BaseController
{
    /**
     * 客户积分明细详情
     * @return JsonResponse
     * @throws \Exception
     */
    public function userScoreDetail()
    {
        $params = $this->request->all();

        $params['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('查看客户积分明细');

        $result = (new Score())->userScoreDetail($params);

        return success($result);
    }
}