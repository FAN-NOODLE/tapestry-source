<?php

namespace App\Models;

class UserCouponTmp extends Base
{
    const COUPON_LIMIT_NO = 1; // 无限制优惠券

    const COUPON_LIMIT_TINE = 2; // 固定时间段

    const COUPON_HAVE_EXPIRED = 0; // 已过期

    const COUPON_WAIT_USE = 1; // 待使用

    const COUPON_HAVE_USED = 3; // 已使用

    /**
     * 表名
     * @var string
     */
    protected $table = 'user_coupon_tmp';

    /**
     * 主键
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'shop_user_coupon_id',
        'enterprise_id',
        'user_id',
        'union_id',
        'open_id',
        'effective',
        'start_time',
        'end_time',
        'status',
        'wx_coupon_id',
        'receive_time',
        'use_time',
        'order_no',
        'order_amount'
    ];

    /**
     * @var array
     */
    public static $selectField = ['id', 'shop_user_coupon_id', 'enterprise_id', 'user_id', 'union_id', 'open_id', 'effective',
                                  'start_time', 'end_time', 'status', 'wx_coupon_id', 'receive_time', 'use_time', 'order_no', 'order_amount'];


    /**
     * 用户优惠券
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     */
    public function coupon()
    {
        return $this->belongsTo(WxCoupon::class, 'wx_coupon_id', 'shop_wx_coupon_id');
    }

    /**
     * 优惠券订单
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     */
    public function order()
    {
        return $this->belongsTo(Order::class, 'order_no', 'order_no');
    }


    /**
     * 格式优惠券
     *
     * @param $coupons
     *
     * @return array
     */
    public static function format($coupons)
    {
        $list = [];
        foreach ($coupons as $userCoupon) {
            $coupon = $userCoupon->coupon;
            $order  = $userCoupon->order;
            switch ($userCoupon->effective) {
                case self::COUPON_LIMIT_NO:
                    $item = [
                        'id'                      => $userCoupon->id,
                        'status'                  => $userCoupon->status,
                        'order_no'                => $userCoupon->order_no ?: '',
                        'order_at'                => (self::COUPON_HAVE_USED == $userCoupon->status && $order) ? date('Y.m.d H:i', $order->order_at) : '',
                        'effective'               => $userCoupon->effective,
                        'use_time'                => self::COUPON_HAVE_USED == $userCoupon->status ? date('Y.m.d H:i', $userCoupon->use_time) : '',
                        'order_amount'            => (self::COUPON_HAVE_USED == $userCoupon->status ? ($userCoupon->order_amount / 100) : 0) . '元',
                        'receive_time'            => date('Y.m.d H:i', $userCoupon->receive_time),
                        'preferential_rule'       => '满' . ($coupon->coupon_minimum / 100) . '减' . ($coupon->coupon_value / 100),
                        'preferential_price'      => ($coupon->coupon_value / 100) . '元',
                        'preferential_time_limit' => '长期有效',
                    ];
                    break;
                case self::COUPON_LIMIT_TINE;
                    $item = [
                        'id'                      => $userCoupon->id,
                        'status'                  => self::COUPON_WAIT_USE == $userCoupon->status ? (time() >= $userCoupon->end_time ? self::COUPON_HAVE_EXPIRED : $userCoupon->status) : $userCoupon->status,
                        'order_at'                => (self::COUPON_HAVE_USED == $userCoupon->status && $order) ? date('Y.m.d H:i', $order->order_at) : '',
                        'order_no'                => $userCoupon->order_no ?: '',
                        'use_time'                => self::COUPON_HAVE_USED == $userCoupon->status ? date('Y.m.d H:i', $userCoupon->use_time) : '',
                        'order_amount'            => (self::COUPON_HAVE_USED == $userCoupon->status ? ($userCoupon->order_amount / 100) : 0) . '元',
                        'effective'               => $userCoupon->effective,
                        'receive_time'            => date('Y.m.d H:i', $userCoupon->receive_time),
                        'preferential_rule'       => '满' . ($coupon->coupon_minimum / 100) . '减' . ($coupon->coupon_value / 100),
                        'preferential_price'      => ($coupon->coupon_value / 100) . '元',
                        'preferential_time_limit' => date('Y.m.d', $userCoupon->start_time) . '-' . date('Y.m.d', $userCoupon->end_time),
                    ];
                    break;
                default:
                    $item = [];
            }
            $item && $list[] = $item;
        }
        return $list;
    }
}
