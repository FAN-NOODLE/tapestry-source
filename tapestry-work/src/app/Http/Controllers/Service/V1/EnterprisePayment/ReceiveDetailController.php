<?php

namespace App\Http\Controllers\Service\V1\EnterprisePayment;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterprisePayment;
use App\Models\EnterprisePaymentRecord;

class ReceiveDetailController extends BaseController
{

    /**
     * 领取红包详情
     */
    public function index()
    {
        $userId = request()->post('user_id', 0);
        $activityId = request()->post('activity_id', 0);
        $enterpriseId = request()->post('enterprise_id', 0);

        if (empty($userId) || empty($activityId) || empty($enterpriseId)) {
            throw new Error(120001, '参数错误');
        }

        $activity = EnterprisePayment::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('id', $activityId)
            ->first(['name', 'send_name', 'wishing', 'qr_code_url']);

        $record = EnterprisePaymentRecord::query()
            ->where('enterprise_payment_id', $activityId)
            ->where('status', EnterprisePaymentRecord::RECORD_STATUS_SENDED)
            ->where('user_id', $userId)
            ->first(['amount', 'receive_refund_time']);

        $return = [];
        $return['name'] = $activity->name ?? '';
        $return['qr_code_url'] = $activity->qr_code_url ?? '';
        $return['send_name'] = $activity->send_name ?? '';
        $return['wishing'] = $activity->wishing ?? '';
        $return['amount'] = $record->amount ?? 0;
        $return['receive_refund_time'] = $record->receive_refund_time ?? 0;

        return success(['data' => $return]);
    }
}
