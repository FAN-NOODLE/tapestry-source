<?php


namespace App\Models;


class EnterprisePayment extends Base
{
    /**
     * 购买类目，0：红包，1：短信，2：智慧导购，3：SCRM，4：商城，5：综合系统
     */
    const PAYMENT_CATEGORY_RED_PACK = 0;
    const PAYMENT_CATEGORY_SMS = 1;
    const PAYMENT_CATEGORY_ZHDG = 2;
    const PAYMENT_CATEGORY_SCRM = 3;
    const PAYMENT_CATEGORY_MALL = 4;
    const PAYMENT_CATEGORY_SYSTEM = 5;

    /**
     * 收款状态，0：待收款，1：已收款
     */

    const PAYMENT_STATUS_NO = 0;
    const PAYMENT_STATUS_YES = 1;
    /**
     * 支付方式，0：其他，1：现金，2：微信，3：支付宝，4：银联，5：对公转账
     */

    const PAYMENT_TYPE_OTHER = 0;
    const PAYMENT_TYPE_CASH = 1;
    const PAYMENT_TYPE_WX = 2;
    const PAYMENT_TYPE_ZFB = 3;
    const PAYMENT_TYPE_BANK = 4;
    const PAYMENT_TYPE_TRANSFER = 5;

    const PAYMENT_TYPE_ARR = [
        self::PAYMENT_TYPE_OTHER    => '其他',
        self::PAYMENT_TYPE_CASH     => '现金',
        self::PAYMENT_TYPE_WX       => '微信',
        self::PAYMENT_TYPE_ZFB      => '支付宝',
        self::PAYMENT_TYPE_BANK     => '银联',
        self::PAYMENT_TYPE_TRANSFER => '对公转账',
    ];
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_payment';

    protected $guarded = [];

}