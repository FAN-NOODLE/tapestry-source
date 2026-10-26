<?php
/**
 * 提现申请用户认证
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

class CertificationAddController extends Controller
{
    /**
     * 提现申请用户认证
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @author  YangHui
     * @date    2020-10-16 15:20:42
     */
    public function index(Request $request)
    {
        //接收参数
        $param = [
            'real_name' => $request->input('real_name'),
            'mobile' => $request->input('mobile'),
            'id_type' => $request->input('id_type'),
            'id_code' => $request->input('id_code'),
            'bank_card_username' => $request->input('bank_card_username'),
            'bank_card_code' => $request->input('bank_card_code'),
            'bank_name' => $request->input('bank_name'),
            'agree_deal' => $request->input('agree_deal')
        ];
        $storeId = $this->request->post('store_id', 0);
        $enterpriseId = Auth::user()['enterprise_id'];
        $workerDetail = $this->storeWorkerDetail($enterpriseId, $storeId);

        $param['enterprise_id'] = $enterpriseId;
        $param['worker_id'] = $workerDetail['id'];
        $param['mem_type'] = $this->userType;

        Log::info('用户认证添加-请求参数', [$param]);

        $result = (new Withdrawal())->certificationAdd($param);

        return success($result['data'] ?? $result);
    }
}
