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

class CashCoupon extends Base
{
    const IS_UN_SALE = 1; // 未上架

    const IS_ON_SALE = 2; // 上架中

    const IS_OFF_SALE = 3; // 已下架

    const NOT_DELETED = 1; // 未删除

    const DELETED = 2; // 已删除

    /**
     * 主键
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * 表名
     * @var string
     */
    protected $table = 'cash_coupon';

    /**
     * 白名单
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'cash_coupon_name',
        'number',
        'amount',
        'quantity',
        'is_sale',
        'send_name',
        'wishing',
        'remark',
        'is_deleted'
    ];
}