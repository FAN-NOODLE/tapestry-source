<?php

namespace App\Http\Controllers\V1\AppCenter\CommunityDiversion;

use Illuminate\Http\JsonResponse;
use CrsSdk\Shop\CommunityDiversion;
use App\Http\Controllers\V1\BaseController;

class SaveController extends BaseController
{
    /**
     * 保存配置
     * @return JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();
        $params = array_merge([
            'enterprise_id' => $this->enterprise['id'],
        ], $input);

        (new CommunityDiversion())->save($params);

        return success();
    }
}
