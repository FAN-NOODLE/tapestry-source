<?php

namespace App\Models\ScrmEnterprise;

use App\Models\Scrm\Members;

class ExternalPayRefundOrder extends ScrmEnterpriseInit
{
    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'external_pay_refund_order';
    /**
     * @var array
     */
    protected $guarded = [];

    const REFUND_STATUS = [
        0 => "已申请退款",
        1 => "退款处理中",
        2 => "退款成功",
        3 => "退款关闭",
        4 => "退款异常",
        5 => "审批中",
        6 => "审批失败",
        7 => "审批取消"
    ];

    /**
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     */
    public function order()
    {
        return $this->belongsTo(ExternalPayOrder::class, 'out_trade_no', 'out_trade_no');
    }

    /**
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     */
    public function member()
    {
        return $this->belongsTo(Members::class,'refund_member_id','id');
    }
}