<?php
/**
 * Created by PhpStorm.
 * User: wangce
 * Date: 2018/11/15
 * Time: 2019-02-20 09:48
 */

namespace App\Models;

class Role extends Base
{

    /**
     * 角色表
     *
     * @var string
     */
    protected $table = 'role';

    /**
     * 类型:自然流量
     */
    const TYPE_OTHER = 1;

    /**
     * 类型:门店
     */
    const TYPE_STORE = 2;

    /**
     * 类型:自然流量+门店
     */
    const TYPE_ALL = 3;

    /**
     * @var array
     */
    protected $fillable = [
        'enterprise_id',
        'name',
        'describe',
        'shop_id',
        'available',
        'other'
    ];

    /**
     * 关联角色门店关系表 一对多
     *
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     */
    public function stores()
    {
        return $this->hasMany('App\Models\RoleStore', 'role_id', 'id');
    }

}
