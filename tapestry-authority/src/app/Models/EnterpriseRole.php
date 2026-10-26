<?php

namespace App\Models;

class EnterpriseRole extends Base
{
    /**
     * 功能权限
     */
    const ROLE_TYPE_POWER = 1;

    /**
     * 数据权限
     */
    const ROLE_TYPE_DATA = 2;


    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_role';

    protected $fillable = [
        'id',
        'enterprise_id',
        'operator_id',
        'role_name',
        'type',
        'auth_rule',
        'identity_type',
        'created_at',
        'updated_at',
        'deleted_at'
    ];

    public function roleAdmins()
    {
        return $this->hasMany(EnterpriseAdminRole::class, 'role_id', 'id');
    }

    public function enterpriseAdmin()
    {
        return $this->belongsTo(EnterpriseAdmin::class, 'operator_id', 'id')
            ->select(['id', 'username']);
    }
}
