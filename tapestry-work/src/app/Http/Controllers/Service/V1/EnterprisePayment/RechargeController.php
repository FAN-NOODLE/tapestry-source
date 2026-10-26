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
use App\Models\EnterprisePaymentBudget;
use App\Services\EnterprisePaymentBudgetService;
use Illuminate\Http\Request;

class RechargeController extends BaseController
{
    /**
     * 充值接口
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->post('enterprise_id');
        $amount = $request->post('amount', 0);
        $voucherUrl = $request->post('voucher_url', '');
        $incomeType = $request->post('income_type', EnterprisePaymentBudget::INCOME_TYPE_IN);

        if (empty($enterpriseId)) {
            throw new Error(110003, '企业ID不能为空');
        }

        if (empty($amount)) {
            throw new Error(110003, '充值金额不能为空');
        }

        if (!in_array($incomeType, [EnterprisePaymentBudget::INCOME_TYPE_IN, EnterprisePaymentBudget::INCOME_TYPE_OUT])) {
            throw new Error(110003, '充值类型只能是充值或扣款');
        }

        $info = (new EnterprisePaymentBudgetService())->getBudgetInfo($enterpriseId);
        // 当前余额
        $balance = $info->balance ?? 0;

        if ($incomeType == EnterprisePaymentBudget::INCOME_TYPE_IN) {
            // 充值
            $balance += $amount;
        } else {
            // 扣款
            $balance -= $amount;
        }

        $balance = ($balance < 0) ? 0 : $balance;

        EnterprisePaymentBudget::query()
            ->create([
                'enterprise_id' => $enterpriseId,
                'amount' => $amount,
                'balance' => $balance,
                'income_type' => $incomeType,
                'voucher_url' => $voucherUrl
            ]);

        return success();
    }
}
