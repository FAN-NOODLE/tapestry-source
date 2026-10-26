<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\Service\V1\EnterprisePayment;

use App\Http\Controllers\Service\V1\BaseController;
use App\Exceptions\Error;
use App\Models\EnterprisePaymentRecord;
use App\Services\EnterprisePaymentBudgetService;

class BudgetTotalController extends BaseController
{
    /**
     * 企业充值对账统计
     */
    public function index()
    {
        $enterpriseId = $this->request->input('enterprise_id', 0);

        if (empty($enterpriseId)) {
            throw new Error(110003, '企业ID不能为空');
        }

        // 获取当前充值账户信息
        $info = (new EnterprisePaymentBudgetService())->getBudgetInfo($enterpriseId);

        // 获取发送金额
        $record = EnterprisePaymentRecord::query()
            ->where('enterprise_id', $enterpriseId)
            ->get(['amount'])
            ->toArray();

        $sendMoney = 0;

        if (!empty($record)) {
            $sendMoney = array_sum(array_column($record, 'amount'));
        }

        $data = [
            'surplus_money' => $info->balance ?? 0,
            'send_num'      => count($record),
            'send_money'    => $sendMoney,
        ];

        return success([
            'data' => $data,
        ]);
    }
}
