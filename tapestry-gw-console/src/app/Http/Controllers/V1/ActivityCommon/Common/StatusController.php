<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/19 0019
 * Time: 20:28
 */

namespace App\Http\Controllers\V1\ActivityCommon\Common;

use App\Http\Controllers\V1\ActivityCommon\ActivityBaseController;

class StatusController extends ActivityBaseController
{
    function index()
    {
        $params = $this->params;

        if ($params['status'] == self::STATUS_UP) { // 上架
            $params['method'] = 'status-up';
            $script = '上架活动';
        } else if ($params['status'] == self::STATUS_DOWN) { // 下架
            $params['method'] = 'status-down';
            $script = '下架活动';
        } else if (in_array($params['status'], [self::STATUS_SUSPEND, self::STATUS_RECOVER])) {
            $params['method'] = 'suspend';
            $script = '暂停或恢复活动';
        }

        try {
            $this->activityApi->getActivity($params);
        } catch (\Exception $e) {
            $this->notice($e, $params);
        }

        $this->writeOperatorLog($script.':'.json_encode($params, JSON_UNESCAPED_UNICODE));

        return success([]);
    }

}
