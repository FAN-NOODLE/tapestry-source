<?php

namespace App\Http\Controllers\V1\AppCenter\CommunityDiversion;

use Illuminate\Http\JsonResponse;
use CrsSdk\Shop\CommunityDiversion;
use App\Http\Controllers\V1\BaseController;

class GetConfigController extends BaseController
{
    /**
     * 获取社群导流配置
     * @return JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $data = (new CommunityDiversion())->getConfig(['enterprise_id' => $this->enterpriseId]);

        return success($data);
    }
}
