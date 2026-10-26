<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/15
 * Time: 4:05 PM
 */

namespace App\Models;

class EnterpriseAppInstall extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_app_install';

    /**
     * 允许写入的字段
     *
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'app_id',
        'identity'
    ];
}
