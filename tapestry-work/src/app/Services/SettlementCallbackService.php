<?php


namespace App\Services;

use App\Exceptions\Error;
use App\Library\Balance\Wetax\EncryptUtil;
use App\Library\Balance\Wetax\ResponseReader;
use App\Library\Balance\Wetax\SettlementConfig;
use Illuminate\Support\Facades\Log;

/**
 * 资金结算回调
 * Class SettlementCallbackService
 * @package App\Services
 */
class SettlementCallbackService
{

    /**
     * 商户结算平台secret
     * @var string
     */
    private $config;

    public function __construct()
    {
        $this->config = SettlementConfig::getInstance();
    }

    /**
     * 解密回调数据并获取data参数
     * @param string $encryptString 加密的数据串
     * @return array
     */
    private function getCallbackData(string $encryptString): array
    {
        $decryptedData = null;
        try {
            $decryptedData = EncryptUtil::decrypt($encryptString, $this->config->getAppSecret());
            Log::info("SettlementCallback decryptedData:" . $decryptedData);
            return json_decode($decryptedData, true);
        } catch (\Exception $e) {
            Log::error("settlement callback data decrypt error", ['exception' => $e->getTraceAsString()]);
            throw new Error(1100015, "回调数据解析失败", 500, $e);
        }
    }

    /**
     * 结算单状态变更回调
     *
     * data参数：
     * [
     *     {
     *         "account_num": "ZD20190729ca7792", // 账单号
     *         "order_random_code": "100001_1564390605", // 客户订单号
     *         "settlement_code": "JS2019072932389f", // 结算单号
     *         "project_code": null, // 项目编码
     *         "is_verification": false, // 结算单校验是否成功
     *         "verification_error_code": "[11,13]", // 错误码
     *         "status": 100, //结算单状态
     *         "service_amount": 0.09, // 服务费
     *         "settle_amount": 1.0, // 绩效费
     *         "total_money": 1.09, // 合计金额
     *         "fail_reason": null, // 失败原因
     *         "hangup_flag": false, // 挂单标识: true:挂单
     *         "hangup_msg": null, // 挂单原因
     *         "hangup_begintime": null, // 挂单时间
     *         "pay_time": null, // 支付时间(扣除商户余额的时间)
     *         "handle_time": null // 处理时间(打款成功或失败的时间)
     *         "surplus_money": 0.0, // 本月剩余免增值税金额(挂单时有值，否则 为 0)
     *     }
     * ]
     *
     * 挂单：挂单指的结算单被暂时挂起，未进入打款。多商户给用户的结算总额超过用户当月结算阶梯费率标准 99399 元时，结算单将挂起
     * 后进入排队序列，次月自动处理;如不希望下月自动处理，可在商户结算后台【财务管理】找到 对应结算单-进行【退款】处理;
     *
     * @param string $callbackData
     */
    public function settleBalanceCallback(string $callbackData)
    {
        Log::info("settleBalanceCallback callbackData: " . $callbackData);
        $data = $this->getCallbackData($callbackData);
        if (!is_array($data)) {
            Log::error("settleBalanceCallback data type error");
            throw new Error(1100016, "结算单状态变更回调数据类型错误");
        }

        // 处理结算回调数据
        $settlementService = new SettlementService();
        foreach ($data as $settlementInfo) {
            $settlementService->processSettlementInfo($settlementInfo);
        }
    }

    /**
     * 结算单退款回调
     *
     * data参数：
     * {
     *     "code":0,"
     *     "msg":"处理成功",
     *     "request_id":"REDACTED_FOR_SOURCE_SHARING",
     *     "appkey":"REDACTED_FOR_SOURCE_SHARING",
     *     "data": {
     *         "settlement_code":"JS19BR19A690E9F9", // 结算流水号
     *         "order_random_code":"09708757-7ea1-4fda", // 客户订单号
     *         "refund_merchant_amount":54736.84, // 余额退回金额
     *         "refund_service_amount":263.16,  // 服务费抵扣券退回金额
     *         "change_code": "FW19BRAA9A200255", // 服务费变更单单流水号
     *         "change_merchant_amount":4263.16, // 变更单余额退回金额
     *         "change_service_amount":236.84  // 变更单抵扣券退回金额
     *      }
     *   }
     * @param string $callbackData
     */
    public function refundBalanceCallback(string $callbackData)
    {
        Log::info("refundBalanceCallback callbackData: " . $callbackData);
        $data = $this->getCallbackData($callbackData);
        if (empty($data['order_random_code']) || empty($data['settlement_code'])) {
            throw new Error(1100017, "客户订单号或结算流水号为空");
        }

        // 更新结算单退款记录为退款成功
        (new SettlementRefundRecordService())->updateRefundSucceed($data);
    }


}