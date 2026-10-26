<?php
/**
 * Created by PhpStorm.
 * User: wangce
 * Date: 2018/11/15
 * Time: 2019-02-20 09:48
 */

namespace App\Models;

class OperationLog extends Base
{
    /**
     * 服务权限类型 0:全部
     */
    const SERVE_TYPE_ALL = 0;

    /**
     * 服务权限类型 1:shop
     */
    const SERVE_TYPE_SHOP = 1;

    /**
     * 服务权限类型 2:crs
     */
    const SERVE_TYPE_CRS = 2;

    /**
     * 操作日志记录表
     *
     * @var string
     */
    protected $table = 'operation_log';

    /**
     * @var array
     */
    protected $fillable = [
        'enterprise_id',
        'admin_id',
        'serve_type',
        'admin_name',
        'module',
        'page',
        'content',
        'ip',
        'http_method',
        'http_params',
        'user_agent'
    ];

}
