<?php
/**
 * Created by PhpStorm.
 * User: wangce
 * Date: 2018/11/15
 * Time: 2019-02-20 09:48
 */

namespace App\Models;

class RolePermission extends Base
{

    /**
     * 权限类型 1:基础
     */
    const TYPE_BASE = 1;

    /**
     * 权限类型 2:应用
     */
    const TYPE_APP = 2;

    /**
     * 角色权限表
     *
     * @var string
     */
    protected $table = 'role_permission';

    /**
     * @var array
     */
    protected $fillable = [
        'role_id',
        'bap_id',
        'type',
        'enterprise_id'
    ];
}
