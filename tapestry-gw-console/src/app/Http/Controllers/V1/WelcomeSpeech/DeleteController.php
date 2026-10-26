<?php

namespace App\Http\Controllers\V1\WelcomeSpeech;

use App\Exceptions\Error;
use Illuminate\Http\Request;
use CrsSdk\Work\WelcomeSpeech;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class DeleteController extends BaseController
{
    /**
     * 欢迎语删除
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['id']);
        if (empty($params['id'])) {
            throw new Error(20012101, 'id不能为空');
        }

        (new WelcomeSpeech())->delete($params);

        return success();
    }
}
