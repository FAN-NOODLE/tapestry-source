<?php
/**
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月17日16:36:43
 */
namespace App\Services;

use App\Common\Constant;
use App\Models\PartnerMember;
use App\Models\ShoppingGuideCashWithdrawal;
use App\Models\StoreWorkerCertification;
use App\Models\WithdrawalApplication;
use App\Models\WithdrawalSetting;
use App\Exceptions\Error;
use Illuminate\Support\Facades\DB;

class WithdrawalServices extends BaseService
{
    /**
     * 获取用户认证信息
     * @param int|string $workerId
     * @return array
     */
    public function getUserCertificationInfo($params): array
    {

        //WORK_ID查询
        $certList = StoreWorkerCertification::query()
            ->where(function ($query) use ($params) {
                $query->where('worker_id', $params['worker_id'])
                    ->orWhere('mobile', $params['mobile'] ?? '-1');
            })
            ->get();
        //如果为空，直接返回
        if ($certList->isEmpty()) {
            return [];
        }
        $info = [];
        //检查当前企业他是否已经认证
        $certArr = $certList->toArray();
        $certArr = array_column($certArr, null, 'enterprise_id');
        //如果不存在当前企业下的用户信息，则重新复制一份
        if (!isset($certArr[ $params['enterprise_id'] ])
            || $certArr[ $params['enterprise_id'] ]['worker_id'] != $params['worker_id']
        ) {
            $data = current($certArr);
            $data['id'] = idGenerate();
            $data['enterprise_id'] = $params['enterprise_id'];
            $data['worker_id'] = $params['worker_id'];
            $data['mem_type'] = $params['mem_type'];
            $data['created_at'] = time();
            $data['updated_at'] = time();
            StoreWorkerCertification::query()->insert($data);
            $info = StoreWorkerCertification::query()->where('id', $data['id'])->first();
        } else {
            $info = StoreWorkerCertification::query()
                ->where('id', $certArr[$params['enterprise_id']]['id'])
                ->where('mem_type', $params['mem_type'])
                ->first();
        }
        if (empty($info)) {
            return [];
        }
        //如果认证状态为默认，但是有数据，则进行一次验证
        if ($info->auth_status == StoreWorkerCertification::AUTH_STATUS_DEFAULT) {
            $result = (new SettlementService())->signAgreement(
                $info->real_name,
                $info->id_type,
                $info->id_code,
                $info->mobile,
                $this->getWithdrawalConfig($params['enterprise_id'])['transfer_no']??''
            );
            if ($result) {
                $info->auth_status = StoreWorkerCertification::AUTH_STATUS_PASS;
                $info->status = StoreWorkerCertification::STATUS_ENABLED;
                $info->save();
            } else {
                return [];
            }
        }
        //如果认证结果与状态依然不符合，返回空,否则返回数据
        if ($info->auth_status == StoreWorkerCertification::AUTH_STATUS_PASS
            && $info->status == StoreWorkerCertification::STATUS_ENABLED
        ) {
            $info = $info->toArray();

            //当日已申请提现或者已提现金额
            $startTime = strtotime(date('Y-m-d'));
            $endTime = strtotime($startTime . ' 23:59:59');
            $info['today_withdrawaled'] = WithdrawalApplication::query()
                ->where('worker_id', $params['worker_id'])
                ->where('enterprise_id', $params['enterprise_id'])
                ->whereBetween('request_handle_at', [$startTime, $endTime])
                ->whereIn('status', [
                    WithdrawalApplication::WAIT_APPLICATION_STATUS,
                    WithdrawalApplication::WITHDRAW_SUCCESS_STATUS,
                    WithdrawalApplication::WITHDRAW_IN_STATUS
                ])
                ->sum('application_amount');

            return $info;
        }
        return [];
    }

    /**
     * 获取认证文件地址信息
     * @param $params 参数
     * @return string
     */
    public function getUserCertificationFileUrl($params): string
    {
        $info = $this->getUserCertificationInfo($params);
        if (empty($info)) {
            return '';
        }
        $result = (new SettlementService())->queryAgreement(
            $info['real_name'],
            $info['id_code'],
            $this->getWithdrawalConfig($params['enterprise_id'])['transfer_no']??''
        );
        if (isset($result['is_signed']) && $result['is_signed'] == true) {
            return $result['agreement_pic'];
        } else {
            return '';
        }
    }

