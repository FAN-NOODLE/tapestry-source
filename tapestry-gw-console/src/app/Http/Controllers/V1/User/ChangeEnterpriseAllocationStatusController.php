<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\User;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\EnterpriseExtend;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class ChangeEnterpriseAllocationStatusController extends BaseController
{
    /**
     * 更新企业分配客户状态
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        // 分配状态：1、开启，2、关闭
        $status = trim($request->post('status', 0));
        // 红包开关 1、开启，2、关闭
        $couponStatus = trim($request->post('coupon_status', 0));

        if (empty($status) && empty($couponStatus)) {
            error(320107);
        }

        $returnData = (new EnterpriseExtend)->changeAllocationStatus(['status' => $status, 'coupon_status' => $couponStatus, 'enterprise_id' => $this->enterpriseId]);

        $switchStatus = $status ?: $couponStatus;

        $switch = $switchStatus == 1 ? '开启' : '关闭';

        if ($couponStatus) {
            $this->writeOperatorLog('更新企业红包应用状态:' . $switch);
        } else {
            $this->writeOperatorLog('更新企业分配客户状态:' . $switch);
        }

        return success($returnData);
    }
}
