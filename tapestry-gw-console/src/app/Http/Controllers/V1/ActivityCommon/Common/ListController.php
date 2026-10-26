<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/19 0019
 * Time: 20:28
 */

namespace App\Http\Controllers\V1\ActivityCommon\Common;

use App\Http\Controllers\V1\ActivityCommon\ActivityBaseController;
use App\Traits\ActivityTrait;

class ListController extends ActivityBaseController
{
    use ActivityTrait;

    function index()
    {
        $params = $this->params;
        try {
            $activityList = $this->activityApi->getActivity($params);
        } catch (\Exception $e) {
            $this->notice($e, $params);
        }
        // 组装数据
        foreach ($activityList['list'] as &$activity) {
            $activity['activity_spu'] = $activity['spu_list'];
            foreach ($activity['activity_spu'] as &$spu) {
                unset($spu['path']);
            }
            $this->ActivityArrayToJson($activity, 'activity_sku', 'spu_name');
            $this->ActivityArrayToJson($activity, 'activity_sku', 'sku_name', true);
            $this->ActivityArrayToJson($activity, 'spu_list', 'qr_code');
            $this->ActivityArrayToJson($activity, 'spu_list', 'spu_id');
            $this->ActivityArrayToJson($activity, 'spu_list', 'path', true);
        }
        unset($activity);
        return success($activityList);
    }

}
