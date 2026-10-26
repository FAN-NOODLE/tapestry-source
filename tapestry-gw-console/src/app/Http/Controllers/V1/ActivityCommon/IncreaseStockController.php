<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/13 0013
 * Time: 14:32
 */

namespace App\Http\Controllers\V1\ActivityCommon;

use App\Http\Controllers\V1\BaseController;
use App\Exceptions\Error;

class IncreaseStockController extends ActivityBaseController
{

    public function index()
    {
        $param = $this->request->input();
        $param['enterprise_id'] = $this->enterpriseId;

        try {
            $this->activityApi->increaseStock($param);
        } catch (\Exception $e) {

            if (in_array($param['type'],[self::TYPE_GROUP,self::TYPE_HAGGLE])){
                throw new Error($e->getCode(), $e->getMessage(),$e->data['data'] ?? []);
            }

            $this->notice($e,$param,'increase stock fail');
        }

        if ($param['type'] == self::TYPE_HELP_GIFT) {
            $this->type = "help";
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
