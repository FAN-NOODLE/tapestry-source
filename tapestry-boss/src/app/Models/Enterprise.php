<?php

namespace App\Models;

class Enterprise extends Base
{
    /**
     * 企业商城版本 1:标准版（旧版）,2:中心商城版(有门店), 3:中心商城(无门店), 4:直营连锁门店
     */
    const NORMAL_VERSION = 1;
    const VERSION_CENTER_SHOP_STORE = 2;
    const VERSION_CENTER_SHOP = 3;
    const VERSION_DIRECT_CHAIN_STORE = 4;

    const STATUS_ENABLE = 1;
    const STATUS_DISABLE = 2;

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise';

    protected $guarded = [];

    public static function searchEnterprise($id)
    {
        return self::query()->where('id', '=', $id)->first();
    }

    /**
     * 关联虚拟开放平台（一对一）
     *
     * @return \Illuminate\Database\Eloquent\Relations\HasOne
     */
    public function virtualOpenPlatform()
    {
        return $this->hasOne(EnterpriseVirtualOpenPlatform::class, 'enterprise_id');
    }

    /**
     * 关联第三方授权（一对多）
     *
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     */
    public function thirdPartAccounts()
    {
        return $this->hasMany(EnterpriseThirdPartAccount::class, 'enterprise_id');
    }

    public function getEditAtAttribute($value)
    {
        return date('Y-m-d H:i:s', $value);
    }

    public function getCreatedAtAttribute($value)
    {
        return date('Y-m-d H:i:s', $value);
    }

    public function getEnterpriseById($id)
    {
        return $this->getModelArray([
            'id' => $id
        ]);
    }

    public function getEnterpriseByHash($hash)
    {
        return $this->getModelArray([
            'hash' => $hash
        ]);
    }
}