    /**
     * 获取用户提现记录列表
     * @param $params
     * @return array
     */
    public function getUserWithdrawalLogs($params): array
    {
        $params['page'] = $params['page'] ?? 1;
        $params['size'] = $params['size'] ?? 20;

        $result = [];
        $total = WithdrawalApplication::query()
            ->where('worker_id', $params['worker_id'])
            ->where('enterprise_id', $params['enterprise_id'])
            ->count();
        if ($total > 0) {
            $result = WithdrawalApplication::query()
                ->where('worker_id', $params['worker_id'])
                ->where('enterprise_id', $params['enterprise_id'])
                ->orderBy('created_at', 'desc')
                ->forPage($params['page'], $params['size'])
                ->get();
        }

        $result = empty($result) ? [] : $result->toArray();

        return [
            'list' => $result,
            'total' => intval($total),
            'page'  => intval($params['page']),
            'size'  => intval($params['size'])
        ];
    }

    /**
     * 获取当前用户的信息
     * @param $params
     * @return array|mixed
     * @throws \Exception
     */
    public function getCurrentWorkerInfo($params)
    {
        $workerCondition = [
            'enterprise_id' => $params['enterprise_id'],
            'ids' => [$params['worker_id']]
        ];

        //如果是合伙人，则返回合伙人信息，如果是导购刚去查询导购信息
        $workerData = (new PartnerMemberService())->getPartnerInfo([
            'enterprise_id' => $params['enterprise_id'],
            'id' => $params['worker_id']
        ]);
        if (empty($workerData)) {
            $workerData = StoreWorkerService::list($workerCondition)['list'];
            $workerData = !empty($workerData) ? current($workerData) : [];
        } else {
            $workerData['type'] = Constant::STORE_WORKER_TYPE_PARTNER;
        }

        return $workerData;
    }

    /**
     * 提现审核
     * @param $params['ids] id数组
     * @param $params['action] 动作
     * @param $params['reason] 驳回原因
     * @return int
     */
    public function withdrawalAudit($params, $enterpriseHash): array
    {
        if (!isset($params['action']) || !in_array($params['action'], ['pass', 'refuse'])) {
            return false;
        }
        $withdrawalLogs = WithdrawalApplication::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->whereIn('id', $params['ids'])
            ->where('status', WithdrawalApplication::WAIT_APPLICATION_STATUS)
            ->get();
        if (empty($withdrawalLogs)) {
            throw new Error(1100002, '未查询到对应的可操作数据');
        }
        //查询用户身份验证信息
        $certificationInfoArr = StoreWorkerCertification::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->whereIn('worker_id', array_column($withdrawalLogs->toArray(), 'worker_id'))
            ->where('auth_status', StoreWorkerCertification::AUTH_STATUS_PASS)
            ->where('status', StoreWorkerCertification::STATUS_ENABLED)
            ->get()
            ->toArray();
        $certificationInfoArr = array_column($certificationInfoArr, null, 'worker_id');

        $returnArr = [];
        $withdrawalConfig = $this->getWithdrawalConfig($params['enterprise_id']);

        //如果是拒绝，则改变状态，并写原因
        if ($params['action'] == 'refuse') {
            $returnArr = $this->auditRefuse($withdrawalLogs, $params);
        }

        //如果是通过，打包数据，通过第三方接口提交提现申请，然后状态到达转账号
        if ($params['action'] == 'pass') {
            $balanceArr = [];
            $withdrawalLogArr = [];
            foreach ($withdrawalLogs as $info) {
                //如果对用户信息审核不通过则直接返回
                if (!isset($certificationInfoArr[$info->worker_id])) {
                    $returnArr = $this->auditFail([$info], $params, '身份认证验证失败');
                    continue;
                } else {
                    $withdrawalLogArr[] = $info;
                    //整合需要提交给第三方打款平台的数据
                    $balanceArr[] = [
                        'project_code' => $withdrawalConfig['transfer_no']??'',
                        'order_random_code' => strval($info->id),
                        'name' => strval($certificationInfoArr[$info->worker_id]['bank_card_username']),
                        'certificate_type' => intval($certificationInfoArr[$info->worker_id]['id_type']),
                        'certificate_num' => strval($certificationInfoArr[$info->worker_id]['id_code']),
                        'phone_num' => strval($certificationInfoArr[$info->worker_id]['mobile']),
                        'bank_name' => strval($certificationInfoArr[$info->worker_id]['bank_name']),
                        'bank_card_num' => strval($certificationInfoArr[$info->worker_id]['bank_card_code']),
                        'settle_amount' => intval($info->application_amount)
                    ];
                }
            }
            if (!empty($balanceArr)) {
                try {
                    (new SettlementService())
                        ->batchSettleBalance($enterpriseHash,$withdrawalConfig['transfer_no']??'', $balanceArr);
                    $returnArr = $this->auditPass($withdrawalLogArr, $params);
                } catch (\Exception $e) {
                    //$returnArr = $this->auditFail($withdrawalLogArr, $params, $e->getMessage());
                    throw new Error(1100010, "批量结算发起失败", 500, $e);
                }
            }
        }

        return [
            'success_ids' => $returnArr['successIds']??[],
            'error_ids' => $returnArr['errorIds']??[]
        ];
    }

