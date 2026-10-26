<?php


namespace App\Http\Controllers\V1\Goods;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use Illuminate\Support\Facades\Log;

/**
 * 获取指定商品的商品详情（包括失效的商品）
 *
 * Class GetDeleteGoodsDetailController
 *
 * @package App\Http\Controllers\V1\Goods
 */
class GetDeleteGoodsDetailController extends BaseController
{
    public function index()
    {
        $params = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;

        try {
            $data = (new Goods())->getDeleteGoodsDetail($params);
        } catch (\Exception $e) {
            Log::error('复制商品获取商品详情失败' . $e->getMessage(), $params);
            notice(450000, [$e->getMessage()]);
        }

        return success($data);
    }
}