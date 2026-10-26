<?php


namespace App\Http\Controllers\V1\Video;


use App\Http\Controllers\Controller;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\MicroGoodsV1\Video;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class GetCategoryFromIdController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $params = $this->getParams();
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new Video())->getVideoCategoryFromPid($params));
    }

    /**
     * @return array
     * @throws ValidationException
     */
    public function getParams(): array
    {
        return $this->validate($this->request, [
            'pid' => 'required|integer'
        ], [
            'pid.*' => '参数[pid]不合法pid不合法'
        ]);
    }
}