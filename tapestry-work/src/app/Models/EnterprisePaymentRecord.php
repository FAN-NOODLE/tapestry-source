<?php


namespace App\Models;


class EnterprisePaymentRecord extends Base
{

    /**
     * 发送领取状态：1、发送中，2、已领取，3、已退回
     */
    const RECORD_STATUS_SENDING = 1;
    const RECORD_STATUS_SENDED = 2;
    const RECORD_STATUS_BACK = 3;

    protected $table = 'enterprise_payment_record';

    protected $guarded = [];

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