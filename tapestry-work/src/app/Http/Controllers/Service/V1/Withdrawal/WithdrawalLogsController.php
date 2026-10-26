<?php
/**
 * 提现 - 记录列表
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月16日16:56:4
 */

namespace App\Http\Controllers\Service\V1\Withdrawal;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\ShoppingGuideCashWithdrawal;
use App\Models\StoreWorkerCertification;
use Illuminate\Http\Request;
use App\Services\WithdrawalServices;

class WithdrawalLogsController extends BaseController
{
    /**
     * 提现 - 记录列表
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   Yang Hui
     * @date     2020年10月16日17:1:28
     */
    public function index(Request $request)
    {
        $params = $request->only([
            'enterprise_id',
            'worker_id',
            'page',
            'size'
        ]);
        $rules = [
            'enterprise_id' => 'required|integer',
            'worker_id' => 'required|integer'
        ];
        $message = [
            'enterprise_id.required' => '企业ID不存在',
            'enterprise_id.integer'  => '企业ID数据类型有误',
            'worker_id.required'      => '导购ID不存在',
            'worker_id.integer' => '导购ID数据类型错误'
        ];
        $this->validate(request(), $rules, $message);

        $params['page'] = $params['page'] ?? 1;
        $params['size'] = $params['size'] ?? 20;

        $data = (new WithdrawalServices())->getUserWithdrawalLogs($params);

        return success($data);
    }
}
