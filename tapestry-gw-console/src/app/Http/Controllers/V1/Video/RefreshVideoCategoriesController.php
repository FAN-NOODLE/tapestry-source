<?php

namespace App\Http\Controllers\V1\Video;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Video;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class RefreshVideoCategoriesController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $params = $this->getParams();
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new Video())->refreshVideoCategories($params));
    }

    /**
     * 获取参数
     * @return array
     * @throws ValidationException
     */
    public function getParams(): array
    {
        return $this->validate($this->request, [
            'category_ids' => 'array',
        ], [
            'category_ids.*' => '参数[category_ids]不合法',
        ]);
    }
}