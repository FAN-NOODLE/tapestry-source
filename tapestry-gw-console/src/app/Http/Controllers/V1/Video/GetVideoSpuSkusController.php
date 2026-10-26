<?php


namespace App\Http\Controllers\V1\Video;


use App\Http\Controllers\Controller;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\MicroGoodsV1\Video;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class GetVideoSpuSkusController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $params = $this->getParams();
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new Video())->getVideoSpuSkus($params));
    }

    /**
     * 获取参数
     * @return array
     * @throws ValidationException
     */
    public function getParams(): array
    {
        return $this->validate($this->request, [
            'spu_id' => 'required|integer'
        ], [
            'spu_id.*' => '参数[spu_id]不合法'
        ]);
    }
}