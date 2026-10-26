<?php


namespace App\Models;


use Illuminate\Support\Facades\DB;

class SettlementRecord extends Base
{
    protected $table = 'settlement_record';

    protected $fillable = [
        'order_random_code',
        'settlement_code',
        'account_num',
        'project_code',
        'name',
        'certificate_type',
        'certificate_num',
        'phone_num',
        'bank_name',
        'bank_card_num',
        'settle_amount',
        'status',
        'is_verification',
        'verification_error_code',
        'fail_reason',
        'hangup_flag',
        'hangup_begintime',
        'hangup_msg',
        'pay_time',
        'handle_time'
    ];

    public function addAll(Array $data)
    {
        $rs = DB::table($this->getTable())->insert($data);
        return $rs;
    }
}