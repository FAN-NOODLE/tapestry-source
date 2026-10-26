<?php

namespace App\Models\ScrmEnterprise;

use App\Models\Scrm\Members;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;

class ExternalPayOrder extends ScrmEnterpriseInit
{

    /**
     * 发货订单
     */
    const DELIVERY_TYPE_ON = 2;
    /**
     * 自提订单
     */
    const DELIVERY_TYPE_SELF = 1;
    /**
     * 已完成订单
     */
    const ORDER_STATE_SUCCESS = 3;
    /**
     * 已发货订单
     */
    const ORDER_STATE_DELIVER = 2;
    /**
     * 订单状态
     */
    const ORDER_STATE = [
        1 => "待发货",
        2 => "已发货",
        3 => "已完成",
        4 => "已退款"
    ];
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
    protected $table = 'external_pay_order';
    /**
     * @var array
     */
    protected $guarded = [];
    /**
     * @var string[]
     */
    protected $hidden = [
        'created_at',
        'updated_at',
        'deleted_at'
    ];

    /**
     * 关联用户信息
     *
     * @return BelongsTo
     */
    public function user()
    {
        return $this->belongsTo(WxUserCorpExternal::class, 'external_userid', 'openid');
    }

    /**
     * 关联退款信息
     *
     * @return HasMany
     */
    public function refund(): HasMany
    {
        return $this->hasMany(ExternalPayRefundOrder::class, 'out_refund_no', 'out_trade_no');
    }

    /**
     * 关联导购
     * @return BelongsTo
     */
    public function member()
    {
        return $this->belongsTo(Members::class, 'member_id', 'id');
    }

    /**
     * @param $key
     *
     * @return void
     */
    public function getContactInfoAttribute($key)
    {
        return $key ? json_decode($key, true) : [];
    }

    /**
     * @param $key
     *
     * @return void
     */
    public function getCommodityListAttribute($key)
    {
        return $key ? json_decode($key, true) : [];
    }

    /**
     * @param $key
     *
     * @return array|mixed
     */
    public function getLogisticsTrackAttribute($key)
    {
        return $key ? json_decode($key, true) : [];
    }
}