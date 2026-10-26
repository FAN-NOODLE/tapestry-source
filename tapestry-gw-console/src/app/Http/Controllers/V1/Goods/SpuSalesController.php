<?php


namespace App\Http\Controllers\V1\Goods;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;

class SpuSalesController extends BaseController
{
    public function index()
    {
        $params = $this->request->only(['goods_id', 'source_id']);

        $data = [
            'enterprise_id' => $this->enterpriseId,
            'source_id'     => empty($params['source_id']) ? $this->enterpriseId : $params['source_id'],
            'goods_id'      => $params['goods_id'] ?? [],
            'is_spu'        => 1,
            'is_month'      => 0
        ];

        $result = (new Goods())->getGoodsSales($data);

        return success($result);
    }
}