<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Relations\HasOne;

class ActivityRecord extends Base
{
    // 中奖状态，0未中奖,1已中奖
    const StatusYes = 1;
    const StatusNo  = 0;
    // 是否核销，0未核销,1已核销
    const IsCheckYes = 1;
    const IsCheckNo  = 0;


    public $primaryKey = 'id';


    protected $table    = 'activity_record';
    protected $fillable = [
        'id',
        'enterprise_id',
        'activity_id',
        'unionid',
        'openid',
        'user_name',
        'qr_code',
        'activity_code',
        'prize_id',
        'status',
        'is_check',
        'check_time',
        'member_id',
        'member_name',
        'version'
    ];

    public function activity(): HasOne
    {
        return $this->hasOne(Activity::class, 'id', 'activity_id');
    }

    public function prize(): HasOne
    {
        return $this->hasOne(ActivityPrize::class, 'id', 'prize_id');
    }
}