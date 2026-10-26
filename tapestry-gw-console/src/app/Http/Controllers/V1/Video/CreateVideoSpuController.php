<?php


namespace App\Http\Controllers\V1\Video;


use App\Http\Controllers\Controller;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\MicroGoodsV1\Video;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class CreateVideoSpuController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $params = $this->getParams();
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new Video())->createVideoSpu($params));
    }

    /**
     * 获取参数
     * @return array
     * @throws ValidationException
     */
    public function getParams(): array
    {
        return $this->validate($this->request, [
            'brand_id' => 'required|integer',
            'spu_ids' => 'required|array',
            'category_id' => 'required|integer',
            'qualifications' => 'array',
        ], [
            'brand_id.*' => '参数[brand_id]不合法',
            'spu_ids.*' => '参数[spu_ids]不合法',
            'category_id.*' => '参数[category_id]不合法',
            'qualifications.*' => '参数[certification]不合法',
        ]);
    }
}