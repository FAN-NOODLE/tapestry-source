<?php
/**
 * 提现申请 - 添加记录
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

class WithdrawalAddController extends Controller
{
    /**
     * 提现申请 - 添加记录
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @author  YangHui
     * @date    2020-10-16 15:20:42
     */
    public function index(Request $request)
    {
        //接收参数
        $param = [
            'application_amount' => $request->input('application_amount'),
        ];

        $storeId = $this->request->post('store_id', 0);
        $enterpriseId = Auth::user()['enterprise_id'];
        $workerDetail = $this->storeWorkerDetail($enterpriseId, $storeId);

        $param['worker_id'] = $workerDetail['id'];
        $param['enterprise_id'] = $enterpriseId;
        $param['mem_type'] = $this->userType;

        Log::info('用户认证添加-请求参数', [$param]);

        $result = (new Withdrawal())->withdrawalAdd($param);

        return success($result['data'] ?? $result);
    }
}
