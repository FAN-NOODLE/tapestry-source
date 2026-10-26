<?php


namespace App\Models;


class EnterpriseAdminRole extends Base
{
    protected $table = 'enterprise_admin_role';

    protected $fillable = [
        'id',
        'enterprise_id',
        'enterprise_admin_id',
        'role_id',
        'created_at',
        'updated_at',
        'deleted_at'
    ];

}