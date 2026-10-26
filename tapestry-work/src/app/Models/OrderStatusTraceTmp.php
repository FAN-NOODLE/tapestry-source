<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class OrderStatusTraceTmp extends Base
{
    protected $table = 'order_status_trace_tmp';

    protected $primaryKey = 'id';

    protected $fillable = [
        'id',
        'enterprise_id',
        'order_no',
        'event',
        'event_desc',
        'from',
        'to',
        'change_at',
        'created_at',
        'updated_at',
    ];

    /**
     * @var array
     */
    public static $selectField = ['id', 'enterprise_id', 'order_no', 'event', 'event_desc', 'from', 'to', 'change_at'];
}
