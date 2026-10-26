<?php


namespace App\Http\Controllers\V1\Video;


use App\Http\Controllers\Controller;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\MicroGoodsV1\Video;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class GetVideoSkuBaseController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $params = $this->getParams();
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new Video())->getVideoSkuBase($params));
    }

    /**
     * 获取参数
     * @return array
     * @throws ValidationException
     */
    public function getParams(): array
    {
        return $this->validate($this->request, [
            'sku_ids' => 'required|array',
        ], [
            'sku_ids.*' => '参数[sku_ids]不合法'
        ]);
    }
}