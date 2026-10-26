<?php


namespace App\Http\Controllers\App\V1\Settlement;


use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\SettlementCallbackService;
use App\Services\SettlementService;

class TestController extends BaseController
{
    /**
     * @var SettlementService
     */
    private $settlementService;

    private $settlementCallbackService;

    public function __construct(SettlementService $settlementService, SettlementCallbackService $settlementCallbackService)
    {
        $this->settlementService = $settlementService;
        $this->settlementCallbackService = $settlementCallbackService;
        $testHeader = request()->header('Settlement-Test', '');
        if (empty($testHeader)) {
            throw new Error(110007, "非法请求", 401);
        }
    }

    public function getProjectByCode()
    {
        $projectCode = request()->input('project_code', '');
        $project = $this->settlementService->getProjectByCode($projectCode);
        return success([
            'project' => $project
        ]);
    }

    public function getProjectList()
    {
        $projectName = request()->input('project_name', '');
        $projectList = $this->settlementService->getProjectList($projectName);
        return success([
            'projectList' => $projectList
        ]);
    }

    public function signAgreement()
    {
        $projectCode = request()->input('project_code', '');
        $name = request()->input('name', '');
        $certificateType = request()->input('certificate_type', 0);;
        $certificateNum = request()->input('certificate_num', '');
        $phoneNum = request()->input('phone_num', '');
        $result = $this->settlementService->signAgreement($name, $certificateType, $certificateNum, $phoneNum, $projectCode);
        return success([
            'result' => $result
        ]);
    }

    public function queryAgreement()
    {
        $projectCode = request()->input('project_code', '');
        $name = request()->input('name', '');
        $certificateNum = request()->input('certificate_num', '');
        $result = $this->settlementService->queryAgreement($name, $certificateNum, $projectCode);
        return success([
            'result' => $result
        ]);
    }

    public function settleBalance()
    {
        $enterpriseHash = request()->input('enterprise_hash', '');
        $projectCode = request()->input('project_code', '');
        $orderRandomCode = request()->input('order_random_code', '');
        $name = request()->input('name', '');
        $certificateType = request()->input('certificate_type', '');
        $certificateNum = request()->input('certificate_num', '');
        $phoneNum = request()->input('phone_num', '');
        $bankName = request()->input('bank_name', '');
        $bankCardNum = request()->input('bank_card_num', '');
        $settleAmount = request()->input('settle_amount', '');

        $this->settlementService->settleBalance($enterpriseHash, $projectCode, $orderRandomCode, $name, $certificateType,
            $certificateNum, $phoneNum, $bankName, $bankCardNum, $settleAmount);
        return success([]);
    }

    public function getSettleBalance()
    {
        $projectCode = request()->input('project_code', '');
        $orderRandomCode = request()->input('order_random_code', '');
        $settlementCode = request()->input('settlement_code', '');
        $settlementInfo = $this->settlementService->getSettleBalance($projectCode, $orderRandomCode, $settlementCode);
        return success([
            'settlement_info' => $settlementInfo
        ]);
    }

    public function syncSettlementInfo()
    {
        $projectCode = request()->input('project_code', '');
        $orderRandomCode = request()->input('order_random_code', '');
        $settlementCode = request()->input('settlement_code', '');
        $settlementInfo = $this->settlementService->syncSettlementInfo($projectCode, $orderRandomCode, $settlementCode);
        return success([
            'settlement_info' => $settlementInfo
        ]);
    }

    public function settlementCallback()
    {
        $callbackData = request()->input('callback_data');
        $this->settlementCallbackService->settleBalanceCallback($callbackData);
        return success([]);
    }

    public function settlementRefundCallback()
    {
        $callbackData = request()->input('callback_data');
        $this->settlementCallbackService->refundBalanceCallback($callbackData);
        return success([]);
    }
}