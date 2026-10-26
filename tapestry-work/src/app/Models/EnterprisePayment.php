<?php


namespace App\Models;


class EnterprisePayment extends Base
{

    // 红包类型 1=固定金额 2=随机金额
    const FIXED_AMOUNT_TYPE = 1;
    const RANDOM_AMOUNT_TYPE = 2;

    /**
     * 上架状态：1=未上架, 2=上架, 3=已下架,4=已结束
     */
    const IS_SALE_DEFAULT = 1;
    const IS_SALE_UP = 2;
    const IS_SALE_DOWN = 3;
    const IS_SALE_OVER = 4;

    /**
     * 删除状态 1 未删除 2 已删除
     */
    const IS_DELETED_NO = 1;
    const IS_DELETED_YES = 2;

    /**
     * 企业现有客户是否可以领取：1=否, 2=是
     */
    const IS_OLD_USER_RECEIVE_NO = 1;
    const IS_OLD_USER_RECEIVE_YES = 2;

    const RANDOM_LIST_KEY = "enterprise_payment_random_";
    const SEND_PACK_KEY = "enterprise_payment_send_pack_";
    const SEND_PACK_BY_USER_KEY = "enterprise_payment_send_pack_user_";

    protected $table = 'enterprise_payment';

    protected $guarded = [];

}