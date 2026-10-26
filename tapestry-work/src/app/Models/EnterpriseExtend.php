<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/10
 * Time: 16:28
 */

namespace App\Models;

class EnterpriseExtend extends Base
{
    const OPEN  = 1;
    const CLOSE = 2;

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'enterprise_extend';

    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'allocation_status',
        'cash_coupon_status',
    ];

    /**
     * @var array
     */
    public static $selectField = ['id', 'enterprise_id', 'allocation_status', 'cash_coupon_status'];
}
