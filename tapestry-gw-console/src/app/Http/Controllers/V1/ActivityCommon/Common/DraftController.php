<?php


namespace App\Http\Controllers\V1\ActivityCommon\Common;


use App\Http\Controllers\V1\ActivityCommon\ActivityBaseController;

class DraftController extends ActivityBaseController
{
    function index()
    {
        $params = $this->params;
        // 检验好友助力邮费状态开关 1：开启  2：关闭
        if ($params['type'] == 'help' && $params['activity_type'] == 6) {
            if (!in_array($params['postage_status'], [1,3])) {
                notice(920002);
            }
            $this->typeCode = 6;
        } elseif ($params['type'] == 'help' && $params['activity_type'] == 7) {
            $params['postage_status'] = 2; // 买赠
            $this->typeCode = 7;
        }
        // 检验好友助力邮费
        if ($params['type'] == 'help' && $params['postage'] <= 0 && $params['activity_type'] == 6) {
            notice(920001);
        }

        $params['activity_ids'] = $params['activity_ids'] ?? [];
        try {
            $activityInfo = $this->activityApi->getActivity($params);
        } catch (\Exception $e) {
            $this->notice($e, $params);
        }

        return success($activityInfo);
    }

}