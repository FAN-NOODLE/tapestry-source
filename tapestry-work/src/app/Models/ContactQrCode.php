<?php

namespace App\Models;

class ContactQrCode extends Base
{
    /**
     * 数据来源
     */
    const SOURCE_TYPE_DEFAULT = 0; //默认正常创建
    const SOURCE_TYPE_AUTO = 1; //内购员工自建
    const SOURCE_TYPE_AUTO_GUIDE = 3; //支付既导购

    /**
     * 联系人数量类型
     */
    const CODE_TYPE_ONE = 1; //单人
    const CODE_TYPE_MANY = 2; //多人

    /**
     * 是否需要验证
     */
    const VERIFY_NO = 1; //需要
    const VERIFY_YES = 0; //不需要

    /**
     * 表名
     * @var string
     */
    protected $table = 'contact_qr_code';

    /**
     * 主键
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'channel',
        'name',
        'enterprise_id',
        'type',
        'is_verify',
        'official_url',
        'transform_url',
        'uc_id',
        'config_id',
        'channel_id',
        'source_type'
    ];

    /**
     * 包含一个或多个员工参数
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     * @author: DLong
     */
    public function params()
    {
        return $this->hasMany(ContactQrCodeParams::class, 'contact_qr_code_id')
            ->select(['contact_qr_code_id', 'object']);
    }

    public function channel()
    {
        return $this->hasOne(Channel::class, 'id', 'channel_id');
    }
}
