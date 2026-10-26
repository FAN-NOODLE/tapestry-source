<?php
/*
* +------------------------------------+
* | Created: phpstorm                  
* +------------------------------------+
* | Author : wyg                       
* +------------------------------------+
* | Date   : 2020/2/7                   
* +------------------------------------+
* | Time   : 14:52                   
* +------------------------------------+
* | Email  : <wangyongguo@vchangyi.com>
* +------------------------------------+
*/

namespace App\Models;

class CashCouponRecord extends Base
{
    const SENDING = 1; // 发送中
    const RECEIVED = 2; // 已领取
    const REFUND = 3; // 已退回

    /**
     * 发送中
     */
    const IS_SEND  = 1;

    /**
     * 已领取
     */
    const IS_RECEIVED  = 2;

    /**
     * 已退回
     */
    const IS_BACK  = 3;

    /**
     * 主键
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * 表名
     * @var string
     */
    protected $table = 'cash_coupon_record';

    /**
     * 白名单
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'cash_coupon_id',
        'bill_no',
        'detail_id',
        'mem_uid',
        'ec_userid',
        'open_id',
        'send_time',
        'receive_refund_time',
        'status'
    ];

    /**
     * 生成单号
     *
     * @return string
     */
    public static function generateBillNo()
    {
        $billNo = self::randLetter(2) . self::randNumber(12) . self::randLetter(1);

        /* 到数据库里查找是否已存在 */
        $record = self::where(['bill_no' => $billNo]);

        if ($record->count()) {
            /* 如果有重复的，则重新生成 */
            return self::generateBillNo();
        }

        return $billNo;
    }

    /**
     * 随机数生成
     *
     * @param $count
     * @return false|string
     */
    private static function randLetter($count)
    {
        $str = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
        $str = str_shuffle($str);
        $str = substr($str, 0, $count);
        return $str;
    }

    /**
     * 随机数生成
     *
     * @param $count
     * @return false|string
     */
    private static function randNumber($count)
    {
        $str = '123456789';
        $str = str_shuffle($str);
        $str = substr($str, 0, $count);
        return $str;
    }
}
