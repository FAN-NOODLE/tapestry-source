<?php


namespace App\Http\Controllers\V1\Video;


use App\Http\Controllers\Controller;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\MicroGoodsV1\Video;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class GetCategoryTreeController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new Video())->getVideoCategoryTree($params));
    }
}