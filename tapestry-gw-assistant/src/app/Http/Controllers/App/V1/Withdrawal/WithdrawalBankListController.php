<?php
/**
 *提现 - 银行列表
 * Created by PhpStorm.
 * User: Yang Hui
 * DateTime: 2020年10月16日15:16:42
 */

namespace App\Http\Controllers\App\V1\Withdrawal;

use App\Http\Controllers\App\V1\BaseController as Controller;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use CrsSdk\Work\Withdrawal;

class WithdrawalBankListController extends Controller
{
    /**
     * 提现 - 银行列表
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @author  YangHui
     * @date    2020-10-16 15:20:42
     */
    public function index(Request $request)
    {
        $param = $request->input();

        $param['enterprise_id'] = $this->enterpriseId;

        $result = (new Withdrawal())->withdrawalBankList($param);

        return success($result['data'] ?? $result);
    }
}
