<?php
/*
 * Created by PhpStorm.
 * User: wangce
 * Date Time: 2019/2/25 9:38
 */

namespace App\Models;


class ApiPower extends Base
{

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'api_power';

    /**
     * @var array
     */
    protected $fillable = [
        'pid',
        'bap_id',
        'route',
        'describe'
    ];

}
