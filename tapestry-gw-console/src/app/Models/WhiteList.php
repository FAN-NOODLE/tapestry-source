<?php
/*
 * Created by PhpStorm.
 * User: wangce
 * Date Time: 2019/2/25 9:38
 */

namespace App\Models;


class WhiteList extends Base
{
    /**
     * 先享卡
     */
    const TYPE_CARD = 1;

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'white_list';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'type',
        'created_at',
        'updated_at',
        'deleted_at',
    ];

}
