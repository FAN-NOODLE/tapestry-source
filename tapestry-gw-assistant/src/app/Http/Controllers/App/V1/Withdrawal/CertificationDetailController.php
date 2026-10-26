<?php
/**
 * 提现申请用户认证信息详情
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

class CertificationDetailController extends Controller
{
    /**
     * 提现申请用户认证信息详情
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

        //接收参数
        $param = [
            'worker_id' => $workerDetail['id']
        ];
        $param['enterprise_id'] = $enterpriseId;
        $param['mem_type'] = $this->userType;

        $result = (new Withdrawal())->certificationDetail($param);

        return success($result['data'] ?? $result);
    }
    /**
     * 提现认证 - 详情
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   Yang Hui
     * @date     2020年10月16日17:1:28
     */
    public function certificationFileUrl(Request $request)
    {
        //接收参数
        $param = [
            'id' => 'id',
        ];
        $param['mem_type'] = $this->userType;
        $param['enterprise_id'] = $this->enterpriseId;

        $storeId = $this->request->post('store_id', 0);
        $enterpriseId = Auth::user()['enterprise_id'];
        $workerDetail = $this->storeWorkerDetail($enterpriseId, $storeId);
        $param['worker_id'] = $workerDetail['id'];

        $result = (new Withdrawal())->certificationFileUrl($param);

        return success($result['data'] ?? $result);
    }
}
