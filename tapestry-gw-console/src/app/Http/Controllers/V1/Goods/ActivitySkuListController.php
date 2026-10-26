<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/12 0012
 * Time: 14:55
 */

namespace App\Http\Controllers\V1\Goods;


use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Activity as GoodsActivity;

class ActivitySkuListController extends BaseController
{
    /**
     * 售卖类型:线上
     */
    const SALE_TYPE_ONLINE = 1;

    /**
     * 售卖类型:线下
     */
    const SALE_TYPE_OFFLINE = 2;

    /**
     * 售卖类型:线上+线下
     */
    const SALE_TYPE_ALL = 3;

    public function index()
    {
        $param = $this->request->input();
        $search['enterprise_id'] = $this->enterpriseId;
        $search['spu_ids'] = (array)$param['spu_id'];
        // 商品的售卖方式（1:线上,2:线下,3:线上+线下），默认没有限制
        // 这里排除线下的商品
        if (isset($param['sale_type']) && $param['sale_type'] == 2) {
            $search['sale_type'] = [2, 3];
        } else {
            $search['sale_type'] = [1, 3];
        }

        try {
            $skuList = (new GoodsActivity())->getSkuBySpuIds($search);
        } catch (\Exception $e) {
            throw new Error($e->getCode(),$e->getMessage());
        }

        if (!$skuList) {
            return success([]);
        }
        $SkuListNew = [];
        foreach ($skuList as $key => $sku) {
            $skuId = $sku['id'];
            $SkuListNew[$key]['sku_id'] = $skuId;
            $SkuListNew[$key]['spu_id'] = $sku['spu_id'];
            $SkuListNew[$key]['name'] = $sku['sku_name'];
            $SkuListNew[$key]['sku_picture'] = $sku['picture_url'];
            $SkuListNew[$key]['sign_price'] = $sku['sign_price'] ?? 0;
            $SkuListNew[$key]['stock'] = $sku['stock'];
            $SkuListNew[$key]['online_stock'] = $sku['online_stock'];
            $SkuListNew[$key]['offline_stock'] = $sku['offline_stock'];
            $SkuListNew[$key]['material_id'] = $sku['material_id'];
            $SkuListNew[$key]['source_id'] = $sku['source_id'];
            $SkuListNew[$key]['spu_name'] = $sku['spu_name'];
            $SkuListNew[$key]['attribute_arr'] = $sku['attribute_arr'];
            // 自提，取线下价格
            if ($sku['sale_type'] == self::SALE_TYPE_OFFLINE) {
                $SkuListNew[$key]['sell_price'] = $sku['off_line_price'];
            } else {
                // 快递，和快递+自提，取线上价格
                $SkuListNew[$key]['sell_price'] = $sku['on_line_price'];
            }
        }
        return success($SkuListNew);
    }

}
