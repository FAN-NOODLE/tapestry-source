<?php

namespace App\Models;

class EnterpriseThirdPartAccount extends Base
{

    /**
     * 授权类型:1.服务号 2.小程序
     */
    const OFFICIAL_ACCOUNT = 1;
    const MINI_PROGRAM_ACCOUNT = 2;

    /**
     *  授权状态：1.有效 2.无效
     */
    const STATUS_ACTIVE = 1;
    const STATUS_NEGATIVE = 2;

    /**
     * 小程序类型：1.商城小程序，2.企业微信小程序
     */
    const TYPE_SHOP = 1;
    const TYPE_QY_WX = 2;


    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_third_part_account';

    protected $guarded = [];

    public function listMiniPrograms($enterpriseIds)
    {
        return self::query()
                ->whereIn('enterprise_id', $enterpriseIds)
                ->where('status', self::STATUS_ACTIVE)
                ->where('platform', self::MINI_PROGRAM_ACCOUNT)
                ->where('type', self::TYPE_SHOP)
                ->get()
                ->toArray() ?? [];
    }

    public function getMiniProgram($enterpriseId)
    {
        $obj = self::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('status', self::STATUS_ACTIVE)
            ->where('platform', self::MINI_PROGRAM_ACCOUNT)
            ->where('type', self::TYPE_SHOP)
            ->first();

        return empty($obj) ? [] : $obj->toArray();
    }

    public static function existsMiniProgram($enterpriseId)
    {
        return self::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('status', self::STATUS_ACTIVE)
            ->where('platform', self::MINI_PROGRAM_ACCOUNT)
            ->where('type', self::TYPE_SHOP)
            ->exists();
    }
}
