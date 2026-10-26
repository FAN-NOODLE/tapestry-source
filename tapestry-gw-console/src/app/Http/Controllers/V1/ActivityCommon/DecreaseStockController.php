<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/13 0013
 * Time: 14:32
 */

namespace App\Http\Controllers\V1\ActivityCommon;


class DecreaseStockController extends ActivityBaseController
{

    public function index()
    {
        $param = $this->params;
        $param['type'] = $this->typeCode; //此处类型为数字

        try {
            $this->activityApi->decreaseStock($param);
        } catch (\Exception $e) {
            $this->notice($e,$param,'decrease stock fail');
        }
        $activityInfo = $this->getActivityInfo($param['rule_id']);
        $skuStock = [];

        foreach ($activityInfo['activity_sku'] as $sku) {
            if ($sku['sku_id'] == $param['sku_id']) {
                $skuStock['sku_id'] = $param['sku_id'];
                $skuStock['stock'] = $sku['stock'];
                $skuStock['activity_stock'] = $sku['activity_stock'];
                break;
            }
        }
        return success($skuStock);
    }

}
