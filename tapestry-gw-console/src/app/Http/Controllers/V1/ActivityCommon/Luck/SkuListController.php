<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/19 0019
 * Time: 20:28
 */

namespace App\Http\Controllers\V1\ActivityCommon\Luck;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroActivity\Activity;
use CrsSdk\Shop\Template;

class SkuListController extends BaseController
{
    /**
     * sku列表
     *
     * @return \Illuminate\Http\JsonResponse
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/6/15 11:20
     */
    public function index()
    {

        $params = $this->request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['source_id'] = $this->enterpriseId;
        //商品的状态（1:待上架,2:上架中,3:已下架,4:已暂停），默认没有限制 例：[1, 3]
        // 过滤仅自提商品
        $params['sale_type'] = [1, 3];
        $params['need_sign_price'] = 1;
        // sku列表数据
        $skuList = (new \CrsSdk\MicroGoodsV1\Activity())->getSkuList($params);


        if (empty($skuList['list'])) {
            return success($skuList);
        }
        // 获取sku占用列表数据
        $params['type'] = $params['occupy_type'] ?? '';
        $params['start_at'] = $params['start_at'] ?? 0;
        $params['end_at'] = $params['end_at'] ?? 0;
        $params['rule_id'] = $params['rule_id'] ?? 0;
        $occupySpuList = (new Activity())->activitySpu($params);

        $resList = [];

        foreach ($skuList['list'] as $key => $value) {
            $skuAttr = $value['attribute_arr'] ?? [];
            $skuAttrStr = implode(' ',$skuAttr);

            $newValue = [
                'market_price' => $value['sign_price'],
                'sell_price'   => $value['on_line_price'],
                'sku_id'       => $value['id'],
                'spu_id'       => $value['spu_id'],
                'sku_name'     => $value['sku_name'],
                'sku_img'      => $value['picture_url'],
                'material_stock' => $value['stock'],
                'material_id'    => $value['material_id'],
                'sku_attr'       => $skuAttrStr,
                'occupy'         => 0,
            ];
            if (in_array($value['spu_id'], $occupySpuList['spu_ids'])) {
                $newValue['occupy'] = 1;
            }
            $resList[] = $newValue;
        }

        $skuList['list'] = $resList;
        return success($skuList);
    }

}
