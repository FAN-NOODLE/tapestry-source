<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Relations\HasOne;

class EnterpriseThirdPartAccount extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_third_part_account';

    /**
     * 小程序的发版类型:1.稳定版 2.单独的项目
     *
     */
    const RELEASE_STABLE = 1;
    const RELEASE_PROJECT = 2;

    /**
     * 授权类型:1.服务号 2.小程序
     *
     */
    const OFFICIAL_ACCOUNT = 1;
    const MINI_PROGRAM_ACCOUNT = 2;

    /**
     * 微信载体类型
     */
    const MINI_PROGRAM = 'MINI_PROGRAM'; // 小程序
    const FU_WU_HAO = 'FU_WU_HAO';    // 服务号
    const DING_YUE_HAO = 'DING_YUE_HAO'; // 订阅号

    /**
     *  授权类型：1.手动授权 2.第三方授权
     *
     */
    const AUTHORIZE_TYPE_CUSTOMER = 1;
    const AUTHORIZE_TYPE_THIRD = 2;

    /**
     *  授权状态：1.有效 2.无效
     */
    const STATUS_ACTIVE = 1;
    const STATUS_NEGATIVE = 2;

    /**
     * 小程序类型：1.商城小程序，2.企业微信小程序 3.公众号
     */
    const TYPE_SHOP = 1;
    const TYPE_QY_WX = 2;
    const TYPE_OFFICIAL = 3;

    /**
     * 认证状态: 1.已认证 2.未认证
     */
    const VERIFIED = 1;
    const UNVERIFIED = 2;

    /**
     * 小程序版本标识集合
     */
    const MINI_PROGRAM_VERSION = [
        self::TYPE_SHOP_SHORT,
        self::TYPE_QY_WX_SHORT
    ];

    /**
     * 商城小程序标识
     */
    const TYPE_SHOP_SHORT = 's';

    /**
     * 企业微信小程序标识
     */
    const TYPE_QY_WX_SHORT = 'q';

    /**
     * 允许小程序类型
     */
    const ALL_ALLOW_TYPE = [
        self::TYPE_SHOP,
        self::TYPE_QY_WX,
    ];

    /**
     * 小程序审核状态
     */
    const AUDIT_STATUS = [
        1 => 'PENDING', // 等待审核
        2 => 'AUDITING', // 审核中
        3 => 'AUDITED', // 审核通过
        4 => 'REFUSED', // 审核失败
        5 => 'WITHDRAW' // 撤审
    ];

    /**
     * 小程序上线状态
     */
    const ONLINE_STATUS = [
        1 => 'PENDING', // 等待
        2 => 'ONLINE', // 上线
        3 => 'OFFLINE' // 下线
    ];

    /**
     * @var array
     */
    protected $guarded = [];

    /**
     * @param $enterpriseId
     * @return EnterpriseThirdPartAccount|Builder|Model|object|null
     */
    public static function getMiniProgram($enterpriseId)
    {
        return self::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('status', self::STATUS_ACTIVE)
            ->where('platform', self::MINI_PROGRAM_ACCOUNT)
            ->where('type', self::TYPE_SHOP)
            ->first();
    }

    /**
     * 小程序最近审核版本
     * @return HasOne
     */
    public function latestCodeLibrary()
    {
        return $this->hasOne(EnterpriseMiniProgramCodeLibrary::class, 'account_id', 'id')
            ->orderByDesc('id');
    }

    /**
     * 当前在用小程序版本
     * @return HasOne
     */
    public function usedCodeLibrary()
    {
        return $this->hasOne(EnterpriseMiniProgramCodeLibrary::class, 'account_id', 'id')
            ->where('using', EnterpriseMiniProgramCodeLibrary::USING_STATUS_TRUE);
    }

    public function enterprise()
    {
        return $this->belongsTo(Enterprise::class);
    }
}
