<?php

namespace App\Http\Controllers\V1\Score;

use CrsSdk\Shop\Score;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class UserScoreListController extends BaseController
{
    /**
     * 客户积分明细列表
     * @return JsonResponse
     * @throws \Exception
     */
    public function userScoreList()
    {
        $params = $this->request->all();

        $params['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('查看客户积分明细列表');

        $result = (new Score())->userScoreList($params);

        return success($result);
    }
}