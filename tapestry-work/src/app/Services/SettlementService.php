<?php


namespace App\Services;

use App\Exceptions\Error;
use App\Library\Balance\Wetax\BalanceSettlement;
use App\Library\Balance\Wetax\SettlementException;
use App\Library\Balance\Wetax\VerificationErrorCode;
use App\Models\SettlementRecord;
use App\Models\StoreWorkerCertification;
use App\Models\WithdrawalApplication;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;
use Illuminate\Validation\Rule;

/**
 * 资金结算服务
 * Class SettleService
 * @package App\Services
 */
class SettlementService extends BaseService
{
    /**
     * 结算服务
     * @var BalanceSettlement
     */
    private $balanceSettlement;

    public function __construct()
    {
        $this->balanceSettlement = new BalanceSettlement();
    }

    /**
     * 授权签约，返回结果为true表示签约成功，失败会抛出异常
     *
     * @param string $projectCode 项目编码
     * @param string $name 用户姓名
     * @param int $certificateType 证件类型(默认为 1)
     * @param string $certificateNum 证件号
     * @param string $phoneNum 手机号
     * @return bool
     */
    public function signAgreement(string $name, int $certificateType, string $certificateNum, string $phoneNum,
                                  string $projectCode): bool
    {
        Log::info("SettlementService.signAgreement", func_get_args());
        $validator = Validator::make([
            'project_code' => $projectCode,
            'name' => $name,
            'certificate_type' => $certificateType,
            'certificate_num' => $certificateNum,
            'phone_num' => $phoneNum,
        ], [
            'project_code' => 'required|string',
            'name' => 'required|string',
            'certificate_type' => 'required|integer',
            'certificate_num' => 'required|string',
            'phone_num' => 'required|string'
        ]);
        if ($validator->fails()) {
            Log::info("SettlementService.signAgreement params validate failed", ['errors' => $validator->errors()->all()]);
            throw new Error(1100011, "授权签约参数验证失败");
        }

        try {
            return $this->balanceSettlement->signAgreement($projectCode, $name, $certificateType, $certificateNum, $phoneNum);
        } catch (SettlementException $e) {
            throw new Error(1100011, $e->getFailReason(), 500, $e);
        }
    }

    /**
     * 获取签约结果
     * 响应结果：
     * [
     *     'is_signed' => true, // 签约结果，true 已签约，false 未签约
     *     'agreement_pic' => '' // 签约生成 PDF 格式文件链接(10 分钟有效期)
     * ]
     *
     * @param string $projectCode 项目编码
     * @param string $name 用户姓名
     * @param string $certificateNum 身份证号
     * @return array
     */
    public function queryAgreement(string $name, string $certificateNum, string $projectCode): array
    {
        Log::info("SettlementService.queryAgreement", func_get_args());
        $validator = Validator::make([
            'project_code' => $projectCode,
            'name' => $name,
            'certificate_num' => $certificateNum,
        ], [
            'project_code' => 'required|string',
            'name' => 'required|string',
            'certificate_num' => 'required|string',
        ]);
        if ($validator->fails()) {
            Log::info("SettlementService.queryAgreement params validate failed", ['errors' => $validator->errors()->all()]);
            throw new Error(1100012, "获取签约结果参数验证失败");
        }

        try {
            return $this->balanceSettlement->queryAgreement($projectCode, $name, $certificateNum);
        } catch (SettlementException $e) {
            throw new Error(1100012, $e->getFailReason(), 500, $e);
        }
    }

    /**
     * 获取项目列表
     * 响应结果：
     * [
     *     {
     *         'project_code': "XM19BK0A8658AEA0", // 项目编号
     *         'project_name': "Test 项目", // 项目名称
     *         'create_time': "2019-11-20 20:46:06", // 创建时间
     *         'allow_merchant_balance': false // 是否允许使用商户余额
     *     }
     * ]
     *
     * @param string $projectName
     * @return array
     */
    public function getProjectList(string $projectName)
    {
        Log::info("SettlementService.getProjectList, projectName: " . $projectName);
        $projectName = trim($projectName);
        if (empty($projectName)) {
            throw new Error(1100013, '项目名称不能为空');
        }
        try {
            return $this->balanceSettlement->getProjectList($projectName);
        } catch (SettlementException $e) {
            throw new Error(1100013, $e->getFailReason(), 500, $e);
        }
    }

