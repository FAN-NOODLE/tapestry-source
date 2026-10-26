<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/9/10
 * Time: 11:48
 */

namespace App\Models;

class CrsOrderStatusTrace extends Base
{
    protected $table = 'crs_order_status_trace';

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
