<?php


namespace App\Http\Controllers\V1\Video;


use App\Http\Controllers\Controller;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\MicroGoodsV1\Video;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class CreateVideoCategoryController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $params = $this->getParams();
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new Video())->handlerVideoCategory($params));
    }

    /**
     * 获取参数
     * @return array
     * @throws ValidationException
     */
    public function getParams(): array
    {
        return $this->validate($this->request, [
            'id' => 'integer',
            'license' => 'string',
            'brand_id' => 'integer',
            'level1' => 'integer',
            'level2' => 'integer',
            'level3' => 'integer',
            'certificate' => 'array',
        ], [
            'id.*' => '参数[id]不合法',
            'license.*' => '参数[license]不合法',
            'brand_id.*' => '参数[brand_id]不合法',
            'level1.*' => '参数[level1]不合法',
            'level2.*' => '参数[level2]不合法',
            'level3.*' => '参数[level3]不合法',
            'certificate.*' => '参数[certificate]不合法',
        ]);
    }
}