<?php


namespace App\Services;


use App\Models\SettlementRefundRecord;
use Illuminate\Support\Facades\Log;

class SettlementRefundRecordService extends BaseService
{


    /**
     * 添加结算单退款记录
     * @param string $projectCode
     * @param string $orderRandomCode
     * @param string $settlementCode
     * @param bool $isInitiated 结算单退款请求是否发送成功，成功会设置退款发起时间
     */
    public function addRefundRecord(string $projectCode, string $orderRandomCode, string $settlementCode,
                                    bool $isInitiated = true)
    {
        $record = [
            'order_random_code' => $orderRandomCode,
            'settlement_code' => $settlementCode,
            'project_code' => $projectCode
        ];
        if ($isInitiated) {
            $record['refund_status'] = SettlementRefundRecord::REFUND_STATUS_PROCESSING;
            $record['refund_initiated_at'] = time();
        } else {
            $record['refund_status'] = SettlementRefundRecord::REFUND_STATUS_FAILED;
            $record['refund_initiated_at'] = 0;
            $record['refund_failed_type'] = SettlementRefundRecord::REFUND_FAILED_TYPE_INITIATED_FAIL;
        }
        SettlementRefundRecord::query()->create($record);
    }

    /**
     * 更新结算退款记录为退款成功
     * @param array $refundData
     */
    public function updateRefundSucceed(array $refundData)
    {
        $toIntAmount = function($floatAmount): int {
            if (empty($floatAmount)) {
                return 0;
            }
            return (int) ($floatAmount * 100);
        };
        $refundRecordQueryWhere = [
            'order_random_code' => $refundData['order_random_code'],
            'settlement_code' => $refundData['settlement_code']
        ];
        $refundRecord = SettlementRefundRecord::query()->where($refundRecordQueryWhere)->first();

        if ($refundRecord != null) {
            SettlementRefundRecord::query()->where('id', $refundRecord['id'])->update([
                'refund_merchant_amount' => $toIntAmount($refundData['refund_merchant_amount'] ?? 0),
                'refund_service_amount' => $toIntAmount($refundData['refund_service_amount'] ?? 0),
                'change_code' => $refundData['change_code']  ?? '',
                'change_merchant_amount' => $toIntAmount($refundData['change_merchant_amount'] ?? 0),
                'change_service_amount' => $toIntAmount($refundData['change_service_amount'] ?? 0),
                'refund_status' => SettlementRefundRecord::REFUND_STATUS_SUCCEED,
                'refund_succeed_at' => time(),
            ]);
        } else {
            Log::error("结算单退款记录不存在", $refundRecordQueryWhere);
        }
    }

}