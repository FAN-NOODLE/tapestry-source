<?php

namespace App\Http\Controllers\V1\Score;

use CrsSdk\Shop\Score;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class DeleteController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->all();

        $params['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('删除积分规则');

        $result = (new  Score())->delete($params);

        return success($result);
    }
}