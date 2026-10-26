<?php

namespace App\Http\Controllers\V1\Score;

use CrsSdk\Shop\Score;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class SaveDeadlineController extends BaseController
{
    /**
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->all();

        $params['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('设置积分有效期');

        $result = (new Score())->saveDeadline($params);

        return success($result);
    }
}