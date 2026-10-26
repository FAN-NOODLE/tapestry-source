<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/19 0019
 * Time: 20:28
 */

namespace App\Http\Controllers\V1\ActivityCommon\Common;

use App\Http\Controllers\V1\ActivityCommon\ActivityBaseController;

class CopyDetailController extends ActivityBaseController
{

    function index()
    {
        $params = $this->params;
        $activityInfo = [];
        try {
            $activityInfo = $this->activityApi->getActivity($params);
        } catch (\Exception $e) {
            $this->notice($e, $params);
        }
        if (isset($activityInfo['activity_sku'])) {
            if ($params['type'] == 'help') {
                // 如果是好友助力返回字段转化
                foreach ($activityInfo['activity_sku'] as $value) {
                    $skuData[] = [
                        'id'                  => $value['sku_id'] ?? 0,
                        'sku_name'            => $value['name'] ?? '',
                        'spu_id'              => $value['spu_id'] ?? 0,
                        'enterprise_id'       => $value['enterprise_id'] ?? 0,
                        'material_id'         => $value['material_id'] ?? 0,
                        'stock'               => $value['stock'] ?? 0,
                        'attribute_arr'       => empty($value['attribute_item_value']) ? [] : explode(' ', $value['attribute_item_value']),
                        'on_line_price'       => $value['sell_price'] ?? 0,
                        'sell_price'          => $value['sell_price'] ?? 0,
                        'activity_price'      => $value['activity_price'] ?? 0,
                        'activity_stock_set'  => $value['activity_stock'] ?? 0,
                        'activity_stock'      => $value['activity_stock'] ?? 0,
                        'activity_stock_type' => $value['activity_stock_type'] ?? 0,
                        'picture_url'         => $value['sku_picture'] ?? '',
                    ];
                    $activityInfo['activity_sku'] = $skuData;
                }
            }

            if (in_array($params['type'], ['flash-sale', 'buyout-price'])){
                $spuData = [];
                foreach (collect($activityInfo['activity_sku'] ?? [])->groupBy('spu_id')->toArray() as $index => $item) {
                    $spuData[] = [
                        'spu_id'   => $index,
                        'spu_name' => current($item)['spu_name'] ?? '',
                        'url' => current($item)['picture_url'],
                        'sku_list' => $item
                    ];
                }
                $activityInfo['activity_sku'] = $spuData;
            }
        }

        unset($sku);
        return success($activityInfo);
    }

}
