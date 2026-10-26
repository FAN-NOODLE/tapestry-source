<?php

namespace App\Http\Controllers\Service\V1\Store;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Micros\Common;
use App\Models\StoreDepartment;

class CloseOrOpenDepartmentController extends BaseController
{

    /**
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2022/2/7 15:45
     */
    public function index(): array
    {
        $id = $this->request->input('id', 0);
        $isClosed = $this->request->input('is_closed', 0);
        $enterpriseId = $this->request->input('enterprise_id', 0);
        if (!$id) {
            throw new Error(600016, '参数错误');
        }
        $info = StoreDepartment::query()->find($id);
        if (!empty($info)) {
            $info->closed_at = $isClosed ? time() : 0;
            $info->save();
            //通知管理层组织架构门店闭店或者重新打开
            Common::microScrm()->setApi("micro-scrm/v1/dashboard/area/store-close-or-open")->post([
                'third_id' => $info->wechat_department_id,
                'is_closed' => (bool)$isClosed,
                'enterprise_id' => $enterpriseId
            ]);
        }
        return success([]);
    }
}
