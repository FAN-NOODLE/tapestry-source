<?php
/**
 * 提现申请用户认证审核
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

class CertificationAuditController extends Controller
{
    /**
     * 提现申请用户认证审核
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @author  YangHui
     * @date    2020-10-16 15:20:42
     */
    public function index(Request $request)
    {
        //接收参数
        $param = [
            'id' => $request->input('id'),
        ];
        $storeId = $this->request->post('store_id', 0);
        $enterpriseId = Auth::user()['enterprise_id'];
        $workerDetail = $this->storeWorkerDetail($enterpriseId, $storeId);

        $param['worker_id'] = $workerDetail['id'];
        $param['enterprise_id'] = $this->enterpriseId;
        $param['mem_type'] = $this->userType;

        $result = (new Withdrawal())->certificationAudit($param);

        return success($result['data'] ?? $result);
    }
}
