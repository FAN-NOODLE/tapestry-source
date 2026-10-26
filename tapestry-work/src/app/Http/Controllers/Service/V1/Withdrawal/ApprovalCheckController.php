<?php
/**
 * 提现 - 审批
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月16日16:56:4
 */

namespace App\Http\Controllers\Service\V1\Withdrawal;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\ShoppingGuideCashWithdrawal;
use App\Models\StoreWorkerCertification;
use App\Models\WithdrawalApplication;
use App\Services\EnterpriseService;
use Illuminate\Http\Request;
use App\Services\WithdrawalServices;
use Illuminate\Validation\Rule;

class ApprovalCheckController extends BaseController
{
    /**
     * 提现 - 审批
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   Yang Hui
     * @date     2020年10月16日17:1:28
     */
    public function index(Request $request)
    {
        $params = $request->only([
            'ids',
            'enterprise_id',
            'admin_id',
            'action',
            'reason'
        ]);
        $rules = [
            'ids' => 'required|array',
            'enterprise_id' => 'required|integer',
            'admin_id' => 'required|integer',
            'action' => [
                'required',
                Rule::in([
                    'pass', //通过
                    'refuse' //拒绝
                ])
            ]
        ];
        $message = [
            'ids.required' => '参数id不存在',
            'ids.array' => '参数id格式错误',
            'enterprise_id.required' => '企业ID不存在',
            'enterprise_id.integer'  => '企业ID数据类型有误',
            'admin_id.required'      => '管理员ID不存在',
            'admin_id.integer' => '导购ID数据类型错误',
            'action.required' => '请求执行的动作错误',
            'action.in' => '无法识别请求的执行动作'
        ];
        $this->validate(request(), $rules, $message);

        // 查询企业信息
        $enterpriseData = (new EnterpriseService())->getDetailByEnterpriseId($params['enterprise_id']);
        $result = (new WithdrawalServices())->withdrawalAudit($params, $enterpriseData['hash']);

        $returnData = [
            'data' => $result
        ];
        return success($returnData);
    }
}
