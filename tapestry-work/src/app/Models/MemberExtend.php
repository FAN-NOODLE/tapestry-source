<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/10
 * Time: 16:28
 */

namespace App\Models;

class MemberExtend extends Base
{
    const OPEN          = 1;
    const CLOSE         = 2;
    const FORCE_OPEN    = 3;
    const FOREVER_CLOSE = 4;

    const AUTHORITY_POSITIVE = 1; // 有权限
    const AUTHORITY_NEGATIVE = 2; // 无权限

    const GUIDE = '导购';

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'member_extend';

    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'mem_uid',
        'role',
        'guide_id',
        'avatar',
        'allocation_status',
        'trade_guide_status',
        'limit_user_count',
        'config_id',
        'contact_qr_code',
    ];

    /**
     * @var array
     */
    public static $selectField = ['id', 'enterprise_id', 'mem_uid', 'guide_id', 'avatar', 'allocation_status', 'limit_user_count', 'config_id', 'contact_qr_code'];
}
