<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class OrderTmp extends Base
{
    const STATUS_ALIAS_ARRAY = [
        20  => '待支付',
        21  => '已取消',
        23  => '获取支付凭证',
        25  => '等待微信支付成功通知',
        30  => '待确认',
        40  => '待发货',
        41  => '售后中(待发货)',
        60  => '待收货',
        61  => '售后中(待收货)',
        80  => '已签收',
        81  => '售后中(退货退款)',
        82  => '售后中(仅退款)',
        83  => '部分退款',
        84  => '售后中(部分退款)',
        100 => '已完成',
        200 => '退款中',
        201 => '退款失败'
    ];

    // 已取消
    const STATUS_CANCEL_ARRAY = [21];

    // 待支付
    const STATUS_WAIT_ARRAY = [20, 23];

    // 已支付
    const STATUS_PAY_ARRAY = [25, 30, 40, 41, 60, 61, 81, 82, 83, 84, 200, 201];

    // 退款/售后
    const STATUS_FINISH_ARRAY = [80, 100];

    // 已成交
    const STATUS_TRADED = [25, 30, 40, 41, 60, 61, 81, 82, 83, 84, 200, 201, 80, 100];

    // 其他
    const STATUS_ALIAS_OTHER = '其他';

    // 其他
    const STATUS_ALIAS_CANCEL= '已取消';

    // 其他
    const STATUS_ALIAS_WAIT = '待付款';

    // 其他
    const STATUS_ALIAS_PAY = '已付款';

    // 其他
    const STATUS_ALIAS_FINISH= '退款/售后';

    // 待付款
    const STATUS_WAIT = 1;

    // 已付款
    const STATUS_PAY = 2;

    // 退款/售后
    const STATUS_FINISH = 3;

    // 已读
    const READ_STATUS_YES = 2;

    // 未读
    const READ_STATUS_NO = 1;

    protected $table = 'order_tmp';

    protected $primaryKey = 'id';

    protected $fillable = [
        'id',
        'user_id',
        'union_id',
        'open_id',
        'enterprise_id',
        'order_no',
        'status',
        'order_type',
        'total_amount',
        'preferential_price',
        'final_amount',
        'payment_channel',
        'pay_at',
        'refund_at',
        'order_at',
        'delivery_name',
        'delivery_mobile',
        'delivery_address',
        'remarks',
        'province',
        'city',
        'area',
        'logistic_identity',
        'logistic_no',
        'send_at',
        'sign_at',
        'cash_fee',
        'coupon_fee',
        'coupon_count',
        'other_fee',
        'read_status',
        'created_at',
        'updated_at',
    ];

    public function goods()
    {
        return $this->hasMany(OrderGoods::class, 'order_no', 'order_no');
    }

    public function user()
    {
        return $this->belongsTo(User::class, 'user_id', 'id');
    }

    /**
     * @var array
     */
    public static $selectField = ['id', 'user_id', 'union_id', 'open_id', 'enterprise_id', 'order_no', 'status', 'order_type', 'total_amount', 'preferential_price',
        'final_amount', 'payment_channel', 'pay_at', 'refund_at', 'order_at', 'delivery_name', 'delivery_mobile', 'delivery_address',
        'remarks', 'province', 'city', 'area', 'logistic_identity', 'logistic_no', 'send_at', 'sign_at', 'cash_fee', 'coupon_fee',
        'coupon_count', 'other_fee'];

    /**
     * 订单状态格式化
     *
     * @param $status
     *
     * @return string
     */
    public static function getStatusAlias($status)
    {
        $alias = self::STATUS_ALIAS_OTHER;

        if(in_array($status, self::STATUS_CANCEL_ARRAY)) {
            $alias = self::STATUS_ALIAS_CANCEL;
        }
        if(in_array($status, self::STATUS_WAIT_ARRAY)) {
            $alias = self::STATUS_ALIAS_WAIT;
        }
        if(in_array($status, self::STATUS_PAY_ARRAY)) {
            $alias = self::STATUS_ALIAS_PAY;
        }
        if(in_array($status, self::STATUS_FINISH_ARRAY)) {
            $alias = self::STATUS_ALIAS_FINISH;
        }
        return $alias;
    }
}
