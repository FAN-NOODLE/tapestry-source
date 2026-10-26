<?php

namespace App\Http\Controllers\V1\WelcomeSpeech;

use Illuminate\Http\Request;
use CrsSdk\Work\WelcomeSpeech;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class DetailController extends BaseController
{
    /**
     * 欢迎语详情
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['id', 'is_default']);

        $params['enterprise_id'] = $this->enterpriseId;

        $detail = (new WelcomeSpeech())->detail($params);

        return success($detail['data']);
    }
}
