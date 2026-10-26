<?php
/**
 * GoodsSalesController.php
 * Created on 2019-08-02 15:14
 * Created by chengyujia@vchangyi.com
 */

namespace App\Http\Controllers\V1\Goods;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use Illuminate\Http\Request;

class GoodsSalesController extends BaseController
{
    /**
     * 活动商品列表
     *
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $search                 = $request->input('search', []);
        $input['page']          = $request->input('page', Constant::PAGE);
        $input['size']          = $request->input('size', Constant::LIMIT);
        $input['enterprise_id'] = $this->enterpriseId;
        $input['source_ids']    = [$this->enterpriseId];
        $input['search']        = $search['spu_name'] ?? '';
        $input['search_type']   = 5;

        $goodsList = (new Goods())->spuList($input);

        foreach ($goodsList['list'] as &$value) {
            $value['spu_id'] = $value['id'];
            if (isset($value['spu_source']))
                $value['spu_picture'] = $value['spu_source']['url'] ?? '';
            unset($value['id'], $value['spu_source']);
        }

        return success($goodsList);
    }
}
