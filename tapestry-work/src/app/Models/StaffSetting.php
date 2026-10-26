<?php

namespace App\Models;

class StaffSetting extends Base
{
    /**
     * binding_type
     */
    const FOREVER = 1;
    const EXPIRE  = 2;

    /**
     * change_owner
     */
    const NOT_ALLOW = 1;
    const ALLOW     = 2;

    /**
     * commission_type
     */
    const HAS_COMMISSION     = 1;
    const NOT_HAS_COMMISSION = 2;

    /**
     * settle_type
     */
    const TRADE_SETTLE     = 1;
    const POST_SALE_SETTLE = 2;

    /**
     * type, 1 内购  2 导购
     */
    const STAFF_TYPE = 1;
    const GUIDE_TYPE = 2;

    /**
     * fission_type 客户分享是否统计导购业绩 1:不计算业绩 2:计算业绩
     * @var string
     */
     const FISSION_TYPE_FAIL = 1;
     const FISSION_TYPE_TRUE = 2;

    protected $table = 'staff_setting';

    protected $guarded = [];

    protected $primaryKey = 'id';

    protected $fillable = [
        'enterprise_id',
        'binding_type',
        'duration_time',
        'change_owner',
        'protected_time',
        'commission_type',
        'settle_type',
        'achievement_desc',
        'type',
        'fission_type',
        'created_at',
        'updated_at',
    ];
}
