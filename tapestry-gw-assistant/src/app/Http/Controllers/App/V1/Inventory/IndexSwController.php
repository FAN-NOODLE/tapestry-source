<?php

namespace App\Http\Controllers\App\V1\Inventory;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\MicroScrm\InventorySwScrm;
use Illuminate\Http\JsonResponse;

class IndexSwController extends BaseController
{

    /**
     * 三方商品库存查询
     *
     * @return JsonResponse
     */
    public function index(): JsonResponse
    {
        $params = request()->all();

        $params['enterprise_id'] = $this->enterpriseId;
        if (empty($params['store_code'])) {
            $params['store_code'] = $this->memberAuthInfo['store_code'];
        }
        $params['member_id'] = $this->memberAuthInfo['member_id'];

        $result = InventorySwScrm::list($params);

        return success($result);
    }

    /**
     * 库存查询的历史搜索记录
     * @return JsonResponse
     */
    public function searchKeywords(): JsonResponse
    {
        $params = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];

        $result = InventorySwScrm::searchKeywords($params);

        return success($result?: []);

    }

    /**
     * 删除库存查询的历史搜索记录
     * @return JsonResponse
     */
    public function delKeywords(): JsonResponse
    {
        $params = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];

        $result = InventorySwScrm::delKeywords($params);

        return success($result?: []);

    }
}