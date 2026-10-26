<?php
/**
 * 提现 - 申请
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月16日16:56:4
 */

namespace App\Http\Controllers\Service\V1\Withdrawal;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Enterprise;
use App\Models\DistributionGoods;
use App\Models\StoreWorkerCertification;
use Illuminate\Http\Request;
use App\Services\WithdrawalServices;
use App\Models\ShoppingGuideCashWithdrawal;
use Illuminate\Support\Facades\DB;
use App\Models\WithdrawalApplication;
use Illuminate\Support\Facades\Redis;

class WithdrawalAddController extends BaseController
{
    //最高提现金额：分
    const MIN_APPLICATION_AMOUNT = 100;
    //最高提现金额：分
    const MAX_APPLICATION_AMOUNT = 5000 * 100;
    /**
     * 提现 - 申请
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
            'application_amount',
            'mem_type'
        ]);

        $rules = [
            'enterprise_id' => 'required|integer',
            'worker_id' => 'required|integer',
            'application_amount' => 'required|integer|min:' .
                self::MIN_APPLICATION_AMOUNT . '|max:' . self::MAX_APPLICATION_AMOUNT
        ];

        $message = [
            'enterprise_id.required' => '企业ID不存在',
            'enterprise_id.integer'  => '企业ID数据类型有误',
            'worker_id.required'      => '导购ID不存在',
            'worker_id.integer' => '导购ID数据类型错误',
            'application_amount.required' => '请输入提现金额',
            'application_amount.integer' => '提现金额格式错误，单位分制',
            'application_amount.min' => '提现金额不能小于' . intval(self::MIN_APPLICATION_AMOUNT/100) . '元',
            'application_amount.max' => '提现金额不能超' . intval(self::MAX_APPLICATION_AMOUNT/100) . '元'
        ];
        $this->validate(request(), $rules, $message);

        $memType = $params['mem_type'];
        if (!in_array($memType, [1, 2])) {
            $memType = 1;
        }

        //这里需要加锁，防止重复提交
        $lockKey = 'work_withdrawal_add_'.$params['worker_id'];
        // 查询锁
        $getCache = Redis::get($lockKey);
        // 检查是否已经超出时间 超出时间释放锁（防止死锁）
        if (!empty($getCache) && $getCache < time()) {
            Redis::del([$lockKey]);
        }

        //验证用户是否通过了身份验证，并获取验证信息
        $certificationInfo = (new WithdrawalServices())->getUserCertificationInfo($params);
        //如果未认证或者未认证通过，则提示认证
        if (empty($certificationInfo)) {
            return error(1100004);
        }

        //获取可提现信息 available_cash_amount
        $balanceInfo = ShoppingGuideCashWithdrawal::query()
            ->where('worker_id', $params['worker_id'])
            ->where('enterprise_id', $params['enterprise_id'])
            ->first();
        if (empty($balanceInfo) || $balanceInfo->available_cash_amount < $params['application_amount']) {
            return error(1100005);
        }

        //检验当天可提现金额
        if ($certificationInfo['today_withdrawaled']
            && ($certificationInfo['today_withdrawaled'] + $params['application_amount']) > (5000 * 100)) {
            return error(1100010);
        }

        // 开始加锁
        $setCache = Redis::setnx($lockKey, time() + 60);
        // 如果已经发起直接跳出
        if (!$setCache) {
            throw new Error(1100007, '提交已在处理中，请勿重复提交');
        }

        //开始事务，写提现申请表，更改提现结算表
        DB::beginTransaction();
        try {
            $updateResult = ShoppingGuideCashWithdrawal::query()
                ->where('worker_id', $params['worker_id'])
                ->where('enterprise_id', $params['enterprise_id'])
                ->where('available_cash_amount', $balanceInfo->available_cash_amount)
                ->update([
                    'available_cash_amount' => $balanceInfo->available_cash_amount-$params['application_amount'],
                    'frozen_amount' => $balanceInfo->frozen_amount + $params['application_amount']
                ]);
            if (!$updateResult) {
                throw new Error(1100006, '');
            }

            $application = new WithdrawalApplication();
            $application->id = idGenerate();
            $application->enterprise_id = $params['enterprise_id'];
            $application->worker_id = $params['worker_id'];
            $application->application_amount = $params['application_amount'];
            $application->request_initiated_at = time();
            $application->status = WithdrawalApplication::WAIT_APPLICATION_STATUS;
            $application->app_type = $memType;
            $application->save();

            DB::commit();
            Redis::del([$lockKey]);
            return success();
        } catch (\Exception $e) {
            DB::rollBack();
            Redis::del([$lockKey]);
            return error(1100006);
        } finally {
            Redis::del([$lockKey]);
        }
    }
}
