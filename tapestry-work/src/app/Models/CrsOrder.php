<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/9/10
 * Time: 11:47
 */

namespace App\Models;

class CrsOrder extends Base
{
    const STATUS_ALIAS_ARRAY = [
        0  => "全部",
        20 => "待支付",
        21 => "已取消",
        23 => "待付款",
        25 => "待确认",
        30 => "待确认",
        40 => "待发货",
        41 => "售后中",
        60 => "待收货",
        61 => "售后中",
        80 => "已签收",
        81 => "售后中",
        82 => "售后中",
        83 => "部分退款",
        84 => "售后中",
        100 => "已完成",
        200 => "退款中",
        201 => "退款失败",
        205 => "全部退款"
    ];

    const STATUS_COMPLETE = 205;
    const STATUS_PRE_DELIVER = 40;
    const PART_REFUND = 83;
    const SIGNED_IN = 80;


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
    const STATUS_ALIAS_PAY = '待发货';

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

    const HAO_WU_JI_ORDER = 2;

    protected $table = 'crs_order';

    protected $primaryKey = 'id';

    protected $fillable = [
        'id',
        'sync_source',
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
        'logistics_price',
        'third_payment_no',
        'store_staff_id',
        'created_at',
        'updated_at',
    ];


    public function goods()
    {
        return $this->hasMany(CrsOrderGoods::class, 'order_no', 'order_no');
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
        $alias = self::STATUS_ALIAS_ARRAY[$status] ?? self::STATUS_ALIAS_OTHER;
        return $alias;
    }

    /**
     * getOrder
     * @desc
     * @param $sourceType
     * @param $orderNo
     *
     * @return mixed
     * @date 2020-09-08
     */
    public static function getOrder($orderNo)
    {
        return self::query()
            ->where('order_no', $orderNo)
            ->first();
    }
}
