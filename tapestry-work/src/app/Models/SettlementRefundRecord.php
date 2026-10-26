<?php


namespace App\Models;

/**
 * 结算单退款记录
 * Class SettlementRefundRecord
 * @package App\Models
 */
class SettlementRefundRecord extends Base
{

    // 退款单状态-退款中
    const REFUND_STATUS_PROCESSING = 1;

    // 退款单状态-退款成功
    const REFUND_STATUS_SUCCEED = 2;

    // 退款单状态-退款失败
    const REFUND_STATUS_FAILED = 3;


    // 退款失败类型-发起退款失败
    const REFUND_FAILED_TYPE_INITIATED_FAIL = 1;

    // 退款失败类型-达到重试次数
    const REFUND_FAILED_TYPE_OVER_RETRIES = 2;

    protected $table = 'settlement_refund_record';

    protected $fillable = [
        'order_random_code',
        'settlement_code',
        'project_code',
        'refund_merchant_amount',
        'refund_service_amount',
        'change_code',
        'change_merchant_amount',
        'change_service_amount',
        'refund_status',
        'refund_initiated_at',
        'refund_succeed_at',
        'refund_failed_type',
        'retry_num',
    ];

}