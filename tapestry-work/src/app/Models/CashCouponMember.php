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


class CashCouponMember extends Base
{
    const MEMBER_IS_NORMAL = 0; // 员工正常

    const MEMBER_IS_LEAVE = 1; // 员工离职

    /**
     * 主键
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * 表名
     * @var string
     */
    protected $table = 'cash_coupon_member';

    /**
     * 白名单
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'cash_coupon_id',
        'object',
        'is_leave'
    ];
}