<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Relations\HasOne;

class EnterpriseAdmin extends Base
{
    /**
     * 是否超管  1-是  2-不是
     */
    const IS_SUPER = 1;
    const NO_SUPER = 2;

    /**
     * 账号状态 1-开启 2-禁用
     */
    const STATUS_OPEN = 1;
    const STATUS_CLOSE = 2;

    /**
     * 表名称
     * @var string
     */
    protected $table = 'enterprise_admin';

    protected $guarded = [];

    /**
     * 超级管理员
     */
    const SUPER_ADMIN = 0;

    public static function updateAvailableSystem(int $enterpriseId, array $data)
    {
        self::query()
            ->where('enterprise_id', '=', $enterpriseId)
            ->where('is_super', '=', self::IS_SUPER)
            ->update($data);
    }

    /**
     * 关联企业信息（一对一）
     * @return HasOne
     */
    public function enterprise()
    {
        return $this->hasOne(Enterprise::class, 'id', 'enterprise_id');
    }

    public function getUpdatedAtAttribute($value)
    {
        return date('Y-m-d H:i:s', $value);
    }

    public function getCreatedAtAttribute($value)
    {
        return date('Y-m-d H:i:s', $value);
    }
}