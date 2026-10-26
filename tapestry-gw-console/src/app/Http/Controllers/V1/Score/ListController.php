<?php

namespace App\Http\Controllers\V1\Score;

use CrsSdk\Shop\Score;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $params = [
            'enterprise_id' => $this->enterpriseId
        ];

        $result = (new Score())->lists($params);

        return success($result);
    }
}