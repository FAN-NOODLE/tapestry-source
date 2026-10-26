<?php

namespace App\Http\Controllers\V1\Score;

use CrsSdk\Shop\Score;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class CountController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $result = (new Score())->count(['enterprise_id' => $this->enterpriseId]);
        return success($result);
    }
}