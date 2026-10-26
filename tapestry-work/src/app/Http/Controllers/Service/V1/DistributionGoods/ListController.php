<?php

namespace App\Http\Controllers\Service\V1\DistributionGoods;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\DistributionGoodsService;
use Illuminate\Http\Request;

class ListController extends BaseController
{
    /**
     * 分销商品列表
     *
     * @param Request $request
     * @param DistributionGoodsService $distributionGoodsService
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request, DistributionGoodsService $distributionGoodsService)
    {
        $enterpriseId   = $request->get('enterprise_id');
        $search         = $this->getParams();
        $data           = $distributionGoodsService->goodsList($enterpriseId, $search);
        return success($data);
    }

    /**
     * 获取参数
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    private function getParams(): array
    {
        $rules = [
            'enterprise_id'         => 'required',
            'page'                  => 'integer',
            'limit'                 => 'integer',
            'spu_name'              => 'max:128',
            'sku_id'                => 'integer'
        ];

        $input = $this->validate($this->request, $rules);
        $input['page']           = $input['page'] ?? Constant::PAGE;
        $input['limit']          = $input['limit'] ?? Constant::SIZE;
        return $input;
    }
}
