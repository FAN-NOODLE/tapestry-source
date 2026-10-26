<?php

namespace App\Http\Controllers\V1\Score;

use CrsSdk\Shop\Score;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class SaveController extends BaseController
{
    /**
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 接值
        $params = $request->all();

        $params['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('添加积分规则');

        $result = (new Score())->save($params);

        return success($result);
    }
}