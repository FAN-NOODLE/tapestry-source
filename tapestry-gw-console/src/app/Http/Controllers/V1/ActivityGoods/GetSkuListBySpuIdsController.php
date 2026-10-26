<?php


namespace App\Http\Controllers\V1\ActivityGoods;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Activity;

class GetSkuListBySpuIdsController extends BaseController
{
    /**
     * 根据spuIds获取skuList
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $input    = $this->request->all();
        $minPrice = $input['min_price'] ?? 0;

        $skuList = (new Activity())->getSkuBySpuIds($input);

        // 如果需要过滤最小价格
        if (!empty($minPrice)) {
            foreach ($skuList as $key => $item) {
                if ($item['on_line_price'] < $minPrice) {
                    unset($skuList[$key]);
                }
            }
            sort($skuList);
        }

        return success($skuList);
    }
}