<?php


namespace App\Http\Controllers\V1\Video;


use App\Http\Controllers\Controller;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\MicroGoodsV1\Video;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class HandlerVideoSpuController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $params = $this->getParams();
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new Video())->handlerVideoSpu($params));
    }

    /**
     * 获取参数
     * @return array
     * @throws ValidationException
     */
    public function getParams(): array
    {
        return $this->validate($this->request, [
            'id' => 'required|integer',
            'brand_id' => 'integer',
            'status' => 'required|string',
            'spu_id' => 'integer',
            'qualifications' => 'array',
            'skus' => 'array',
            'product_id' => 'integer',
            'path' => 'string',
            'out_product_id' => 'string',
            'category_id' => 'integer|required_if:status,UPDATE',
        ], [
            'id.*' => '参数[id]不合法',
            'brand_id.*' => '参数[brand_id]不合法',
            'status.*' => '参数[status]不合法',
            'spu_id.*' => '参数[spu_id]不合法',
            'qualifications.*' => '参数[qualifications]不合法',
            'skus.*' => '参数[skus]不合法',
            'path.*' => '参数[path]不合法',
            'product_id.*' => '参数[product_id]不合法',
            'out_product_id.*' => '参数[out_product_id]不合法',
            'category_id.integer' => '参数[category_id]不合法',
            'category_id.required_if' => '参数[category_id]不能为空',
        ]);
    }
}