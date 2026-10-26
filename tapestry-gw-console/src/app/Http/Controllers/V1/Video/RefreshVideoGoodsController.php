<?php

namespace App\Http\Controllers\V1\Video;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Video;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class RefreshVideoGoodsController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $params = $this->getParams();
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new Video())->refreshVideoGoods($params));
    }

    /**
     * 获取参数
     * @return array
     * @throws ValidationException
     */
    public function getParams(): array
    {
        return $this->validate($this->request, [
            'spu_ids' => 'array',
        ], [
            'spu_ids.*' => '参数[spu_ids]不合法',
        ]);
    }
}