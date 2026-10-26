<?php

namespace App\Models\ScrmEnterprise;


class EnterpriseUser extends ScrmEnterpriseInit
{

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'enterprise_user';


    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $guarded = [];


    public static function updateEnterpriseUserById(int $enterpriseId, int $Id, array $updateData)
    {
        return static::query()
            ->setConn($enterpriseId)
            ->where('id', $Id)
            ->withTrashed()
            ->update($updateData);
    }

}

