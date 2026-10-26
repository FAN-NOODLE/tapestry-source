<?php
/**
 * 提现申请 - 提现记录列表
 * Created by PhpStorm.
 * User: Yang Hui
 * DateTime: 2020年10月16日15:16:42
 */

namespace App\Http\Controllers\App\V1\Withdrawal;

use App\Http\Controllers\App\V1\BaseController as Controller;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use CrsSdk\Work\Withdrawal;
use Illuminate\Support\Facades\Auth;

class WithdrawalLogsController extends Controller
{
    /**
     * 提现申请 - 提现记录列表
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @author  YangHui
     * @date    2020-10-16 15:20:42
     */
    public function index(Request $request)
    {
        $storeId = $this->request->post('store_id', 0);
        $enterpriseId = Auth::user()['enterprise_id'];
        $workerDetail = $this->storeWorkerDetail($enterpriseId, $storeId);

        $param = $request->input();
        $param['worker_id'] = $workerDetail['id'];
        $param['enterprise_id'] = $enterpriseId;
        $param['mem_type'] = $this->userType;

        $result = (new Withdrawal())->withdrawalLogs($param);

        return success($result['data'] ?? $result);
    }
}
