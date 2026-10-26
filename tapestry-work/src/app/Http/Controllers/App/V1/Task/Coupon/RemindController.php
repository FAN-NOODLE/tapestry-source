<?php

namespace App\Http\Controllers\App\V1\Task\Coupon;

use App\Common\Constant;
use App\Http\Controllers\App\V1\Task\TaskBaseController;
use App\Models\Task;
use Illuminate\Http\Request;

class RemindController extends TaskBaseController
{

    /**
     * 优惠券提醒
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {

        $customerId = $request->get('ec_userid', '');

        // 用户ID
        if (empty($customerId)) {
            error(200011);
        }

        // 写入优惠券提醒记录
        $this->saveLog(0, $customerId, Task::MEMBER_COUPON, Constant::TASK_SHOW_STATUS);

        return success(['data' => []]);
    }
}

