<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/19 0019
 * Time: 20:28
 */

namespace App\Http\Controllers\V1\ActivityCommon\Common;

use App\Http\Controllers\V1\ActivityCommon\ActivityBaseController;

class DeleteController extends ActivityBaseController
{

    function index()
    {
        $params = $this->params;

        // 删除活动
        try {
            $this->activityApi->getActivity($params);
        } catch (\Exception $e) {
            $this->notice($e, $params);
        }

        $this->writeOperatorLog("删除活动:". json_encode($params, JSON_UNESCAPED_UNICODE));

        return success([]);
    }

}