    /**
     * 驳回提现请求
     * @param $withdrawalLogs
     * @param $params
     * @return array
     */
    public function auditRefuse($withdrawalLogs, $params)
    {
        if (!isset($params['reason']) || empty($params['reason'])) {
            throw new Error(1100009, '请填写驳回原因');
        }

        //执行成功的ID数组
        $successIds = [];
        //执行失败的id数组
        $errorIds = [];

        DB::beginTransaction();
        try {
            foreach ($withdrawalLogs as $info) {
                //$info->enterprise_id = $params['enterprise_id'];
                $info->operator_id = $params['admin_id'] ?? 0;
                $info->status = WithdrawalApplication::WITHDRAW_REJECT_STATUS;
                $info->reason = $params['reason'];
                $info->request_handle_at = time();
                //保存成功后，将款项退回
                if ($info->save()) {
                    $this->changeUserAmount(
                        $params['enterprise_id'],
                        $info->worker_id,
                        $info->application_amount,
                        false
                    );

                    $successIds[] = $info->id;
                } else {
                    $errorIds[] = $info->id;
                }
            }
            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
        }


        return [
            'successIds' => $successIds,
            'errorIds' => $errorIds
        ];
    }

    /**
     * 打款失败
     * @param $withdrawalLogs
     * @param $params
     * @param $reason
     */
    public function auditFail($withdrawalLogs, $params, $reason = '打款失败')
    {
        //执行成功的ID数组
        $successIds = [];
        //执行失败的id数组
        $errorIds = [];
        DB::beginTransaction();
        try {
            foreach ($withdrawalLogs as $info) {
                //$info->enterprise_id = $params['enterprise_id'];
                $info->operator_id = $params['admin_id']??0;
                $info->status = WithdrawalApplication::WITHDRAW_FAIL_STATUS;
                $info->reason = $reason;
                $info->request_handle_at = time();
                //保存成功后，将款项退回
                if ($info->save()) {
                    $this->changeUserAmount(
                        $params['enterprise_id'],
                        $info->worker_id,
                        $info->application_amount,
                        false
                    );

                    $successIds[] = $info->id;
                } else {
                    $errorIds[] = $info->id;
                }
            }
            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
        }
        return [
            'successIds' => $successIds,
            'errorIds' => $errorIds
        ];
    }

    /**
     * 审核通过，进入到提现中
     * @param $withdrawalLogs
     * @param $params
     * @return array
     */
    public function auditPass($withdrawalLogs, $params)
    {
        //执行成功的ID数组
        $successIds = [];
        //执行失败的id数组
        $errorIds = [];
        DB::beginTransaction();
        try {
            foreach ($withdrawalLogs as $info) {
                //$info->enterprise_id = $params['enterprise_id'];
                $info->operator_id = $params['admin_id']??0;
                $info->status = WithdrawalApplication::WITHDRAW_IN_STATUS;
                $info->request_handle_at = time();
                if ($info->save()) {
                    /*$this->changeUserAmount(
                        $params['enterprise_id'],
                        $info->worker_id,
                        $info->application_amount,
                        true
                    );*/

                    $successIds[] = $info->id;
                } else {
                    $errorIds[] = $info->id;
                }
            }
            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
        }
        return [
            'successIds' => $successIds,
            'errorIds' => $errorIds
        ];
    }

    /**
     * 改变用户的可提现金额
     * @param $enterpriseId
     * @param $workerId 用户ID
     * @param $amount 金额
     * @param bool $successed 成功或者回退金额
     * @return bool
     */
    public function changeUserAmount($enterpriseId, $workerId, $amount, $successd = false)
    {
        $userCash = ShoppingGuideCashWithdrawal::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('worker_id', $workerId)
            ->where('frozen_amount', '>=', $amount)
            ->first();
        //如果用户账号不存在，返回false
        if (empty($userCash)) {
            return false;
        }
        //如果是成功，判断冻结金额是否能满足回退条件
        if ($successd) {
            $userCash->frozen_amount -= $amount;
            $userCash->cash_withdrawal += $amount;
        } else {
            $userCash->frozen_amount -= $amount;
            $userCash->available_cash_amount += $amount;
        }

        return $userCash->save();
    }

    /**
     * 获取提现配置信息
     * @param $enterpriseId
     * @return array
     */
    public function getWithdrawalConfig($enterpriseId)
    {
        $withdrawalConfig = WithdrawalSetting::query()
            ->where('enterprise_id', $enterpriseId)
            ->first();
        return empty($withdrawalConfig) ? [] : $withdrawalConfig->toArray();
    }
}
