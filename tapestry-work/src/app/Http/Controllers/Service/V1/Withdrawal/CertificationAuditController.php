<?php
/**
 * 提现认证 - 审核
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月16日16:56:4
 */

namespace App\Http\Controllers\Service\V1\Withdrawal;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Enterprise;
use App\Models\DistributionGoods;
use App\Models\StoreWorkerCertification;
use Illuminate\Http\Request;
use App\Services\WithdrawalServices;
use Illuminate\Support\Facades\Log;

class CertificationAuditController extends BaseController
{
    /**
     * 提现认证 - 审核
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   Yang Hui
     * @date     2020年10月16日17:1:28
     */
    public function index(Request $request)
    {
        $params = $request->input();
        $result = false;
        $message = "";
        $certification = [];
        try {
            $certification = (new WithdrawalServices())->getUserCertificationInfo($params);
        } catch (\Exception $e) {
            Log::info($params['worker_id'] . '用户身份认证-' . $e->getMessage());
            $message = $e->getMessage();
            $certification = [];
        }
        if (empty($certification)) {
            return success([
                'data' => [
                    "status" => boolval($result),
                    "message" => empty($message) ? '未查询到认证信息' : $message
                ]
            ]);
        }
        //进行信息验证操作,如果为true，则认证通过
        if ($certification['auth_status'] == StoreWorkerCertification::AUTH_STATUS_PASS) {
            $result = true;
            $message = "认证成功";
        } else {
            $result = false;
            $message = empty($message) ? '认证失败' : $message;
        }
        return success([
            'data' => [
                "status" => boolval($result),
                "message" => $message
            ]
        ]);
    }
}
