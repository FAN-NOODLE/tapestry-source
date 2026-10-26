<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\Service\V1\EnterprisePayment;

use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\EnterprisePaymentRequest;
use App\Exceptions\Error;
use App\Models\EnterprisePayment;
use App\Models\EnterprisePaymentBudget;
use App\Services\EnterprisePaymentBudgetService;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class SaveController extends BaseController
{
    /**
     * 编辑企业微信红包
     *
     * @return array
     */
    public function index()
    {
        $params = $this->request->input();
        // 验证数据
        (new EnterprisePaymentRequest())->save($params);

        // 验证金额
        if (empty($params['enterprise_payment_amount'])) {
            throw new Error(110007, '预计红包消费总金额不能为0');
        }

        $detail = EnterprisePayment::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('id', $params['id'])
            ->first();

        if (empty($detail)) {
            throw new Error(110009, '活动信息不存在');
        }

        // 如果活动非未上架或者活动已经删除则不能更新
        if ($detail->is_sale != EnterprisePayment::IS_SALE_DEFAULT || EnterprisePayment::IS_DELETED_YES == $detail->is_deleted) {
            throw new Error(110010, '活动不能被编辑');
        }

        // 初始化
        $incomeType = EnterprisePaymentBudget::INCOME_TYPE_OUT;
        // 如果原始金额比更改后金额大
        if ($detail->enterprise_payment_amount > $params['enterprise_payment_amount']) {
            $amount     = $detail->enterprise_payment_amount - $params['enterprise_payment_amount'];
            $incomeType = EnterprisePaymentBudget::INCOME_TYPE_IN;
        } else {
            $amount = $params['enterprise_payment_amount'] - $detail->enterprise_payment_amount;
        }

        // 获取当前充值账户信息
        $info       = (new EnterprisePaymentBudgetService())->getBudgetInfo($params['enterprise_id']);
        $oldBalance = $info->balance ?? 0;
        if ($oldBalance < $amount) {
            throw new Error(110008, '预计红包消费总金额不能大于当前账户余额');
        }

        // 添加
        DB::beginTransaction();
        try {

            // 如果有金额需要扣款或者返款
            if ($amount > 0) {
                // 默认类型为活动扣款
                $type    = EnterprisePaymentBudget::EVENT_DEDUCTION_TYPE;
                $balance = $oldBalance - $amount;
                if ($incomeType == EnterprisePaymentBudget::INCOME_TYPE_IN) {
                    $balance = $oldBalance + $amount;
                    $type    = EnterprisePaymentBudget::ACTIVITY_REFUND_TYPE;
                }

                $BudgetData = [
                    'enterprise_id' => $params['enterprise_id'],
                    'number'        => $detail->number,
                    'amount'        => $amount,
                    'balance'       => $balance,
                    'income_type'   => $incomeType,
                    'type'          => $type
                ];

                EnterprisePaymentBudget::query()
                    ->create($BudgetData);
            }

            // 更新活动
            EnterprisePayment::query()
                ->where('id', $params['id'])
                ->update([
                    'enterprise_id'             => $params['enterprise_id'],
                    'name'                      => $params['name'],
                    'type'                      => $params['type'],
                    'min_amount'                => $params['min_amount'] ?? 0,
                    'max_amount'                => $params['max_amount'] ?? 0,
                    'fixed_amount'              => $params['fixed_amount'],
                    'enterprise_payment_amount' => $params['enterprise_payment_amount'],
                    'quantity'                  => $params['quantity'],
                    'is_old_user_receive'       => $params['is_old_user_receive'],
                    'qr_code_url'               => $params['qr_code_url'],
                    'send_name'                 => $params['send_name'],
                    'wishing'                   => $params['wishing'],
                    'remark'                    => $params['remark'],
                ]);

            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            Log::info('编辑红包活动失败: ' . $e->getMessage());
            notice('100000', '编辑活动失败');

        } catch (\Throwable $e) {
            DB::rollBack();
            Log::info('编辑红包活动失败: ' . $e->getMessage());
            notice('100000', '编辑活动失败');
        }

        return success(['data' => []]);
    }
}
