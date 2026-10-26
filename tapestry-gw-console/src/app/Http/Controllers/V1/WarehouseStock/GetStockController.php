<?php


namespace App\Http\Controllers\V1\WarehouseStock;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\WarehouseStock;

class GetStockController extends BaseController
{
    /**
     * 获取物料的剩余库存
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        $stock = (new WarehouseStock())->getStock($input);

        return success($stock);
    }
}