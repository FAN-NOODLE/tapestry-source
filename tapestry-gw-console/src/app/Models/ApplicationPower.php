<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/20
 * Time: 16:11
 */
namespace App\Models;

class ApplicationPower extends Base
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
     * 是否可用 1:可用
     */
    const AVAILABLE_IS = 1;

    /**
     * 是否可用 0:不可用
     */
    const AVAILABLE_DISABLED = 0;

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'application_power';

    /**
     * @var array
     */
    protected $fillable = [
        'serve_type',
        'name',
        'alias',
        'component',
        'path',
        'parent',
        'available',
        'enterprise_id'
    ];
}
