<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/23
 * Time: 17:35
 */

namespace App\Models;

class CommissionRecord extends Base
{
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'commission_record';

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
        'mem_uid',
        'amount_price',
        'order_id',
        'charge_price',
        'user_id'
    ];
}
