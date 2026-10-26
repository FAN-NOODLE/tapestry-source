<?php

namespace App\Console\Commands;

use App\Models\CashCouponRecord;
use CrsSdk\MicroPay\CashCoupon as CashCouponSdk;
use Illuminate\Console\Command;

class AutoQueryCashCouponInfo extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:query_cash_coupon_info';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '定时查询红包发送记录';

    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     */
    public function handle()
    {
        $cashCouponRecord = CashCouponRecord::query()
            ->where('status', CashCouponRecord::SENDING)
            ->get()
            ->toArray();

        if ($cashCouponRecord) {
            foreach ($cashCouponRecord as $value) {
                $params = [
                    'enterprise_id' => $value['enterprise_id'],
                    'bill_no' => $value['bill_no'],
                ];

                $result = (new CashCouponSdk())->getRedPackInfo($params);

                if (!empty($result['data'])) {
                    switch ($result['data']['status']) {
                        // 返回状态有5个，其中SENDING:发放中，SENT:已发放待领取，RFUND_ING:退款中不做处理，
                        // RECEIVED:已领取变更状态为已领取，FAILED：发放失败， REFUND:已退款变更状态为已退款
                        case 'SENDING':
                        case 'SENT':
                        case 'RFUND_ING':
                            break;
                        case 'RECEIVED':
                            $this->receiveAct($value, $result['data']);
                            break;
                        case 'FAILED':
                        case 'REFUND':
                            $this->refundAct($value, $result['data']);
                            break;
                    }
                }
            }
        }
    }

    /**
     * 已领取
     *
     * @param $cashCouponRecord
     * @param $result
     */
    private function receiveAct($cashCouponRecord, $result)
    {
        CashCouponRecord::query()
            ->where('status', CashCouponRecord::SENDING)
            ->where('id', $cashCouponRecord['id'])
            ->update([
                'status' => CashCouponRecord::RECEIVED,
                'send_time' => strtotime($result['send_time']),
                'receive_refund_time' => strtotime($result['receive_time']),
            ]);
    }

    /**
     * 已退款
     *
     * @param $cashCouponRecord
     * @param $result
     */
    private function refundAct($cashCouponRecord, $result)
    {
        CashCouponRecord::query()
            ->where('status', CashCouponRecord::SENDING)
            ->where('id', $cashCouponRecord['id'])
            ->update([
                'status' => CashCouponRecord::REFUND,
                'send_time' => strtotime($result['send_time']),
                'receive_refund_time' => strtotime($result['refund_time']),
            ]);
    }
}
