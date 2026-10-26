<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/19 0019
 * Time: 20:28
 */

namespace App\Http\Controllers\V1\ActivityCommon\Common;

use App\Exceptions\Error;
use App\Http\Controllers\V1\ActivityCommon\ActivityBaseController;

class AddController extends ActivityBaseController
{

    const UP_TYPE_SELF =1;
    const UP_TYPE_CRON =2;

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
        $params['up_type'] = in_array($params['up_type']  ?? self::UP_TYPE_SELF, [self::UP_TYPE_SELF, self::UP_TYPE_CRON]) ? $params['up_type'] ?? self::UP_TYPE_SELF : self::UP_TYPE_SELF;
        try {
            $activityInfo = $this->activityApi->getActivity($params);
        } catch (\Exception $e) {;
            $occupySpuData = [];
            if ($e->getCode() == 200071) {
                $occupySpuData = $e->data['data'];
            }
            $this->notice($e, $params, null, ['occupy_spu' => $occupySpuData]);
        }

        return success($activityInfo);
    }

}
