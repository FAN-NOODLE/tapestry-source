<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/13 0013
 * Time: 14:35
 */

namespace App\Http\Controllers\V1\ActivityCommon;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroActivity\Activity;

class ActivityListController extends BaseController
{

    /**
     * 选择活动列表
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $param = $this->request->input();
        $param['enterprise_id'] = $this->enterpriseId;
        // 只获取已上架 和 已结束的活动
        $param['activity_status'] = [ActivityBaseController::STATUS_UP, ActivityBaseController::STATUS_OVER];
        try {
            $list = (new Activity)->activityList($param);
        } catch (\Exception $e) {
            throw new Error($e->getCode(), $e->getMessage());
        }

        return success($list ?? []);
    }

}
