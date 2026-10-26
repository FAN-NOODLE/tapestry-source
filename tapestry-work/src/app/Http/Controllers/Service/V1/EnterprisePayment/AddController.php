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

class AddController extends BaseController
{
    /**
     * 添加企业微信红包
     *
     * @return array
     */
    public function index()
    {
        $params = $this->request->input();
        // 验证数据
        (new EnterprisePaymentRequest())->add($params);

        // 验证金额
        if (empty($params['enterprise_payment_amount'])) {
            throw new Error(110007, '预计红包消费总金额不能为0');
        }

        // 获取当前充值账户信息
        $info    = (new EnterprisePaymentBudgetService())->getBudgetInfo($params['enterprise_id']);
        $balance = $info->balance ?? 0;
        if ($balance < $params['enterprise_payment_amount']) {
            throw new Error(110008, '预计红包消费总金额不能大于当前账户余额');
        }

        $model = EnterprisePayment::query();

        $count = $model->where('enterprise_id', $params['enterprise_id'])
            ->count();

        // 添加
        DB::beginTransaction();
        try {
            // 获取版本号
            $number = 'YLHB' . sprintf("%05d", $count + 1);
            // 如果有金额则需要扣款
            if ($params['enterprise_payment_amount'] > 0) {
                $BudgetData = [
                    'enterprise_id' => $params['enterprise_id'],
                    'number'        => $number,
                    'amount'        => $params['enterprise_payment_amount'],
                    'balance'       => $balance - $params['enterprise_payment_amount'],
                    'income_type'   => EnterprisePaymentBudget::INCOME_TYPE_OUT,
                    'type'          => EnterprisePaymentBudget::EVENT_DEDUCTION_TYPE
                ];

                EnterprisePaymentBudget::query()
                    ->create($BudgetData);
            }

            // 创建活动
            $model->create([
                'enterprise_id'             => $params['enterprise_id'],
                'name'                      => $params['name'],
                'type'                      => $params['type'],
                'number'                    => $number,
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
            Log::info('创建红包活动失败: ' . $e->getMessage());
            notice('100000', '创建红包活动失败');

        } catch (\Throwable $e) {
            DB::rollBack();
            Log::info('创建红包活动失败: ' . $e->getMessage());
            notice('100000', '创建红包活动失败');
        }

        return success(['data' => []]);
    }
}