    /**
     * 根据项目编码获取项目信息
     * 响应结果：
     * [
     *     'project_code': "XM19BK0A8658AEA0", // 项目编号
     *     'project_name': "Test 项目", // 项目名称
     *     'create_time': "2019-11-20 20:46:06", // 创建时间
     *     'allow_merchant_balance': false // 是否允许使用商户余额
     * ]
     *
     * @param string $projectCode
     * @return array
     */
    public function getProjectByCode(string $projectCode)
    {
        Log::info("SettlementService.getProjectList, projectCode: " . $projectCode);
        $projectCode = trim($projectCode);
        if (empty($projectCode)) {
            throw new Error(1100014, '项目编码不能为空');
        }

        try {
            return $this->balanceSettlement->getProjectByCode($projectCode);
        } catch (SettlementException $e) {
            throw new Error(1100014, $e->getFailReason(), 500, $e);
        }
    }

    /**
     * 发起结算单，结算单状态变更会异步通知
     * @param string $enterpriseHash
     * @param string $projectCode
     * @param string $orderRandomCode
     * @param string $name
     * @param int $certificateType
     * @param string $certificateNum
     * @param string $phoneNum
     * @param string $bankName
     * @param string $bankCardNum
     * @param int $settleAmount
     */
    public function settleBalance(string $enterpriseHash,
                                  string $projectCode,
                                  string $orderRandomCode,
                                  string $name,
                                  int $certificateType,
                                  string $certificateNum,
                                  string $phoneNum,
                                  string $bankName,
                                  string $bankCardNum,
                                  int $settleAmount)
    {
        Log::info("SettlementService.settleBalance", func_get_args());

        $settleData = [
            'order_random_code' => $orderRandomCode,
            'project_code' => $projectCode,
            'name' => $name,
            'certificate_type' => $certificateType,
            'certificate_num' => $certificateNum,
            'phone_num' => $phoneNum,
            'bank_name' => $bankName,
            'bank_card_num' => $bankCardNum,
            'settle_amount' => $settleAmount
        ];
        $settleDataList = [$settleData];
        $this->settleValidator($settleDataList);

        $paymentWay = BalanceSettlement::PAYMENT_WAY_BANK_CARD;
        $floatSettleAmount = (float) ($settleAmount / 100);
        try {
            $this->balanceSettlement->settleBalance($projectCode, $orderRandomCode, $name, $certificateType,
                $certificateNum, $phoneNum, $floatSettleAmount, $paymentWay, $bankName, $bankCardNum);
        } catch (\Exception $e) {
            Log::error("SettlementService.settleBalance fail", [
                'code' => $e->getCode(),
                'message' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
            throw new Error(1100014, "结算发起失败", 500, $e);
        }

        SettlementRecord::query()->create($settleData);
    }

    /**
     * 批量结算
     * @param string $enterpriseHash 企业hash
     * @param string $projectCode
     * @param array $settleDataList
     */
    public function batchSettleBalance(string $enterpriseHash, string $projectCode, array $settleDataList)
    {
        Log::info("SettlementService.batchSettleBalance", func_get_args());
        $this->settleValidator($settleDataList);

        $realSettleInfoList = [];
        foreach ($settleDataList as $value) {
            $settleInfo['project_code'] = $value['project_code'];
            $settleInfo['order_random_code'] = $value['order_random_code'];
            $settleInfo['name'] = $value['name'];
            $settleInfo['certificate_type'] = $value['certificate_type'];
            $settleInfo['certificate_num'] = $value['certificate_num'];
            $settleInfo['phone_num'] = $value['phone_num'];
            $settleInfo['bank_name'] = $value['bank_name'];
            $settleInfo['bankcard_num'] = $value['bank_card_num'];
            $settleInfo['payment_way'] = BalanceSettlement::PAYMENT_WAY_BANK_CARD;
            $settleInfo['settle_amount'] = (float) ($value['settle_amount'] / 100);
            $realSettleInfoList[] = $settleInfo;
        }

        try {
            $this->balanceSettlement->batchSettleBalance($enterpriseHash, $projectCode, $realSettleInfoList);
        } catch (\Exception $e) {
            Log::error("SettlementService.batchSettleBalance fail", [
                'code' => $e->getCode(),
                'message' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
            throw new Error(1100014, "批量结算发起失败", 500, $e);
        }

        $time = time();
        foreach ($settleDataList as &$value) {
            $value['created_at'] = $time;
        }
        DB::table((new SettlementRecord())->getTable())->insert($settleDataList);
    }

    /**
     * 结算数据验证
     * @param array $settleDataList
     */
    private function settleValidator(array $settleDataList)
    {
        $data = [
            'list' => $settleDataList
        ];
        $rules = [
            'list' => 'required|array',
            'list.*.project_code' => 'required|string',
            'list.*.order_random_code' => 'required|string',
            'list.*.name' => 'required|string',
            'list.*.certificate_type' => ['required', Rule::in([
                StoreWorkerCertification::ID_TYPE_MAINLAND_RESIDENT,
                StoreWorkerCertification::ID_TYPE_HK_MA_PASS,
                StoreWorkerCertification::ID_TYPE_HK_MA_RESIDENT,
                StoreWorkerCertification::ID_TYPE_TW_PASS,
                StoreWorkerCertification::ID_TYPE_TW_RESIDENT,
            ])],
            'list.*.certificate_num' => 'required|string',
            'list.*.phone_num' => 'required|string',
            'list.*.bank_name' => 'required|string',
            'list.*.bank_card_num' => 'required|string',
            'list.*.settle_amount' => 'required|integer|min:100',
        ];
        $validator = Validator::make($data, $rules);
        if ($validator->fails()) {
            Log::info("SettlementService.settleBalance params validate failed",
                ['errors' => $validator->errors()->all()]);
            throw new Error(1100014, "结算单数据验证失败");
        }
    }

    /**
     * 结算单发起退款，并添加数据库退款记录
     * @param string $enterpriseHash
     * @param string $projectCode
     * @param string $orderRandomCode
     * @param string $settlementCode
     */
    public function refundBalance(string $enterpriseHash, string $projectCode, string $orderRandomCode, string $settlementCode)
    {
        Log::info("SettlementService.refundBalance", func_get_args());
        $validator = Validator::make([
            'project_code' => $projectCode,
            'order_random_code' => $orderRandomCode,
            'settlement_code' => $settlementCode,
        ], [
            'project_code' => 'required|string',
            'order_random_code' => 'required|string',
            'settlement_code' => 'required|string',
        ]);
        if ($validator->fails()) {
            Log::info("SettlementService.refundBalance params validate failed",
                ['errors' => $validator->errors()->all()]);
            throw new Error(1100020, "结算单退款参数验证失败");
        }

        $refundRecordService = new SettlementRefundRecordService();
        try {
            $this->balanceSettlement->refundBalance($enterpriseHash, $projectCode, $orderRandomCode, $settlementCode);
            $refundRecordService->addRefundRecord($projectCode, $orderRandomCode, $settlementCode);
        } catch (\Exception $e) {
            Log::error("settlement refund initiate failed", [
                'code' => $e->getCode(),
                'message' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
            $refundRecordService->addRefundRecord($projectCode, $orderRandomCode, $settlementCode, false);
        }
    }

    /**
     * 结算单信息查询
     * @param string $projectCode
     * @param string $orderRandomCode
     * @param string $settlementCode
     * @return array
     */
    public function getSettleBalance(string $projectCode, string $orderRandomCode, string $settlementCode = '')
    {
        Log::info("SettlementService.getBalance", func_get_args());

        $validator = Validator::make([
            'project_code' => $projectCode,
            'order_random_code' => $orderRandomCode,
            'settlement_code' => $settlementCode,
        ], [
            'project_code' => 'required|string',
            'order_random_code' => 'required|string',
            'settlement_code' => 'string',
        ]);
        if ($validator->fails()) {
            Log::info("SettlementService.getSettleBalance params validate failed",
                ['errors' => $validator->errors()->all()]);
            throw new Error(1100019, "结算单查询参数验证失败");
        }
        try {
            return $this->balanceSettlement->getBalance($projectCode, $orderRandomCode, $settlementCode);
        } catch (\Exception $e) {
            Log::info("SettlementService.getBalance", [
                'code' => $e->getCode(),
                'message' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
            throw new Error(1100019, "结算单查询失败", 500, $e);
        }
    }

    /**
     * 同步结算信息
     * @param string $projectCode
     * @param string $orderRandomCode
     * @param string $settlementCode
     */
    public function syncSettlementInfo(string $projectCode, string $orderRandomCode, string $settlementCode = '')
    {
        Log::info("SettlementService.syncSettlementInfo", func_get_args());

        $settlementInfo = $this->getSettleBalance($projectCode, $orderRandomCode, $settlementCode);

        return $settlementInfo;
    }

    /**
     * 处理结算单状态信息
     * @param array $info
     */
    public function processSettlementInfo(array $info)
    {
        if (empty($info['order_random_code']) || empty($info['settlement_code']) || empty($info['project_code'])) {
            Log::error("settleBalanceCallback callbackData main params not exists", ['info' => $info]);
            return;
        }
        $isRefund = false;
        $enterpriseHash = '';
        try {
            // 更新结算单信息
            Log::info("settleBalanceCallback settlement record update");
            $settlementRecord = SettlementRecord::query()->where('order_random_code', $info['order_random_code'])->first();
            if ($settlementRecord != null) {
                $setData = [];
                isset($info['settlement_code']) && $setData['settlement_code'] = (string) $info['settlement_code'];
                isset($info['account_num']) && $setData['account_num'] = (string) $info['account_num'];
                isset($info['status']) && $setData['status'] = (int) $info['status'];
                isset($info['is_verification']) && $setData['is_verification'] = (int) $info['is_verification'];
                isset($info['verification_error_code']) && $setData['verification_error_code'] = (string) $info['verification_error_code'];
                isset($info['fail_reason']) && $setData['fail_reason'] = (string) $info['fail_reason'];
                isset($info['hangup_flag']) && $setData['hangup_flag'] = (int) $info['hangup_flag'];
                isset($info['hangup_begintime']) && $setData['hangup_begintime'] = (string) $info['hangup_begintime'];
                isset($info['hangup_msg']) && $setData['hangup_msg'] = (string) $info['hangup_msg'];
                isset($info['pay_time']) && $setData['pay_time'] = (string) $info['pay_time'];
                isset($info['handle_time']) && $setData['handle_time'] = (string) $info['handle_time'];
                if (!empty($setData)) {
                    SettlementRecord::query()->where('id', $settlementRecord['id'])->update($setData);
                }
                Log::info("settleBalanceCallback settlement record update succeed");
            } else {
                Log::error("settleBalanceCallback settlement record not exists",
                    ['order_random_code' => $info['order_random_code']]);
            }


            // 处理结算单状态，更新提现状态
            $withdrawalSetData = [
                'status' => null,
                'reason' => ""
            ];
            $withdrawalIsSucceed = null;
            if (!empty($info['status'])) {
                if (100 == $info['status']) { // 结算单未审核
                    if (false === ($info['is_verification'] ?? null)) { // 结算单创建失败（审核不通过）
                        // 更新提现失败，和失败原因
                        $withdrawalSetData['status'] = WithdrawalApplication::WITHDRAW_FAIL_STATUS;
                        $withdrawalSetData['reason'] = VerificationErrorCode::getErrorByCodeList(
                            json_decode($info["verification_error_code"] ?? "", true));
                        $withdrawalIsSucceed  = false;
                    }
                } else if (1000 == $info['status']) { // 打款成功
                    // 更新提现成功
                    $withdrawalSetData['status'] = WithdrawalApplication::WITHDRAW_SUCCESS_STATUS;
                    $withdrawalIsSucceed = true;
                } else if (1004 == $info['status']) { // 打款失败
                    // 更新提现失败和失败原因
                    $withdrawalSetData['status'] = WithdrawalApplication::WITHDRAW_FAIL_STATUS;
                    $withdrawalSetData['reason'] = !empty($info['fail_reason']) ? $info['fail_reason'] : '打款失败';
                    $withdrawalIsSucceed = false;
                    $isRefund = true;
                }
            }

            Log::info("settleBalanceCallback withdrawal status update", $withdrawalSetData);
            $withdrawalId = (int)$info['order_random_code'];
            $withdrawal = WithdrawalApplication::query()->where('id', $withdrawalId)->first();
            if ($withdrawal == null) {
                Log::error("settleBalanceCallback withdrawal not exists", ['withdrawalId' => $withdrawalId]);
            }

            if ($withdrawalSetData['status'] != null && $withdrawal != null) {
                $enterpriseData = (new EnterpriseService())->getDetailByEnterpriseId($withdrawal['enterprise_id']);
                if (!empty($enterpriseData)) {
                    $enterpriseHash = $enterpriseData['hash'];
                }
                DB::beginTransaction();
                try {
                    WithdrawalApplication::query()->where('id', $withdrawalId)->update($withdrawalSetData);

                    // 更新导购可提现金额表
                    if (null !== $withdrawalIsSucceed) {
                        (new WithdrawalServices())->changeUserAmount($withdrawal['enterprise_id'],
                            $withdrawal['worker_id'], $withdrawal['application_amount'], $withdrawalIsSucceed);
                    }

                    Log::info("settleBalanceCallback withdrawal status update succeed");
                    DB::commit();
                } catch (\Exception $e) {
                    DB::rollBack();
                    throw $e;
                }
            }
            Log::info("settleBalanceCallback process succeed");
        } catch (\Exception $e) {
            Log::error("settleBalanceCallback process failed", [
                'code' => $e->getCode(),
                'message' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
        }

        // 发起结算单退款
        if ($isRefund) {
            if (!empty($enterpriseHash)) {
                $this->refundBalance($enterpriseHash, $info['project_code'], $info['order_random_code'], $info['settlement_code']);
            } else {
                Log::error("企业hash不存在", ['enterpriseHash' => $enterpriseHash]);
            }
        }
    }

}