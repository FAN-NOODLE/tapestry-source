<?php
/**
 * Created by PhpStorm.
 * Date: 2018/11/15
 * Time: 4:05 PM
 */

namespace App\Models;

class EnterpriseVirtualOpenPlatform extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_virtual_open_platform';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'open_app_id'
    ];
}
